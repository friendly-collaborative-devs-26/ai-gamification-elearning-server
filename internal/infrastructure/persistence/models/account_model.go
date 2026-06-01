package models

import "time"

type AccountModel struct {
	UserID            uint      `gorm:"not null;index"`
	Provider          string    `gorm:"not null"`
	ProviderAccountID string    `gorm:"not null"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`

	User UserModel `gorm:"foreignKey:UserID"`
}

func (AccountModel) TableName() string { return "accounts" }
