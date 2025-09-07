package repositories

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/utils"
	"time"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	BaseRepository[models.User]
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByRole(role string) ([]models.User, error)
	UpdateLastLogin(userID uint) error
	ChangePassword(userID uint, hashedPassword string) error
}

// userRepository 用户仓储实现
type userRepository struct {
	*baseRepository[models.User]
	db     *gorm.DB
	logger utils.Logger
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB, logger utils.Logger) UserRepository {
	return &userRepository{
		baseRepository: NewBaseRepository[models.User](db, logger).(*baseRepository[models.User]),
		db:             db,
		logger:         logger,
	}
}

// GetByUsername 根据用户名获取用户
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where(&models.User{Username: username}).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("根据用户名获取用户失败", utils.ErrorField(err), utils.String("username", username))
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where(&models.User{Email: email}).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Error("根据邮箱获取用户失败", utils.ErrorField(err), utils.String("email", email))
		return nil, err
	}
	return &user, nil
}

// GetByRole 根据角色获取用户列表
func (r *userRepository) GetByRole(role string) ([]models.User, error) {
	var users []models.User
	err := r.db.Where("role = ?", role).Find(&users).Error
	if err != nil {
		r.logger.Error("根据角色获取用户列表失败", utils.ErrorField(err), utils.String("role", role))
		return nil, err
	}
	return users, nil
}

// UpdateLastLogin 更新最后登录时间
func (r *userRepository) UpdateLastLogin(userID uint) error {
	now := time.Now()
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", now).Error
	if err != nil {
		r.logger.Error("更新用户最后登录时间失败", utils.ErrorField(err), utils.Int("user_id", int(userID)))
		return err
	}
	return nil
}

// ChangePassword 修改密码
func (r *userRepository) ChangePassword(userID uint, hashedPassword string) error {
	err := r.db.Model(&models.User{}).Where("id = ?", userID).Update("password_hash", hashedPassword).Error
	if err != nil {
		r.logger.Error("修改用户密码失败", utils.ErrorField(err), utils.Int("user_id", int(userID)))
		return err
	}
	return nil
}
