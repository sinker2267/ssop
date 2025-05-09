package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sinker/ssop/internal/config"
	"github.com/sinker/ssop/internal/models"
	"github.com/sinker/ssop/internal/repository"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/utils"
)

// 定义错误
var (
	ErrInvalidCredentials   = errors.New("无效的用户名或密码")
	ErrUserAlreadyExists    = errors.New("用户已存在")
	ErrInvalidToken         = errors.New("无效的令牌")
	ErrTokenExpired         = errors.New("令牌已过期")
	ErrTokenBlacklisted    = errors.New("令牌已被撤销")
	ErrInvalidRefreshToken  = errors.New("无效的刷新令牌")
	ErrRefreshTokenExpired  = errors.New("刷新令牌已过期")
	ErrFailedToGenerateToken = errors.New("生成令牌失败")
)

// 定义声明
type Claims struct {
	UserID     string   `json:"userId"`
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

// TokenPair 令牌对
type TokenPair struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"`
	UserID       string `json:"userId"`
	Username     string `json:"username"`
	Role         string `json:"role"`
}

// AuthService 认证服务接口
type AuthService interface {
	Register(registerDTO RegisterDTO) (*models.User, error)
	Login(loginDTO LoginDTO) (*TokenPair, error)
	RefreshToken(refreshToken string) (*TokenPair, error)
	VerifyToken(tokenString string) (*Claims, error)
	GuestLogin() (*TokenPair, error)
	Logout(tokenString string) error
}

// authService 认证服务实现
type authService struct {
	userRepo     repository.UserRepository
	jwtConfig    config.JWTConfig
	tokenService TokenService
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo repository.UserRepository, jwtConfig config.JWTConfig, tokenService TokenService) AuthService {
	return &authService{
		userRepo:     userRepo,
		jwtConfig:    jwtConfig,
		tokenService: tokenService,
	}
}

// RegisterDTO 注册数据传输对象
type RegisterDTO struct {
	Username     string `json:"username" binding:"required,min=3,max=50"`
	Password     string `json:"password" binding:"required,min=6,max=50"`
	Email        string `json:"email" binding:"required,email"`
	FullName     string `json:"fullName" binding:"required"`
	Organization string `json:"organization" binding:"required"`
}

// Register 用户注册
func (s *authService) Register(registerDTO RegisterDTO) (*models.User, error) {
	// 检查用户名是否已存在
	existingUser, _ := s.userRepo.FindByUsername(registerDTO.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("用户名已存在")
	}

	// 检查邮箱是否已存在
	existingEmail, _ := s.userRepo.FindByEmail(registerDTO.Email)
	if existingEmail != nil {
		return nil, fmt.Errorf("邮箱已被注册")
	}

	// 创建新用户
	user := &models.User{
		Username:     registerDTO.Username,
		Password:     registerDTO.Password,
		Email:        registerDTO.Email,
		FullName:     registerDTO.FullName,
		Organization: registerDTO.Organization,
		Role:         "student", // 默认角色
	}

	// 保存用户
	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LoginDTO 登录数据传输对象
type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (s *authService) Login(loginDTO LoginDTO) (*TokenPair, error) {
	// 查找用户
	user, err := s.userRepo.FindByUsername(loginDTO.Username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// 验证密码
	if !user.VerifyPassword(loginDTO.Password) {
		return nil, ErrInvalidCredentials
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		return nil, err
	}

	// 生成令牌
	return s.generateTokenPair(user)
}

// RefreshToken 刷新令牌
func (s *authService) RefreshToken(refreshToken string) (*TokenPair, error) {
	// 解析刷新令牌
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, ErrInvalidRefreshToken
		}
		return nil, ErrInvalidRefreshToken
	}

	if !token.Valid {
		return nil, ErrInvalidRefreshToken
	}

	// 验证令牌是否过期
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrRefreshTokenExpired
	}

	// 检查令牌是否在黑名单中
	isBlacklisted, err := s.tokenService.IsBlacklisted(refreshToken)
	if err != nil {
		logger.Error("检查刷新令牌是否在黑名单中失败", "error", err)
		return nil, err
	}
	if isBlacklisted {
		return nil, ErrTokenBlacklisted
	}

	// 获取用户
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// 将旧的刷新令牌加入黑名单
	if err := s.tokenService.AddToBlacklist(refreshToken, claims); err != nil {
		logger.Error("将旧刷新令牌加入黑名单失败", "error", err)
		// 继续处理，不中断流程
	}

	// 生成新令牌
	return s.generateTokenPair(user)
}

// VerifyToken 验证令牌
func (s *authService) VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.Secret), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// 验证令牌是否过期
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}

	// 检查令牌是否在黑名单中
	isBlacklisted, err := s.tokenService.IsBlacklisted(tokenString)
	if err != nil {
		logger.Error("检查令牌是否在黑名单中失败", "error", err)
		return nil, err
	}
	if isBlacklisted {
		return nil, ErrTokenBlacklisted
	}

	return claims, nil
}

// GuestLogin 游客登录
func (s *authService) GuestLogin() (*TokenPair, error) {
	// 创建游客用户
	guestUser := &models.User{
		ID:       "guest_" + utils.GenerateID(""),
		Username: "guest_" + utils.GenerateID(""),
		Role:     "guest",
	}

	// 生成令牌
	return s.generateTokenPair(guestUser)
}

// Logout 用户登出
func (s *authService) Logout(tokenString string) error {
	// 验证令牌
	claims, err := s.VerifyToken(tokenString)
	if err != nil {
		// 如果令牌已过期或无效，不需要加入黑名单
		if err == ErrTokenExpired || err == ErrInvalidToken {
			return nil
		}
		// 如果令牌已经在黑名单中，直接返回成功
		if err == ErrTokenBlacklisted {
			return nil
		}
		return err
	}

	// 将令牌加入黑名单
	return s.tokenService.AddToBlacklist(tokenString, claims)
}

// generateTokenPair 生成令牌对
func (s *authService) generateTokenPair(user *models.User) (*TokenPair, error) {
	// 获取权限
	permissions := user.UserPermissions()

	// 设置过期时间
	accessTokenExp := time.Now().Add(s.jwtConfig.AccessTokenExp)
	refreshTokenExp := time.Now().Add(s.jwtConfig.RefreshTokenExp)

	// 创建访问令牌声明
	accessClaims := &Claims{
		UserID:     user.ID,
		Username:   user.Username,
		Role:       user.Role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ssop",
			Subject:   user.ID,
		},
	}

	// 创建刷新令牌声明
	refreshClaims := &Claims{
		UserID:     user.ID,
		Username:   user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ssop",
			Subject:   user.ID,
		},
	}

	// 生成访问令牌
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, ErrFailedToGenerateToken
	}

	// 生成刷新令牌
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, ErrFailedToGenerateToken
	}

	// 返回令牌对
	return &TokenPair{
		Token:        accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.jwtConfig.AccessTokenExp.Seconds()),
		UserID:       user.ID,
		Username:     user.Username,
		Role:         user.Role,
	}, nil
} 