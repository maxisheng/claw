package service

import (
	"cms-backend/middleware"
	"cms-backend/models"
	"cms-backend/repository"
	"errors"
	"time"
)

// AdminService 管理员业务逻辑层
type AdminService struct {
	adminRepo *repository.AdminRepository
}

// NewAdminService 创建管理员服务
func NewAdminService(adminRepo *repository.AdminRepository) *AdminService {
	return &AdminService{adminRepo: adminRepo}
}

// Login 管理员登录
func (s *AdminService) Login(username, password string) (*models.Admin, string, error) {
	// 查找管理员
	admin, err := s.adminRepo.FindByUsername(username)
	if err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if !middleware.CheckPassword(admin.Password, password) {
		return nil, "", errors.New("用户名或密码错误")
	}

	// 检查状态
	if admin.Status != models.StatusActive {
		return nil, "", errors.New("账号已被禁用")
	}

	// 生成 Token
	token, err := middleware.GenerateToken(admin.ID, admin.Username, string(admin.Role))
	if err != nil {
		return nil, "", errors.New("生成 token 失败")
	}

	// 更新最后登录时间
	_ = s.adminRepo.UpdateLastLogin(admin.ID)

	return admin, token, nil
}

// GetByID 根据 ID 获取管理员
func (s *AdminService) GetByID(id uint) (*models.Admin, error) {
	return s.adminRepo.FindByID(id)
}

// GetMe 获取当前管理员信息
func (s *AdminService) GetMe(id uint) (*models.Admin, error) {
	return s.adminRepo.FindByID(id)
}

// UpdateProfile 更新个人资料
func (s *AdminService) UpdateProfile(id uint, email, avatar string) (*models.Admin, error) {
	admin, err := s.adminRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("管理员不存在")
	}

	if email != "" {
		admin.Email = email
	}
	if avatar != "" {
		admin.Avatar = avatar
	}

	err = s.adminRepo.Update(admin)
	if err != nil {
		return nil, errors.New("更新失败")
	}

	return admin, nil
}

// ChangePassword 修改密码
func (s *AdminService) ChangePassword(id uint, oldPassword, newPassword string) error {
	admin, err := s.adminRepo.FindByID(id)
	if err != nil {
		return errors.New("管理员不存在")
	}

	// 验证旧密码
	if !middleware.CheckPassword(admin.Password, oldPassword) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := middleware.HashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	err = s.adminRepo.UpdatePassword(id, hashedPassword)
	if err != nil {
		return errors.New("更新密码失败")
	}

	return nil
}

// GetAll 获取所有管理员
func (s *AdminService) GetAll() ([]models.Admin, error) {
	return s.adminRepo.FindAll()
}

// Create 创建管理员
func (s *AdminService) Create(username, password, email string, role models.AdminRole) (*models.Admin, error) {
	// 检查用户名是否存在
	if s.adminRepo.ExistsByUsername(username) {
		return nil, errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := middleware.HashPassword(password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	admin := &models.Admin{
		Username: username,
		Password: hashedPassword,
		Email:    email,
		Role:     role,
		Status:   models.StatusActive,
	}

	if admin.Role == "" {
		admin.Role = models.RoleEditor
	}

	err = s.adminRepo.Create(admin)
	if err != nil {
		return nil, errors.New("创建失败")
	}

	return admin, nil
}

// Update 更新管理员
func (s *AdminService) Update(id uint, email string, role models.AdminRole, status models.AdminStatus) (*models.Admin, error) {
	admin, err := s.adminRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("管理员不存在")
	}

	if email != "" {
		admin.Email = email
	}
	if role != "" {
		admin.Role = role
	}
	if status != "" {
		admin.Status = status
	}

	err = s.adminRepo.Update(admin)
	if err != nil {
		return nil, errors.New("更新失败")
	}

	return admin, nil
}

// Delete 删除管理员
func (s *AdminService) Delete(id, currentAdminID uint) error {
	// 不能删除自己
	if id == currentAdminID {
		return errors.New("不能删除自己")
	}

	_, err := s.adminRepo.FindByID(id)
	if err != nil {
		return errors.New("管理员不存在")
	}

	return s.adminRepo.Delete(id)
}

// GetStats 获取统计数据
func (s *AdminService) GetStats() (map[string]interface{}, error) {
	admins, _ := s.adminRepo.FindAll()
	return map[string]interface{}{
		"total": len(admins),
	}, nil
}

// UpdateLastLogin 更新最后登录时间（公开方法）
func (s *AdminService) UpdateLastLogin(id uint) {
	_ = s.adminRepo.UpdateLastLogin(id)
	now := time.Now()
	_ = s.adminRepo.Update(&models.Admin{
		ID:        id,
		LastLogin: &now,
	})
}
