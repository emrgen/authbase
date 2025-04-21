package service

import (
	"context"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/permission"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/pkg/tester"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testAdminProjectName        = "test-admin-project"
	testAdminProjectVisibleName = "Test Admin Project"
	testAdminProjectEmail       = "test-admin@mail.com"
)

var testAdminProjectPassword = "test-admin-password"

type authbase struct {
	adminProject *v1.Project
	accessKey    *v1.AccessKey
	perm         permission.AuthBasePermission
	provider     store.Provider
	redis        *cache.Redis
	ctx          context.Context
}

func createAdminProject(t *testing.T) *authbase {
	db := tester.TestDB()
	provider := store.NewDefaultProvider(store.NewGormStore(db))
	keyProvider := x.NewUnverifiedKeyProvider()
	redis := tester.TestRedis()
	adminProjectService := NewAdminProjectService(provider, redis)
	verifier := x.NewUnverifiedVerifier()
	accessTokenService := NewAccessKeyService(permission.NewNullAuthbasePermission(), provider, redis, keyProvider, verifier)

	ctx := context.TODO()

	res, err := adminProjectService.CreateAdminProject(ctx, &v1.CreateAdminProjectRequest{
		Name:        testAdminProjectName,
		VisibleName: testAdminProjectVisibleName,
		Email:       testAdminProjectEmail,
		Password:    &testAdminProjectPassword,
	})
	assert.NoError(t, err, "error creating admin project, %v", err)

	ctx = context.WithValue(ctx, x.ProjectIDKey, uuid.MustParse(res.Project.Id))
	ctx = context.WithValue(ctx, x.PoolIDKey, uuid.MustParse(res.Account.PoolId))
	ctx = context.WithValue(ctx, x.AccountIDKey, uuid.MustParse(res.Account.Id))

	access, err := accessTokenService.CreateAccessKey(ctx, &v1.CreateAccessKeyRequest{
		//PoolId:   res.Project.PoolId,
		Email:    testAdminProjectEmail,
		Password: testAdminProjectPassword,
	})
	assert.NotNil(t, access, "access key is nil")

	perm := permission.NewStoreBasedPermission(provider)

	return &authbase{
		adminProject: res.Project,
		accessKey:    access.Token,
		perm:         perm,
		provider:     provider,
		redis:        redis,
		ctx:          ctx,
	}
}

func TestAdminProjectService_CreateAdminProject(t *testing.T) {
	tester.RemoveDBFile()
	tester.Setup()
	createAdminProject(t)
}
