package models

import "time"

type PermissionModel struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

func (PermissionModel) TableName() string { return "permissions" }
