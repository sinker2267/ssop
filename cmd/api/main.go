package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sinker/ssop/internal/application/service"
	"github.com/sinker/ssop/internal/domain/entity"
	"github.com/sinker/ssop/internal/infrastructure/middleware"
	"github.com/sinker/ssop/internal/infrastructure/persistence"
	"github.com/sinker/ssop/internal/interfaces/api"
	"github.com/sinker/ssop/pkg/auth"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 初始化数据库连接
	dsn := "root:password@tcp(127.0.0.1:3306)/ssop?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移表结构
	if err := autoMigrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化JWT配置
	jwtConfig := auth.DefaultJWTConfig

	// 初始化仓储
	userRepo := persistence.NewUserRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, jwtConfig)

	// 初始化处理器
	authHandler := api.NewAuthHandler(authService)

	// 初始化路由
	r := gin.Default()

	// 注册中间件
	r.Use(middleware.CORS())
	r.Use(middleware.ErrorHandler())

	// 注册路由
	apiV1 := r.Group("/api/v1")
	authHandler.Register(apiV1)

	// 添加需要认证的路由
	protected := apiV1.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtConfig))
	{
		// 这里添加需要认证的路由
		protected.GET("/auth/profile", authHandler.GetProfile)
	}

	// 启动服务
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// 自动迁移表结构
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// 在这里添加需要迁移的实体
		&entity.User{},
	)
} 