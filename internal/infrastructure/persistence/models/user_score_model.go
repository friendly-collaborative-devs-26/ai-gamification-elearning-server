package models

import "time"

type UserScoreModel struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null;index"`
	BlockID    uint      `gorm:"not null;index"`
	TotalScore int       `gorm:"default:0"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`

	User  UserModel  `gorm:"foreignKey:UserID"`
	Block BlockModel `gorm:"foreignKey:BlockID"`
}

func (UserScoreModel) TableName() string { return "user_scores" }
