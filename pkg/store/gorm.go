package store

import (
	"context"
	"errors"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

var _ AuthBaseStore = new(GormStore)

type GormStore struct {
	db *gorm.DB
}

func (g *GormStore) GetUserCount(ctx context.Context, projectID uuid.UUID) (uint32, error) {
	var count int64
	g.db.Model(&model.User{}).Where("project_id = ?", projectID).Count(&count)

	return uint32(count), nil
}

func (g *GormStore) GetMemberCount(ctx context.Context, projectID uuid.UUID) (uint32, error) {
	var count int64
	g.db.Model(&model.ProjectMember{}).Where("project_id = ?", projectID).Count(&count)

	return uint32(count), nil
}

func (g *GormStore) ListProjectMembersUsers(ctx context.Context, orgID uuid.UUID, userIDs []uuid.UUID) ([]*model.ProjectMember, error) {
	var permissions []*model.ProjectMember
	err := g.db.Find(&permissions, "project_id = ? AND user_id IN ?", orgID, userIDs).Error
	return permissions, err
}

func (g *GormStore) GetMasterProject(ctx context.Context) (*model.Project, error) {
	var org model.Project
	err := g.db.Where("master = ?", true).First(&org).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrProjectNotFound
	}

	return &org, err
}

// DeleteSessionByUserID expire and delete all sessions for a user which not deleted or expired already
func (g *GormStore) DeleteSessionByUserID(ctx context.Context, userID uuid.UUID) error {
	return g.db.Model(&model.Session{}).
		Where("user_id = ? AND expired_at > ?", userID, time.Now()).
		Update("expired_at", time.Now()).
		Error
}

func (g *GormStore) ListSessions(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]*model.Session, error) {
	var sessions []*model.Session
	err := g.db.Limit(perPage).Offset(page*perPage).Preload("User").Find(&sessions, "project_id = ?", orgID).Error
	return sessions, err
}

func (g *GormStore) ListActiveSessions(ctx context.Context, userID uuid.UUID) ([]*model.Session, error) {
	var sessions []*model.Session
	err := g.db.Find(&sessions, "user_id = ? AND expired_at > ?)", userID, time.Now()).Error
	return sessions, err
}

func (g *GormStore) CreateSession(ctx context.Context, session *model.Session) error {
	return g.db.Create(session).Error
}

func (g *GormStore) DeleteSession(ctx context.Context, id uuid.UUID) error {
	session := model.Session{ID: id.String()}
	return g.db.Delete(&session).Error
}

func (g *GormStore) CreateVerificationCode(ctx context.Context, code *model.VerificationCode) error {
	return g.db.Create(code).Error
}

func (g *GormStore) GetVerificationCode(ctx context.Context, code string) (*model.VerificationCode, error) {
	var vc model.VerificationCode
	err := g.db.Where("code = ?", code).First(&vc).Error
	return &vc, err
}

func (g *GormStore) DeleteVerificationCode(ctx context.Context, code string) error {
	// hard delete the verification code
	return g.db.Delete(&model.VerificationCode{Code: code}).Error
}

func (g *GormStore) CreateToken(ctx context.Context, token *model.Token) error {
	return g.db.Create(token).Error
}

func (g *GormStore) GetTokenByID(ctx context.Context, id uuid.UUID) (*model.Token, error) {
	var token model.Token
	err := g.db.Where("id = ?", id).First(&token).Error
	return &token, err
}

func (g *GormStore) ListUserTokens(ctx context.Context, orgID, userID uuid.UUID, page, perPage int) ([]*model.Token, int, error) {
	var tokens []*model.Token
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Token{}).Count(&total).Error; err != nil {
			return err
		}
		err := tx.Limit(perPage).Offset(page*perPage).Find(&tokens, "project_id = ? AND user_id = ?", orgID, userID).Error

		return err
	})

	return tokens, int(total), err
}

func (g *GormStore) DeleteToken(ctx context.Context, id uuid.UUID) error {
	return g.db.Delete(&model.Token{ID: id.String()}).Error
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (g *GormStore) CreateUser(ctx context.Context, user *model.User) error {
	return g.db.Create(user).Error
}

func (g *GormStore) GetUserByEmail(ctx context.Context, orgID uuid.UUID, email string) (*model.User, error) {
	var user model.User
	err := g.db.Find(&user, "project_id = ? AND email = ?", orgID, email).Error
	return &user, err
}

func (g *GormStore) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := g.db.Where("id = ?", id.String()).Preload("Project").First(&user).Error
	return &user, err
}

