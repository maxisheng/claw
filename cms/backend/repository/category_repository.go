package repository

import (
	"cms-backend/models"
	"gorm.io/gorm"
)

// CategoryRepository 分类数据访问层
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓库
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// FindByID 根据 ID 查找
func (r *CategoryRepository) FindByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

// FindAll 查找所有分类
func (r *CategoryRepository) FindAll() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("sort_order ASC, id DESC").Find(&categories).Error
	return categories, err
}

// Create 创建分类
func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

// Update 更新分类
func (r *CategoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.Category{}, id).Error
}

// ExistsBySlug 检查 slug 是否存在
func (r *CategoryRepository) ExistsBySlug(slug string, excludeID ...uint) bool {
	query := r.db.Where("slug = ?", slug)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	var count int64
	query.Model(&models.Category{}).Count(&count)
	return count > 0
}

// Count 统计分类数量
func (r *CategoryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Category{}).Count(&count).Error
	return count, err
}
