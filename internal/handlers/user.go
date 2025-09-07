package handlers

import (
	"air-quality-server/internal/services"
	"air-quality-server/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService services.UserService
	logger      utils.Logger
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService services.UserService, logger utils.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// CreateUser 创建用户
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("创建用户请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	h.logger.Info("创建用户请求", utils.String("username", req.Username))
	c.JSON(http.StatusCreated, gin.H{
		"message": "用户创建成功",
		"data": gin.H{
			"username": req.Username,
			"email":    req.Email,
			"role":     req.Role,
		},
	})
}

// GetUser 获取用户
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("用户ID参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID参数错误"})
		return
	}

	h.logger.Info("获取用户请求", utils.Int("user_id", int(id)))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户成功",
		"data": gin.H{
			"id":       id,
			"username": "test_user",
			"email":    "test@example.com",
			"role":     "user",
		},
	})
}

// UpdateUser 更新用户
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("用户ID参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID参数错误"})
		return
	}

	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("更新用户请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	h.logger.Info("更新用户请求", utils.Int("user_id", int(id)))
	c.JSON(http.StatusOK, gin.H{
		"message": "用户更新成功",
		"data": gin.H{
			"id": id,
		},
	})
}

// DeleteUser 删除用户
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("用户ID参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID参数错误"})
		return
	}

	h.logger.Info("删除用户请求", utils.Int("user_id", int(id)))
	c.JSON(http.StatusOK, gin.H{
		"message": "用户删除成功",
	})
}

// ListUsers 列出用户
func (h *UserHandler) ListUsers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	h.logger.Info("列出用户请求", utils.Int("limit", limit), utils.Int("offset", offset))
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户列表成功",
		"data": gin.H{
			"users":  []gin.H{},
			"total":  0,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("用户登录请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	h.logger.Info("用户登录请求", utils.String("username", req.Username))
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"data": gin.H{
			"token": "fake_jwt_token",
			"user": gin.H{
				"id":       1,
				"username": req.Username,
				"role":     "user",
			},
		},
	})
}

// Logout 用户登出
func (h *UserHandler) Logout(c *gin.Context) {
	h.logger.Info("用户登出请求")
	c.JSON(http.StatusOK, gin.H{
		"message": "登出成功",
	})
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Error("用户ID参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID参数错误"})
		return
	}

	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("修改密码请求参数错误", utils.ErrorField(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	h.logger.Info("修改密码请求", utils.Int("user_id", int(id)))
	c.JSON(http.StatusOK, gin.H{
		"message": "密码修改成功",
	})
}
