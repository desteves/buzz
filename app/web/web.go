// Package web provides functionality for web applications.
package web

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

type BuzzForm struct {
	A string `json:"a"`
	P string `json:"p"`
}

func CleanInput(input string) string {
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

func UpdateFormContent(node *html.Node, targetID string, newContent string) {
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
		UpdateFormContent(child, targetID, newContent)
	}
}

// UpdateDivContent finds the div with a given ID and updates its content
func UpdateDivContent(node *html.Node, targetID string, newContent string) {
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
		UpdateDivContent(child, targetID, newContent)
	}
}
