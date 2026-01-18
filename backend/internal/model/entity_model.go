package model

import (
	"time"

	"gorm.io/gorm"
)

type Entity struct {
	// gorm.Model
	// Table string `gorm:"-"`
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
