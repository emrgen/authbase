package model

type Permission struct {
	OrganizationID string `gorm:"primaryKey;not null"`
	UserID         string `gorm:"primaryKey;not null"`
	Permission     uint32 `gorm:"not null;default:0"`

	Organization *Organization `gorm:"foreignKey:OrganizationID;OnDelete:CASCADE"` // delete permissions when organization is deleted
	User         *User         `gorm:"foreignKey:UserID;OnDelete:CASCADE"`         // delete permissions when user is deleted
}

func (Permission) TableName() string {
	return tableName("permissions")
}
