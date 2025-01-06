package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	x "github.com/emrgen/authbase/x"
	"github.com/google/uuid"
)

// NewClientService creates a new ClientService.
func NewClientService(perm permission.MemberPermission, store store.Provider, cache *cache.Redis) *ClientService {
	return &ClientService{store: store, perm: perm}
}

var _ v1.ClientServiceServer = new(ClientService)

type ClientService struct {
	perm  permission.MemberPermission
	store store.Provider
	v1.UnimplementedClientServiceServer
}

func (c *ClientService) CreateClient(ctx context.Context, request *v1.CreateClientRequest) (*v1.CreateClientResponse, error) {
	userID, err := x.GetAuthbaseUserID(ctx)
	if err != nil {
		return nil, err
	}

	// check if the user has the permission to create a client
	err = c.perm.CheckProjectPermission(ctx, userID, "write")
	if err != nil {
		return nil, err
	}

	// create the client
	as, err := store.GetProjectStore(ctx, c.store)
	if err != nil {
		return nil, err
	}

	secret := x.Keygen()

	client := model.Client{
		ID:        uuid.New().String(),
		ProjectID: request.GetProjectId(),
		Name:      request.GetName(),
		Secret:    secret,
	}
	err = as.CreateClient(ctx, &client)
	if err != nil {
		return nil, err
	}

	return &v1.CreateClientResponse{
		Client: &v1.Client{
			Id:           client.ID,
			ProjectId:    request.GetProjectId(),
			ClientSecret: secret,
		},
	}, nil
}

func (c *ClientService) GetClient(ctx context.Context, request *v1.GetClientRequest) (*v1.GetClientResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClientService) ListClients(ctx context.Context, request *v1.ListClientsRequest) (*v1.ListClientsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *ClientService) UpdateClient(ctx context.Context, request *v1.UpdateClientRequest) (*v1.UpdateClientResponse, error) {
	var err error
	userID, err := uuid.Parse(request.GetClientId())
	if err != nil {
		return nil, err
	}

	as, err := store.GetProjectStore(ctx, c.store)
	if err != nil {
		return nil, err
	}

	member, err := as.GetAccountByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	orgID, err := uuid.Parse(member.ProjectID)
	if err != nil {
		return nil, err
	}

	err = c.perm.CheckProjectPermission(ctx, orgID, "write")
	if err != nil {
		return nil, err
	}

	// update the member and the permission
	err = as.Transaction(func(tx store.AuthBaseStore) error {
		perm, err := tx.GetProjectMemberByID(ctx, orgID, userID)
		if err != nil {
			return err
		}

		if request.GetPermission() != v1.Permission_UNKNOWN {
			perm.Permission = uint32(request.GetPermission())
		}

		err = tx.UpdateAccount(ctx, member)
		if err != nil {
			return err
		}

		err = tx.UpdateProjectMember(ctx, perm)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.UpdateClientResponse{
		Client: &v1.Client{
			Id: member.ID,
		},
	}, nil
}

func (c *ClientService) DeleteClient(ctx context.Context, request *v1.DeleteClientRequest) (*v1.DeleteClientResponse, error) {
	//TODO implement me
	panic("implement me")
}
