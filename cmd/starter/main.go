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
	fmt.Println("Starting node synchronization...")

	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: No .env file found")
	}

	database.Connect()

	apiKey := os.Getenv("UPTIME_API_KEY")
	if apiKey == "" {
		fmt.Println("Error: UPTIME_API_KEY environment variable is not set")
		return
	}

	req, err := http.NewRequest("GET", "https://api.aminh.pro/aio/v2/asc/uptime/list", nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("API returned status code: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	var apiRes ApiResponse
	if err := json.Unmarshal(body, &apiRes); err != nil {
		fmt.Printf("Error parsing JSON response: %v\n", err)
		return
	}

	if !apiRes.Success {
		fmt.Println("API request failed (success=false)")
		return
	}

	successCount := 0
	for _, url := range apiRes.Data {
		if strings.TrimSpace(url) == "" {
			continue
		}

		var node models.Node
		err := database.DB.Where("url = ?", url).First(&node).Error
		if err != nil {
			if strings.Contains(err.Error(), "record not found") {
				if createErr := database.DB.Create(&models.Node{URL: url}).Error; createErr != nil {
					fmt.Printf("Error creating node for URL %s: %v\n", url, createErr)
				} else {
					successCount++
				}
			} else {
				fmt.Printf("Error checking node for URL %s: %v\n", url, err)
			}
		}
	}

	fmt.Printf("Node synchronization completed. Added %d new nodes.\n", successCount)
}
