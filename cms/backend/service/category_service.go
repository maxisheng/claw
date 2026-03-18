package service

import (
	"cms-backend/models"
	"cms-backend/repository"
	"errors"
)

// CategoryService 分类业务逻辑层
type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

// NewCategoryService 创建分类服务
func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// GetAll 获取所有分类
func (s *CategoryService) GetAll() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

// GetByID 根据 ID 获取分类
func (s *CategoryService) GetByID(id uint) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

// Create 创建分类
func (s *CategoryService) Create(req *models.CategoryRequest) (*models.Category, error) {
	// 检查 slug 是否存在
	if req.Slug != "" && s.categoryRepo.ExistsBySlug(req.Slug) {
		return nil, errors.New("URL 标识已存在")
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		ParentID:    req.ParentID,
	}

	if category.Slug == "" {
		category.Slug = req.Name // 简单处理，实际应该生成 slug
	}

	err := s.categoryRepo.Create(category)
	if err != nil {
		return nil, errors.New("创建失败")
	}

	return category, nil
}

// Update 更新分类
func (s *CategoryService) Update(id uint, req *models.CategoryRequest) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("分类不存在")
	}

	// 检查 slug 是否被其他分类使用
	if req.Slug != "" && req.Slug != category.Slug && s.categoryRepo.ExistsBySlug(req.Slug, id) {
		return nil, errors.New("URL 标识已存在")
	}

	category.Name = req.Name
	if req.Slug != "" {
		category.Slug = req.Slug
	}
	category.Description = req.Description
	category.SortOrder = req.SortOrder
	category.ParentID = req.ParentID

	err = s.categoryRepo.Update(category)
	if err != nil {
		return nil, errors.New("更新失败")
	}

	return category, nil
}

// Delete 删除分类
func (s *CategoryService) Delete(id uint) error {
	_, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("分类不存在")
	}

	return s.categoryRepo.Delete(id)
}

// GetStats 获取分类统计数据
func (s *CategoryService) GetStats() (map[string]interface{}, error) {
	count, err := s.categoryRepo.Count()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total": count,
	}, nil
}
