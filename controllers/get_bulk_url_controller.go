package controllers

import (
	"log"
	"os"

	"uptime/database"
	"uptime/models"

	"github.com/gofiber/fiber/v2"
)

type BulkURLResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func GetBulkURL(c *fiber.Ctx) error {
	key := c.Get("Authorization")
	apiKey := os.Getenv("UPTIME_API_KEY")
	if key != apiKey {
		return c.Status(401).JSON(BulkURLResponse{
			Code:    401,
			Msg:     "Token expired",
			Success: false,
			Data:    nil,
		})
	}
	
	var body struct {
		URLs []string `json:"urls"`
	}
	if err := c.BodyParser(&body); err != nil || len(body.URLs) == 0 {
		return c.Status(422).JSON(BulkURLResponse{
			Code:    422,
			Msg:     "URL is empty",
			Success: false,
			Data:    nil,
		})
	}

	var nodes []models.Node
	if err := database.DB.Preload("Histories").Where("url IN ?", body.URLs).Order("id desc").Find(&nodes).Error; err != nil {
		log.Println("Database error:", err)
		return c.Status(500).JSON(BulkURLResponse{
			Code:    500,
			Msg:     "Database error",
			Success: false,
			Data:    nil,
		})
	}

	// آماده‌سازی خروجی مشابه لاراول
	type HistoryResponse struct {
		ID        uint     `json:"id"`
		NodeID    uint     `json:"node_id"`
		Delay     *float64 `json:"delay,omitempty"`
		Status    *uint    `json:"status,omitempty"`
		Up        int      `json:"up"`
		Suspended int      `json:"suspended"`
		Exception *string  `json:"exception"`
		CreatedAt int64    `json:"created_at"`
		UpdatedAt int64    `json:"updated_at"`
	}

	type NodeResponse struct {
		ID        uint              `json:"id"`
		URL       string            `json:"url"`
		Histories []HistoryResponse `json:"histories"`
	}

	result := make([]NodeResponse, len(nodes))
	for i, n := range nodes {
		histories := make([]HistoryResponse, len(n.Histories))
		for j, h := range n.Histories {
			histories[j] = HistoryResponse{
				ID:        h.ID,
				NodeID:    h.NodeID,
				Delay:     h.Delay,
				Status:    h.Status,
				Up:        boolToInt(h.Up),
				Suspended: boolToInt(h.Suspended),
				Exception: h.Exception,
				CreatedAt: h.CreatedAt.Unix(),
				UpdatedAt: h.UpdatedAt.Unix(),
			}
		}
		result[i] = NodeResponse{
			ID:        n.ID,
			URL:       n.URL,
			Histories: histories,
		}
	}

	return c.JSON(BulkURLResponse{
		Code:    200,
		Msg:     "urls report",
		Success: true,
		Data: map[string]interface{}{
			"urls": result,
		},
	})
}
