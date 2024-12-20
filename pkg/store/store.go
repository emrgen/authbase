package store

import (
	"context"
	"errors"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/google/uuid"
)

var (
	ErrOrganizationExists = errors.New("organization already exists")
)

// AuthBaseStore is the interface for interacting with the database.
type AuthBaseStore interface {
	UserStore
	SessionStore
	OrganizationStore
	PermissionStore
	ProviderStore
	RefreshTokenStore
	TokenStore
	VerificationCodeStore
	Migrate() error
	Transaction(func(AuthBaseStore) error) error
}

// OrganizationStore is the interface for interacting with the organization database.
type OrganizationStore interface {
	// CreateOrganization creates a new organization in the database.
	CreateOrganization(ctx context.Context, org *model.Organization) error
	// GetOrganizationByName retrieves an organization by its name.
	GetOrganizationByName(ctx context.Context, name string) (*model.Organization, error)
	// GetOrganizationByID retrieves an organization by its ID.
	GetOrganizationByID(ctx context.Context, id uuid.UUID) (*model.Organization, error)
	// ListOrganizations retrieves a list of organizations.
	ListOrganizations(ctx context.Context, page, perPage int) ([]*model.Organization, int, error)
	// UpdateOrganization updates an organization in the database.
	UpdateOrganization(ctx context.Context, org *model.Organization) error
	// DeleteOrganization deletes an organization from the database.
	DeleteOrganization(ctx context.Context, id uuid.UUID) error
}

// UserStore is the interface for interacting with the user database.
// User can be a member or an end user.
type UserStore interface {
	// CreateUser creates a new user in the database.
	CreateUser(ctx context.Context, user *model.User) error
	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	// GetUserByID retrieves a user by their ID.
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	// UpdateUser updates a user in the database.
	UpdateUser(ctx context.Context, user *model.User) error
	// DeleteUser deletes a user from the database.
	DeleteUser(ctx context.Context, id uuid.UUID) error
	// ListUsersByOrg retrieves a list of users by organization.
	ListUsersByOrg(ctx context.Context, member bool, orgID uuid.UUID, page, perPage int) ([]*model.User, int, error)
	// DisableUser disables a user in the database.
	DisableUser(ctx context.Context, id uuid.UUID) error
	// EnableUser enables a user in the database.
	EnableUser(ctx context.Context, id uuid.UUID) error
	// VerifyUser verifies a user in the database.
	VerifyUser(ctx context.Context, id uuid.UUID) error
	// UserExists checks if a user exists in the database.
	UserExists(ctx context.Context, orgID uuid.UUID, username, email string) ([]*model.User, error)
}

type SessionStore interface {
	// CreateSession creates a new session in the database.
	CreateSession(ctx context.Context, session *model.Session) error
	// ListSessions retrieves a list of sessions.
	ListSessions(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]*model.Session, error)
	// DeleteSession deletes a session from the database.
	DeleteSession(ctx context.Context, id uuid.UUID) error
	// DeleteSessionByUserID deletes a session from the database by user ID.
	DeleteSessionByUserID(ctx context.Context, userID uuid.UUID) error
}

// PermissionStore is the interface for interacting with the permission database.
type PermissionStore interface {
	// CreatePermission creates a new permission in the database.
	CreatePermission(ctx context.Context, permission *model.Permission) error
	// GetPermissionByID retrieves a permission by its ID.
	GetPermissionByID(ctx context.Context, orgID, userID uuid.UUID) (*model.Permission, error)
	// ListPermissions retrieves a list of permissions.
	ListPermissions(ctx context.Context, page, perPage int) ([]*model.Permission, error)
	// UpdatePermission updates a permission in the database.
	UpdatePermission(ctx context.Context, permission *model.Permission) error
	// DeletePermission deletes a permission from the database.
	DeletePermission(ctx context.Context, orgID, userID uuid.UUID) error
}

// ProviderStore is the interface for interacting with the provider database.
type ProviderStore interface {
	// CreateProvider creates a new provider in the database.
	CreateProvider(ctx context.Context, provider *model.Provider) error
	// GetProviderByID retrieves a provider by its ID.
	GetProviderByID(ctx context.Context, id uuid.UUID) (*model.Provider, error)
	// ListProviders retrieves a list of providers.
	ListProviders(ctx context.Context, page, perPage int) ([]*model.Provider, error)
	// UpdateProvider updates a provider in the database.
	UpdateProvider(ctx context.Context, provider *model.Provider) error
	// DeleteProvider deletes a provider from the database.
	DeleteProvider(ctx context.Context, id uuid.UUID) error
}

// RefreshTokenStore is the interface for interacting with the refresh token database.
type RefreshTokenStore interface {
	// CreateRefreshToken creates a new refresh token in the database.
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	// GetRefreshTokenByID retrieves a
	GetRefreshTokenByID(ctx context.Context, token string) (*model.RefreshToken, error)
	// ListRefreshTokens retrieves a list of refresh tokens.
	ListRefreshTokens(ctx context.Context, page, perPage int) ([]*model.RefreshToken, error)
	// UpdateRefreshToken updates a refresh token in the database.
	UpdateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	// DeleteRefreshToken deletes a refresh token from the database.
	DeleteRefreshToken(ctx context.Context, token string) error
}

// TokenStore is the interface for interacting with the token database.
type TokenStore interface {
	// CreateToken creates a new token in the database.
	CreateToken(ctx context.Context, token *model.Token) error
	// GetTokenByID retrieves a token by its ID.
	GetTokenByID(ctx context.Context, id uuid.UUID) (*model.Token, error)
	// ListUserTokens retrieves a list of tokens by user.
	ListUserTokens(ctx context.Context, orgID, userID uuid.UUID, page, perPage int) ([]*model.Token, int, error)
	// DeleteToken updates a token in the database.
	DeleteToken(ctx context.Context, id uuid.UUID) error
}

type VerificationCodeStore interface {
	// CreateVerificationCode creates a new code
	CreateVerificationCode(ctx context.Context, code *model.VerificationCode) error
	// GetVerificationCode retrieves a code by hash id
	GetVerificationCode(ctx context.Context, code string) (*model.VerificationCode, error)
	// DeleteVerificationCode deletes the code
	DeleteVerificationCode(ctx context.Context, code string) error
}
