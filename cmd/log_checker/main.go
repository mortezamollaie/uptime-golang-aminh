package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"

	"uptime/database"
	"uptime/models"
)

func main() {
	fmt.Println("Starting log cleanup process...")
	
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	database.Connect()

	cutoffDate := time.Now().AddDate(0, 0, -31)
	fmt.Printf("Deleting logs older than: %s\n", cutoffDate.Format("2006-01-02 15:04:05"))

	res := database.DB.Where("created_at < ?", cutoffDate).Delete(&models.NodeLog{})

	if res.Error != nil {
		fmt.Printf("Error deleting old logs: %v\n", res.Error)
		return
	}

	fmt.Printf("Successfully deleted %d old log record(s) from node_logs table.\n", res.RowsAffected)
}
