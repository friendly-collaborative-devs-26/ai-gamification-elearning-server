package models

import "time"

type BuyedStoreItemModel struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UserID      uint      `gorm:"not null;index"`
	ItemID      uint      `gorm:"not null;index"`
	PurchasedAt time.Time `gorm:"autoCreateTime"`

	User      UserModel      `gorm:"foreignKey:UserID"`
	StoreItem StoreItemModel `gorm:"foreignKey:ItemID"`
}

func (BuyedStoreItemModel) TableName() string { return "buyed_store_items" }
