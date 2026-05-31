package models

import (
	"encoding/json"
	"time"
)

type StoreItemModel struct {
	ID            uint   `gorm:"primaryKey;autoIncrement"`
	Name          string `gorm:"not null"`
	Description   string
	Price         int             `gorm:"not null"`
	ItemType      string          `gorm:"not null"`
	Configuration json.RawMessage `gorm:"type:jsonb"`
	IsActive      bool            `gorm:"default:true"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
}

func (StoreItemModel) TableName() string { return "store_items" }
