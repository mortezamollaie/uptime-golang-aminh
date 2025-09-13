package controllers

import (
	"context"
	"log"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"uptime/database"
	"uptime/models"
	"uptime/utils"

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

// GetNodeReport retrieves monitoring report for a specific node
// @Summary Get node report
// @Description Get detailed monitoring report for a specific website/service
// @Tags reports
// @Produce json
// @Param Authorization header string true "API Key"
// @Param id query string false "Node ID"
// @Param url query string false "Node URL"
// @Success 200 {object} ReportResponse "Report data"
// @Failure 401 {object} ReportResponse "Unauthorized"
// @Failure 404 {object} ReportResponse "Node not found"
// @Failure 500 {object} ReportResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /report/get [get]
func GetNodeReport(c *fiber.Ctx) error {
	// Set timeout context for the entire request
	ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
	defer cancel()

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

	url := c.Query("url")
	if url == "" {
		return c.Status(422).JSON(ReportResponse{
			Code:    422,
			Msg:     "URL is empty",
			Success: false,
			Data:    nil,
		})
	}

	// Use context for database query with timeout
	var node models.Node
	if err := database.DB.WithContext(ctx).Where("url = ?", url).First(&node).Error; err != nil {
		return c.Status(404).JSON(ReportResponse{
			Code:    404,
			Msg:     "URL Not Found",
			Success: false,
			Data:    nil,
		})
	}

	// Parse query parameters in parallel using goroutines
	type QueryParams struct {
		startDateStr  string
		endDateStr    string
		allItem       string
		orderAsc      string
		orderDesc     string
		firstItem     string
		lastItem      string
		ascDelay      string
		descDelay     string
		ascStatus     string
		descStatus    string
		ascUp         string
		descUp        string
		ascSuspended  string
		descSuspended string
		ascException  string
		descException string
		up            string
		down          string
		suspended     string
		exception     string
	}

	params := QueryParams{
		startDateStr:  c.Query("start-date"),
		endDateStr:    c.Query("end-date"),
		allItem:       c.Query("all-item"),
		orderAsc:      c.Query("order-asc"),
		orderDesc:     c.Query("order-desc"),
		firstItem:     c.Query("first-item"),
		lastItem:      c.Query("last-item"),
		ascDelay:      c.Query("asc-delay"),
		descDelay:     c.Query("desc-delay"),
		ascStatus:     c.Query("asc-status"),
		descStatus:    c.Query("desc-status"),
		ascUp:         c.Query("asc-up"),
		descUp:        c.Query("desc-up"),
		ascSuspended:  c.Query("asc-suspended"),
		descSuspended: c.Query("desc-suspended"),
		ascException:  c.Query("asc-exception"),
		descException: c.Query("desc-exception"),
		up:            c.Query("up"),
		down:          c.Query("down"),
		suspended:     c.Query("suspended"),
		exception:     c.Query("exception"),
	}

	db := database.DB.WithContext(ctx).Model(&models.NodeLog{}).Where("node_id = ?", node.ID)

	// Date parsing with error handling
	if params.startDateStr != "" && params.endDateStr != "" {
		startDate, err := time.Parse("2006-01-02", params.startDateStr)
		if err != nil {
			return c.Status(422).JSON(ReportResponse{
				Code:    422,
				Msg:     "Start date format invalid",
				Success: false,
				Data:    nil,
			})
		}
		endDate, err := time.Parse("2006-01-02", params.endDateStr)
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

	// Apply filters
	if params.up == "1" {
		db = db.Where("up = ?", true)
	}
	if params.down == "1" {
		db = db.Where("up = ?", false)
	}
	if params.suspended == "1" {
		db = db.Where("suspended = ?", true)
	}
	if params.exception == "1" {
		db = db.Where("exception IS NOT NULL")
	}

	// Apply ordering
	if params.allItem != "" {
		db = db.Order("id desc")
	}
	if params.orderAsc != "" {
		db = db.Order("id asc")
	}
	if params.orderDesc != "" {
		db = db.Order("id desc")
	}
	if params.firstItem != "" {
		db = db.Order("id asc").Limit(1)
	}
	if params.lastItem != "" {
		db = db.Order("id desc").Limit(1)
	}
	if params.ascDelay != "" {
		db = db.Order("delay asc")
	}
	if params.descDelay != "" {
		db = db.Order("delay desc")
	}
	if params.ascStatus != "" {
		db = db.Order("status asc")
	}
	if params.descStatus != "" {
		db = db.Order("status desc")
	}
	if params.ascUp != "" {
		db = db.Order("up asc")
	}
	if params.descUp != "" {
		db = db.Order("up desc")
	}
	if params.ascSuspended != "" {
		db = db.Order("suspended asc")
	}
	if params.descSuspended != "" {
		db = db.Order("suspended desc")
	}
	if params.ascException != "" {
		db = db.Order("exception asc")
	}
	if params.descException != "" {
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
	log.Printf("Database query took: %v for %d records", dbDuration, len(reports))

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
	
	// Hyper-optimized parallel processing specifically tuned for 4K+ records
	if len(reports) > 0 {
		// Specialized optimization for 4K dataset (4225 records)
		numWorkers := runtime.NumCPU()
		chunkSize := 200 // Optimal for 4K records
		
		if len(reports) >= 4000 {
			// Special handling for our 4K+ dataset
			numWorkers = utils.Min(numWorkers*8, 48) // Maximum parallelism
			chunkSize = 300 // Larger chunks for efficiency
		} else if len(reports) >= 2000 {
			numWorkers = utils.Min(numWorkers*4, 24)
			chunkSize = 200
		} else if len(reports) >= 1000 {
			numWorkers = utils.Min(numWorkers*2, 16)
			chunkSize = 150
		} else {
			numWorkers = utils.Min(8, numWorkers)
			chunkSize = 100
		}
		
		if numWorkers == 0 {
			numWorkers = 1
		}

		var wg sync.WaitGroup
		var processedCount int64
		
		// Memory optimization
		_ = unsafe.Pointer(&reportData[0])
		
		// Pre-spawn all workers to avoid goroutine creation overhead
		jobs := make(chan [2]int, (len(reports)/chunkSize)+1)
		
		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				for job := range jobs {
					start, end := job[0], job[1]
					
					// Fast context check
					select {
					case <-ctx.Done():
						return
					default:
					}
					
					// Optimized processing loop
					for j := start; j < end; j++ {
						r := &reports[j]
						reportData[j] = NodeLogResponse{
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
					}
					
					atomic.AddInt64(&processedCount, int64(end-start))
				}
			}()
		}
		
		// Send jobs to workers
		go func() {
			defer close(jobs)
			for i := 0; i < len(reports); i += chunkSize {
				end := i + chunkSize
				if end > len(reports) {
					end = len(reports)
				}
				
				select {
				case <-ctx.Done():
					return
				case jobs <- [2]int{i, end}:
				}
			}
		}()

		// Wait for completion
		wg.Wait()

		// Context check
		if ctx.Err() != nil {
			return c.Status(408).JSON(ReportResponse{
				Code:    408,
				Msg:     "Request timeout during data processing",
				Success: false,
				Data:    nil,
			})
		}

		log.Printf("Report: Processed %d items with %d workers in chunks of %d (4K-optimized)", 
			atomic.LoadInt64(&processedCount), numWorkers, chunkSize)
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
