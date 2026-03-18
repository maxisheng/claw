package middleware

import (
	"cms-backend/models"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var jwtKey = []byte("cms-secret-key-change-in-production-2024")

// Claims JWT Claims
type Claims struct {
	AdminID  uint   `json:"admin_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateToken 生成 JWT Token
func GenerateToken(adminID uint, username string, role string) (string, error) {
	claims := Claims{
		AdminID:  adminID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

// AuthMiddleware JWT 认证中间件
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"code": 401, "message": "Missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid authorization header"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"code": 401, "message": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 验证管理员是否存在且状态正常
		var admin models.Admin
		if err := db.First(&admin, claims.AdminID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(401, gin.H{"code": 401, "message": "Admin not found"})
				c.Abort()
				return
			}
			c.JSON(500, gin.H{"code": 500, "message": "Database error"})
			c.Abort()
			return
		}

		if admin.Status != "active" {
			c.JSON(403, gin.H{"code": 403, "message": "Admin account is disabled"})
			c.Abort()
			return
		}

		// 设置上下文
		c.Set("adminID", claims.AdminID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole 角色权限中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(403, gin.H{"code": 403, "message": "Role not found"})
			c.Abort()
			return
		}

		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}

		c.JSON(403, gin.H{"code": 403, "message": "Insufficient permissions"})
		c.Abort()
	}
}
