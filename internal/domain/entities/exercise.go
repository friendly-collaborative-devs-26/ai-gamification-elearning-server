package entities

import (
	"encoding/json"
	"time"
)

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type Exercise struct {
	ID            uint
	BlockID       uint
	CreatorID     uint
	Title         string
	Description   string
	ExerciseURL   string
	Configuration json.RawMessage
	Difficulty    Difficulty
	Score         int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
