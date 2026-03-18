package routes

import (
	"cms-backend/middleware"
	"cms-backend/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	api := r.Group("/api")
	{
		// 公开路由 - 登录
		api.POST("/login", func(c *gin.Context) { loginHandler(c, db) })

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(db))
		{
			// 管理员相关
			auth.GET("/admin/me", func(c *gin.Context) { getMeHandler(c, db) })
			auth.PUT("/admin/me", func(c *gin.Context) { updateProfileHandler(c, db) })
			auth.PUT("/admin/change-password", func(c *gin.Context) { changePasswordHandler(c, db) })
			auth.POST("/admin/logout", func(c *gin.Context) { logoutHandler(c) })

			// 管理员管理（仅超级管理员）
			auth.GET("/admins", middleware.RequireRole("super_admin", "admin"), func(c *gin.Context) { listAdminsHandler(c, db) })
			auth.POST("/admins", middleware.RequireRole("super_admin"), func(c *gin.Context) { createAdminHandler(c, db) })
			auth.PUT("/admins/:id", middleware.RequireRole("super_admin"), func(c *gin.Context) { updateAdminHandler(c, db) })
			auth.DELETE("/admins/:id", middleware.RequireRole("super_admin"), func(c *gin.Context) { deleteAdminHandler(c, db) })

			// 文章管理
			auth.GET("/articles", func(c *gin.Context) { listArticlesHandler(c, db) })
			auth.GET("/articles/:id", func(c *gin.Context) { getArticleHandler(c, db) })
			auth.POST("/articles", func(c *gin.Context) { createArticleHandler(c, db) })
			auth.PUT("/articles/:id", func(c *gin.Context) { updateArticleHandler(c, db) })
			auth.DELETE("/articles/:id", func(c *gin.Context) { deleteArticleHandler(c, db) })

			// 分类管理
			auth.GET("/categories", func(c *gin.Context) { listCategoriesHandler(c, db) })
			auth.POST("/categories", func(c *gin.Context) { createCategoryHandler(c, db) })
			auth.PUT("/categories/:id", func(c *gin.Context) { updateCategoryHandler(c, db) })
			auth.DELETE("/categories/:id", func(c *gin.Context) { deleteCategoryHandler(c, db) })
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})
}

// loginHandler 登录处理
func loginHandler(c *gin.Context, db *gorm.DB) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var admin models.Admin
	if err := db.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid username or password"})
		return
	}

	// 验证密码
	if !middleware.CheckPassword(admin.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid username or password"})
		return
	}

	// 检查状态
	if admin.Status != models.StatusActive {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "Account is disabled"})
		return
	}

	// 生成 Token
	token, err := middleware.GenerateToken(admin.ID, admin.Username, string(admin.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to generate token"})
		return
	}

	// 更新最后登录时间
	now := time.Now()
	db.Model(&admin).Update("last_login", now)
	admin.LastLogin = &now

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data": models.LoginResponse{
			Token: token,
			Admin: admin,
		},
	})
}

// logoutHandler 登出处理
func logoutHandler(c *gin.Context) {
	// JWT 是无状态的，客户端只需删除 token 即可
	// 这里可以记录日志或加入黑名单（可选）
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Logout successful",
	})
}

// getMeHandler 获取当前管理员信息
func getMeHandler(c *gin.Context, db *gorm.DB) {
	adminID, _ := c.Get("adminID")
	var admin models.Admin
	if err := db.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Admin not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    admin.ToResponse(),
	})
}

// updateProfileHandler 更新个人资料
func updateProfileHandler(c *gin.Context, db *gorm.DB) {
	adminID, _ := c.Get("adminID")
	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var admin models.Admin
	if err := db.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Admin not found"})
		return
	}

	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	db.Model(&admin).Updates(updates)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Profile updated",
		"data":    admin.ToResponse(),
	})
}

// changePasswordHandler 修改密码
func changePasswordHandler(c *gin.Context, db *gorm.DB) {
	adminID, _ := c.Get("adminID")
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var admin models.Admin
	if err := db.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Admin not found"})
		return
	}

	// 验证旧密码
	if !middleware.CheckPassword(admin.Password, req.OldPassword) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Old password is incorrect"})
		return
	}

	// 加密新密码
	hashedPassword, err := middleware.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Failed to hash password"})
		return
	}

	db.Model(&admin).Update("password", hashedPassword)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Password changed successfully",
	})
}

