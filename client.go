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
	v1.UserServiceClient
	v1.MemberServiceClient
	v1.PermissionServiceClient
	v1.AuthServiceClient
	v1.OauthServiceClient
	gopackv1.TokenServiceClient
	v1.OfflineTokenServiceClient
	io.Closer
}

type client struct {
	conn *grpc.ClientConn
	v1.AdminProjectServiceClient
	v1.ProjectServiceClient
	v1.UserServiceClient
	v1.MemberServiceClient
	v1.PermissionServiceClient
	v1.AuthServiceClient
	v1.OauthServiceClient
	gopackv1.TokenServiceClient
	v1.OfflineTokenServiceClient
}

func NewClient(port string) (Client, error) {
	conn, err := grpc.NewClient(":4000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &client{
		conn:                      conn,
		ProjectServiceClient:      v1.NewProjectServiceClient(conn),
		UserServiceClient:         v1.NewUserServiceClient(conn),
		MemberServiceClient:       v1.NewMemberServiceClient(conn),
		PermissionServiceClient:   v1.NewPermissionServiceClient(conn),
		AuthServiceClient:         v1.NewAuthServiceClient(conn),
		OauthServiceClient:        v1.NewOauthServiceClient(conn),
		TokenServiceClient:        gopackv1.NewTokenServiceClient(conn),
		AdminProjectServiceClient: v1.NewAdminProjectServiceClient(conn),
		OfflineTokenServiceClient: v1.NewOfflineTokenServiceClient(conn),
	}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}
