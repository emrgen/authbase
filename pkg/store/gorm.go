package store

import (
	"context"
	"errors"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// NewGormStore creates a new GormStore.
func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

var _ AuthBaseStore = new(GormStore)

// GormStore is a Gorm implementation of the store.
type GormStore struct {
	db *gorm.DB
}

func (g *GormStore) ListRolesByNames(ctx context.Context, poolID uuid.UUID, names []string) ([]*model.Role, error) {
	var roles []*model.Role
	err := g.db.Find(&roles, "pool_id = ? AND name IN ?", poolID.String(), names).Error
	return roles, err
}

func (g *GormStore) CreateRole(ctx context.Context, role *model.Role) error {
	return g.db.Create(role).Error
}

func (g *GormStore) GetRole(ctx context.Context, poolID uuid.UUID, name string) (*model.Role, error) {
	var role model.Role
	err := g.db.Where("name = ? AND pool_id = ?", name, poolID.String()).First(&role).Error
	if role.Name == "" {
		return nil, ErrRoleNotFound
	}

	return &role, err
}

func (g *GormStore) ListRoles(ctx context.Context, poolID uuid.UUID, page, perPage int) ([]*model.Role, int, error) {
	var roles []*model.Role
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Role{}).Where("pool_id = ?", poolID.String()).Count(&total).Error; err != nil {
			return err
		}
		return tx.Limit(perPage).Offset(page*perPage).Find(&roles, "pool_id = ?", poolID.String()).Error
	})

	return roles, int(total), err

}

func (g *GormStore) UpdateRole(ctx context.Context, role *model.Role) error {
	return g.db.Save(role).Error
}

func (g *GormStore) DeleteRole(ctx context.Context, poolID uuid.UUID, name string) error {
	role := model.Role{Name: name, PoolID: poolID.String()}
	return g.db.Delete(&role).Error
}

func (g *GormStore) CreateGroup(ctx context.Context, group *model.Group) error {
	return g.db.Create(group).Error
}

func (g *GormStore) GetGroup(ctx context.Context, id uuid.UUID) (*model.Group, error) {
	var group model.Group
	err := g.db.Where("id = ?", id).Preload("Roles").First(&group).Error
	return &group, err
}

func (g *GormStore) ListGroups(ctx context.Context, poolID uuid.UUID, page, perPage int) ([]*model.Group, int, error) {
	var groups []*model.Group
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Group{}).Preload("Roles").Where("pool_id = ?", poolID.String()).Count(&total).Error; err != nil {
			return err
		}
		return tx.Limit(perPage).Offset(page*perPage).Preload("Roles").Find(&groups, "pool_id = ?", poolID.String()).Error
	})

	return groups, int(total), err
}

func (g *GormStore) UpdateGroup(ctx context.Context, group *model.Group) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(group).Association("Roles").Replace(group.Roles); err != nil {
			return err
		}
		return tx.Save(group).Error
	})
}

func (g *GormStore) DeleteGroup(ctx context.Context, id uuid.UUID) error {
	group := model.Group{ID: id.String()}
	return g.db.Delete(&group).Error
}

func (g *GormStore) AddGroupMember(ctx context.Context, member *model.GroupMember) error {
	return g.db.Create(member).Error
}

func (g *GormStore) ListGroupMemberByAccount(ctx context.Context, accountID uuid.UUID) ([]*model.GroupMember, error) {
	var groups []*model.GroupMember
	err := g.db.Where("account_id = ?", accountID.String()).Preload("Group.Roles").Find(&groups).Error
	return groups, err
}

func (g *GormStore) RemoveGroupMember(ctx context.Context, groupID, accountID uuid.UUID) error {
	member := model.GroupMember{GroupID: groupID.String(), AccountID: accountID.String()}
	return g.db.Delete(&member).Error
}

func (g *GormStore) ListGroupMembers(ctx context.Context, groupID uuid.UUID, page, perPage int) ([]*model.GroupMember, int, error) {
	var members []*model.GroupMember
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.GroupMember{}).Preload("Account").Where("group_id = ?", groupID).Count(&total).Error; err != nil {
			return err
		}
		return tx.Limit(perPage).Offset(page*perPage).Preload("Account").Find(&members, "group_id = ?", groupID).Error
	})

	return members, int(total), err
}

func (g *GormStore) AddPoolMember(ctx context.Context, member *model.PoolMember) error {
	return g.db.Create(member).Error
}

func (g *GormStore) GetPoolMember(ctx context.Context, poolID, accountID uuid.UUID) (*model.PoolMember, error) {
	var member model.PoolMember
	err := g.db.Where("pool_id = ? AND account_id = ?", poolID, accountID).First(&member).Error
	return &member, err
}

func (g *GormStore) ListPoolMembers(ctx context.Context, poolID uuid.UUID, page, perPage int) ([]*model.PoolMember, int, error) {
	var members []*model.PoolMember
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.PoolMember{}).Where("pool_id = ?", poolID).Count(&total).Error; err != nil {
			return err
		}
		return tx.Limit(perPage).Offset(page*perPage).Find(&members, "pool_id = ?", poolID).Error
	})

	return members, int(total), err
}

