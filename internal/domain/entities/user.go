package entities

import "time"

type User struct {
	ID        uint
	Username  string
	Email     string
	Password  *string
	Image     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
