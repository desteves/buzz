package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func oauthGoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func oauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// fmt.Println(r.FormValue("code"))
	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....
	// fmt.Printf("UserInfo: %s\n", data)
	w.Write(data)

	// // Read the HTML template file
	// tmplFile, err := ioutil.ReadFile("./templates/welcome.html")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// // Parse the template
	// t := template.Must(template.New("html").Parse(string(tmplFile)))

	// UserInfo: {
	// 	"id": "105630908885014381855",
	// 	"email": "diana@pulumi.com",
	// 	"verified_email": true,
	// 	"name": "Diana Esteves",
	// 	"given_name": "Diana",
	// 	"family_name": "Esteves",
	// 	"picture": "https://lh3.googleusercontent.com/a/ACg8ocIeQugr0fxzdXtYdxVezUcoVwXrJFis_WuRz-hI28zbvQ=s96-c",
	// 	"locale": "en",
	// 	"hd": "pulumi.com"
	//   }

	// Create a data structure to pass values to the template
	// Execute the template with the data and write the result to the ResponseWriter
	// err = t.Execute(w, data)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
