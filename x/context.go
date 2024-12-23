package x

import (
	"context"
	"github.com/google/uuid"
)

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value("userID").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrUserNotFoundInContext
	}

	return userID, nil
}

func GetOrganizationID(ctx context.Context) (uuid.UUID, error) {
	organizationID, ok := ctx.Value("organizationID").(uuid.UUID)
	if !ok {
		return uuid.UUID{}, ErrOrganizationNotFoundInContext
	}

	return organizationID, nil
}

func GetOAuthState(ctx context.Context) (string, error) {
	state, ok := ctx.Value("oauthstate").(string)
	if !ok {
		return "", ErrOAuthStateNotFoundInContext
	}

	return state, nil
}
