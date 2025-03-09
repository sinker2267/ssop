package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户实体
type User struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:50;not null;uniqueIndex"`
	Password  string    `gorm:"size:100;not null"`
	Email     string    `gorm:"size:100;uniqueIndex"`
	Phone     string    `gorm:"size:20;uniqueIndex"`
	LastLogin time.Time
	IsActive  bool      `gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser 创建新用户
func NewUser(username, password, email, phone string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:  username,
		Password:  string(hashedPassword),
		Email:     email,
		Phone:     phone,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// CheckPassword 验证密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UpdatePassword 更新密码
func (u *User) UpdatePassword(newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateProfile 更新用户信息
func (u *User) UpdateProfile(email, phone string) {
	u.Email = email
	u.Phone = phone
	u.UpdatedAt = time.Now()
}

// RecordLogin 记录登录时间
func (u *User) RecordLogin() {
	u.LastLogin = time.Now()
} 