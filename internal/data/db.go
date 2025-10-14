package data

import (
	"fmt"
	"log"
	"tds_server/internal/config"
	"tds_server/internal/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) error {
	// 构建连接字符串
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=prefer TimeZone=%s",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.DbName, cfg.DB.Port, cfg.DB.TimeZone,
	)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层 *sql.DB 来设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB from gorm DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移（只创建缺失表/列，不删除）
	err = DB.AutoMigrate(&model.UserToken{})
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	log.Println("Database initialization finished")
	return nil
}
