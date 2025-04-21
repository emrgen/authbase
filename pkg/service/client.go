package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/secret"
	"github.com/emrgen/authbase/pkg/store"
	x "github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewClientService creates a new ClientService.
func NewClientService(perm permission.MemberPermission, store store.Provider, secrets secret.Store) *ClientService {
	return &ClientService{store: store, perm: perm, secret: secrets}
}

var _ v1.ClientServiceServer = new(ClientService)

// ClientService is the service for managing clients for a account pool.
type ClientService struct {
	perm   permission.MemberPermission
	store  store.Provider
	secret secret.Store
	v1.UnimplementedClientServiceServer
}

// CreateClient creates a new client for the given pool.
func (c *ClientService) CreateClient(ctx context.Context, request *v1.CreateClientRequest) (*v1.CreateClientResponse, error) {
	as, err := store.GetProjectStore(ctx, c.store)
	if err != nil {
		return nil, err
	}

	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return nil, err
	}

	// check if the user has permission to create a client for the pool
	err = c.perm.CheckProjectPermission(ctx, accountID, permission.ProjectPermissionWrite)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.GetPoolId())
	if err != nil {
		return nil, err
	}
	pool, err := as.GetPoolByID(ctx, poolID)
	if err != nil {
		return nil, err
	}

	client := model.Client{
		ID:          uuid.New().String(),
		PoolID:      pool.ID,
		Name:        request.GetName(),
		CreatedByID: accountID.String(),
	}
	err = as.CreateClient(ctx, &client)
	if err != nil {
		return nil, err
	}

	account, err := as.GetAccountByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// TODO: should we return the client secret to the user?
	return &v1.CreateClientResponse{
		Client: &v1.Client{
			Id:     client.ID,
			PoolId: request.GetPoolId(),
			Name:   client.Name,
			CreatedByUser: &v1.Account{
				Id:          accountID.String(),
				VisibleName: account.VisibleName,
			},
		},
	}, nil
}

func (c *ClientService) GetClient(ctx context.Context, request *v1.GetClientRequest) (*v1.GetClientResponse, error) {
	as, err := store.GetProjectStore(ctx, c.store)
	if err != nil {
		return nil, err
	}

	clientID, err := uuid.Parse(request.GetClientId())
	if err != nil {
		return nil, err
	}

	client, err := as.GetClientByID(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return &v1.GetClientResponse{
		Client: &v1.Client{
			Id:     client.ID,
			PoolId: client.PoolID,
			Name:   client.Name,
			CreatedByUser: &v1.Account{
				Id:          client.CreatedByID,
				VisibleName: client.CreatedByAccount.VisibleName,
			},
		},
	}, nil

}

func (c *ClientService) ListClients(ctx context.Context, request *v1.ListClientsRequest) (*v1.ListClientsResponse, error) {
	as, err := store.GetProjectStore(ctx, c.store)
	if err != nil {
		return nil, err
	}

	poolID, err := uuid.Parse(request.GetPoolId())
	if err != nil {
		return nil, err
	}

	page := x.GetPageFromRequest(request)
	clients, total, err := as.ListClients(ctx, poolID, int(page.Page), int(page.Size))
	if err != nil {
		return nil, err
	}

	var clientProtos []*v1.Client
	for _, client := range clients {
		clientProtos = append(clientProtos, &v1.Client{
			Id:        client.ID,
			PoolId:    client.PoolID,
			Name:      client.Name,
			CreatedAt: timestamppb.New(client.CreatedAt),
		})
	}

	return &v1.ListClientsResponse{
		Clients: clientProtos,
		Meta: &v1.Meta{
			Total: int32(total),
			Page:  page.Page,
			Size:  page.Size,
		},
	}, nil
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
	as, err := store.GetProjectStore(ctx, c.store)
	if err != nil {
		return nil, err
	}

	clientID, err := uuid.Parse(request.GetClientId())
	if err != nil {
		return nil, err
	}

	err = as.DeleteClient(ctx, clientID)
	if err != nil {
		return nil, err
	}

	return &v1.DeleteClientResponse{}, nil

}