func (g *GormStore) UpdateUser(ctx context.Context, user *model.User) error {
	return g.db.Save(user).Error
}

func (g *GormStore) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user := model.User{ID: id.String()}
	return g.db.Delete(&user).Error
}

func (g *GormStore) ListUsersByOrg(ctx context.Context, member bool, orgID uuid.UUID, page, perPage int) ([]*model.User, int, error) {
	var users []*model.User
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if member {
			if err := tx.Model(&model.User{}).Where("project_id = ? AND member = ?", orgID.String(), member).Count(&total).Error; err != nil {
				return err
			}

			members, err := g.ListProjectMembers(ctx, orgID, page, perPage)
			if err != nil {
				return err
			}

			for _, member := range members {
				users = append(users, member.User)
			}

			return nil
		} else {
			if err := tx.Model(&model.User{}).Where("project_id = ?", orgID.String()).Count(&total).Error; err != nil {
				return err
			}
			return g.db.Where("project_id = ?", orgID).Limit(perPage).Offset(page * perPage).Find(&users).Error
		}
	})

	return users, int(total), err
}

func (g *GormStore) DisableUser(ctx context.Context, id uuid.UUID) error {
	user := model.User{ID: id.String()}
	return g.db.Model(&user).Update("disabled", true).Update("disabled_at", gorm.Expr("NOW()")).Error
}

func (g *GormStore) EnableUser(ctx context.Context, id uuid.UUID) error {
	user := model.User{ID: id.String()}
	return g.db.Model(&user).Update("disabled", false).Update("disabled_at", nil).Error
}

func (g *GormStore) VerifyUser(ctx context.Context, id uuid.UUID) error {
	user := model.User{ID: id.String()}
	return g.db.Model(&user).Update("verified", true).Update("verified_at", gorm.Expr("NOW()")).Error
}

func (g *GormStore) UserExists(ctx context.Context, orgID uuid.UUID, username, email string) ([]*model.User, error) {
	var users []*model.User
	err := g.db.Where("project_id = ? AND (username = ? OR email = ?)", orgID, username, email).Find(&users).Error
	return users, err
}

func (g *GormStore) CreateProject(ctx context.Context, org *model.Project) error {
	err := g.db.Create(org).Error

	if err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return ErrProjectExists
		}
	}

	return err
}

// GetProjectByName retrieves an organization by name
// TODO: heavily rate limited and should be used with caution,
// maybe use reCAPTCHA to verify the user is not a bot
func (g *GormStore) GetProjectByName(ctx context.Context, name string) (*model.Project, error) {
	var org model.Project
	err := g.db.Where("name = ?", name).First(&org).Error
	return &org, err
}

func (g *GormStore) GetProjectByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	var org model.Project
	err := g.db.Where("id = ?", id).First(&org).Error
	return &org, err
}

func (g *GormStore) ListProjects(ctx context.Context, page, perPage int) ([]*model.Project, int, error) {
	var orgs []*model.Project
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Project{}).Count(&total).Error; err != nil {
			return err
		}
		return g.db.Limit(perPage).Offset(page * perPage).Find(&orgs).Error
	})

	return orgs, int(total), err
}

func (g *GormStore) UpdateProject(ctx context.Context, org *model.Project) error {
	return g.db.Save(org).Error
}

// DeleteProject deletes an organization from the database
func (g *GormStore) DeleteProject(ctx context.Context, id uuid.UUID) error {
	org := model.Project{ID: id.String()}
	return g.db.Delete(&org).Error
}

func (g *GormStore) CreateKeypair(ctx context.Context, keypair *model.Keypair) error {
	// NOTE: we should only have one keypair per project, so we can safely delete all existing keypairs and create a new one
	// this will cause all the existing tokens to be invalidated and the users will have to re-authenticate
	err := g.Transaction(func(tx AuthBaseStore) error {
		if err := tx.(*GormStore).db.Where("project_id = ?", keypair.ProjectID).Delete(&model.Keypair{}).Error; err != nil {
			return err
		}

		return tx.(*GormStore).db.Create(keypair).Error
	})

	return err
}

