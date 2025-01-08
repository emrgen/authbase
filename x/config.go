package x

import (
	"context"
	"google.golang.org/grpc/metadata"
	"os"
)

type AuthbaseConfig struct {
	AccessKey string
}

func ConfigFromEnv() (*AuthbaseConfig, error) {
	// load tiny config
	authbaseKey := os.Getenv("AUTHBASE_KEY")

	return &AuthbaseConfig{
		AccessKey: authbaseKey,
	}, nil
}

// IntoContext creates a new context with the TinyAPIKey in the metadata.
func (p *AuthbaseConfig) IntoContext() context.Context {
	// create a new context
	ctx := context.Background()

	md := metadata.New(map[string]string{"Authorization": "Bearer " + p.AccessKey})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	return ctx
}
