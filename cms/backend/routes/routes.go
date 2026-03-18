package routes

import (
	"cms-backend/handlers"
	"cms-backend/middleware"
	"cms-backend/repository"
	"cms-backend/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRoutes 设置所有路由
func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// 初始化 Repository
	adminRepo := repository.NewAdminRepository(db)
	articleRepo := repository.NewArticleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	// 初始化 Service
	adminService := service.NewAdminService(adminRepo)
	articleService := service.NewArticleService(articleRepo)
	categoryService := service.NewCategoryService(categoryRepo)

	// 初始化 Handler
	adminHandler := handlers.NewAdminHandler(adminService)
	articleHandler := handlers.NewArticleHandler(articleService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// API 路由组
	api := r.Group("/api")
	{
		// 公开路由 - 登录
		api.POST("/login", adminHandler.Login)

		// 需要认证的路由
		auth := api.Group("")
		auth.Use(middleware.AuthMiddleware(db))
		{
			// 管理员相关
			auth.GET("/admin/me", adminHandler.GetMe)
			auth.PUT("/admin/me", adminHandler.UpdateProfile)
			auth.PUT("/admin/change-password", adminHandler.ChangePassword)
			auth.POST("/admin/logout", adminHandler.Logout)

			// 管理员管理（仅超级管理员）
			auth.GET("/admins", middleware.RequireRole("super_admin", "admin"), adminHandler.ListAdmins)
			auth.POST("/admins", middleware.RequireRole("super_admin"), adminHandler.CreateAdmin)
			auth.PUT("/admins/:id", middleware.RequireRole("super_admin"), adminHandler.UpdateAdmin)
			auth.DELETE("/admins/:id", middleware.RequireRole("super_admin"), adminHandler.DeleteAdmin)

			// 文章管理
			auth.GET("/articles", articleHandler.ListArticles)
			auth.GET("/articles/:id", articleHandler.GetArticle)
			auth.POST("/articles", articleHandler.CreateArticle)
			auth.PUT("/articles/:id", articleHandler.UpdateArticle)
			auth.DELETE("/articles/:id", articleHandler.DeleteArticle)

			// 分类管理
			auth.GET("/categories", categoryHandler.ListCategories)
			auth.POST("/categories", categoryHandler.CreateCategory)
			auth.PUT("/categories/:id", categoryHandler.UpdateCategory)
			auth.DELETE("/categories/:id", categoryHandler.DeleteCategory)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
			"version":   "1.0.0",
		})
	})
}
