package entities

import "time"

type Permission struct {
	ID          uint
	Name        string
	Description string
	CreatedAt   time.Time
}
