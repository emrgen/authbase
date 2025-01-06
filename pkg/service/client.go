package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
)

func NewClientService(perm permission.MemberPermission, store store.Provider, cache *cache.Redis) *ClientService {
	return &ClientService{store: store}
}

var _ v1.ClientServiceServer = new(ClientService)

type ClientService struct {
	store store.Provider
	v1.UnimplementedClientServiceServer
}

func (c *ClientService) CreateClient(ctx context.Context, request *v1.CreateClientRequest) (*v1.CreateClientResponse, error) {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

func (c *ClientService) DeleteClient(ctx context.Context, request *v1.DeleteClientRequest) (*v1.DeleteClientResponse, error) {
	//TODO implement me
	panic("implement me")
}
