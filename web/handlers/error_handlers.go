package handlers

import (
	"air-quality-server/internal/utils"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ErrorHandler 错误处理器
type ErrorHandler struct {
	logger utils.Logger
}

// NewErrorHandler 创建错误处理器
func NewErrorHandler(logger utils.Logger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Error     string            `json:"error"`
	Message   string            `json:"message"`
	Code      int               `json:"code"`
	RequestID string            `json:"request_id,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
}

// WebError 页面错误结构
type WebError struct {
	Code       string
	Message    string
	Details    string
	RequestID  string
	StatusCode int
}

// HandleError 处理错误并返回适当的响应
func (h *ErrorHandler) HandleError(c *gin.Context, err error, statusCode int) {
	requestID := h.getRequestID(c)

	// 记录错误日志
	h.logger.Error("请求处理错误",
		utils.ErrorField(err),
		utils.String("request_id", requestID),
		utils.String("path", c.Request.URL.Path),
		utils.String("method", c.Request.Method),
		utils.Int("status_code", statusCode))

	// 根据请求类型返回不同的响应
	if h.isAPIRequest(c) {
		h.handleAPIError(c, err, statusCode, requestID)
	} else {
		h.handleWebError(c, err, statusCode, requestID)
	}
}

// HandlePanic 处理panic恢复
func (h *ErrorHandler) HandlePanic(c *gin.Context) {
	if r := recover(); r != nil {
		requestID := h.getRequestID(c)

		// 记录panic日志
		h.logger.Error("请求处理panic",
			utils.String("panic", fmt.Sprintf("%v", r)),
			utils.String("stack", string(debug.Stack())),
			utils.String("request_id", requestID),
			utils.String("path", c.Request.URL.Path),
			utils.String("method", c.Request.Method))

		// 根据请求类型返回不同的响应
		if h.isAPIRequest(c) {
			h.handleAPIError(c, fmt.Errorf("internal server error"), http.StatusInternalServerError, requestID)
		} else {
			h.handleWebError(c, fmt.Errorf("internal server error"), http.StatusInternalServerError, requestID)
		}
	}
}

// handleAPIError 处理API错误
func (h *ErrorHandler) handleAPIError(c *gin.Context, err error, statusCode int, requestID string) {
	response := ErrorResponse{
		Error:     http.StatusText(statusCode),
		Message:   err.Error(),
		Code:      statusCode,
		RequestID: requestID,
		Timestamp: time.Now(),
	}

	c.JSON(statusCode, response)
}

// handleWebError 处理Web页面错误
func (h *ErrorHandler) handleWebError(c *gin.Context, err error, statusCode int, requestID string) {
	errorCode := fmt.Sprintf("%d", statusCode)
	errorMessage := err.Error()

	// 根据状态码设置默认错误信息
	switch statusCode {
	case http.StatusNotFound:
		errorCode = "404"
		errorMessage = "页面未找到"
	case http.StatusInternalServerError:
		errorCode = "500"
		errorMessage = "服务器内部错误"
	case http.StatusForbidden:
		errorCode = "403"
		errorMessage = "访问被拒绝"
	case http.StatusBadRequest:
		errorCode = "400"
		errorMessage = "请求参数错误"
	}

	webError := WebError{
		Code:       errorCode,
		Message:    errorMessage,
		Details:    err.Error(),
		RequestID:  requestID,
		StatusCode: statusCode,
	}

	c.HTML(statusCode, "base.html", gin.H{
		"Title":        "错误页面",
		"CurrentPage":  "error",
		"ErrorCode":    webError.Code,
		"ErrorMessage": webError.Message,
		"ErrorDetails": webError.Details,
		"RequestID":    webError.RequestID,
	})
}

// isAPIRequest 判断是否为API请求
func (h *ErrorHandler) isAPIRequest(c *gin.Context) bool {
	path := c.Request.URL.Path
	return len(path) > 4 && path[:4] == "/api" || len(path) > 8 && path[:8] == "/web/api"
}

// getRequestID 获取请求ID
func (h *ErrorHandler) getRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return uuid.New().String()
}

// ValidationError 验证错误结构
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrors 验证错误集合
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// HandleValidationError 处理验证错误
func (h *ErrorHandler) HandleValidationError(c *gin.Context, errors []ValidationError) {
	requestID := h.getRequestID(c)

	h.logger.Warn("请求验证失败",
		utils.String("request_id", requestID),
		utils.String("path", c.Request.URL.Path),
		utils.String("method", c.Request.Method),
		utils.Any("validation_errors", errors))

	if h.isAPIRequest(c) {
		response := ErrorResponse{
			Error:     "Validation Failed",
			Message:   "请求参数验证失败",
			Code:      http.StatusBadRequest,
			RequestID: requestID,
			Timestamp: time.Now(),
		}

		// 添加验证错误详情
		details := make(map[string]string)
		for _, err := range errors {
			details[err.Field] = err.Message
		}
		response.Details = details

		c.JSON(http.StatusBadRequest, response)
	} else {
		// Web页面验证错误处理
		c.HTML(http.StatusBadRequest, "base.html", gin.H{
			"Title":        "参数错误",
			"CurrentPage":  "error",
			"ErrorCode":    "400",
			"ErrorMessage": "请求参数验证失败",
			"ErrorDetails": "请检查您输入的信息是否正确",
			"RequestID":    requestID,
		})
	}
}

// NotFoundHandler 404处理器
func (h *ErrorHandler) NotFoundHandler(c *gin.Context) {
	h.HandleError(c, fmt.Errorf("page not found: %s", c.Request.URL.Path), http.StatusNotFound)
}

// MethodNotAllowedHandler 405处理器
func (h *ErrorHandler) MethodNotAllowedHandler(c *gin.Context) {
	h.HandleError(c, fmt.Errorf("method not allowed: %s", c.Request.Method), http.StatusMethodNotAllowed)
}

// RecoveryMiddleware 恢复中间件
func (h *ErrorHandler) RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer h.HandlePanic(c)
		c.Next()
	}
}

// RequestIDMiddleware 请求ID中间件
func (h *ErrorHandler) RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// ErrorLoggerMiddleware 错误日志中间件
func (h *ErrorHandler) ErrorLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 记录错误状态码
		if c.Writer.Status() >= 400 {
			requestID := h.getRequestID(c)
			h.logger.Error("HTTP请求错误",
				utils.Int("status_code", c.Writer.Status()),
				utils.String("request_id", requestID),
				utils.String("path", c.Request.URL.Path),
				utils.String("method", c.Request.Method),
				utils.String("client_ip", c.ClientIP()),
				utils.String("user_agent", c.Request.UserAgent()))
		}
	}
}
