// Package oauth provides authentication functionality.
package oauth

import (
	"net/http"

	"golang.org/x/oauth2"
)

type Data struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	Hd            string `json:"hd"`
	AccessToken   string `json:"accessToken"`
}

type Provider interface {
	OAuthConfig() *oauth2.Config
	ValidateToken(token string) (bool, error)
	GenerateCookie(w http.ResponseWriter) string
	GetData(token string) (Data, error)
}
