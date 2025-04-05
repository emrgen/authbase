package oauth

import (
	"context"
	"errors"
	"golang.org/x/oauth2"
)

// Provider is an interface for OAuth providers
type Provider interface {
	// GetName returns the name of the provider
	GetName() string
	// GetType returns the type of the provider
	GetType() string
	// GetToken returns the token from the provider
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
}

// GetProvider returns the provider based on the name
func GetProvider(name string, config oauth2.Config) (Provider, error) {
	switch name {
	case "google":
		return NewGoogleProvider(config), nil
	default:
		return nil, ErrUnsupportedProvider
	}
}

// ErrUnsupportedProvider is returned when the provider is not supported
var ErrUnsupportedProvider = errors.New("unsupported provider")
