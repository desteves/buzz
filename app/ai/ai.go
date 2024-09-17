// Package ai provides functionality for AI operations.
package ai

type Provider interface {
	GenerateResponse(prompt string) (string, error)
}
