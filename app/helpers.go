package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/net/html"
)

// Helper functions

func getEnvOrDefault(envVarName, defaultValue string) string {
	value := os.Getenv(envVarName)
	if value == "" {
		log.Printf("Environment variable %s not set, using default value %s", envVarName, defaultValue)
		return defaultValue
	}
	return value
}

func parseResponse(resp *genai.GenerateContentResponse) (string, error) {

	var formattedContent string

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				formattedContent += fmt.Sprintf("%s\n", part)
			}
		} else {
			return "", fmt.Errorf("no content found in response")
		}
	}
	return formattedContent, nil
}

func cleanInput(input string) string {
	// Define a regex pattern to match non-alphanumeric characters and diacritics
	pattern := regexp.MustCompile(`[^a-zA-Z\\s]+`)

	// Replace non-alphanumeric characters and diacritics with an empty string
	cleaned := pattern.ReplaceAllStringFunc(input, func(s string) string {
		var result []rune
		for _, r := range s {
			// Check if the rune is alphanumeric or whitespace
			if unicode.IsLetter(r) || unicode.IsSpace(r) {
				result = append(result, r)
			}
		}
		return string(result)
	})

	// Remove extra whitespaces
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return cleaned
}

func updateFormContent(node *html.Node, targetID string, newContent string) {
	if node.Type == html.ElementNode && node.Data == "input" {
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == targetID {
				// Set the new default value
				for j, innerAttr := range node.Attr {
					if innerAttr.Key == "value" {
						node.Attr[j].Val = newContent // Update existing value
						return
					}
				}
				// If 'value' attribute doesn't exist, add it
				node.Attr = append(node.Attr, html.Attribute{
					Key: "value",
					Val: newContent,
				})
				return
			}
		}
	}

	// Recurse through child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		updateFormContent(child, targetID, newContent)
	}
}

// Find the div with a given ID and update its content
func updateDivContent(node *html.Node, targetID string, newContent string) {
	if node.Type == html.ElementNode && node.Data == "div" {
		// Check if the div has the desired ID
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == targetID {
				// Set the new content
				node.FirstChild = &html.Node{
					Type: html.TextNode,
					Data: newContent,
				}
				return
			}
		}
	}

	// Recurse through child nodes
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		updateDivContent(child, targetID, newContent)
	}
}

func validateToken(accessToken string) (bool, error) {
	// Construct the tokeninfo URL with the access token
	tokeninfoURL := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/tokeninfo?access_token=%s", accessToken)

	// Make a GET request to the tokeninfo endpoint
	resp, err := http.Get(tokeninfoURL)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
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

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) (Profile, error) {
	// Use code to get token and get user info from Google.
	host := getEnvOrDefault("REDIR", "http://localhost:8000")
	redirectUrl := host + "/app"
	googleOauthConfig.RedirectURL = redirectUrl
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return Profile{}, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	// Extract the OAuth2 token
	accessToken := token.AccessToken

	response, err := http.Get(OAUTH_API + token.AccessToken)
	if err != nil {
		return Profile{}, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return Profile{}, fmt.Errorf("failed read response: %s", err.Error())
	}
	myProfile := Profile{}
	err = json.Unmarshal(contents, &myProfile)
	if err != nil {
		return Profile{}, fmt.Errorf("invalid response: %s", err.Error())
	}
	myProfile.AccessToken = accessToken
	return myProfile, nil
}
