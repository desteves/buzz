// Package ai provides functionality for AI operations.
package ai

import (
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

func ParseResponse(resp *genai.GenerateContentResponse) (string, error) {

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
