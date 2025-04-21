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
	"github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
	"os"
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

	// check if the master project already exists
	project, err := as.GetMasterProject(ctx)
	if !errors.Is(err, store.ErrMasterProjectNotFound) && err != nil {
		return nil, err
	}
	err = nil // ignore the error if the project is not found
	if project != nil {
		return nil, x.ErrProjectExists
	}

	account := model.Account{
		ID:            uuid.New().String(),
		Email:         request.GetEmail(),
		VisibleName:   request.GetVisibleName(),
		SassAdmin:     true,
		ProjectMember: true,
	}

	project = &model.Project{
		ID:      uuid.New().String(),
		Name:    request.GetName(),
		OwnerID: account.ID,
		Master:  true,
	}
	account.ProjectID = project.ID

	perm := model.ProjectMember{
		ProjectID:  project.ID,
		AccountID:  account.ID,
		Permission: uint32(v1.Permission_OWNER),
	}

	pool := model.Pool{
		ID:        uuid.New().String(),
		ProjectID: project.ID,
		Name:      "default",
		Default:   true,
	}
	account.PoolID = pool.ID
	project.PoolID = pool.ID

	poolMember := model.PoolMember{
		AccountID:  account.ID,
		PoolID:     pool.ID,
		Permission: uint32(v1.Permission_OWNER),
	}

	client := model.Client{
		ID:          request.GetClientId(),
		PoolID:      pool.ID,
		Name:        "default",
		CreatedByID: account.ID,
		Default:     true,
	}

	// Create project and account in a transaction
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		logrus.Infof("account: %v", account)

		err := tx.CreateProject(ctx, project)
		if err != nil {
			return err
		}

		//// if password is provided, clientSecretHash it and provider it
		password := request.GetPassword()
		if password != "" {
			secret := x.Keygen()
			hash := x.HashPassword(password, secret)
			account.PasswordHash = string(hash)
			account.Salt = secret
		}

		err = tx.CreateAccount(ctx, &account)
		if err != nil {
			return err
		}

		err = tx.CreateProjectMember(ctx, &perm)
		if err != nil {
			return err
		}

		err = tx.CreatePool(ctx, &pool)
		if err != nil {
			return err
		}

		err = tx.AddPoolMember(ctx, &poolMember)
		if err != nil {
			return err
		}

		err = tx.CreateClient(ctx, &client)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// if the mail server is configured, send a verification email anyway
	// verification email will be sent if the account is created successfully
	//if request.GetVerifyEmail() {
	//	verificationCode := x.GenerateVerificationCode()
	//	err := a.cache.Set("email:"+account.Email, verificationCode, time.Hour)
	//	if err != nil {
	//		logrus.Errorf("failed to set verification code in cache: %v", err)
	//		return nil, err
	//	}
	//}

	return &v1.CreateAdminProjectResponse{
		Project: &v1.Project{
			Id:     project.ID,
			Name:   project.Name,
			PoolId: pool.ID,
		},
		Account: &v1.Account{
			Id:        account.ID,
			Email:     account.Email,
			ProjectId: account.ProjectID,
			PoolId:    pool.ID,
		},
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
