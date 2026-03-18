package models

import (
	"time"

	"gorm.io/gorm"
)

// AdminRole 管理员角色
type AdminRole string

const (
	RoleSuperAdmin AdminRole = "super_admin" // 超级管理员
	RoleAdmin      AdminRole = "admin"       // 普通管理员
	RoleEditor     AdminRole = "editor"      // 编辑
)

// AdminStatus 管理员状态
type AdminStatus string

const (
	StatusActive   AdminStatus = "active"   // 活跃
	StatusInactive AdminStatus = "inactive" // 禁用
)

// Admin 管理员模型
type Admin struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"`
	Email     string         `gorm:"size:100" json:"email"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Role      AdminRole      `gorm:"size:20;default:editor" json:"role"`
	Status    AdminStatus    `gorm:"size:20;default:active" json:"status"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	Admin Admin  `json:"admin"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Email  string `json:"email" binding:"omitempty,email"`
	Avatar string `json:"avatar"`
}

// AdminResponse 管理员响应（不包含密码）
type AdminResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Avatar    string     `json:"avatar"`
	Role      AdminRole  `json:"role"`
	Status    AdminStatus `json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
}

// ToResponse 转换为响应对象
func (a *Admin) ToResponse() AdminResponse {
	return AdminResponse{
		ID:        a.ID,
		Username:  a.Username,
		Email:     a.Email,
		Avatar:    a.Avatar,
		Role:      a.Role,
		Status:    a.Status,
		LastLogin: a.LastLogin,
		CreatedAt: a.CreatedAt,
	}
}
