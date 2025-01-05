package model

type ProjectMember struct {
	ProjectID  string `gorm:"primaryKey;not null"`
	UserID     string `gorm:"primaryKey;not null"`
	Permission uint32 `gorm:"not null;default:0"`

	Project *Project `gorm:"foreignKey:ProjectID;OnDelete:CASCADE"` // delete permissions when project is deleted
	User    *User    `gorm:"foreignKey:UserID;OnDelete:CASCADE"`    // delete permissions when user is deleted
}

func (ProjectMember) TableName() string {
	return tableName("project_members")
}
