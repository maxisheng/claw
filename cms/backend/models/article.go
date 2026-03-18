package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:200" json:"title" binding:"required"`
	Content   string         `gorm:"type:text" json:"content" binding:"required"`
	Summary   string         `gorm:"size:500" json:"summary"`
	Cover     string         `gorm:"size:500" json:"cover"`
	Slug      string         `gorm:"uniqueIndex;size:200" json:"slug"`
	Status    string         `gorm:"size:20;default:draft" json:"status"` // draft, published, archived
	CategoryID uint          `json:"category_id"`
	Category  Category       `gorm:"foreignKey:CategoryID" json:"category"`
	AuthorID  uint           `json:"author_id"`
	Author    User           `gorm:"foreignKey:AuthorID" json:"author"`
	ViewCount int            `gorm:"default:0" json:"view_count"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ArticleRequest struct {
	Title      string `json:"title" binding:"required"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary"`
	Cover      string `json:"cover"`
	Slug       string `json:"slug"`
	Status     string `json:"status"`
	CategoryID uint   `json:"category_id"`
}
