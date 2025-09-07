package database

import (
	"fmt"
	"uptime/config"
	"uptime/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	dsn := config.AppConfig.Database.DSN

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // کم کردن لاگ‌های اضافی
	})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	fmt.Println("✅ Connected to MySQL!")

	if err := DB.AutoMigrate(&models.Node{}, &models.NodeLog{}, &models.History{}); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	fmt.Println("✅ Tables migrated successfully!")
}
