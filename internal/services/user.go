package services

import (
	"air-quality-server/internal/models"
	"air-quality-server/internal/repositories"
	"air-quality-server/internal/utils"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id uint) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id uint) error
	ListUsers(ctx context.Context, limit, offset int) ([]models.User, error)
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	UpdateLastLogin(ctx context.Context, userID uint) error
}

// userService 用户服务实现
type userService struct {
	userRepo repositories.UserRepository
	logger   utils.Logger
}

// NewUserService 创建用户服务
func NewUserService(userRepo repositories.UserRepository, logger utils.Logger) UserService {
	return &userService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// CreateUser 创建用户
func (s *userService) CreateUser(ctx context.Context, user *models.User) error {
	// 检查用户名是否已存在
	if existingUser, _ := s.userRepo.GetByUsername(user.Username); existingUser != nil {
		return errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if existingUser, _ := s.userRepo.GetByEmail(user.Email); existingUser != nil {
		return errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("密码加密失败", utils.ErrorField(err))
		return err
	}
	user.PasswordHash = string(hashedPassword)

	// 设置默认值
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.Status = "active"

	if err := s.userRepo.Create(ctx, user); err != nil {
		s.logger.Error("创建用户失败", utils.ErrorField(err))
		return err
	}

	s.logger.Info("用户创建成功", utils.String("username", user.Username))
	return nil
}

// GetUser 获取用户
func (s *userService) GetUser(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("获取用户失败", utils.ErrorField(err), utils.Int("user_id", int(id)))
		return nil, err
	}
	return user, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		s.logger.Error("根据用户名获取用户失败", utils.ErrorField(err), utils.String("username", username))
		return nil, err
	}
	return user, nil
}

// GetUserByEmail 根据邮箱获取用户
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		s.logger.Error("根据邮箱获取用户失败", utils.ErrorField(err), utils.String("email", email))
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	now := time.Now()
	user.UpdatedAt = now

	// 使用结构体更新，只更新非零值字段
	updateData := &models.User{
		Username:  user.Username,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		UpdatedAt: now,
	}

	if err := s.userRepo.Update(ctx, user.ID, updateData); err != nil {
		s.logger.Error("更新用户失败", utils.ErrorField(err), utils.Int("user_id", int(user.ID)))
		return err
	}

	s.logger.Info("用户更新成功", utils.Int("user_id", int(user.ID)))
	return nil
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.logger.Error("删除用户失败", utils.ErrorField(err), utils.Int("user_id", int(id)))
		return err
	}

	s.logger.Info("用户删除成功", utils.Int("user_id", int(id)))
	return nil
}

// ListUsers 列出用户
func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	request := &repositories.ListRequest{
		Page:     (offset / limit) + 1,
		PageSize: limit,
	}
	response, err := s.userRepo.List(ctx, request)
	if err != nil {
		s.logger.Error("列出用户失败", utils.ErrorField(err))
		return nil, err
	}
	users := response.Data
	return users, nil
}

// AuthenticateUser 用户认证
func (s *userService) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		s.logger.Error("用户认证失败", utils.ErrorField(err), utils.String("username", username))
		return nil, errors.New("用户名或密码错误")
	}

	if user == nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != "active" {
		return nil, errors.New("用户账户已被禁用")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		s.logger.Error("密码验证失败", utils.ErrorField(err), utils.String("username", username))
		return nil, errors.New("用户名或密码错误")
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(uint(user.ID)); err != nil {
		s.logger.Error("更新最后登录时间失败", utils.ErrorField(err), utils.Int("user_id", int(user.ID)))
	}

	s.logger.Info("用户认证成功", utils.String("username", username))
	return user, nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.Error("获取用户失败", utils.ErrorField(err), utils.Int("user_id", int(userID)))
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("新密码加密失败", utils.ErrorField(err))
		return err
	}

	// 更新密码
	if err := s.userRepo.ChangePassword(userID, string(hashedPassword)); err != nil {
		s.logger.Error("修改密码失败", utils.ErrorField(err), utils.Int("user_id", int(userID)))
		return err
	}

	s.logger.Info("密码修改成功", utils.Int("user_id", int(userID)))
	return nil
}

// UpdateLastLogin 更新最后登录时间
func (s *userService) UpdateLastLogin(ctx context.Context, userID uint) error {
	if err := s.userRepo.UpdateLastLogin(userID); err != nil {
		s.logger.Error("更新最后登录时间失败", utils.ErrorField(err), utils.Int("user_id", int(userID)))
		return err
	}
	return nil
}
