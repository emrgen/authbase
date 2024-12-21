package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"time"
)

var _ v1.AdminOrganizationServiceServer = (*AdminOrganizationService)(nil)

type AdminOrganizationService struct {
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedAdminOrganizationServiceServer
}

// NewAdminOrganizationService creates a new admin organization service
func NewAdminOrganizationService(store store.Provider, cache *cache.Redis) v1.AdminOrganizationServiceServer {
	return &AdminOrganizationService{store: store, cache: cache}
}

// CreateAdminOrganization creates a new organization
func (a *AdminOrganizationService) CreateAdminOrganization(ctx context.Context, request *v1.CreateAdminOrganizationRequest) (*v1.CreateAdminOrganizationResponse, error) {
	as := a.store.Default()

	user := model.User{
		ID:        uuid.New().String(),
		Email:     request.GetEmail(),
		Username:  request.GetUsername(),
		SassAdmin: true,
	}

	org := model.Organization{
		ID:      uuid.New().String(),
		Name:    request.GetName(),
		OwnerID: user.ID,
		Master:  true,
	}
	user.OrganizationID = org.ID

	// Create organization and user in a transaction
	err := as.Transaction(func(tx store.AuthBaseStore) error {
		err := tx.CreateOrganization(ctx, &org)
		if err != nil {
			return err
		}

		// if the mail server is configured, send a verification email anyway
		// verification email will be sent if the user is created successfully
		if request.GetVerifyEmail() {
			verificationCode := x.GenerateCode()
			err := a.cache.Set("email:"+user.Email, verificationCode, time.Hour)
			if err != nil {
				return err
			}
			defer func() {
				if err == nil {
					// send verification email
				}
			}()
		}

		// if password is provided, hash it and store it
		password := request.GetPassword()
		if password != "" {
			secret := x.Keygen()
			hash, err := x.HashPassword(password, secret)
			if err != nil {
				return err
			}

			user.Password = string(hash)
			user.Salt = secret
		}

		err = tx.CreateUser(ctx, &user)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateAdminOrganizationResponse{
		Id: org.ID,
	}, nil
}

// CreateMigration creates a new migration for the project
func (a *AdminOrganizationService) CreateMigration(ctx context.Context, request *v1.CreateMigrationRequest) (*v1.CreateMigrationResponse, error) {
	projectID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	as, err := a.store.Provide(projectID)
	if err != nil {
		return nil, err
	}

	err = as.Migrate()
	if err != nil {
		return nil, err
	}

	return &v1.CreateMigrationResponse{
		Message: "Migration successful",
	}, nil
}
