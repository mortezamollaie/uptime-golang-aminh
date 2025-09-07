package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		DSN string
	}
	Server struct {
		Port string
	}
	UptimeChecker struct {
		CheckInterval   time.Duration
		RequestTimeout  time.Duration
		MaxWorkers      int
	}
	API struct {
		Key string
	}
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}
	
	AppConfig = &Config{}
	
	// Database config
	AppConfig.Database.DSN = getEnv("MYSQL_DSN", "root:@tcp(127.0.0.1:3306)/uptime_db?charset=utf8mb4&parseTime=True&loc=Local")
	
	// Server config
	AppConfig.Server.Port = getEnv("PORT", "3000")
	
	// Uptime checker config
	checkIntervalStr := getEnv("CHECK_INTERVAL", "2m")
	if checkInterval, err := time.ParseDuration(checkIntervalStr); err == nil {
		AppConfig.UptimeChecker.CheckInterval = checkInterval
	} else {
		AppConfig.UptimeChecker.CheckInterval = 2 * time.Minute
	}
	
	timeoutStr := getEnv("REQUEST_TIMEOUT", "45s")
	if timeout, err := time.ParseDuration(timeoutStr); err == nil {
		AppConfig.UptimeChecker.RequestTimeout = timeout
	} else {
		AppConfig.UptimeChecker.RequestTimeout = 45 * time.Second
	}
	
	maxWorkersStr := getEnv("MAX_WORKERS", "50")
	if maxWorkers, err := strconv.Atoi(maxWorkersStr); err == nil {
		AppConfig.UptimeChecker.MaxWorkers = maxWorkers
	} else {
		AppConfig.UptimeChecker.MaxWorkers = 50
	}
	
	// API config
	AppConfig.API.Key = getEnv("UPTIME_API_KEY", "")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