func (g *GormStore) UpdatePoolMember(ctx context.Context, member *model.PoolMember) error {
	return g.db.Save(member).Error
}

func (g *GormStore) RemovePoolMember(ctx context.Context, poolID, accountID uuid.UUID) error {
	member := model.PoolMember{PoolID: poolID.String(), AccountID: accountID.String()}
	return g.db.Delete(&member).Error
}

func (g *GormStore) CreatePool(ctx context.Context, pool *model.Pool) error {
	return g.db.Create(pool).Error
}

func (g *GormStore) GetMasterPool(ctx context.Context, projectID uuid.UUID) (*model.Pool, error) {
	var pool model.Pool
	err := g.db.Where("project_id = ? AND master = ?", projectID, true).First(&pool).Error
	return &pool, err
}

func (g *GormStore) GetPoolByID(ctx context.Context, id uuid.UUID) (*model.Pool, error) {
	var pool model.Pool
	err := g.db.Where("id = ?", id).First(&pool).Error
	return &pool, err
}

func (g *GormStore) ListPools(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.Pool, int, error) {
	var pools []*model.Pool
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Pool{}).Where("project_id = ?", projectID).Count(&total).Error; err != nil {
			return err
		}
		return tx.Limit(perPage).Offset(page*perPage).Find(&pools, "project_id = ?", projectID).Error
	})

	return pools, int(total), err
}

func (g *GormStore) UpdatePool(ctx context.Context, pool *model.Pool) error {
	return g.db.Save(pool).Error
}

func (g *GormStore) DeletePool(ctx context.Context, id uuid.UUID) error {
	pool := model.Pool{ID: id.String()}
	return g.db.Delete(&pool).Error
}

func (g *GormStore) GetAccountCount(ctx context.Context, projectID uuid.UUID) (uint32, error) {
	var count int64
	g.db.Model(&model.Account{}).Where("project_id = ?", projectID).Count(&count)

	return uint32(count), nil
}

func (g *GormStore) GetMemberCount(ctx context.Context, projectID uuid.UUID) (uint32, error) {
	var count int64
	g.db.Model(&model.ProjectMember{}).Where("project_id = ?", projectID).Count(&count)

	return uint32(count), nil
}

func (g *GormStore) ListProjectMembersByAccountIDs(ctx context.Context, projectID uuid.UUID, accountIDs []uuid.UUID) ([]*model.ProjectMember, error) {
	var permissions []*model.ProjectMember
	err := g.db.Find(&permissions, "project_id = ? AND user_id IN ?", projectID, accountIDs).Error
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

func (g *GormStore) CreateClient(ctx context.Context, client *model.Client) error {
	return g.db.Create(client).Error
}

func (g *GormStore) GetClientByID(ctx context.Context, id uuid.UUID) (*model.Client, error) {
	var client model.Client
	err := g.db.Where("id = ?", id).Preload("Pool").Preload("CreatedByAccount").First(&client).Error
	return &client, err
}

func (g *GormStore) ListClients(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.Client, int, error) {
	var clients []*model.Client
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Client{}).Where("pool_id = ?", projectID).Count(&total).Error; err != nil {
			return err
		}
		return tx.Limit(perPage).Offset(page*perPage).Find(&clients, "pool_id = ?", projectID).Error
	})

	return clients, int(total), err
}

func (g *GormStore) UpdateClient(ctx context.Context, client *model.Client) error {
	return g.db.Save(client).Error
}

func (g *GormStore) DeleteClient(ctx context.Context, id uuid.UUID) error {
	client := model.Client{ID: id.String()}
	return g.db.Delete(&client).Error
}

// DeleteSessionByAccountID expire and delete all sessions for a user which not deleted or expired already
func (g *GormStore) DeleteSessionByAccountID(ctx context.Context, userID uuid.UUID) error {
	return g.db.Model(&model.Session{}).
		Where("user_id = ? AND expired_at > ?", userID, time.Now()).
		Update("expired_at", time.Now()).
		Error
}

func (g *GormStore) ListActiveAccounts(ctx context.Context, poolID uuid.UUID, page, perPage int) ([]*model.Session, error) {
	var sessions []*model.Session
	err := g.db.Limit(perPage).Offset(page*perPage).Select("DISTINCT account_id").Preload("Account").Find(&sessions, "pool_id = ?", poolID.String()).Error
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

func (g *GormStore) CreateAccessKey(ctx context.Context, token *model.AccessKey) error {
	return g.db.Create(token).Error
}

func (g *GormStore) GetAccessKeyByID(ctx context.Context, id uuid.UUID) (*model.AccessKey, error) {
	var token model.AccessKey
	err := g.db.Where("id = ?", id).First(&token).Error
	return &token, err
}

func (g *GormStore) ListAccountAccessKeys(ctx context.Context, orgID, userID uuid.UUID, page, perPage int) ([]*model.AccessKey, int, error) {
	var tokens []*model.AccessKey
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.AccessKey{}).Count(&total).Error; err != nil {
			return err
		}
		err := tx.Limit(perPage).Offset(page*perPage).Find(&tokens, "project_id = ? AND account_id = ?", orgID, userID).Error

		return err
	})

	return tokens, int(total), err
}

