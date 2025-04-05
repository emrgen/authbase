package oauth

import (
	"context"
	"golang.org/x/oauth2"
)

var _ Provider = new(GoogleProvider)

// GoogleProvider is an implementation of the Provider interface for Google OAuth2
type GoogleProvider struct {
	config oauth2.Config
}

// NewGoogleProvider creates a new GoogleProvider instance
func NewGoogleProvider(config oauth2.Config) *GoogleProvider {
	return &GoogleProvider{
		config: config,
	}
}

func (g *GoogleProvider) GetName() string {
	return "google"
}

func (g *GoogleProvider) GetType() string {
	return "oauth2"
}

func (g *GoogleProvider) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	return g.config.Exchange(ctx, code)
}
