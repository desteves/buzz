package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Google struct {
	OAuthAPI string
}

func NewGoogleProvider() *Google {
	return &Google{
		OAuthAPI: "https://www.googleapis.com/oauth2/v3/userinfo?access_token=",
	}
}

// TODO - clean this up inside interface class

// getEnv function to return the value of an environment variable or a default value if not set
func getEnv(envVarName string, defaultValue string) string {
	value, exists := os.LookupEnv(envVarName)
	if !exists {
		return defaultValue
	}
	return value
}

// var googleOauthConfig = &oauth2.Config{

func (g *Google) OAuthConfig() *oauth2.Config {

	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
		RedirectURL:  getEnv("REDIR", "http://localhost:8000/app"),
	}
}

func (g *Google) ValidateToken(token string) (bool, error) {
	// Construct the tokeninfo URL with the access token
	tokeninfoURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=%s", token)

	// Make a GET request to the tokeninfo endpoint
	resp, err := http.Get(tokeninfoURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	// Parse the JSON response
	var tokenInfo map[string]interface{}
	err = json.Unmarshal(body, &tokenInfo)
	if err != nil {
		return false, err
	}

	// Check if the token is valid
	if resp.StatusCode == http.StatusOK {
		return true, nil
	} else {
		errorDescription := tokenInfo["error_description"].(string)
		return false, fmt.Errorf("token validation failed: %s", errorDescription)
	}
}

func (g *Google) GenerateCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func (g *Google) GetData(code string) (Data, error) {
	// Use code to get token and get user info from Google.

	data := Data{}
	token, err := g.OAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		return data, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	// Extract the OAuth2 token
	accessToken := token.AccessToken

	response, err := http.Get(g.OAuthAPI + token.AccessToken)
	if err != nil {
		return data, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return data, fmt.Errorf("failed read response: %s", err.Error())
	}
	err = json.Unmarshal(contents, &data)
	if err != nil {
		return data, fmt.Errorf("invalid response: %s", err.Error())
	}
	data.AccessToken = accessToken
	return data, nil
}
