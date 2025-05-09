package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sinker/ssop/internal/config"
	"github.com/sinker/ssop/internal/handlers"
	"github.com/sinker/ssop/internal/middleware"
	"github.com/sinker/ssop/internal/models"
	"github.com/sinker/ssop/internal/repository"
	"github.com/sinker/ssop/internal/services"
	"github.com/sinker/ssop/pkg/logger"
	"github.com/sinker/ssop/pkg/redis"
	"github.com/sinker/ssop/pkg/utils"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 初始化配置
	cfg := config.LoadConfig()

	// 初始化日志
	logger.InitLogger(cfg.LogLevel)

	// 设置gin模式
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库
	db, err := models.InitDB(cfg.DBConfig)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}

	// 初始化Redis
	if err := redis.InitRedis(cfg.RedisConfig); err != nil {
		logger.Fatal("Failed to connect to Redis", "error", err)
	}
	// 确保Redis连接在程序退出时关闭
	defer func() {
		if err := redis.Close(); err != nil {
			logger.Error("Failed to close Redis connection", "error", err)
		}
	}()

	// 确保存储目录存在
	ensureStorageDirs(cfg.StorageConfig)

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	datasetRepo := repository.NewDatasetRepository(db)
	analysisRepo := repository.NewAnalysisRepository(db)
	systemRepo := repository.NewSystemRepository(db)

	// 初始化服务
	tokenService := services.NewTokenService()
	authService := services.NewAuthService(userRepo, cfg.JWTConfig, tokenService)
	userService := services.NewUserService(userRepo)
	datasetService := services.NewDatasetService(datasetRepo, cfg.StorageConfig.DatasetDir)
	analysisService := services.NewAnalysisService(analysisRepo, datasetRepo, cfg.StorageConfig.AnalysisDir)
	systemService := services.NewSystemService(systemRepo)

	// 初始化路由
	router := gin.Default()
	
	// 注册中间件
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogger())

	// API版本
	v1 := router.Group("/api/v1")
	
	// 注册路由
	authMiddleware := middleware.AuthMiddleware(authService)
	handlers.RegisterAuthRoutes(v1, authService)
	handlers.RegisterUserRoutes(v1, userService, authMiddleware)
	handlers.RegisterDatasetRoutes(v1, datasetService, authMiddleware)
	handlers.RegisterAnalysisRoutes(v1, analysisService, authMiddleware)
	handlers.RegisterSystemRoutes(v1, systemService, authMiddleware)

	// 创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: router,
	}

	// 优雅关闭
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", "error", err)
		}
	}()

	logger.Info("Server started", "port", cfg.Port)

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}

// 确保存储目录存在
func ensureStorageDirs(cfg config.StorageConfig) {
	dirs := []string{
		cfg.BaseDir,
		cfg.DatasetDir,
		cfg.AnalysisDir,
	}

	for _, dir := range dirs {
		if err := utils.EnsureDir(dir); err != nil {
			logger.Error("Failed to create directory", "path", dir, "error", err)
		}
	}
} 