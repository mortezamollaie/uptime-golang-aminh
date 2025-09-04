package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"

	"uptime/database"
	"uptime/models"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	database.Connect()

	cutoffDate := time.Now().AddDate(0, 0, -31)

	res := database.DB.Where("created_at < ?", cutoffDate).Delete(&models.NodeLog{})

	if res.Error != nil {
		fmt.Println("Error deleting old logs:", res.Error)
		return
	}

	fmt.Printf("%d old log(s) deleted from node_logs table.\n", res.RowsAffected)
}
