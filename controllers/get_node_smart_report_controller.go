package controllers

import (
	"context"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"uptime/database"
	"uptime/models"
	"uptime/utils"

	"github.com/gofiber/fiber/v2"
)

func GetNodeSmartReport(c *fiber.Ctx) error {
	// Set timeout context for the entire request
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

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
	if err := database.DB.WithContext(ctx).Where("url = ?", url).First(&node).Error; err != nil {
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

	db := database.DB.WithContext(ctx).Model(&models.NodeLog{}).Where("node_id = ?", node.ID)

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
	
	// Measure database query time
	dbStartTime := time.Now()
	if err := db.Find(&reports).Error; err != nil {
		log.Println("Database error:", err)
		return c.Status(500).JSON(ReportResponse{
			Code:    500,
			Msg:     "Database error",
			Success: false,
			Data:    nil,
		})
	}
	dbDuration := time.Since(dbStartTime)
	log.Printf("Smart Query database took: %v for %d records", dbDuration, len(reports))

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
	var downCount int64 // Use atomic counter for thread-safe counting
	
	// Optimize parallel processing with worker pool pattern
	if len(reports) > 0 {
		// Use worker pool for better resource management
		numWorkers := runtime.NumCPU()
		if len(reports) < 1000 {
			// For small datasets, use fewer workers to avoid overhead
			numWorkers = utils.Min(numWorkers, utils.Max(1, len(reports)/100))
		} else {
			// For large datasets, use more workers but cap at CPU count * 2
			numWorkers = utils.Min(numWorkers*2, 16)
		}
		
		if numWorkers == 0 {
			numWorkers = 1
		}

		// Create jobs channel
		jobs := make(chan int, len(reports))
		var wg sync.WaitGroup
		var processedCount int64

		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				localDownCount := int64(0)
				
				for {
					select {
					case <-ctx.Done():
						// Exit early if context is cancelled
						atomic.AddInt64(&downCount, localDownCount)
						return
					case idx, ok := <-jobs:
						if !ok {
							atomic.AddInt64(&downCount, localDownCount)
							return
						}
						
						// Process single item
						r := reports[idx]
						if !r.Up {
							localDownCount++
						}
						
						reportData[idx] = NodeLogResponse{
							ID:        r.ID,
							NodeID:    r.NodeID,
							Delay:     r.Delay,
							Status:    r.Status,
							Up:        boolToInt(r.Up),
							Suspended: boolToInt(r.Suspended),
							Exception: r.Exception,
							CreatedAt: r.CreatedAt.Unix(),
							UpdatedAt: r.UpdatedAt.Unix(),
						}
						
						atomic.AddInt64(&processedCount, 1)
					}
				}
			}()
		}

		// Send jobs to workers
		go func() {
			defer close(jobs)
			for i := 0; i < len(reports); i++ {
				select {
				case <-ctx.Done():
					return
				case jobs <- i:
				}
			}
		}()

		// Wait for all workers to finish
		wg.Wait()

		// Check if context was cancelled during processing
		if ctx.Err() != nil {
			return c.Status(408).JSON(ReportResponse{
				Code:    408,
				Msg:     "Request timeout during data processing",
				Success: false,
				Data:    nil,
			})
		}

		log.Printf("Smart Report: Processed %d items with %d workers, found %d down entries", 
			atomic.LoadInt64(&processedCount), numWorkers, atomic.LoadInt64(&downCount))
	}

	return c.JSON(ReportResponse{
		Code:    200,
		Msg:     "url report",
		Success: true,
		Data: map[string]interface{}{
			"reports":    reportData,
			"log_count":  len(reportData),
			"down_count": downCount,
		},
	})
}
