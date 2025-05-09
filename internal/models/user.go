package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           string    `json:"userId" gorm:"primaryKey;type:varchar(32)"`
	Username     string    `json:"username" gorm:"type:varchar(50);uniqueIndex"`
	Password     string    `json:"-" gorm:"type:varchar(255)"`
	Email        string    `json:"email" gorm:"type:varchar(100);uniqueIndex"`
	FullName     string    `json:"fullName" gorm:"type:varchar(100)"`
	Organization string    `json:"organization" gorm:"type:varchar(100)"`
	Role         string    `json:"role" gorm:"type:varchar(20);default:'student'"`
	CreatedAt    *time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdatedAt    *time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	LastLoginAt  *time.Time `json:"lastLoginTime"`
}

// TableName 表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前的钩子函数
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 生成密码哈希
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UpdatePassword 更新密码
func (u *User) UpdatePassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// UserPermissions 获取用户权限
func (u *User) UserPermissions() []string {
	switch u.Role {
	case "admin":
		return []string{"user:read", "user:write", "data:read", "data:write", "data:delete", "analysis:use", "system:admin"}
	case "researcher":
		return []string{"user:read", "data:read", "data:write", "analysis:use"}
	case "student":
		return []string{"user:read", "data:read", "analysis:use"}
	case "guest":
		return []string{"user:read", "data:read"}
	default:
		return []string{}
	}
} 