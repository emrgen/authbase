package service

import (
	"context"
	"encoding/json"
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

var _ v1.OrganizationServiceServer = new(OrganizationService)

type OrganizationService struct {
	perm  permission.MemberPermission
	store store.Provider
	cache *cache.Redis
	v1.UnimplementedOrganizationServiceServer
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(perm permission.MemberPermission, store store.Provider, cache *cache.Redis) *OrganizationService {
	return &OrganizationService{perm: perm, store: store, cache: cache}
}

func (o *OrganizationService) CreateOrganization(ctx context.Context, request *v1.CreateOrganizationRequest) (*v1.CreateOrganizationResponse, error) {
	var err error
	err = o.perm.CheckMasterOrganizationPermission(ctx, "write")
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

	org := model.Organization{
		ID:      uuid.New().String(),
		Name:    request.GetName(),
		OwnerID: user.ID,
	}
	user.OrganizationID = org.ID

	perm := model.Permission{
		OrganizationID: org.ID,
		UserID:         user.ID,
		Permission:     uint32(v1.Permission_ADMIN),
	}

	// if this is the first organization, make the organization is the master organization
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		_, total, _ := tx.ListOrganizations(ctx, 1, 1)
		if total == 0 {
			org.Master = true
			user.SassAdmin = true
		}

		err := tx.CreateOrganization(ctx, &org)
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

		err = tx.CreatePermission(ctx, &perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.CreateOrganizationResponse{
		Id:   org.ID,
		Name: org.Name,
	}, nil
}

// GetOrganizationId gets the organization ID, given the name
func (o *OrganizationService) GetOrganizationId(ctx context.Context, request *v1.GetOrganizationIdRequest) (*v1.GetOrganizationIdResponse, error) {
	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	org, err := as.GetOrganizationByName(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &v1.GetOrganizationIdResponse{
		Id:   org.ID,
		Name: org.Name,
	}, nil
}

func (o *OrganizationService) GetOrganization(ctx context.Context, request *v1.GetOrganizationRequest) (*v1.GetOrganizationResponse, error) {
	var err error

	err = o.perm.CheckMasterOrganizationPermission(ctx, "read")
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

	org, err := as.GetOrganizationByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.GetOrganizationResponse{
		Organization: &v1.Organization{
			Id:        org.ID,
			Name:      org.Name,
			OwnerId:   org.OwnerID,
			CreatedAt: timestamppb.New(org.CreatedAt),
			UpdatedAt: timestamppb.New(org.UpdatedAt),
		},
	}, nil
}

func (o *OrganizationService) ListOrganizations(ctx context.Context, request *v1.ListOrganizationsRequest) (*v1.ListOrganizationsResponse, error) {
	var err error

	err = o.perm.CheckMasterOrganizationPermission(ctx, "read")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}
	page := utils.GetPage(request)

	orgs, total, err := as.ListOrganizations(ctx, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var organizations []*v1.Organization
	for _, org := range orgs {
		organizations = append(organizations, &v1.Organization{
			Id:        org.ID,
			Name:      org.Name,
			OwnerId:   org.OwnerID,
			Master:    org.Master,
			CreatedAt: timestamppb.New(org.CreatedAt),
			UpdatedAt: timestamppb.New(org.UpdatedAt),
		})
	}

	return &v1.ListOrganizationsResponse{
		Organizations: organizations,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil
}

func (o *OrganizationService) UpdateOrganization(ctx context.Context, request *v1.UpdateOrganizationRequest) (*v1.UpdateOrganizationResponse, error) {
	var err error
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	err = o.perm.CheckOrganizationPermission(ctx, id, "write")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		org, err := tx.GetOrganizationByID(ctx, id)
		if err != nil {
			return err
		}

		org.Name = request.GetName()

		err = tx.UpdateOrganization(ctx, org)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateOrganizationResponse{}, nil
}

func (o *OrganizationService) DeleteOrganization(ctx context.Context, request *v1.DeleteOrganizationRequest) (*v1.DeleteOrganizationResponse, error) {

	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}
	err = o.perm.CheckOrganizationPermission(ctx, id, "write")
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, o.store)
	if err != nil {
		return nil, err
	}

	err = as.DeleteOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteOrganizationResponse{}, nil
}

func (o *OrganizationService) AddOauthProvider(ctx context.Context, request *v1.AddOauthProviderRequest) (*v1.AddOauthProviderResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = o.perm.CheckOrganizationPermission(ctx, orgID, "write")
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

	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	providerModel := model.OauthProvider{
		ID:             uuid.New().String(),
		OrganizationID: orgID.String(),
		Config:         string(data),
	}

	err = as.CreateOauthProvider(ctx, &providerModel)
	if err != nil {
		return nil, err
	}

	return &v1.AddOauthProviderResponse{
		Message: "Oauth provider added successfully",
	}, nil
}

func (o *OrganizationService) GetOauthProvider(ctx context.Context, request *v1.GetOauthProviderRequest) (*v1.GetOauthProviderResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = o.perm.CheckOrganizationPermission(ctx, orgID, "read")
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

	providerConfig := make(map[string]interface{})
	err = json.Unmarshal([]byte(provider.Config), &providerConfig)
	if err != nil {
		return nil, err
	}

	clientID := providerConfig["client_id"].(string)
	clientSecret := providerConfig["client_secret"].(string)

	return &v1.GetOauthProviderResponse{
		Provider: &v1.OAuthProvider{
			Id:           provider.ID,
			Name:         providerConfig["provider"].(string),
			ClientId:     clientID,
			ClientSecret: clientSecret,
		},
	}, nil
}

func (o *OrganizationService) ListOauthProviders(ctx context.Context, request *v1.ListOauthProvidersRequest) (*v1.ListOauthProvidersResponse, error) {
	orgID, err := uuid.Parse(request.GetOrganizationId())
	if err != nil {
		return nil, err
	}

	err = o.perm.CheckOrganizationPermission(ctx, orgID, "read")
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
		providerConfig := make(map[string]interface{})
		err = json.Unmarshal([]byte(provider.Config), &providerConfig)
		if err != nil {
			return nil, err
		}

		oauthProviders = append(oauthProviders, &v1.OAuthProvider{
			Id:       provider.ID,
			Name:     providerConfig["provider"].(string),
			ClientId: providerConfig["client_id"].(string),
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
//	organization_id: "organization_id",
//	provider: "Google",
//	client_id: "client_id",
//	client_secret: "client_secret",
func (o *OrganizationService) UpdateOauthProvider(ctx context.Context, request *v1.UpdateOauthProviderRequest) (*v1.UpdateOauthProviderResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteOauthProvider deletes the oauth provider information.
// The provider ID is required to delete the provider information.
// Example:
//
//	id: "provider_id"
func (o *OrganizationService) DeleteOauthProvider(ctx context.Context, request *v1.DeleteOauthProviderRequest) (*v1.DeleteOauthProviderResponse, error) {
	//TODO implement me
	panic("implement me")
}
