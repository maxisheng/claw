package models

import (
	"time"

	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Slug        string         `gorm:"size:100;uniqueIndex" json:"slug"`
	Description string         `gorm:"size:500" json:"description"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	ParentID    *uint          `json:"parent_id,omitempty"`
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// CategoryRequest 分类请求
type CategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	ParentID    *uint  `json:"parent_id,omitempty"`
}
