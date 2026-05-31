package entities

import "time"

type BlockType string

const (
	BlockTypeRoadMap BlockType = "road_map"
	BlockTypeCourse  BlockType = "course"
	BlockTypeModule  BlockType = "module"
)

type Block struct {
	ID          uint
	Name        string
	Description string
	Type        BlockType
	ParentID    *uint
	CreatorID   uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
