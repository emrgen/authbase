package permission

import (
	"context"
	"github.com/google/uuid"
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
	CreateMemberPermission(ctx context.Context, subjectID uuid.UUID) error
	DeleteMemberPermission(ctx context.Context, subjectID uuid.UUID) error
	CheckMemberPermission(ctx context.Context, subjectID uuid.UUID, objectID uuid.UUID, objectType string, relation string) (bool, error)
}

type UserPermission interface {
	CreateUserPermission(ctx context.Context, subjectID uuid.UUID) error
	DeleteUserPermission(ctx context.Context, subjectID uuid.UUID) error
	CheckUserPermission(ctx context.Context, subjectID uuid.UUID, objectID uuid.UUID, objectType string, relation string) (bool, error)
}

type OrganizationPermission interface {
	CreateOrganizationPermission(ctx context.Context, subjectID uuid.UUID) error
	DeleteOrganizationPermission(ctx context.Context, subjectID uuid.UUID) error
	CheckOrganizationPermission(ctx context.Context, subjectID uuid.UUID, objectID uuid.UUID, objectType string, relation string) (bool, error)
}

// AuthBasePermission is an interface representing the authbase permissions that are used in the service layer
type AuthBasePermission interface {
	MemberPermission
	UserPermission
	OrganizationPermission
}

// AuthZedPermission is a struct that implements the AuthBasePermission interface
// it delegates the permission checks to the Zed package
type AuthZedPermission struct {
}

func (a *AuthZedPermission) CreateUserPermission(ctx context.Context, subjectID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZedPermission) DeleteUserPermission(ctx context.Context, subjectID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZedPermission) CheckUserPermission(ctx context.Context, subjectID uuid.UUID, objectID uuid.UUID, objectType string, relation string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

var _ AuthBasePermission = new(AuthZedPermission)

func NewAuthZedPermission() *AuthZedPermission {
	return &AuthZedPermission{}
}

func (a *AuthZedPermission) CreateMemberPermission(ctx context.Context, subjectID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZedPermission) DeleteMemberPermission(ctx context.Context, subjectID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZedPermission) CheckMemberPermission(ctx context.Context, subjectID uuid.UUID, objectID uuid.UUID, objectType string, relation string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZedPermission) CreateOrganizationPermission(ctx context.Context, subjectID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZedPermission) DeleteOrganizationPermission(ctx context.Context, subjectID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (a *AuthZedPermission) CheckOrganizationPermission(ctx context.Context, subjectID uuid.UUID, objectID uuid.UUID, objectType string, relation string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
