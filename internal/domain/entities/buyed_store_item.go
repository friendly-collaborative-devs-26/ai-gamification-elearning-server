package entities

import (
	"time"
)

type BuyedStoreItem struct {
	ID          uint
	UserID      uint
	ItemID      uint
	PurchasedAt time.Time
}
