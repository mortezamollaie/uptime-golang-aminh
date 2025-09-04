package database

import (
	"fmt"
	"os"
	"uptime/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:@tcp(127.0.0.1:3306)/uptime_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %v", err))
	}

	fmt.Println("✅ Connected to MySQL!")

	if err := DB.AutoMigrate(&models.Node{}, &models.NodeLog{}, &models.History{}); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	fmt.Println("✅ Tables migrated successfully!")
}
