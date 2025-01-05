package store

import (
	"context"
	"errors"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/google/uuid"
)

var (
	ErrProjectExists           = errors.New("project already exists")
	ErrProjectNotFound         = errors.New("project not found")
	ErrPermissionNotFound      = errors.New("permission not found")
	ErrPermissionAlreadyExists = errors.New("permission already exists")
)

// AuthBaseStore is the interface for interacting with the database.
type AuthBaseStore interface {
	UserStore
	SessionStore
	ProjectStore
	ProjectMemberStore
	ProviderStore
	RefreshTokenStore
	TokenStore
	VerificationCodeStore
	Migrate() error
	Transaction(func(AuthBaseStore) error) error
}

// ProjectStore is the interface for interacting with the project database.
type ProjectStore interface {
	// CreateProject creates a new project in the database.
	CreateProject(ctx context.Context, org *model.Project) error
	// GetProjectByName retrieves an project by its name.
	GetProjectByName(ctx context.Context, name string) (*model.Project, error)
	// GetProjectByID retrieves an project by its ID.
	GetProjectByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	// GetMasterProject retrieves the master project.
	GetMasterProject(ctx context.Context) (*model.Project, error)
	// ListProjects retrieves a list of projects.
	ListProjects(ctx context.Context, page, perPage int) ([]*model.Project, int, error)
	// UpdateProject updates an project in the database.
	UpdateProject(ctx context.Context, org *model.Project) error
	// DeleteProject deletes an project from the database.
	DeleteProject(ctx context.Context, id uuid.UUID) error
	// CreateKeypair creates a new keypair in the database.
	CreateKeypair(ctx context.Context, keypair *model.Keypair) error
	// GetKeypair retrieves a keypair by its ID.
	GetKeypair(ctx context.Context, id uuid.UUID) (*model.Keypair, error)
}

// UserStore is the interface for interacting with the user database.
// User can be a member or an end user.
type UserStore interface {
	// CreateUser creates a new user in the database.
	CreateUser(ctx context.Context, user *model.User) error
	// GetUserByEmail retrieves a user by their email address.
	GetUserByEmail(ctx context.Context, orgID uuid.UUID, email string) (*model.User, error)
	// GetUserByID retrieves a user by their ID.
	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	// UpdateUser updates a user in the database.
	UpdateUser(ctx context.Context, user *model.User) error
	// DeleteUser deletes a user from the database.
	DeleteUser(ctx context.Context, id uuid.UUID) error
	// ListUsersByOrg retrieves a list of users by project.
	ListUsersByOrg(ctx context.Context, member bool, projectID uuid.UUID, page, perPage int) ([]*model.User, int, error)
	// DisableUser disables a user in the database.
	DisableUser(ctx context.Context, id uuid.UUID) error
	// EnableUser enables a user in the database.
	EnableUser(ctx context.Context, id uuid.UUID) error
	// VerifyUser verifies a user in the database.
	VerifyUser(ctx context.Context, id uuid.UUID) error
	// UserExists checks if a user exists in the database.
	UserExists(ctx context.Context, projectID uuid.UUID, username, email string) ([]*model.User, error)
	// GetUserCount retrieves the number of users in a project.
	GetUserCount(ctx context.Context, projectID uuid.UUID) (uint32, error)
}

type SessionStore interface {
	// CreateSession creates a new session in the database.
	CreateSession(ctx context.Context, session *model.Session) error
	// ListSessions retrieves a list of sessions.
	ListSessions(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]*model.Session, error)
	// ListActiveSessions retrieves a list of active sessions.
	ListActiveSessions(ctx context.Context, userID uuid.UUID) ([]*model.Session, error)
	// DeleteSession deletes a session from the database.
	DeleteSession(ctx context.Context, sessionID uuid.UUID) error
	// DeleteSessionByUserID deletes a session from the database by user ID.
	DeleteSessionByUserID(ctx context.Context, userID uuid.UUID) error
}

// ProjectMemberStore is the interface for interacting with the permission database.
type ProjectMemberStore interface {
	// CreateProjectMember creates a new permission in the database.
	CreateProjectMember(ctx context.Context, permission *model.ProjectMember) error
	// GetProjectMemberByID retrieves a permission by its ID.
	GetProjectMemberByID(ctx context.Context, orgID, userID uuid.UUID) (*model.ProjectMember, error)
	// ListProjectMembers retrieves a list of permissions.
	ListProjectMembers(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.ProjectMember, error)
	// ListProjectMembersUsers retrieves a list of permissions by project.
	ListProjectMembersUsers(ctx context.Context, orgID uuid.UUID, userIDs []uuid.UUID) ([]*model.ProjectMember, error)
	// UpdateProjectMember updates a permission in the database.
	UpdateProjectMember(ctx context.Context, permission *model.ProjectMember) error
	// DeleteProjectMember deletes a permission from the database.
	DeleteProjectMember(ctx context.Context, orgID, userID uuid.UUID) error
	// GetMemberCount retrieves the number of members in a project.
	GetMemberCount(ctx context.Context, projectID uuid.UUID) (uint32, error)
}

// ProviderStore is the interface for interacting with the provider database.
type ProviderStore interface {
	// CreateOauthProvider CreateProvider creates a new provider in the database.
	CreateOauthProvider(ctx context.Context, provider *model.OauthProvider) error
	// GetOauthProviderByID retrieves a provider by its ID.
	GetOauthProviderByID(ctx context.Context, id uuid.UUID) (*model.OauthProvider, error)
	// GetOauthProviderByName retrieves a provider by its name.
	GetOauthProviderByName(ctx context.Context, orgID uuid.UUID, provider string) (*model.OauthProvider, error)
	// ListOauthProviders retrieves a list of providers.
	ListOauthProviders(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]*model.OauthProvider, uint32, error)
	// UpdateOauthProvider updates a provider in the database.
	UpdateOauthProvider(ctx context.Context, provider *model.OauthProvider) error
	// DeleteOauthProvider deletes a provider from the database.
	DeleteOauthProvider(ctx context.Context, id uuid.UUID) error
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
