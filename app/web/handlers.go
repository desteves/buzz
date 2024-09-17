package web

import (
	"encoding/json"
	"fmt"
	"log"
	"main/ai"
	"main/oauth"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type MuxHandler struct {
	Handler       *http.ServeMux
	AIProvider    ai.Provider
	OAuthProvider oauth.Provider
}

func NewHandler() MuxHandler {

	var m MuxHandler

	// log.Printf("Adding routes")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static/")))
	mux.HandleFunc("/auth/google/login", m.loginHandler)
	mux.HandleFunc("/app", m.appHandler)
	mux.HandleFunc("/submit", m.submitHandler)

	m.Handler = mux
	m.AIProvider = ai.NewGeminiProvider()
	m.OAuthProvider = oauth.NewGoogleProvider()

	return m
}

func (m *MuxHandler) submitHandler(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	var data BuzzForm
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Get the values from the parsed JSON
	a := data.A
	p := data.P

	input := CleanInput(p)
	// log.Print(input)

	// // Validate the access token
	// Define the validateToken function or import it from another package
	_, err := m.OAuthProvider.ValidateToken(a)
	if err != nil {
		log.Printf("invalid access token: %s", a)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	html, err := m.AIProvider.GenerateResponse(input)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, html)

}

func (m *MuxHandler) loginHandler(w http.ResponseWriter, r *http.Request) {
	oauthState := m.OAuthProvider.GenerateCookie(w)
	u := m.OAuthProvider.OAuthConfig().AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (m *MuxHandler) appHandler(w http.ResponseWriter, r *http.Request) {
	data, err := m.OAuthProvider.GetData(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Step 1: Read the HTML template
	htmlContent, err := os.ReadFile("./static/app.html")
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Step 2: Parse the HTML
	doc, err := html.Parse(strings.NewReader(string(htmlContent)))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Step 3: Update the content of the div with a specific ID
	targetID := "profile"
	newContent := data.Name
	UpdateDivContent(doc, targetID, newContent)

	targetID = "a"
	newContent = data.AccessToken
	UpdateFormContent(doc, targetID, newContent)

	// Step 4: Convert the updated HTML structure back to a string
	var sb strings.Builder
	err = html.Render(&sb, doc)
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}

	// Set the HTTP response headers
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Write the updated HTML content as the HTTP response
	fmt.Fprint(w, sb.String())

}
