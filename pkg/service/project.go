package service

import (
	"context"
	"errors"
	"google.golang.org/grpc/peer"
	"time"

	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/emrgen/authbase/x/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ v1.ProjectServiceServer = new(ProjectService)

type ProjectService struct {
	perm  permission.MemberPermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedProjectServiceServer
}

// NewProjectService creates a new project service
func NewProjectService(perm permission.MemberPermission, store store.Provider, cache *cache.Redis) *ProjectService {
	return &ProjectService{perm: perm, store: store, cache: cache}
}

// CreateProject creates a new project and the owner user
func (o *ProjectService) CreateProject(ctx context.Context, request *v1.CreateProjectRequest) (*v1.CreateProjectResponse, error) {
	var err error

	// TODO: use token with create project permission to reduce the token scope
	// check if the user has permission to create an project
	err = o.perm.CheckMasterProjectPermission(ctx, "write")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	password := request.GetPassword()
	verifyEmail := request.GetVerifyEmail()

	user := model.User{
		ID:       uuid.New().String(),
		Email:    request.GetEmail(),
		Username: request.GetUsername(),
		Member:   true,
	}

	org := model.Project{
		ID:      uuid.New().String(),
		Name:    request.GetName(),
		OwnerID: user.ID,
	}
	user.ProjectID = org.ID

	perm := model.ProjectMember{
		ProjectID:  org.ID,
		UserID:     user.ID,
		Permission: uint32(v1.Permission_ADMIN),
	}

	// if this is the first project, make the project is the master project
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		_, total, _ := tx.ListProjects(ctx, 1, 1)
		if total == 0 {
			org.Master = true
			user.SassAdmin = true
		}

		err := tx.CreateProject(ctx, &org)
		if err != nil {
			return err
		}

		// if password is provided, email verification is not strictly required
		// FIXME: if the mail server config is provider the email verification will fail with error
		if password == "" || verifyEmail {
			verificationCode := x.RefreshToken()
			// save the code to the db
			err := o.cache.Set("email:"+user.Email, verificationCode, time.Hour)
			if err != nil {
				return err
			}

			if password == "" {
				// send email password reset email
				logrus.Infof("reset password code: %s", verificationCode)
			} else {
				// send email verification email
				logrus.Infof("verification code: %s", verificationCode)
			}
		} else if password != "" {
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

	return &v1.CreateProjectResponse{
		Id:   org.ID,
		Name: org.Name,
	}, nil
}

// GetProjectId gets the organization ID, given the name
func (o *ProjectService) GetProjectId(ctx context.Context, request *v1.GetProjectIdRequest) (*v1.GetProjectIdResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, errors.New("failed to get peer")
	}

	// FIXME: rate limit the number of requests from the same IP
	// use exponential backoff for retries to prevent abuse

	logrus.Infof("peer: %v", p.Addr.String())

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	org, err := as.GetProjectByName(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &v1.GetProjectIdResponse{
		Id:   org.ID,
		Name: org.Name,
	}, nil
}

func (o *ProjectService) GetProject(ctx context.Context, request *v1.GetProjectRequest) (*v1.GetProjectResponse, error) {
	var err error

	err = o.perm.CheckMasterProjectPermission(ctx, "read")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	org, err := as.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.GetProjectResponse{
		Project: &v1.Project{
			Id:        org.ID,
			Name:      org.Name,
			OwnerId:   org.OwnerID,
			CreatedAt: timestamppb.New(org.CreatedAt),
			UpdatedAt: timestamppb.New(org.UpdatedAt),
		},
	}, nil
}

func (o *ProjectService) ListProjects(ctx context.Context, request *v1.ListProjectsRequest) (*v1.ListProjectsResponse, error) {
	var err error

	err = o.perm.CheckMasterProjectPermission(ctx, "read")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}
	page := utils.GetPage(request)

	orgs, total, err := as.ListProjects(ctx, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var organizations []*v1.Project
	for _, org := range orgs {
		organizations = append(organizations, &v1.Project{
			Id:        org.ID,
			Name:      org.Name,
			OwnerId:   org.OwnerID,
			Master:    org.Master,
			CreatedAt: timestamppb.New(org.CreatedAt),
			UpdatedAt: timestamppb.New(org.UpdatedAt),
		})
	}

	return &v1.ListProjectsResponse{
		Projects: organizations,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil
}

func (o *ProjectService) UpdateProject(ctx context.Context, request *v1.UpdateProjectRequest) (*v1.UpdateProjectResponse, error) {
	var err error
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	err = o.perm.CheckProjectPermission(ctx, id, "write")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		org, err := tx.GetProjectByID(ctx, id)
		if err != nil {
			return err
		}

		org.Name = request.GetName()

		err = tx.UpdateProject(ctx, org)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateProjectResponse{}, nil
}

func (o *ProjectService) DeleteProject(ctx context.Context, request *v1.DeleteProjectRequest) (*v1.DeleteProjectResponse, error) {

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	err = o.perm.CheckProjectPermission(ctx, id, "write")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	err = as.DeleteProject(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteProjectResponse{}, nil
}

func (o *ProjectService) AddOauthProvider(ctx context.Context, request *v1.AddOauthProviderRequest) (*v1.AddOauthProviderResponse, error) {
	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = o.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	provider := request.GetProvider()

	m := make(map[string]interface{})
	m["provider"] = provider.GetName()
	m["client_id"] = provider.GetClientId()
	m["client_secret"] = provider.GetClientSecret()
	oauthConfig := model.OAuthConfig{
		Provider:     provider.GetName(),
		ClientID:     provider.GetClientId(),
		ClientSecret: provider.GetClientSecret(),
		Scopes:       "openid profile email",
	}

	providerModel := model.OauthProvider{
		ID:        uuid.New().String(),
		ProjectID: orgID.String(),
		Config:    oauthConfig,
	}

	err = as.CreateOauthProvider(ctx, &providerModel)
	if err != nil {
		return nil, err
	}

	return &v1.AddOauthProviderResponse{
		Message: "Oauth provider added successfully",
	}, nil
}

func (o *ProjectService) GetOauthProvider(ctx context.Context, request *v1.GetOauthProviderRequest) (*v1.GetOauthProviderResponse, error) {
	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = o.perm.CheckProjectPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	provider, err := as.GetOauthProviderByName(ctx, orgID, request.GetProvider())
	if err != nil {
		return nil, err
	}

	return &v1.GetOauthProviderResponse{
		Provider: &v1.OAuthProvider{
			Id:           provider.ID,
			Name:         provider.Config.Provider,
			ClientId:     provider.Config.ClientID,
			ClientSecret: provider.Config.ClientSecret,
		},
	}, nil
}

func (o *ProjectService) ListOauthProviders(ctx context.Context, request *v1.ListOauthProvidersRequest) (*v1.ListOauthProvidersResponse, error) {
	orgID, err := uuid.Parse(request.GetProjectId())
	if err != nil {
		return nil, err
	}

	err = o.perm.CheckProjectPermission(ctx, orgID, "read")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	page := utils.GetPage(request)

	providers, total, err := as.ListOauthProviders(ctx, orgID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var oauthProviders []*v1.OAuthProvider
	for _, provider := range providers {
		oauthProviders = append(oauthProviders, &v1.OAuthProvider{
			Id:       provider.ID,
			Name:     provider.Config.Provider,
			ClientId: provider.Config.ClientID,
		})
	}

	return &v1.ListOauthProvidersResponse{
		Providers: oauthProviders,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil
}

// UpdateOauthProvider updates the oauth provider information.
// The provider ID is required to update the provider information.
// Example:
//
//	project_id: "project_id",
//	provider: "Google",
//	client_id: "client_id",
//	client_secret: "client_secret",
func (o *ProjectService) UpdateOauthProvider(ctx context.Context, request *v1.UpdateOauthProviderRequest) (*v1.UpdateOauthProviderResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteOauthProvider deletes the oauth provider information.
// The provider ID is required to delete the provider information.
// Example:
//
//	id: "provider_id"
func (o *ProjectService) DeleteOauthProvider(ctx context.Context, request *v1.DeleteOauthProviderRequest) (*v1.DeleteOauthProviderResponse, error) {
	//TODO implement me
	panic("implement")
}
