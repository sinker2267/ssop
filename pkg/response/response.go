package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Response 标准API响应格式
type Response struct {
	Code      int         `json:"code"`      // 状态码
	Message   string      `json:"message"`   // 状态描述
	Data      interface{} `json:"data"`      // 响应数据
	Timestamp int64       `json:"timestamp"` // 时间戳
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}, msg string) {
	if msg == "" {
		msg = "成功"
	}
	c.JSON(http.StatusOK, Response{
		Code:      200,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code int, msg string) {
	// 确保错误消息不为空
	if msg == "" {
		msg = "操作失败"
	}

	// 获取对应的HTTP状态码
	httpStatus := getHTTPStatus(code)

	c.JSON(httpStatus, Response{
		Code:      code,
		Message:   msg,
		Data:      nil,
		Timestamp: time.Now().UnixMilli(),
	})
}

// 根据业务码获取对应的HTTP状态码
func getHTTPStatus(code int) int {
	switch code {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 500:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

// BadRequest 返回400错误
func BadRequest(c *gin.Context, msg string) {
	Fail(c, 400, msg)
}

// Unauthorized 返回401错误
func Unauthorized(c *gin.Context, msg string) {
	if msg == "" {
		msg = "未授权的访问"
	}
	Fail(c, 401, msg)
}

// Forbidden 返回403错误
func Forbidden(c *gin.Context, msg string) {
	if msg == "" {
		msg = "禁止访问"
	}
	Fail(c, 403, msg)
}

// NotFound 返回404错误
func NotFound(c *gin.Context, msg string) {
	if msg == "" {
		msg = "资源不存在"
	}
	Fail(c, 404, msg)
}

// ServerError 返回500错误
func ServerError(c *gin.Context, msg string) {
	if msg == "" {
		msg = "服务器内部错误"
	}
	Fail(c, 500, msg)
} 