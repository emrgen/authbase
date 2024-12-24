package permission

import (
	"context"
	"errors"

	v1 "github.com/emrgen/authbase/apis/v1"
	"github.com/emrgen/authbase/pkg/store"
	"github.com/emrgen/authbase/x"
	gox "github.com/emrgen/gopack/x"
	"github.com/google/uuid"
)

type ObjectType string

const (
	ObjectTypeUser               ObjectType = "user"
	ObjectTypeOrganization       ObjectType = "organization"
	ObjectTypeProject            ObjectType = "project"
	ObjectTypeMasterOrganization ObjectType = "master/organization"
)

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

type MemberPermission interface {
	//CreatePermission(ctx context.Context, objectType ObjectType, memberID uuid.UUID, relation string) error
	//DeletePermission(ctx context.Context, objectType ObjectType, memberID uuid.UUID, relation string) error

	// CheckMasterOrganizationPermission checks if the user has the permission to perform
	CheckMasterOrganizationPermission(ctx context.Context, relation string) error
	// CheckOrganizationPermission checks if the user has the permission to perform the action on the organization
	CheckOrganizationPermission(ctx context.Context, orgID uuid.UUID, relation string) error
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

type StoreBasedPermission struct {
	store store.Provider
}

func NewStoreBasedPermission(store store.Provider) *StoreBasedPermission {
	return &StoreBasedPermission{store: store}
}

// CheckMasterOrganizationPermission checks if the user has the permission to perform the action on the master organization
func (s *StoreBasedPermission) CheckMasterOrganizationPermission(ctx context.Context, relation string) error {
	userID, err := gox.GetUserID(ctx)
	if err != nil {
		return err
	}

	as, err := store.GetProjectStore(ctx, s.store)
	if err != nil {
		return err
	}
	user, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Disabled {
		return errors.New("user account is disabled")
	}

	if !user.Member {
		return x.ErrNotOrganizationMember
	}

	// if the user is a member of the master organization
	if user.Organization.Master {
		permission, err := as.GetPermissionByID(ctx, uuid.MustParse(user.OrganizationID), userID)
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
	"read":   uint32(v1.Permission_READ),
	"write":  uint32(v1.Permission_WRITE),
	"delete": uint32(v1.Permission_DELETE),
	"admin":  uint32(v1.Permission_ADMIN),
}

// CheckOrganizationPermission checks if the user has the permission to perform the action on the organization
func (s *StoreBasedPermission) CheckOrganizationPermission(ctx context.Context, orgID uuid.UUID, relation string) error {
	userID, err := gox.GetUserID(ctx)
	if err != nil {
		return err
	}

	as, err := store.GetProjectStore(ctx, s.store)
	if err != nil {
		return err
	}

	user, err := as.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Disabled {
		return errors.New("user account is disabled")
	}

	if !user.Member {
		return x.ErrNotOrganizationMember
	}

	err = s.CheckMasterOrganizationPermission(ctx, relation)
	if errors.Is(err, x.ErrUnauthorized) {
		_, err := gox.GetUserID(ctx)
		if err != nil {
			return err
		}

		permission, err := as.GetPermissionByID(ctx, orgID, userID)
		if err != nil {
			return err
		}

		perm := permissionMap[relation]

		// check if the user has the write permission
		if permission.Permission >= perm {
			return nil
		}

		// as the user does not have the permission on the master org and target org
		// return unauthorized
		return x.ErrUnauthorized
	}

	return err
}

// NullAuthbasePermission is a struct that implements the AuthBasePermission interface
type NullAuthbasePermission struct {
}

func NewNullAuthbasePermission() *NullAuthbasePermission {
	return &NullAuthbasePermission{}
}

var _ AuthBasePermission = new(NullAuthbasePermission)

// CheckMasterOrganizationPermission checks if the user has the permission to perform
// for NullAuthbasePermission it always returns nil, meaning the user has the permission
func (n *NullAuthbasePermission) CheckMasterOrganizationPermission(ctx context.Context, relation string) error {
	return nil
}

// CheckOrganizationPermission checks if the user has the permission to perform the action on the organization
// for NullAuthbasePermission it always returns nil, meaning the user has the permission
func (n *NullAuthbasePermission) CheckOrganizationPermission(ctx context.Context, orgID uuid.UUID, relation string) error {
	return nil
}
