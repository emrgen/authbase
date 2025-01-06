package authbase

import (
	v1 "github.com/emrgen/authbase/apis/v1"
	gopackv1 "github.com/emrgen/gopack/apis/v1"
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
	gopackv1.TokenServiceClient
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
	gopackv1.TokenServiceClient
}

func NewClient(port string) (Client, error) {
	conn, err := grpc.NewClient(":4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &client{
		conn:                      conn,
		ProjectServiceClient:      v1.NewProjectServiceClient(conn),
		AuthServiceClient:         v1.NewAuthServiceClient(conn),
		TokenServiceClient:        gopackv1.NewTokenServiceClient(conn),
		AdminProjectServiceClient: v1.NewAdminProjectServiceClient(conn),
		SessionServiceClient:      v1.NewSessionServiceClient(conn),
		AccountServiceClient:      v1.NewAccountServiceClient(conn),
		OAuth2ServiceClient:       v1.NewOAuth2ServiceClient(conn),
		AccessKeyServiceClient:    v1.NewAccessKeyServiceClient(conn),
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
