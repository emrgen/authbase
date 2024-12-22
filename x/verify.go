package x

import (
	"context"
	"errors"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
)

type UserVerifier interface {
	VerifyEmailPassword(ctx context.Context, orgID uuid.UUID, email, password string) (*model.User, error)
	VerifyToken(ctx context.Context, token string) (*Claims, error)
}

type StoreBasedUserVerifier struct {
	store store.Provider
	redis *cache.Redis
}

func NewStoreBasedUserVerifier(store store.Provider, redis *cache.Redis) *StoreBasedUserVerifier {
	return &StoreBasedUserVerifier{
		store: store,
		redis: redis,
	}
}

func (v *StoreBasedUserVerifier) VerifyEmailPassword(ctx context.Context, orgID uuid.UUID, email, password string) (*model.User, error) {
	as, err := store.GetProjectStore(ctx, v.store)
	if err != nil {
		return nil, err
	}

	user, err := as.GetUserByEmail(ctx, orgID, email)
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

func (v *StoreBasedUserVerifier) VerifyToken(ctx context.Context, token string) (*Claims, error) {
	claims, err := VerifyJWTToken(token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
