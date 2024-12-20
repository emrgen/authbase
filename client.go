package authbase

import (
	v1 "github.com/emrgen/authbase/apis/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client interface {
	v1.AdminOrganizationServiceClient
	v1.OrganizationServiceClient
	v1.UserServiceClient
	v1.PermissionServiceClient
	v1.AuthServiceClient
	v1.OauthServiceClient
	v1.TokenServiceClient
}

type client struct {
	v1.AdminOrganizationServiceClient
	v1.OrganizationServiceClient
	v1.UserServiceClient
	v1.PermissionServiceClient
	v1.AuthServiceClient
	v1.OauthServiceClient
	v1.TokenServiceClient
}

func NewClient(port string) (Client, error) {
	conn, err := grpc.NewClient(":4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &client{
		OrganizationServiceClient:      v1.NewOrganizationServiceClient(conn),
		UserServiceClient:              v1.NewUserServiceClient(conn),
		PermissionServiceClient:        v1.NewPermissionServiceClient(conn),
		AuthServiceClient:              v1.NewAuthServiceClient(conn),
		OauthServiceClient:             v1.NewOauthServiceClient(conn),
		TokenServiceClient:             v1.NewTokenServiceClient(conn),
		AdminOrganizationServiceClient: v1.NewAdminOrganizationServiceClient(conn),
	}, nil
}
