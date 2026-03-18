package handlers

import (
	"cms-backend/models"
	"cms-backend/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminHandler 管理员 HTTP 处理器
type AdminHandler struct {
	adminService *service.AdminService
}

// NewAdminHandler 创建管理员处理器
func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
	})
}

// Login 管理员登录
func (h *AdminHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, token, err := h.adminService.Login(req.Username, req.Password)
	if err != nil {
		Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	Success(c, models.LoginResponse{
		Token: token,
		Admin: *admin,
	})
}

// GetMe 获取当前管理员信息
func (h *AdminHandler) GetMe(c *gin.Context) {
	adminID, exists := c.Get("adminID")
	if !exists {
		Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	admin, err := h.adminService.GetMe(adminID.(uint))
	if err != nil {
		Error(c, http.StatusNotFound, err.Error())
		return
	}

	Success(c, admin.ToResponse())
}

// UpdateProfile 更新个人资料
func (h *AdminHandler) UpdateProfile(c *gin.Context) {
	adminID, exists := c.Get("adminID")
	if !exists {
		Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, err := h.adminService.UpdateProfile(adminID.(uint), req.Email, req.Avatar)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, admin.ToResponse())
}

// ChangePassword 修改密码
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	adminID, exists := c.Get("adminID")
	if !exists {
		Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	err := h.adminService.ChangePassword(adminID.(uint), req.OldPassword, req.NewPassword)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, nil)
}

// Logout 登出
func (h *AdminHandler) Logout(c *gin.Context) {
	// JWT 是无状态的，客户端只需删除 token
	Success(c, nil)
}

// ListAdmins 管理员列表
func (h *AdminHandler) ListAdmins(c *gin.Context) {
	admins, err := h.adminService.GetAll()
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response := make([]models.AdminResponse, 0, len(admins))
	for _, admin := range admins {
		response = append(response, admin.ToResponse())
	}

	Success(c, response)
}

// CreateAdmin 创建管理员
func (h *AdminHandler) CreateAdmin(c *gin.Context) {
	var req struct {
		Username string           `json:"username" binding:"required"`
		Password string           `json:"password" binding:"required,min=6"`
		Email    string           `json:"email" binding:"omitempty,email"`
		Role     models.AdminRole `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, err := h.adminService.Create(req.Username, req.Password, req.Email, req.Role)
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, admin.ToResponse())
}

// UpdateAdmin 更新管理员
func (h *AdminHandler) UpdateAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "无效的管理员 ID")
		return
	}

	var req struct {
		Email  string           `json:"email"`
		Role   models.AdminRole `json:"role"`
		Status models.AdminStatus `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	admin, err := h.adminService.Update(uint(id), req.Email, req.Role, req.Status)
	if err != nil {
		Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	Success(c, admin.ToResponse())
}

// DeleteAdmin 删除管理员
func (h *AdminHandler) DeleteAdmin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		Error(c, http.StatusBadRequest, "无效的管理员 ID")
		return
	}

	adminID, exists := c.Get("adminID")
	if !exists {
		Error(c, http.StatusUnauthorized, "未登录")
		return
	}

	err = h.adminService.Delete(uint(id), adminID.(uint))
	if err != nil {
		Error(c, http.StatusBadRequest, err.Error())
		return
	}

	Success(c, nil)
}
