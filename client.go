package authbase

import (
	v1 "github.com/emrgen/authbase/apis/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

type Client interface {
	v1.AdminProjectServiceClient
	v1.ProjectServiceClient
	v1.AccountServiceClient
	v1.ProjectMemberServiceClient
	v1.AuthServiceClient
	v1.OAuth2ServiceClient
	v1.AccessKeyServiceClient
	v1.SessionServiceClient
	v1.ClientServiceClient
	v1.PoolServiceClient
	v1.PoolMemberServiceClient
	v1.TokenServiceClient
	v1.GroupServiceClient
	io.Closer
}

type client struct {
	conn *grpc.ClientConn
	v1.AdminProjectServiceClient
	v1.ProjectServiceClient
	v1.AccountServiceClient
	v1.ProjectMemberServiceClient
	v1.AuthServiceClient
	v1.SessionServiceClient
	v1.AccessKeyServiceClient
	v1.OAuth2ServiceClient
	v1.ClientServiceClient
	v1.PoolServiceClient
	v1.PoolMemberServiceClient
	v1.TokenServiceClient
	v1.GroupServiceClient
}

func NewClient(port string) (Client, error) {
	conn, err := grpc.NewClient(":4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &client{
		conn:                       conn,
		ProjectServiceClient:       v1.NewProjectServiceClient(conn),
		AuthServiceClient:          v1.NewAuthServiceClient(conn),
		TokenServiceClient:         v1.NewTokenServiceClient(conn),
		AdminProjectServiceClient:  v1.NewAdminProjectServiceClient(conn),
		SessionServiceClient:       v1.NewSessionServiceClient(conn),
		AccountServiceClient:       v1.NewAccountServiceClient(conn),
		OAuth2ServiceClient:        v1.NewOAuth2ServiceClient(conn),
		AccessKeyServiceClient:     v1.NewAccessKeyServiceClient(conn),
		ClientServiceClient:        v1.NewClientServiceClient(conn),
		ProjectMemberServiceClient: v1.NewProjectMemberServiceClient(conn),
		PoolServiceClient:          v1.NewPoolServiceClient(conn),
		PoolMemberServiceClient:    v1.NewPoolMemberServiceClient(conn),
		GroupServiceClient:         v1.NewGroupServiceClient(conn),
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
