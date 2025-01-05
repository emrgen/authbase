package service

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"go.uber.org/ratelimit"
	"os"
	"time"
)

var _ v1.AdminProjectServiceServer = (*AdminProjectService)(nil)

type AdminProjectService struct {
	provider store.Provider
	cache    *cache.Redis
	limited  ratelimit.Limiter
	v1.UnimplementedAdminProjectServiceServer
}

// NewAdminProjectService creates a new admin project service
func NewAdminProjectService(store store.Provider, cache *cache.Redis) v1.AdminProjectServiceServer {
	return &AdminProjectService{provider: store, cache: cache, limited: ratelimit.New(1)}
}

// CreateAdminProject creates a new project
func (a *AdminProjectService) CreateAdminProject(ctx context.Context, request *v1.CreateAdminProjectRequest) (*v1.CreateAdminProjectResponse, error) {
	// if the app is running in master mode, return an error as this operation is not allowed
	if os.Getenv("APP_MODE") == "multistore" {
		return nil, x.ErrForbidden
	}

	// rate limit the request
	a.limited.Take()

	as := a.provider.Default()

	// check if the master org already exists
	org, err := as.GetMasterProject(ctx)
	if err != nil && !errors.Is(err, store.ErrProjectNotFound) {
		return nil, err
	}

	if org != nil {
		return nil, x.ErrProjectExists
	}

	user := model.User{
		ID:        uuid.New().String(),
		Email:     request.GetEmail(),
		Username:  request.GetUsername(),
		SassAdmin: true,
		Member:    true,
	}

	org = &model.Project{
		ID:      uuid.New().String(),
		Name:    request.GetName(),
		OwnerID: user.ID,
		Master:  true,
	}
	user.ProjectID = org.ID

	perm := model.ProjectMember{
		ProjectID:  org.ID,
		UserID:     user.ID,
		Permission: uint32(v1.Permission_OWNER),
	}

	// Create project and user in a transaction
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		err := tx.CreateProject(ctx, org)
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

		// if password is provided, hash it and provider it
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

		err = tx.CreateProjectMember(ctx, &perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateAdminProjectResponse{
		Id: org.ID,
	}, nil
}

// CreateMigration creates a new migration for the project
func (a *AdminProjectService) CreateMigration(ctx context.Context, request *v1.CreateMigrationRequest) (*v1.CreateMigrationResponse, error) {
	projectID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	as, err := a.provider.Provide(projectID)
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
