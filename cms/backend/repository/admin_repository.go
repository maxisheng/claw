package repository

import (
	"cms-backend/models"
	"gorm.io/gorm"
)

// AdminRepository 管理员数据访问层
type AdminRepository struct {
	db *gorm.DB
}

// NewAdminRepository 创建管理员仓库
func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// FindByUsername 根据用户名查找
func (r *AdminRepository) FindByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindByID 根据 ID 查找
func (r *AdminRepository) FindByID(id uint) (*models.Admin, error) {
	var admin models.Admin
	err := r.db.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindAll 查找所有管理员
func (r *AdminRepository) FindAll() ([]models.Admin, error) {
	var admins []models.Admin
	err := r.db.Find(&admins).Error
	return admins, err
}

// Create 创建管理员
func (r *AdminRepository) Create(admin *models.Admin) error {
	return r.db.Create(admin).Error
}

// Update 更新管理员
func (r *AdminRepository) Update(admin *models.Admin) error {
	return r.db.Save(admin).Error
}

// Delete 删除管理员
func (r *AdminRepository) Delete(id uint) error {
	return r.db.Delete(&models.Admin{}, id).Error
}

// UpdatePassword 更新密码
func (r *AdminRepository) UpdatePassword(id uint, hashedPassword string) error {
	return r.db.Model(&models.Admin{}).Where("id = ?", id).Update("password", hashedPassword).Error
}

// UpdateLastLogin 更新最后登录时间
func (r *AdminRepository) UpdateLastLogin(id uint) error {
	return r.db.Model(&models.Admin{}).Where("id = ?", id).Update("last_login", gorm.Expr("NOW()")).Error
}

// ExistsByUsername 检查用户名是否存在
func (r *AdminRepository) ExistsByUsername(username string, excludeID ...uint) bool {
	query := r.db.Where("username = ?", username)
	if len(excludeID) > 0 {
		query = query.Where("id != ?", excludeID[0])
	}
	var count int64
	query.Model(&models.Admin{}).Count(&count)
	return count > 0
}
