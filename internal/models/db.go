package models

import (
	"fmt"

	"github.com/sinker/ssop/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(config config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// 自动迁移所有模型
	err = db.AutoMigrate(
		&User{},
		&Dataset{},
		&AnalysisTask{},
		&AnalysisResult{},
		&SystemSetting{},
		&AuditLog{},
	)
	
	return db, err
} 