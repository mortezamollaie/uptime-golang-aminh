package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"

	"uptime/database"
	"uptime/models"
)

type ApiResponse struct {
	Success bool     `json:"success"`
	Data    []string `json:"data"`
}

func main() {
	fmt.Println("Start Progress")

	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	database.Connect()

	apiKey := os.Getenv("UPTIME_API_KEY")
	if apiKey == "" {
		fmt.Println("UPTIME_API_KEY is not set")
		return
	}

	req, err := http.NewRequest("GET", "https://api.aminh.pro/aio/v2/asc/uptime/list", nil)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var apiRes ApiResponse
	if err := json.Unmarshal(body, &apiRes); err != nil {
		fmt.Println("JSON decode error:", err)
		return
	}

	if !apiRes.Success {
		fmt.Println("Request failed (success=false)")
		return
	}

	for _, url := range apiRes.Data {
		if url == "" {
			continue
		}

		var node models.Node
		if err := database.DB.Where("url = ?", url).First(&node).Error; err != nil {
			if err.Error() == "record not found" || strings.Contains(err.Error(), "record not found") {
				database.DB.Create(&models.Node{URL: url})
			} else {
				fmt.Println("Error checking node:", err)
			}
		}
	}

	fmt.Println("End Progress")
}
