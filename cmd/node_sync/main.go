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
	// لود کردن .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	// اتصال دیتابیس
	database.Connect()

	// گرفتن API KEY
	apiKey := os.Getenv("UPTIME_API_KEY")
	if apiKey == "" {
		fmt.Println("UPTIME_API_KEY is not set")
		return
	}

	// ساخت request
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

	syncNodes(apiRes.Data)
}

func syncNodes(urls []string) {
	if len(urls) == 0 {
		fmt.Println("No URLs from API")
		return
	}

	var existingNodes []models.Node
	database.DB.Find(&existingNodes)

	existingUrls := make([]string, len(existingNodes))
	for i, node := range existingNodes {
		existingUrls[i] = node.URL
	}

	toDelete := difference(existingUrls, urls)
	if len(toDelete) > 0 {
		database.DB.Where("url IN ?", toDelete).Delete(&models.Node{})
		fmt.Println("Deleted nodes:", strings.Join(toDelete, ", "))
	}

	toAdd := difference(urls, existingUrls)
	for _, url := range toAdd {
		if url != "" {
			database.DB.Create(&models.Node{URL: url})
		}
	}

	if len(toAdd) > 0 {
		fmt.Println("Added nodes:", strings.Join(toAdd, ", "))
	} else {
		fmt.Println("No new nodes to add.")
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
