package controllers

import (
	"log"
	"os"
	"time"

	"uptime/database"
	"uptime/models"

	"github.com/gofiber/fiber/v2"
)

type ReportResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func GetNodeReport(c *fiber.Ctx) error {
	key := c.Get("Authorization")
	apiKey := os.Getenv("UPTIME_API_KEY")
	log.Println("Header:", key, "Env:", apiKey)

	if key != apiKey {
		return c.Status(401).JSON(ReportResponse{
			Code:    401,
			Msg:     "Token expired",
			Success: false,
			Data:    nil,
		})
	}

	url := c.Query("url")
	if url == "" {
		return c.Status(422).JSON(ReportResponse{
			Code:    422,
			Msg:     "URL is empty",
			Success: false,
			Data:    nil,
		})
	}

	var node models.Node
	if err := database.DB.Where("url = ?", url).First(&node).Error; err != nil {
		return c.Status(404).JSON(ReportResponse{
			Code:    404,
			Msg:     "URL Not Found",
			Success: false,
			Data:    nil,
		})
	}

	startDateStr := c.Query("start-date")
	endDateStr := c.Query("end-date")
	allItem := c.Query("all-item")
	orderAsc := c.Query("order-asc")
	orderDesc := c.Query("order-desc")
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

	db := database.DB.Model(&models.NodeLog{}).Where("node_id = ?", node.ID)

	if startDateStr != "" && endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(422).JSON(ReportResponse{
				Code:    422,
				Msg:     "Start date format invalid",
				Success: false,
				Data:    nil,
			})
		}
		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(422).JSON(ReportResponse{
				Code:    422,
				Msg:     "End date format invalid",
				Success: false,
				Data:    nil,
			})
		}
		db = db.Where("created_at BETWEEN ? AND ?", startDate, endDate.Add(23*time.Hour+59*time.Minute+59*time.Second))
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

	if allItem != "" {
		db = db.Order("id desc")
	}
	if orderAsc != "" {
		db = db.Order("id asc")
	}
	if orderDesc != "" {
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

	var reports []models.NodeLog
	if err := db.Find(&reports).Error; err != nil {
		log.Println("Database error:", err)
		return c.Status(500).JSON(ReportResponse{
			Code:    500,
			Msg:     "Database error",
			Success: false,
			Data:    nil,
		})
	}

	type NodeLogResponse struct {
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

	reportData := make([]NodeLogResponse, len(reports))
	for i, r := range reports {
		exception := r.Exception
		reportData[i] = NodeLogResponse{
			ID:        r.ID,
			NodeID:    r.NodeID,
			Delay:     r.Delay,
			Status:    r.Status,
			Up:        boolToInt(r.Up),
			Suspended: boolToInt(r.Suspended),
			Exception: exception,
			CreatedAt: r.CreatedAt.Unix(),
			UpdatedAt: r.UpdatedAt.Unix(),
		}
	}

	return c.JSON(ReportResponse{
		Code:    200,
		Msg:     "url report",
		Success: true,
		Data: map[string]interface{}{
			"reports": reportData,
		},
	})
}
