package permission

import (
	"context"
	"errors"
	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	"github.com/google/uuid"
)

type ObjectType string

type Permission interface {
	CreatePermission(ctx context.Context, objectID uuid.UUID, objectType string, subjectID uuid.UUID, subjectType string, relation string) error
	DeletePermission(ctx context.Context, objectID uuid.UUID, objectType string, subjectID uuid.UUID, subjectType string, relation string) error
	DeleteSubjectPermissions(ctx context.Context, subjectID uuid.UUID, subjectType string) error
	DeleteObjectPermissions(ctx context.Context, objectID uuid.UUID, objectType string) error
	Check(ctx context.Context, objectID uuid.UUID, objectType string, subjectID uuid.UUID, subjectType string, relation string) (bool, error)
}

type AuthZed struct {
}

func (a *AuthZed) CreatePermission(ctx context.Context, objectID uuid.UUID, objectType string, subjectID uuid.UUID, subjectType string, relation string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZed) DeletePermission(ctx context.Context, objectID uuid.UUID, objectType string, subjectID uuid.UUID, subjectType string, relation string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZed) DeleteSubjectPermissions(ctx context.Context, subjectID uuid.UUID, subjectType string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZed) DeleteObjectPermissions(ctx context.Context, objectID uuid.UUID, objectType string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZed) Check(ctx context.Context, objectID uuid.UUID, objectType string, subjectID uuid.UUID, subjectType string, relation string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

// MemberPermission is an interface representing the member permissions that are used in the service layer
type MemberPermission interface {
	//CreateProjectMember(ctx context.Context, objectType ObjectType, memberID uuid.UUID, relation string) error
	//DeleteProjectMember(ctx context.Context, objectType ObjectType, memberID uuid.UUID, relation string) error

	// CheckMasterProjectPermission checks if the user has the permission to perform
	CheckMasterProjectPermission(ctx context.Context, relation string) error
	// CheckProjectPermission checks if the user has the permission to perform the action on the project
	CheckProjectPermission(ctx context.Context, orgID uuid.UUID, relation string) error
}

// AuthBasePermission is an interface representing the authbase permissions that are used in the service layer
type AuthBasePermission interface {
	MemberPermission
}

// AuthZedPermission is a struct that implements the AuthBasePermission interface
// it delegates the permission checks to the Zed package
type AuthZedPermission struct {
}

func NewAuthZedPermission() *AuthZedPermission {
	return &AuthZedPermission{}
}

var _ AuthBasePermission = new(StoreBasedPermission)

// StoreBasedPermission is a struct that implements the AuthBasePermission interface
type StoreBasedPermission struct {
	store store.Provider
}

func NewStoreBasedPermission(store store.Provider) *StoreBasedPermission {
	return &StoreBasedPermission{store: store}
}

// CheckMasterProjectPermission checks if the user has the permission to perform the action on the master project
func (s *StoreBasedPermission) CheckMasterProjectPermission(ctx context.Context, relation string) error {
	accountID, err := x.GetAuthbaseAccountID(ctx)
	if err != nil {
		return err
	}

	as, err := store.GetProjectStore(ctx, s.store)
	if err != nil {
		return err
	}
	user, err := as.GetAccountByID(ctx, accountID)
	if err != nil {
		return err
	}

	if user.Disabled {
		return errors.New("user account is disabled")
	}

	//if !user.Member {
	//	return x.ErrNotProjectMember
	//}

	// if the user is a member of the master project
	if user.Project.Master {
		permission, err := as.GetProjectMemberByID(ctx, uuid.MustParse(user.ProjectID), accountID)
		if err != nil {
			return err
		}

		perm := permissionMap[relation]
		if permission.Permission >= perm {
			return nil
		}
	}

	return x.ErrUnauthorized
}

var permissionMap = map[string]uint32{
	"unknown": uint32(v1.Permission_UNKNOWN),
	"none":    uint32(v1.Permission_NONE),
	"viewer":  uint32(v1.Permission_VIEWER),
	"admin":   uint32(v1.Permission_OWNER),
	"owner":   uint32(v1.Permission_OWNER),
}

// CheckProjectPermission checks if the user has the permission to perform the action on the project
func (s *StoreBasedPermission) CheckProjectPermission(ctx context.Context, projectID uuid.UUID, relation string) error {
	//scopes, err := x.GetAuthbaseScopes(ctx)
	//if err != nil {
	//	return err
	//}

	return nil
}

// NullAuthbasePermission is a struct that implements the AuthBasePermission interface
type NullAuthbasePermission struct {
}

func NewNullAuthbasePermission() *NullAuthbasePermission {
	return &NullAuthbasePermission{}
}

var _ AuthBasePermission = new(NullAuthbasePermission)

// CheckMasterProjectPermission checks if the user has the permission to perform
// for NullAuthbasePermission it always returns nil, meaning the user has the permission
func (n *NullAuthbasePermission) CheckMasterProjectPermission(ctx context.Context, relation string) error {
	return nil
}

// CheckProjectPermission checks if the user has the permission to perform the action on the project
// for NullAuthbasePermission it always returns nil, meaning the user has the permission
func (n *NullAuthbasePermission) CheckProjectPermission(ctx context.Context, orgID uuid.UUID, relation string) error {
	return nil
}
