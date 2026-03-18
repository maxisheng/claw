package service

import (
	"cms-backend/models"
	"cms-backend/repository"
	"errors"
)

// ArticleService 文章业务逻辑层
type ArticleService struct {
	articleRepo *repository.ArticleRepository
}

// NewArticleService 创建文章服务
func NewArticleService(articleRepo *repository.ArticleRepository) *ArticleService {
	return &ArticleService{articleRepo: articleRepo}
}

// GetAll 获取所有文章
func (s *ArticleService) GetAll() ([]models.Article, error) {
	return s.articleRepo.FindAll()
}

// GetByID 根据 ID 获取文章
func (s *ArticleService) GetByID(id uint) (*models.Article, error) {
	return s.articleRepo.FindByID(id)
}

// Create 创建文章
func (s *ArticleService) Create(req *models.ArticleRequest, authorID uint) (*models.Article, error) {
	// 检查 slug 是否存在
	if s.articleRepo.ExistsBySlug(req.Slug) {
		return nil, errors.New("URL 标识已存在")
	}

	article := &models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		Slug:       req.Slug,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		AuthorID:   authorID,
	}

	if article.Status == "" {
		article.Status = "draft"
	}

	err := s.articleRepo.Create(article)
	if err != nil {
		return nil, errors.New("创建失败")
	}

	// 重新加载关联数据
	return s.articleRepo.FindByID(article.ID)
}

// Update 更新文章
func (s *ArticleService) Update(id uint, req *models.ArticleRequest) (*models.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("文章不存在")
	}

	// 检查 slug 是否被其他文章使用
	if req.Slug != article.Slug && s.articleRepo.ExistsBySlug(req.Slug, id) {
		return nil, errors.New("URL 标识已存在")
	}

	article.Title = req.Title
	article.Content = req.Content
	article.Summary = req.Summary
	article.Cover = req.Cover
	article.Slug = req.Slug
	article.Status = req.Status
	article.CategoryID = req.CategoryID

	err = s.articleRepo.Update(article)
	if err != nil {
		return nil, errors.New("更新失败")
	}

	// 重新加载关联数据
	return s.articleRepo.FindByID(article.ID)
}

// Delete 删除文章
func (s *ArticleService) Delete(id uint) error {
	_, err := s.articleRepo.FindByID(id)
	if err != nil {
		return errors.New("文章不存在")
	}

	return s.articleRepo.Delete(id)
}

// GetStats 获取文章统计数据
func (s *ArticleService) GetStats() (map[string]interface{}, error) {
	count, err := s.articleRepo.Count()
	if err != nil {
		return nil, err
	}

	totalViews, err := s.articleRepo.TotalViews()
	if err != nil {
		return nil, err
	}

	recent, err := s.articleRepo.FindRecent(5)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total":        count,
		"total_views":  totalViews,
		"recent_articles": recent,
	}, nil
}

// IncrementViewCount 增加浏览量
func (s *ArticleService) IncrementViewCount(id uint) error {
	return s.articleRepo.IncrementViewCount(id)
}
