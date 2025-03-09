package service

import (
	"context"
	"time"

	"github.com/sinker/ssop/internal/application/dto"
	"github.com/sinker/ssop/internal/domain/entity"
	"github.com/sinker/ssop/internal/domain/repository"
	"github.com/sinker/ssop/pkg/auth"
	"github.com/sinker/ssop/pkg/errors"
)

// AuthService 认证服务接口
type AuthService interface {
	// Register 用户注册
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error)

	// Login 用户登录
	Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, error)

	// RefreshToken 刷新令牌
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error)

	// GetUserByID 根据ID获取用户
	GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error)
}

// AuthServiceImpl 认证服务实现
type AuthServiceImpl struct {
	userRepo   repository.UserRepository
	jwtConfig  auth.JWTConfig
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository, jwtConfig auth.JWTConfig) AuthService {
	return &AuthServiceImpl{
		userRepo:   userRepo,
		jwtConfig:  jwtConfig,
	}
}

// Register 用户注册
func (s *AuthServiceImpl) Register(ctx context.Context, req dto.RegisterRequest) (*dto.UserResponse, error) {
	// 检查用户名是否存在
	existUser, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if existUser != nil {
		return nil, errors.ErrUserAlreadyExists
	}

	// 检查邮箱是否存在
	existUser, err = s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if existUser != nil {
		return nil, errors.NewAppError(409, "邮箱已被注册", nil)
	}

	// 创建用户
	user, err := entity.NewUser(req.Username, req.Password, req.Email, req.Phone)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// 保存用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.ErrInternalServer
	}

	// 返回用户信息
	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
}

// Login 用户登录
func (s *AuthServiceImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.TokenResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	// 记录登录时间
	user.RecordLogin()
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, errors.ErrInternalServer
	}

	// 生成令牌
	accessToken, refreshToken, err := auth.GenerateToken(user.ID, user.Username, s.jwtConfig)
	if err != nil {
		return nil, errors.ErrInternalServer
	}

	// 返回令牌
	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtConfig.TokenExpiry / time.Second),
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
		},
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthServiceImpl) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest) (*dto.TokenResponse, error) {
	// 刷新令牌
	accessToken, refreshToken, err := auth.RefreshToken(req.RefreshToken, s.jwtConfig)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// 解析令牌获取用户信息
	claims, err := auth.ParseToken(accessToken, s.jwtConfig)
	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	// 获取用户
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// 返回令牌
	return &dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.jwtConfig.TokenExpiry / time.Second),
		User: dto.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
		},
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *AuthServiceImpl) GetUserByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.ErrInternalServer
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// 返回用户信息
	return &dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
	}, nil
} 