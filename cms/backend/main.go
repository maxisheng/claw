package main

import (
	"cms-backend/models"
	"cms-backend/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 从环境变量读取 MySQL 配置，或使用默认值
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		// 默认 MySQL 连接配置
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			getEnv("DB_USER", "root"),
			getEnv("DB_PASSWORD", "root"),
			getEnv("DB_HOST", "localhost"),
			getEnv("DB_PORT", "3306"),
			getEnv("DB_NAME", "cms"),
		)
	}

	// 初始化数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Connected to MySQL database")

	// 自动迁移
	err = db.AutoMigrate(&models.User{}, &models.Article{}, &models.Category{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建默认管理员账户
	var admin models.User
	if db.Where("username = ?", "admin").First(&admin).Error != nil {
		db.Create(&models.User{
			Username: "admin",
			Password: "admin123", // 实际生产环境应该加密
			Role:     "admin",
		})
	}

	// 设置 Gin
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
	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
