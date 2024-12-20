package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/emrgen/authbase/x/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

var _ v1.OrganizationServiceServer = new(OrganizationService)

type OrganizationService struct {
	store store.AuthBaseStore
	cache *cache.Redis
	v1.UnimplementedOrganizationServiceServer
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(store store.AuthBaseStore, cache *cache.Redis) *OrganizationService {
	return &OrganizationService{store: store, cache: cache}
}

func (o *OrganizationService) CreateOrganization(ctx context.Context, request *v1.CreateOrganizationRequest) (*v1.CreateOrganizationResponse, error) {
	password := request.GetPassword()
	verifyEmail := request.GetVerifyEmail()

	user := model.User{
		ID:       uuid.New().String(),
		Email:    request.GetEmail(),
		Username: request.GetUsername(),
	}

	org := model.Organization{
		ID:      uuid.New().String(),
		Name:    request.GetName(),
		OwnerID: user.ID,
	}
	user.OrganizationID = org.ID

	// if this is the first organization, make the organization is the master organization
	err := o.store.Transaction(func(tx store.AuthBaseStore) error {
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
	org, err := o.store.GetOrganizationByName(ctx, request.GetName())
	if err != nil {
		return nil, err
	}

	return &v1.GetOrganizationIdResponse{
		Id:   org.ID,
		Name: org.Name,
	}, nil
}

func (o *OrganizationService) GetOrganization(ctx context.Context, request *v1.GetOrganizationRequest) (*v1.GetOrganizationResponse, error) {
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	org, err := o.store.GetOrganizationByID(ctx, id)
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
	page := utils.GetPage(request)

	orgs, total, err := o.store.ListOrganizations(ctx, int(page.Page), int(page.Size))
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
	id, err := uuid.Parse(request.GetId())
	if err != nil {
		return nil, err
	}

	err = o.store.Transaction(func(tx store.AuthBaseStore) error {
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

	err = o.store.DeleteOrganization(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteOrganizationResponse{}, nil
}

func (o *OrganizationService) AddOauthProvider(ctx context.Context, request *v1.AddOauthProviderRequest) (*v1.AddOauthProviderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OrganizationService) GetOauthProvider(ctx context.Context, request *v1.GetOauthProviderRequest) (*v1.GetOauthProviderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OrganizationService) ListOauthProviders(ctx context.Context, request *v1.ListOauthProvidersRequest) (*v1.ListOauthProvidersResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OrganizationService) UpdateOauthProvider(ctx context.Context, request *v1.UpdateOauthProviderRequest) (*v1.UpdateOauthProviderResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (o *OrganizationService) DeleteOauthProvider(ctx context.Context, request *v1.DeleteOauthProviderRequest) (*v1.DeleteOauthProviderResponse, error) {
	//TODO implement me
	panic("implement me")
}
