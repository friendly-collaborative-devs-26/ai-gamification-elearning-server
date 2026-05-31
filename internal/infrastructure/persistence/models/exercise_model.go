package models

import (
	"encoding/json"
	"time"
)

type ExerciseModel struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	BlockID       uint   `gorm:"not null;index"`
	CreatorID     uint   `gorm:"not null;index"`
	Title         string `gorm:"not null"`
	Description   string
	ExerciseURL   string
	Configuration json.RawMessage `gorm:"type:jsonb"`
	Difficulty    string          `gorm:"not null"`
	Score         int             `gorm:"default:0"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
	UpdatedAt     time.Time       `gorm:"autoUpdateTime"`

	Block   BlockModel `gorm:"foreignKey:BlockID"`
	Creator UserModel  `gorm:"foreignKey:CreatorID"`
}

func (ExerciseModel) TableName() string { return "exercises" }