func (g *GormStore) GetKeypair(ctx context.Context, id uuid.UUID) (*model.Keypair, error) {
	var keypair model.Keypair
	err := g.db.Where("id = ?", id).First(&keypair).Error
	return &keypair, err
}

func (g *GormStore) CreateProjectMember(ctx context.Context, permission *model.ProjectMember) error {
	return g.db.Create(permission).Error
}

func (g *GormStore) GetProjectMemberByID(ctx context.Context, orgID, userID uuid.UUID) (*model.ProjectMember, error) {
	var permission model.ProjectMember
	err := g.db.Where("project_id = ? AND user_id = ?", orgID.String(), userID.String()).First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPermissionNotFound
	}

	return &permission, err
}

func (g *GormStore) ListProjectMembers(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.ProjectMember, error) {
	var permissions []*model.ProjectMember
	err := g.db.Where("project_id = ?", projectID).Preload("User").Limit(perPage).Offset(page * perPage).Order("permission DESC").Find(&permissions).Error
	return permissions, err
}

func (g *GormStore) UpdateProjectMember(ctx context.Context, permission *model.ProjectMember) error {
	return g.db.Save(permission).Error
}

func (g *GormStore) DeleteProjectMember(ctx context.Context, orgID, userID uuid.UUID) error {
	permission := model.ProjectMember{ProjectID: orgID.String(), UserID: userID.String()}
	return g.db.Delete(&permission).Error
}

func (g *GormStore) CreateOauthProvider(ctx context.Context, provider *model.OauthProvider) error {
	return g.db.Create(provider).Error
}

func (g *GormStore) GetOauthProviderByID(ctx context.Context, id uuid.UUID) (*model.OauthProvider, error) {
	var provider model.OauthProvider
	err := g.db.Where("id = ?", id).First(&provider).Error
	return &provider, err
}

// GetOauthProviderByName implements AuthBaseStore.
func (g *GormStore) GetOauthProviderByName(ctx context.Context, orgID uuid.UUID, provider string) (*model.OauthProvider, error) {
	var oauthProvider model.OauthProvider
	err := g.db.Where("project_id = ? AND provider = ?", orgID, provider).First(&oauthProvider).Error
	return &oauthProvider, err
}

func (g *GormStore) ListOauthProviders(ctx context.Context, orgID uuid.UUID, page, perPage int) ([]*model.OauthProvider, uint32, error) {
	var providers []*model.OauthProvider
	err := g.db.Limit(perPage).Offset(page*perPage).Find(&providers, "project_id = ?", orgID).Error
	if err != nil {
		return providers, 0, err
	}

	var total int64
	if err := g.db.Model(&model.OauthProvider{}).Where("project_id = ?", orgID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return providers, uint32(total), nil
}

func (g *GormStore) UpdateOauthProvider(ctx context.Context, provider *model.OauthProvider) error {
	return g.db.Save(provider).Error
}

func (g *GormStore) DeleteOauthProvider(ctx context.Context, id uuid.UUID) error {
	provider := model.OauthProvider{ID: id.String()}
	return g.db.Delete(&provider).Error
}

func (g *GormStore) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	return g.db.Create(token).Error
}

func (g *GormStore) GetRefreshTokenByID(ctx context.Context, refreshToken string) (*model.RefreshToken, error) {
	var token model.RefreshToken
	err := g.db.Where("token = ?", refreshToken).First(&token).Error
	return &token, err
}

func (g *GormStore) ListRefreshTokens(ctx context.Context, page, perPage int) ([]*model.RefreshToken, error) {
	var tokens []*model.RefreshToken
	err := g.db.Limit(perPage).Offset(page * perPage).Find(&tokens).Error
	return tokens, err
}

func (g *GormStore) UpdateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	return g.db.Save(token).Error
}

func (g *GormStore) DeleteRefreshToken(ctx context.Context, token string) error {
	return g.db.Delete(&model.RefreshToken{Token: token}).Error
}

func (g *GormStore) Migrate() error {
	return model.Migrate(g.db)
}

func (g *GormStore) Transaction(f func(AuthBaseStore) error) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		return f(&GormStore{db: tx})
	})
}
