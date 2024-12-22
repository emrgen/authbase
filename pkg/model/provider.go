package model

type OauthProvider struct {
	ID             string `gorm:"unique;not null" json:"id"`
	Name           string `gorm:"primaryKey;not null"`
	OrganizationID string `gorm:"primaryKey;not null"`
	Config         string `gorm:"type:json"` // JSON string
}

func (OauthProvider) TableName() string {
	return tableName("oauth_providers")
}
