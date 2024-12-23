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
	store store.AuthBaseStore
}

func NewStoreBasedPermission(store store.AuthBaseStore) *StoreBasedPermission {
	return &StoreBasedPermission{store: store}
}

// CheckMasterOrganizationPermission checks if the user has the permission to perform the action on the master organization
func (s *StoreBasedPermission) CheckMasterOrganizationPermission(ctx context.Context, relation string) error {
	userID, err := x.GetUserID(ctx)
	if err != nil {
		return err
	}

	user, err := s.store.GetUserByID(ctx, userID)
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
		permission, err := s.store.GetPermissionByID(ctx, uuid.MustParse(user.OrganizationID), userID)
		if err != nil {
			return err
		}

		// check if the user has the write permission
		if relation == "write" {
			if permission.Permission >= uint32(v1.Permission_WRITE) {
				return nil
			}
		}

		// check if the user has the read permission
		// being a member of master org the user has implicit read permission
		if relation == "read" {
			return nil
		}

		if relation == "delete" {
			if permission.Permission >= uint32(v1.Permission_DELETE) {
				return nil
			}
		}
	}

	return x.ErrUnauthorized
}

// CheckOrganizationPermission checks if the user has the permission to perform the action on the organization
func (s *StoreBasedPermission) CheckOrganizationPermission(ctx context.Context, orgID uuid.UUID, relation string) error {
	userID, err := x.GetUserID(ctx)
	if err != nil {
		return err
	}

	user, err := s.store.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.Disabled {
		return errors.New("user account is disabled")
	}

	if !user.Member {
		return x.ErrNotOrganizationMember
	}

	if relation == "write" {
		err := s.CheckMasterOrganizationPermission(ctx, "write")
		if errors.Is(err, x.ErrUnauthorized) {
			_, err := x.GetUserID(ctx)
			if err != nil {
				return err
			}

			permission, err := s.store.GetPermissionByID(ctx, orgID, userID)
			if err != nil {
				return err
			}

			// check if the user has the write permission
			if permission.Permission >= uint32(v1.Permission_WRITE) {
				return nil
			}
		}

		return err
	}

	if relation == "read" {
		err := s.CheckMasterOrganizationPermission(ctx, "read")
		if errors.Is(err, x.ErrUnauthorized) {
			_, err := x.GetUserID(ctx)
			if err != nil {
				return err
			}

			permission, err := s.store.GetPermissionByID(ctx, orgID, userID)
			if err != nil {
				return err
			}

			// check if the user has the read permission
			if permission.Permission >= uint32(v1.Permission_READ) {
				return nil
			}
		}

		return err
	}

	return x.ErrUnauthorized
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
