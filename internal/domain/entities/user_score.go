package entities

import "time"

type UserScore struct {
	ID         uint
	UserID     uint
	BlockID    uint
	TotalScore int
	UpdatedAt  time.Time
}
