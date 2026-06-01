package entities

import "time"

type UserPermission struct {
	ID           uint
	UserID       uint
	PermissionID uint
	CreatedAt    time.Time
}
