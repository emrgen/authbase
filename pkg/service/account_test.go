package service

import (
	"github.com/emrgen/authbase/pkg/tester"
	"testing"
)

func TestAccountService_CreateAccount(t *testing.T) {
	tester.RemoveDBFile()
	tester.Setup()

	//ab, _, project := createProject(t, "test-project-1", "Test Project 1", testEmail1, "password")
	//accountService := NewAccountService(ab.perm, ab.provider, ab.redis)
	//res, err := accountService.CreateAccount(ab.ctx, &v1.CreateAccountRequest{
	//	//PoolId:      project.PoolId,
	//	Username:    "test-username",
	//	VisibleName: "Test User",
	//	Email:       testEmail2,
	//	Password:    "password",
	//})
	//assert.NoError(t, err, "failed to create account")
	//assert.Equal(t, res.Account.Username, "test-username", "account username is not correct")
}