func (g *GormStore) DeleteAccessKey(ctx context.Context, id uuid.UUID) error {
	return g.db.Delete(&model.AccessKey{ID: id.String()}).Error
}

func (g *GormStore) CreateAccount(ctx context.Context, user *model.Account) error {
	return g.db.Create(user).Error
}

func (g *GormStore) GetAccountByEmail(ctx context.Context, poolID uuid.UUID, email string) (*model.Account, error) {
	var user model.Account
	err := g.db.Find(&user, "pool_id = ? AND email = ?", poolID, email).Error
	return &user, err
}

func (g *GormStore) GetAccountByID(ctx context.Context, id uuid.UUID) (*model.Account, error) {
	var user model.Account
	err := g.db.Where("id = ?", id.String()).Preload("Project").First(&user).Error
	return &user, err
}

func (g *GormStore) UpdateAccount(ctx context.Context, user *model.Account) error {
	return g.db.Save(user).Error
}

func (g *GormStore) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	user := model.Account{ID: id.String()}
	return g.db.Delete(&user).Error
}

func (g *GormStore) ListProjectAccounts(ctx context.Context, member bool, projectID uuid.UUID, page, perPage int) ([]*model.Account, int, error) {
	var users []*model.Account
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if member {
			if err := tx.Model(&model.Account{}).Where("project_id = ? AND member = ?", projectID.String(), member).Count(&total).Error; err != nil {
				return err
			}

			members, err := g.ListProjectMembers(ctx, projectID, page, perPage)
			if err != nil {
				return err
			}

			for _, member := range members {
				users = append(users, member.Account)
			}

			return nil
		} else {
			if err := tx.Model(&model.Account{}).Where("project_id = ?", projectID.String()).Count(&total).Error; err != nil {
				return err
			}
			return g.db.Where("project_id = ?", projectID).Limit(perPage).Offset(page * perPage).Find(&users).Error
		}
	})

	return users, int(total), err
}

func (g *GormStore) ListPoolAccounts(ctx context.Context, member bool, poolID uuid.UUID, page, perPage int) ([]*model.Account, int, error) {
	var users []*model.Account
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if member {
			if err := tx.Model(&model.Account{}).Where("pool_id = ? AND member = ?", poolID.String(), member).Count(&total).Error; err != nil {
				return err
			}

			members, err := g.ListProjectMembers(ctx, poolID, page, perPage)
			if err != nil {
				return err
			}

			for _, member := range members {
				users = append(users, member.Account)
			}

			return nil
		} else {
			if err := tx.Model(&model.Account{}).Where("pool_id = ?", poolID.String()).Count(&total).Error; err != nil {
				return err
			}
			return g.db.Where("pool_id = ?", poolID).Limit(perPage).Offset(page * perPage).Find(&users).Error
		}
	})

	return users, int(total), err
}

func (g *GormStore) DisableAccount(ctx context.Context, id uuid.UUID) error {
	user := model.Account{ID: id.String()}
	return g.db.Model(&user).Update("disabled", true).Update("disabled_at", gorm.Expr("NOW()")).Error
}

func (g *GormStore) EnableAccount(ctx context.Context, id uuid.UUID) error {
	user := model.Account{ID: id.String()}
	return g.db.Model(&user).Update("disabled", false).Update("disabled_at", nil).Error
}

func (g *GormStore) VerifyAccount(ctx context.Context, id uuid.UUID) error {
	user := model.Account{ID: id.String()}
	return g.db.Model(&user).Update("verified", true).Update("verified_at", gorm.Expr("NOW()")).Error
}

func (g *GormStore) AccountExists(ctx context.Context, orgID uuid.UUID, username, email string) ([]*model.Account, error) {
	var users []*model.Account
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
		if err := tx.(*GormStore).db.Where("client_id = ?", keypair.ClientID).Delete(&model.Keypair{}).Error; err != nil {
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
	err := g.db.Where("project_id = ? AND account_id = ?", orgID.String(), userID.String()).First(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPermissionNotFound
	}

	return &permission, err
}

func (g *GormStore) ListProjectMembers(ctx context.Context, projectID uuid.UUID, page, perPage int) ([]*model.ProjectMember, error) {
	var permissions []*model.ProjectMember
	err := g.db.Where("project_id = ?", projectID).Preload("Account").Limit(perPage).Offset(page * perPage).Order("permission DESC").Find(&permissions).Error
	return permissions, err
}

func (g *GormStore) UpdateProjectMember(ctx context.Context, permission *model.ProjectMember) error {
	return g.db.Save(permission).Error
}

func (g *GormStore) DeleteProjectMember(ctx context.Context, orgID, userID uuid.UUID) error {
	permission := model.ProjectMember{ProjectID: orgID.String(), AccountID: userID.String()}
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
