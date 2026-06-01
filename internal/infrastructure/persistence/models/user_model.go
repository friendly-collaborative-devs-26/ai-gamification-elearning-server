package models

import "time"

type UserModel struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Username  string `gorm:"uniqueIndex;not null"`
	Email     string `gorm:"uniqueIndex;not null"`
	Password  *string
	Image     *string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string { return "users" }