// listAdminsHandler 管理员列表
func listAdminsHandler(c *gin.Context, db *gorm.DB) {
	var admins []models.Admin
	db.Find(&admins)
	
	response := make([]models.AdminResponse, 0, len(admins))
	for _, admin := range admins {
		response = append(response, admin.ToResponse())
	}
	
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    response,
	})
}

// createAdminHandler 创建管理员
func createAdminHandler(c *gin.Context, db *gorm.DB) {
	var req struct {
		Username string           `json:"username" binding:"required"`
		Password string           `json:"password" binding:"required,min=6"`
		Email    string           `json:"email" binding:"omitempty,email"`
		Role     models.AdminRole `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	// 检查用户名是否存在
	var existing models.Admin
	if db.Where("username = ?", req.Username).First(&existing).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Username already exists"})
		return
	}

	hashedPassword, _ := middleware.HashPassword(req.Password)
	admin := models.Admin{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Role:     req.Role,
		Status:   models.StatusActive,
	}

	if admin.Role == "" {
		admin.Role = models.RoleEditor
	}

	db.Create(&admin)
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Admin created",
		"data":    admin.ToResponse(),
	})
}

// updateAdminHandler 更新管理员
func updateAdminHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var req struct {
		Email  string           `json:"email"`
		Role   models.AdminRole `json:"role"`
		Status models.AdminStatus `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var admin models.Admin
	if err := db.First(&admin, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Admin not found"})
		return
	}

	db.Model(&admin).Updates(req)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Admin updated",
		"data":    admin.ToResponse(),
	})
}

// deleteAdminHandler 删除管理员
func deleteAdminHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	adminID, _ := c.Get("adminID")
	
	// 不能删除自己
	if id == fmt.Sprintf("%d", adminID) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Cannot delete yourself"})
		return
	}
	
	db.Delete(&models.Admin{}, id)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Admin deleted",
	})
}

// 文章相关 handler
func listArticlesHandler(c *gin.Context, db *gorm.DB) {
	var articles []models.Article
	db.Preload("Author").Preload("Category").Order("id DESC").Find(&articles)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    articles,
	})
}

func getArticleHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var article models.Article
	if err := db.Preload("Author").Preload("Category").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    article,
	})
}

func createArticleHandler(c *gin.Context, db *gorm.DB) {
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	adminID, _ := c.Get("adminID")
	article := models.Article{
		Title:      req.Title,
		Content:    req.Content,
		Summary:    req.Summary,
		Cover:      req.Cover,
		Slug:       req.Slug,
		Status:     req.Status,
		CategoryID: req.CategoryID,
		AuthorID:   adminID.(uint),
	}

	if article.Status == "" {
		article.Status = "draft"
	}

	db.Create(&article)
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Article created",
		"data":    article,
	})
}

func updateArticleHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var req models.ArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var article models.Article
	if err := db.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Article not found"})
		return
	}

	db.Model(&article).Updates(req)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Article updated",
		"data":    article,
	})
}

func deleteArticleHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	db.Delete(&models.Article{}, id)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Article deleted",
	})
}

// 分类相关 handler
func listCategoriesHandler(c *gin.Context, db *gorm.DB) {
	var categories []models.Category
	db.Order("sort_order ASC, id DESC").Find(&categories)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    categories,
	})
}

func createCategoryHandler(c *gin.Context, db *gorm.DB) {
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	category := models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		ParentID:    req.ParentID,
	}

	db.Create(&category)
	c.JSON(http.StatusCreated, gin.H{
		"code":    0,
		"message": "Category created",
		"data":    category,
	})
}

func updateCategoryHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var req models.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": err.Error()})
		return
	}

	var category models.Category
	if err := db.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Category not found"})
		return
	}

	db.Model(&category).Updates(req)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Category updated",
		"data":    category,
	})
}

func deleteCategoryHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	db.Delete(&models.Category{}, id)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Category deleted",
	})
}
