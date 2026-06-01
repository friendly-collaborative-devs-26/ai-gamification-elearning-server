package models

import "time"

type BlockTypeModel string

const (
	BlockTypeRoadMap BlockTypeModel = "road_map"
	BlockTypeCourse  BlockTypeModel = "course"
	BlockTypeModule  BlockTypeModel = "module"
)

type BlockModel struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	Type        BlockTypeModel `gorm:"not null"`
	ParentID    *uint          `gorm:"index"`
	CreatorID   uint           `gorm:"not null;index"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`

	Parent  *BlockModel `gorm:"foreignKey:ParentID"`
	Creator UserModel   `gorm:"foreignKey:CreatorID"`
}

func (BlockModel) TableName() string { return "blocks" }
