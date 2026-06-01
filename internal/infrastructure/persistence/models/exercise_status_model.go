package models

import "time"

type ExerciseStatusModel struct {
	ID              uint `gorm:"primaryKey;autoIncrement"`
	UserID          uint `gorm:"not null;index"`
	ExerciseID      uint `gorm:"not null;index"`
	ExerciseForkURL string
	Verdict         string
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`

	User     UserModel     `gorm:"foreignKey:UserID"`
	Exercise ExerciseModel `gorm:"foreignKey:ExerciseID"`
}

func (ExerciseStatusModel) TableName() string { return "exercise_status" }
