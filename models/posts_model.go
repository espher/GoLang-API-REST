package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	UserID      uint
	User        User      `gorm:"foreignKey:UserID"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (p *Post) TableName() string {
	return "posts"
}
