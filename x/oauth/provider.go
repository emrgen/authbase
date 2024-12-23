package oauth

import (
	"context"
	"golang.org/x/oauth2"
)

type Provider interface {
	// GetName returns the name of the provider
	GetName() string
	// GetType returns the type of the provider
	GetType() string
	// GetToken returns the token from the provider
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
}
