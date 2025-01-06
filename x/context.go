package x

import (
	"context"
)

func GetOAuth2State(ctx context.Context) (string, error) {
	state, ok := ctx.Value("oauthstate").(string)
	if !ok {
		return "", ErrOAuthStateNotFoundInContext
	}

	return state, nil
}
