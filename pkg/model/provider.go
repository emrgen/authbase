package model

type OAuthConfig struct {
	Provider     string `json:"provider"`
	RedirectURL  string `json:"redirect_url"`
	CallbackURL  string `json:"callback_url"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scopes       string `json:"scopes"`
}

type OauthProvider struct {
	ID       string      `gorm:"unique;not null" json:"id"`
	Provider string      `gorm:"primaryKey;not null"`
	PoolID   string      `gorm:"primaryKey;not null"`
	Pool     *Pool       `gorm:"foreignKey:PoolID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Config   OAuthConfig `gorm:"embedded;embeddedPrefix:config_"`
}

func (OauthProvider) TableName() string {
	return tableName("oauth_providers")
}
