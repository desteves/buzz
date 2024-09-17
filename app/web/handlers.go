package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/ai"
	"main/auth"
	"net/http"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"golang.org/x/net/html"
	"google.golang.org/api/option"
)

func New() http.Handler {
	log.Printf("Adding routes")
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./static/")))
	mux.HandleFunc("/auth/google/login", loginHandler)
	mux.HandleFunc("/app", appHandler)
	mux.HandleFunc("/submit", submitHandler)

	return mux
}

func submitHandler(w http.ResponseWriter, r *http.Request) {

	// Read the request body
	var data BuzzForm
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Get the values from the parsed JSON
	a := data.A
	p := data.P

	input := cleanInput(p)
	// log.Print(input)

	// // Validate the access token
	// Define the validateToken function or import it from another package
	_, err := auth.ValidateToken(a)
	if err != nil {
		log.Printf("invalid access token: %s", a)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Initialize the AI client
	// new comment
	var ctx = context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-pro")

	resp, err := model.GenerateContent(ctx, genai.Text("Use the NATO phonetic alphabet to spell out this phrase: "+input+". Give me the response as an HTML table such that the first column is the letter and the second column is the word from the NATO phonetic alphabet."))
	if err != nil {
		log.Fatal(err)
	}
	html, err := ai.ParseResponse(resp)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprint(w, html)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	oauthState := auth.GenerateStateOauthCookie(w)
	u := auth.GoogleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	data, err := auth.GetUserDataFromGoogle(r.FormValue("code"))
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
	updateDivContent(doc, targetID, newContent)

	targetID = "a"
	newContent = data.AccessToken
	updateFormContent(doc, targetID, newContent)

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
