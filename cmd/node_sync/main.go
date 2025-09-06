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

	syncNodes(apiRes.Data)
}

func syncNodes(urls []string) {
	if len(urls) == 0 {
		fmt.Println("No URLs received from API")
		return
	}

	var existingNodes []models.Node
	if err := database.DB.Find(&existingNodes).Error; err != nil {
		fmt.Printf("Error fetching existing nodes: %v\n", err)
		return
	}

	existingUrls := make([]string, len(existingNodes))
	for i, node := range existingNodes {
		existingUrls[i] = node.URL
	}

	toDelete := difference(existingUrls, urls)
	if len(toDelete) > 0 {
		if err := database.DB.Where("url IN ?", toDelete).Delete(&models.Node{}).Error; err != nil {
			fmt.Printf("Error deleting nodes: %v\n", err)
		} else {
			fmt.Printf("Deleted %d node(s): %s\n", len(toDelete), strings.Join(toDelete, ", "))
		}
	}

	toAdd := difference(urls, existingUrls)
	var successCount int
	for _, url := range toAdd {
		if strings.TrimSpace(url) != "" {
			if err := database.DB.Create(&models.Node{URL: url}).Error; err != nil {
				fmt.Printf("Error adding node %s: %v\n", url, err)
			} else {
				successCount++
			}
		}
	}

	if successCount > 0 {
		fmt.Printf("Successfully added %d new node(s)\n", successCount)
	} else if len(toAdd) == 0 {
		fmt.Println("No new nodes to add")
	}
}

func difference(a, b []string) []string {
	m := make(map[string]bool)
	for _, item := range b {
		m[item] = true
	}
	var diff []string
	for _, item := range a {
		if !m[item] {
			diff = append(diff, item)
		}
	}
	return diff
}
