package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"uptime/database"
	"uptime/internal/logcleanup"
)

func main() {
	fmt.Println("Starting log cleanup process...")

	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	database.Connect()
	logcleanup.CleanupOldLogs()
}
