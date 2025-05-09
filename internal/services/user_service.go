package services

import (
	"errors"

	"github.com/sinker/ssop/internal/models"
	"github.com/sinker/ssop/internal/repository"
)

// 用户相关错误
var (
	ErrUserNotFound = errors.New("用户不存在")
)

// UserService 用户服务接口
type UserService interface {
	GetUserByID(id string) (*models.User, error)
	GetCurrentUserInfo(id string) (*UserInfoDTO, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// UserInfoDTO 用户信息数据传输对象
type UserInfoDTO struct {
	UserID       string   `json:"userId"`
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	FullName     string   `json:"fullName"`
	Organization string   `json:"organization"`
	Role         string   `json:"role"`
	Permissions  []string `json:"permissions"`
	CreateTime   string   `json:"createTime"`
	LastLoginTime string   `json:"lastLoginTime"`
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetCurrentUserInfo 获取当前用户信息
func (s *userService) GetCurrentUserInfo(id string) (*UserInfoDTO, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// 获取权限
	permissions := user.UserPermissions()

	// 转换为DTO
	userInfo := &UserInfoDTO{
		UserID:        user.ID,
		Username:      user.Username,
		Email:         user.Email,
		FullName:      user.FullName,
		Organization:  user.Organization,
		Role:          user.Role,
		Permissions:   permissions,
		CreateTime:    user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		LastLoginTime: user.LastLoginAt.Format("2006-01-02T15:04:05Z"),
	}

	return userInfo, nil
} 