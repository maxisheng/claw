package main

import (
	"cms-backend/middleware"
	"cms-backend/models"
	"cms-backend/repository"
	"cms-backend/routes"
	"cms-backend/service"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// MySQL 连接配置
	dsn := "root:123456@tcp(localhost:3306)/cms?charset=utf8mb4&parseTime=True&loc=Local"

	// 初始化数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ 连接数据库失败:", err)
	}
	log.Println("✅ 已连接到 MySQL 数据库")

	// 自动迁移数据表
	err = db.AutoMigrate(&models.Admin{}, &models.Article{}, &models.Category{})
	if err != nil {
		log.Fatal("❌ 数据库迁移失败:", err)
	}
	log.Println("✅ 数据表已创建")

	// 创建默认管理员账户
	CreateDefaultAdmin(db)

	// 初始化 Repository
	adminRepo := repository.NewAdminRepository(db)

	// 初始化 Service
	adminService := service.NewAdminService(adminRepo)

	// 设置 Gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 配置 CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 设置路由
	routes.SetupRoutes(r, db)

	// 启动服务
	log.Println("🚀 服务器启动在 :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("❌ 启动服务器失败:", err)
	}

	// 避免未使用变量警告
	_ = adminService
}

// CreateDefaultAdmin 创建默认管理员账户
func CreateDefaultAdmin(db *gorm.DB) {
	var admin models.Admin
	if db.Where("username = ?", "admin").First(&admin).Error != nil {
		// 密码加密
		hashedPassword, _ := middleware.HashPassword("admin123")
		db.Create(&models.Admin{
			Username: "admin",
			Password: hashedPassword,
			Email:    "admin@example.com",
			Role:     models.RoleSuperAdmin,
			Status:   models.StatusActive,
		})
		log.Println("✅ 默认管理员已创建：admin / admin123")
	}
}
