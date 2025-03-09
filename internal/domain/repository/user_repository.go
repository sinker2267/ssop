package repository

import (
	"context"

	"github.com/sinker/ssop/internal/domain/entity"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// GetByID 根据ID获取用户
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	
	// GetByUsername 根据用户名获取用户
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	
	// GetByEmail 根据邮箱获取用户
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	
	// Create 创建用户
	Create(ctx context.Context, user *entity.User) error
	
	// Update 更新用户
	Update(ctx context.Context, user *entity.User) error
	
	// Delete 删除用户
	Delete(ctx context.Context, id uint) error
} 