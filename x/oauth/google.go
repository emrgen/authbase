package oauth

import (
	"context"
	"golang.org/x/oauth2"
)

var _ Provider = new(GoogleProvider)

type GoogleProvider struct {
	config oauth2.Config
}

func NewGoogleProvider() *GoogleProvider {
	return &GoogleProvider{}
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
