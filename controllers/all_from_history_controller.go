package controllers

import (
	"log"
	"os"

	"uptime/database"
	"uptime/models"

	"github.com/gofiber/fiber/v2"
)

func AllFormHistory(c *fiber.Ctx) error {
	key := c.Get("Authorization")
	apiKey := os.Getenv("UPTIME_API_KEY")
	if key != apiKey {
		return c.Status(401).JSON(ReportResponse{
			Code:    401,
			Msg:     "Token expired",
			Success: false,
			Data:    nil,
		})
	}

	allItem := c.Query("all-item")
	firstItem := c.Query("first-item")
	lastItem := c.Query("last-item")
	ascDelay := c.Query("asc-delay")
	descDelay := c.Query("desc-delay")
	ascStatus := c.Query("asc-status")
	descStatus := c.Query("desc-status")
	ascUp := c.Query("asc-up")
	descUp := c.Query("desc-up")
	ascSuspended := c.Query("asc-suspended")
	descSuspended := c.Query("desc-suspended")
	ascException := c.Query("asc-exception")
	descException := c.Query("desc-exception")
	up := c.Query("up")
	down := c.Query("down")
	suspended := c.Query("suspended")
	exception := c.Query("exception")

	db := database.DB.Model(&models.History{})

	if allItem != "" {
		db = db.Order("id desc")
	}
	if firstItem != "" {
		db = db.Order("id asc").Limit(1)
	}
	if lastItem != "" {
		db = db.Order("id desc").Limit(1)
	}
	if ascDelay != "" {
		db = db.Order("delay asc")
	}
	if descDelay != "" {
		db = db.Order("delay desc")
	}
	if ascStatus != "" {
		db = db.Order("status asc")
	}
	if descStatus != "" {
		db = db.Order("status desc")
	}
	if ascUp != "" {
		db = db.Order("up asc")
	}
	if descUp != "" {
		db = db.Order("up desc")
	}
	if ascSuspended != "" {
		db = db.Order("suspended asc")
	}
	if descSuspended != "" {
		db = db.Order("suspended desc")
	}
	if ascException != "" {
		db = db.Order("exception asc")
	}
	if descException != "" {
		db = db.Order("exception desc")
	}

	if up == "1" {
		db = db.Where("up = ?", true)
	}
	if down == "1" {
		db = db.Where("up = ?", false)
	}
	if suspended == "1" {
		db = db.Where("suspended = ?", true)
	}
	if exception == "1" {
		db = db.Where("exception IS NOT NULL")
	}

	var histories []models.History
	if err := db.Find(&histories).Error; err != nil {
		log.Println("Database error:", err)
		return c.Status(500).JSON(ReportResponse{
			Code:    500,
			Msg:     "Database error",
			Success: false,
			Data:    nil,
		})
	}

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

	data := make([]HistoryResponse, len(histories))
	for i, h := range histories {
		data[i] = HistoryResponse{
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

	return c.JSON(ReportResponse{
		Code:    200,
		Msg:     "url report",
		Success: true,
		Data: map[string]interface{}{
			"reports": data,
		},
	})
}
