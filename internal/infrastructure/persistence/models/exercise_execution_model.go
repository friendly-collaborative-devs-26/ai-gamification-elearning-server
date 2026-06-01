package models

import "time"

type ExerciseExecutionStatusModel string

const (
	ExecutionStatusInProgress ExerciseExecutionStatusModel = "in_progress"
	ExecutionStatusCompleted  ExerciseExecutionStatusModel = "completed"
	ExecutionStatusFailed     ExerciseExecutionStatusModel = "failed"
)

type ExerciseExecutionModel struct {
	ID               uint `gorm:"primaryKey;autoIncrement"`
	ExerciseStatusID uint `gorm:"not null;index"`
	ExecutionTime    int
	Status           ExerciseExecutionStatusModel `gorm:"type:varchar(20);not null"`
	CreatedAt        time.Time                    `gorm:"autoCreateTime"`

	ExerciseStatus ExerciseStatusModel `gorm:"foreignKey:ExerciseStatusID"`
}

func (ExerciseExecutionModel) TableName() string { return "exercise_executions" }
