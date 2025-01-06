package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/pkg/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testAdminProjectName        = "test-admin-project"
	testAdminProjectVisibleName = "Test Admin Project"
	testAdminProjectEmail       = "test-admin@mail.com"
)

var testAdminProjectPassword = "test-admin-password"

func createAdminProject(t *testing.T) (*v1.Project, *v1.AccessKey) {
	db := tester.TestDB()
	provider := store.NewDefaultProvider(store.NewGormStore(db))
	redis := tester.TestRedis()
	adminProjectService := NewAdminProjectService(provider, redis)
	accessTokenService := NewAccessKeyService(permission.NewNullAuthbasePermission(), provider, redis)

	res, err := adminProjectService.CreateAdminProject(context.Background(), &v1.CreateAdminProjectRequest{
		Name:        testAdminProjectName,
		VisibleName: testAdminProjectVisibleName,
		Email:       testAdminProjectEmail,
		Password:    &testAdminProjectPassword,
	})
	assert.NoError(t, err, "error creating admin project")

	access, err := accessTokenService.CreateAccessKey(context.Background(), &v1.CreateAccessKeyRequest{
		Email:    testAdminProjectEmail,
		Password: testAdminProjectPassword,
	})
	assert.NotNil(t, access, "access key is nil")

	return res.Project, access.Token
}

func TestAdminProjectService_CreateAdminProject(t *testing.T) {
	tester.Setup()
}
