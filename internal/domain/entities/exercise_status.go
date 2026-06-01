package entities

import "time"

type Verdict string

const (
	VerdictApproved Verdict = "approved"
	VerdictRejected Verdict = "rejected"
)

type ExerciseStatus struct {
	ID              uint
	UserID          uint
	ExerciseID      uint
	ExerciseForkURL string
	Verdict         Verdict
	UpdatedAt       time.Time
}
