package entities

import (
	"encoding/json"
	"time"
)

type StoreItem struct {
	ID            uint
	Name          string
	Description   string
	Price         int
	ItemType      string
	Configuration json.RawMessage
	IsActive      bool
	CreatedAt     time.Time
}
