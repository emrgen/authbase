package x

import (
	"context"
	"errors"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
	"strings"
)

// UserVerifier is an interface to verify the user.
type UserVerifier interface {
	// VerifyEmailPassword verifies the email and password of a user.
	VerifyEmailPassword(ctx context.Context, orgID uuid.UUID, email, password string) (*model.Account, error)
	// VerifyToken verifies the token.
	VerifyToken(ctx context.Context, token string) (*Claims, error)
	// VerifyAccessKey verifies the access key.
	VerifyAccessKey(ctx context.Context, id uuid.UUID, accessKey string) (*Claims, error)
}

// StoreBasedUserVerifier is a user verifier that uses the store to verify the user.
type StoreBasedUserVerifier struct {
	store store.Provider
	redis *cache.Redis
}

// NewStoreBasedUserVerifier creates a new StoreBasedUserVerifier.
func NewStoreBasedUserVerifier(store store.Provider, redis *cache.Redis) *StoreBasedUserVerifier {
	return &StoreBasedUserVerifier{
		store: store,
		redis: redis,
	}
}

// VerifyEmailPassword verifies the email and password of a user.
// It returns the user if the email and password are correct.
func (v *StoreBasedUserVerifier) VerifyEmailPassword(ctx context.Context, projectID uuid.UUID, email, password string) (*model.Account, error) {
	as, err := store.GetProjectStore(ctx, v.store)
	if err != nil {
		return nil, err
	}

	user, err := as.GetAccountByEmail(ctx, projectID, email)
	if err != nil {
		return nil, err
	}

	if user.Disabled {
		return nil, errors.New("user account is disabled")
	}

	ok := CompareHashAndPassword(user.Password, password, user.Salt)
	if !ok {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

// VerifyToken verifies the token.
func (v *StoreBasedUserVerifier) VerifyToken(ctx context.Context, token string) (*Claims, error) {
	claims, err := VerifyJWTToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// VerifyAccessKey verifies the access key.
func (v *StoreBasedUserVerifier) VerifyAccessKey(ctx context.Context, id uuid.UUID, key string) (*Claims, error) {
	as, err := store.GetProjectStore(ctx, v.store)
	if err != nil {
		return nil, err
	}

	accessKey, err := as.GetAccessKeyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if accessKey.Token != key {
		return nil, errors.New("invalid access key")
	}

	claims := &Claims{
		ProjectID: accessKey.ProjectID,
		AccountID: accessKey.AccountID,
	}

	if accessKey.Scopes != "" {
		claims.Scopes = strings.Split(accessKey.Scopes, ",")
	}

	return claims, nil
}

// NoOpUserVerifier is a user verifier that does nothing.
type NoOpUserVerifier struct {
}
