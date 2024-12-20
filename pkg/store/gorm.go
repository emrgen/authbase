package store

import (
	"context"
	"errors"
	"github.com/emrgen/authbase/pkg/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ AuthBaseStore = new(GormStore)

type GormStore struct {
	db *gorm.DB
}

func (g *GormStore) CreateVerificationCode(ctx context.Context, code *model.VerificationCode) error {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) GetVerificationCode(ctx context.Context, code string) (*model.VerificationCode, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GormStore) DeleteVerificationCode(ctx context.Context, code string) error {
	//TODO implement me
	panic("implement me")
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
		err := tx.Limit(perPage).Offset(page*perPage).Find(&tokens, "organization_id = ? AND user_id = ?", orgID, userID).Error

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

func (g *GormStore) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := g.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (g *GormStore) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := g.db.Where("id = ?", id).First(&user).Error
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
		if err := tx.Model(&model.User{}).Where("organization_id = ? AND member = ?", orgID, member).Count(&total).Error; err != nil {
			return err
		}
		return g.db.Where("organization_id = ? AND member = ?", orgID, member).Limit(perPage).Offset(page * perPage).Find(&users).Error
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
	err := g.db.Where("organization_id = ? AND (username = ? OR email = ?)", orgID, username, email).Find(&users).Error
	return users, err
}

func (g *GormStore) CreateOrganization(ctx context.Context, org *model.Organization) error {
	err := g.db.Create(org).Error

	if err != nil {
		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return ErrOrganizationExists
		}
	}

	return err
}

func (g *GormStore) GetOrganizationByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	err := g.db.Where("id = ?", id).First(&org).Error
	return &org, err
}

func (g *GormStore) ListOrganizations(ctx context.Context, page, perPage int) ([]*model.Organization, int, error) {
	var orgs []*model.Organization
	var total int64

	err := g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&model.Organization{}).Count(&total).Error; err != nil {
			return err
		}
		return g.db.Limit(perPage).Offset(page * perPage).Find(&orgs).Error
	})

	return orgs, int(total), err
}

func (g *GormStore) UpdateOrganization(ctx context.Context, org *model.Organization) error {
	return g.db.Save(org).Error
}

func (g *GormStore) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	org := model.Organization{ID: id.String()}
	return g.db.Delete(&org).Error
}

func (g *GormStore) CreatePermission(ctx context.Context, permission *model.Permission) error {
	return g.db.Create(permission).Error
}

func (g *GormStore) GetPermissionByID(ctx context.Context, orgID, userID uuid.UUID) (*model.Permission, error) {
	var permission model.Permission
	err := g.db.Where("organization_id = ? AND user_id = ?", orgID, userID).First(&permission).Error
	return &permission, err
}

func (g *GormStore) ListPermissions(ctx context.Context, page, perPage int) ([]*model.Permission, error) {
	var permissions []*model.Permission
	err := g.db.Limit(perPage).Offset(page * perPage).Find(&permissions).Error
	return permissions, err
}

func (g *GormStore) UpdatePermission(ctx context.Context, permission *model.Permission) error {
	return g.db.Save(permission).Error
}

func (g *GormStore) DeletePermission(ctx context.Context, orgID, userID uuid.UUID) error {
	permission := model.Permission{OrganizationID: orgID.String(), UserID: userID.String()}
	return g.db.Delete(&permission).Error
}

func (g *GormStore) CreateProvider(ctx context.Context, provider *model.Provider) error {
	return g.db.Create(provider).Error
}

func (g *GormStore) GetProviderByID(ctx context.Context, id uuid.UUID) (*model.Provider, error) {
	var provider model.Provider
	err := g.db.Where("id = ?", id).First(&provider).Error
	return &provider, err
}

func (g *GormStore) ListProviders(ctx context.Context, page, perPage int) ([]*model.Provider, error) {
	var providers []*model.Provider
	err := g.db.Limit(perPage).Offset(page * perPage).Find(&providers).Error
	return providers, err
}

func (g *GormStore) UpdateProvider(ctx context.Context, provider *model.Provider) error {
	return g.db.Save(provider).Error
}

func (g *GormStore) DeleteProvider(ctx context.Context, id uuid.UUID) error {
	provider := model.Provider{ID: id.String()}
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
