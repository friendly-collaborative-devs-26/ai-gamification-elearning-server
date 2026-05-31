package entities

import "time"

type ExecutionStatus string

const (
	ExecutionStatusInProgress ExecutionStatus = "in_progress"
	ExecutionStatusCompleted  ExecutionStatus = "completed"
	ExecutionStatusFailed     ExecutionStatus = "failed"
)

type ExerciseExecution struct {
	ID               uint
	ExerciseStatusID uint
	ExecutionTime    int
	Status           ExecutionStatus
	CreatedAt        time.Time
}
