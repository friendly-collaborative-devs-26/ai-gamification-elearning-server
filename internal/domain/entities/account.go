package entities

import "time"

type Account struct {
	UserID            uint
	Provider          string
	ProviderAccountID string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
