package repository

import (
	"cms-backend/models"
	"gorm.io/gorm"
)

// ArticleRepository 文章数据访问层
type ArticleRepository struct {
	db *gorm.DB
}

// NewArticleRepository 创建文章仓库
func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

// FindByID 根据 ID 查找
func (r *ArticleRepository) FindByID(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Author").Preload("Category").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// FindAll 查找所有文章
func (r *ArticleRepository) FindAll() ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Preload("Author").Preload("Category").Order("id DESC").Find(&articles).Error
	return articles, err
}

// FindByStatus 根据状态查找
func (r *ArticleRepository) FindByStatus(status string) ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Preload("Author").Preload("Category").Where("status = ?", status).Order("id DESC").Find(&articles).Error
	return articles, err
}

// Create 创建文章
func (r *ArticleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

// Update 更新文章
func (r *ArticleRepository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

// Delete 删除文章
func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Delete(&models.Article{}, id).Error
}

// IncrementViewCount 增加浏览量
func (r *ArticleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).Update("view_count", gorm.Expr("view_count + 1")).Error
}

// ExistsBySlug 检查 slug 是否存在
func (r *ArticleRepository) ExistsBySlug(slug string, excludeID ...uint) bool {
	query := r.db.Where("slug = ?", slug)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	var count int64
	query.Model(&models.Article{}).Count(&count)
	return count > 0
}

// Count 统计文章数量
func (r *ArticleRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Article{}).Count(&count).Error
	return count, err
}

// TotalViews 统计总浏览量
func (r *ArticleRepository) TotalViews() (int64, error) {
	var total int64
	err := r.db.Model(&models.Article{}).Select("COALESCE(SUM(view_count), 0)").Scan(&total).Error
	return total, err
}

// FindRecent 查找最新文章
func (r *ArticleRepository) FindRecent(limit int) ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Preload("Author").Order("id DESC").Limit(limit).Find(&articles).Error
	return articles, err
}
