package auth

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func NewOauthConfig(issuerURL, clientID, clientSecret, redirectURL string) (*oauth2.Config, *oidc.Provider) {
	provider, err := oidc.NewProvider(context.TODO(), issuerURL)
	if err != nil {
		fmt.Println("Error creating auth provider")
	}
	oauth2Config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  redirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "groups", "openid", "email"},
	}
	return &oauth2Config, provider
}
