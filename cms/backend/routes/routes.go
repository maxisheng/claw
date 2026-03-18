package routes

import (
	"cms-backend/middleware"
	"cms-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// API 路由组
	api := r.Group("/api")
	{
		// 公开路由
		api.POST("/login", func(c *gin.Context) { loginHandler(c, db) })
		api.GET("/articles", func(c *gin.Context) { getArticlesHandler(c, db) })
		api.GET("/articles/:slug", func(c *gin.Context) { getArticleHandler(c, db) })
		api.GET("/categories", func(c *gin.Context) { getCategoriesHandler(c, db) })

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(db))
		{
			// 用户管理
			auth.GET("/user", func(c *gin.Context) { getUserHandler(c, db) })

			// 文章管理
			auth.POST("/articles", func(c *gin.Context) { createArticleHandler(c, db) })
			auth.PUT("/articles/:id", func(c *gin.Context) { updateArticleHandler(c, db) })
			auth.DELETE("/articles/:id", func(c *gin.Context) { deleteArticleHandler(c, db) })

			// 分类管理
			auth.POST("/categories", func(c *gin.Context) { createCategoryHandler(c, db) })
			auth.PUT("/categories/:id", func(c *gin.Context) { updateCategoryHandler(c, db) })
			auth.DELETE("/categories/:id", func(c *gin.Context) { deleteCategoryHandler(c, db) })
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
}

func loginHandler(c *gin.Context, db *gorm.DB) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 简单密码验证（生产环境应该用 bcrypt）
	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// 生成 JWT token（简化版本，实际应该用 JWT）
	token := middleware.GenerateToken(user.ID, user.Username)

	c.JSON(http.StatusOK, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

func getArticlesHandler(c *gin.Context, db *gorm.DB) {
	var articles []models.Article
	db.Preload("Author").Preload("Category").Where("status = ?", "published").Find(&articles)
	c.JSON(http.StatusOK, articles)
}

func getArticleHandler(c *gin.Context, db *gorm.DB) {
	slug := c.Param("slug")
	var article models.Article
	if err := db.Preload("Author").Preload("Category").Where("slug = ?", slug).First(&article).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}
	c.JSON(http.StatusOK, article)
}

func createArticleHandler(c *gin.Context, db *gorm.DB) {
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userID")
	article := models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		Slug:       req.Slug,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		AuthorID:   userID.(uint),
	}

	db.Create(&article)
	c.JSON(http.StatusCreated, article)
}

func updateArticleHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	article.Title = req.Title
	article.Content = req.Content
	article.Summary = req.Summary
	article.Cover = req.Cover
	article.Slug = req.Slug
	article.Status = req.Status
	article.CategoryID = req.CategoryID

	db.Save(&article)
	c.JSON(http.StatusOK, article)
}

func deleteArticleHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	db.Delete(&models.Article{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted"})
}

func getCategoriesHandler(c *gin.Context, db *gorm.DB) {
	var categories []models.Category
	db.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

func createCategoryHandler(c *gin.Context, db *gorm.DB) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	db.Create(&category)
	c.JSON(http.StatusCreated, category)
}

func updateCategoryHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	category.Name = req.Name
	category.Slug = req.Slug
	category.Description = req.Description
	category.ParentID = req.ParentID

	db.Save(&category)
	c.JSON(http.StatusOK, category)
}

func deleteCategoryHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	db.Delete(&models.Category{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}

func getUserHandler(c *gin.Context, db *gorm.DB) {
	userID, _ := c.Get("userID")
	var user models.User
	db.First(&user, userID)
	c.JSON(http.StatusOK, user)
}
