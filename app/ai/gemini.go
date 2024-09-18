package ai

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Gemini struct {
	// Model *genai.GenerativeModel
	Client *genai.Client
}

func NewGeminiProvider() *Gemini {

	// Initialize the AI client
	var ctx = context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	return &Gemini{Client: client}

}

func (g *Gemini) GenerateResponse(prompt string) (string, error) {

	model := g.Client.GenerativeModel("gemini-pro")

	resp, err := model.GenerateContent(context.Background(), genai.Text("Use the NATO phonetic alphabet to spell out this phrase: "+prompt+". Give me the response as an HTML table such that the first column is the letter and the second column is the word from the NATO phonetic alphabet."))
	if err != nil {
		log.Fatal(err)
	}
	html, err := parseResponse(resp)
	if err != nil {
		log.Fatal(err)
	}
	return html, nil

}

func parseResponse(resp *genai.GenerateContentResponse) (string, error) {
	var formattedContent string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				formattedContent += fmt.Sprintf("%s", part)
			}
		} else {
			return "", fmt.Errorf("no content found in response")
		}
	}
	return formattedContent, nil
}
