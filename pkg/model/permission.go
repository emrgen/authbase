package model

type Permission struct {
	OrganizationID string `gorm:"primaryKey;not null"`
	UserID         string `gorm:"primaryKey;not null"`
	Permission     uint32 `gorm:"not null;default:0"`

	Organization *Organization `gorm:"foreignKey:OrganizationID"`
	User         *User         `gorm:"foreignKey:UserID"`
}

func (_ Permission) TableName() string {
	return tableName("permissions")
}
