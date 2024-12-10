package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/google/uuid"
)

var _ v1.AuthServiceServer = new(AuthService)

// AuthService is a service that implements the AuthServiceServer interface
type AuthService struct {
	store store.AuthBaseStore
	cache *cache.Redis
	v1.UnimplementedAuthServiceServer
}

func (a *AuthService) CheckUserAlreadyExists(ctx context.Context, request *v1.CheckUserAlreadyExistsRequest) (*v1.CheckEmailAlreadyExistsResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	users, err := a.store.UserExists(ctx, orgID, request.GetUsername(), request.GetEmail())
	var emailExists bool
	var usernameExists bool
	for _, user := range users {
		if user.Email == request.GetEmail() {
			emailExists = true
		}

		if user.Username == request.GetUsername() {
			usernameExists = true
		}
	}

	return &v1.CheckEmailAlreadyExistsResponse{
		EmailExists:    emailExists,
		UsernameExists: usernameExists,
	}, nil

}

// Register creates a new user if the username and email are unique
func (a *AuthService) Register(ctx context.Context, request *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthService) Login(ctx context.Context, request *v1.LoginRequest) (*v1.LoginResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthService) Refresh(ctx context.Context, request *v1.RefreshRequest) (*v1.RefreshResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthService) VerifyEmail(ctx context.Context, request *v1.VerifyEmailRequest) (*v1.VerifyEmailResponse, error) {
	// check if the email is already verified
	// if it is, return an error
	return nil, nil
}
