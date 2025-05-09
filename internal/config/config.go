package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// Config 应用配置
type Config struct {
	Environment string
	Port        int
	LogLevel    string
	DBConfig    DBConfig
	RedisConfig RedisConfig
	JWTConfig   JWTConfig
	StorageConfig StorageConfig
}

// DBConfig 数据库配置
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret           string
	AccessTokenExp  time.Duration
	RefreshTokenExp time.Duration
}

// StorageConfig 存储配置
type StorageConfig struct {
	BaseDir       string
	DatasetDir    string
	AnalysisDir   string
	MaxUploadSize int64
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	// 获取应用环境
	env := getEnv("APP_ENV", "development")
	
	// 获取端口
	port, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))
	
	// 获取日志级别
	logLevel := getEnv("LOG_LEVEL", "info")
	
	// 获取数据库配置
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	dbConfig := DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     dbPort,
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "ssop"),
	}
	
	// 获取Redis配置
	redisPort, _ := strconv.Atoi(getEnv("REDIS_PORT", "6379"))
	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	redisConfig := RedisConfig{
		Host:     getEnv("REDIS_HOST", "localhost"),
		Port:     redisPort,
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       redisDB,
	}
	
	// 获取JWT配置
	accessExp, _ := strconv.Atoi(getEnv("JWT_ACCESS_EXP", "3600"))
	refreshExp, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXP", "604800"))
	jwtConfig := JWTConfig{
		Secret:          getEnv("JWT_SECRET", "ssop_secret_key"),
		AccessTokenExp:  time.Duration(accessExp) * time.Second,
		RefreshTokenExp: time.Duration(refreshExp) * time.Second,
	}
	
	// 获取存储配置
	baseDir := getEnv("STORAGE_BASE_DIR", "./storage")
	maxUploadSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "1073741824"), 10, 64)
	
	storageConfig := StorageConfig{
		BaseDir:       baseDir,
		DatasetDir:    filepath.Join(baseDir, "datasets"),
		AnalysisDir:   filepath.Join(baseDir, "analysis"),
		MaxUploadSize: maxUploadSize,
	}
	
	return &Config{
		Environment:   env,
		Port:          port,
		LogLevel:      logLevel,
		DBConfig:      dbConfig,
		RedisConfig:   redisConfig,
		JWTConfig:     jwtConfig,
		StorageConfig: storageConfig,
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 