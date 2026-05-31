package models

import "time"

type UserPermissionModel struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	UserID       uint      `gorm:"not null;index"`
	PermissionID uint      `gorm:"not null;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	User       UserModel       `gorm:"foreignKey:UserID"`
	Permission PermissionModel `gorm:"foreignKey:PermissionID"`
}

func (UserPermissionModel) TableName() string { return "user_permissions" }
