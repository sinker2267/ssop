package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 解包错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// Is 判断错误类型
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// NewAppError 创建应用错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "无效的凭证", errors.New("invalid credentials"))
	ErrUserNotFound       = NewAppError(http.StatusNotFound, "用户不存在", errors.New("user not found"))
	ErrUserAlreadyExists  = NewAppError(http.StatusConflict, "用户已存在", errors.New("user already exists"))
	ErrInvalidToken       = NewAppError(http.StatusUnauthorized, "无效的令牌", errors.New("invalid token"))
	ErrExpiredToken       = NewAppError(http.StatusUnauthorized, "令牌已过期", errors.New("token expired"))
	ErrInternalServer     = NewAppError(http.StatusInternalServerError, "服务器内部错误", errors.New("internal server error"))
	ErrBadRequest         = NewAppError(http.StatusBadRequest, "无效的请求", errors.New("bad request"))
	ErrForbidden          = NewAppError(http.StatusForbidden, "禁止访问", errors.New("forbidden"))
) 