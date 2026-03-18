package handlers

import (
	"cms-backend/models"
	"cms-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ArticleHandler 文章 HTTP 处理器
type ArticleHandler struct {
	articleService *service.ArticleService
}

// NewArticleHandler 创建文章处理器
func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: articleService}
}

// ListArticles 文章列表
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	articles, err := h.articleService.GetAll()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	Success(c, articles)
}

// GetArticle 文章详情
func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "无效的文章 ID")
		return
	}

	article, err := h.articleService.GetByID(uint(id))
	if err != nil {
		Error(c, http.StatusNotFound, "文章不存在")
		return
	}

	// 增加浏览量
	_ = h.articleService.IncrementViewCount(uint(id))

	Success(c, article)
}

// CreateArticle 创建文章
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	adminID, exists := c.Get("adminID")
	if !exists {
		Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	article, err := h.articleService.Create(&req, adminID.(uint))
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, article)
}

// UpdateArticle 更新文章
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "无效的文章 ID")
		return
	}

	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	article, err := h.articleService.Update(uint(id), &req)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, article)
}

// DeleteArticle 删除文章
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "无效的文章 ID")
		return
	}

	err = h.articleService.Delete(uint(id))
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, nil)
}
