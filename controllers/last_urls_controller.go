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

func LastURLs(c *fiber.Ctx) error {
	// Set timeout context for the entire request
	ctx, cancel := context.WithTimeout(c.Context(), 45*time.Second)
	defer cancel()

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

	var nodes []models.Node
	if err := database.DB.WithContext(ctx).Preload("Histories").Order("id desc").Find(&nodes).Error; err != nil {
		log.Println("Database error:", err)
		return c.Status(500).JSON(BulkURLResponse{
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

	type NodeResponse struct {
		ID        uint              `json:"id"`
		URL       string            `json:"url"`
		Histories []HistoryResponse `json:"histories"`
	}

	boolToInt := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}

	result := make([]NodeResponse, len(nodes))
	
	// Ultra-optimized parallel processing with chunked workers for all URLs
	if len(nodes) > 0 {
		// Calculate total work (nodes + their histories)
		totalHistories := 0
		for _, node := range nodes {
			totalHistories += len(node.Histories)
		}
		
		// Ultra-optimized worker calculation
		numWorkers := runtime.NumCPU()
		chunkSize := 50 // Optimized chunk size for node processing
		
		if len(nodes) < 50 {
			numWorkers = utils.Min(4, numWorkers)
			chunkSize = 20
		} else if len(nodes) < 200 {
			numWorkers = utils.Min(numWorkers*2, 12)
			chunkSize = 30
		} else {
			// For large datasets (all URLs), use maximum workers
			numWorkers = utils.Min(numWorkers*4, 24)
			chunkSize = 80
		}
		
		if numWorkers == 0 {
			numWorkers = 1
		}

		var wg sync.WaitGroup
		var processedCount int64
		
		// Pre-allocate memory for better performance
		_ = unsafe.Pointer(&result[0]) // Memory hint for better cache locality
		
		// Process in ultra-optimized chunks using semaphore
		semaphore := make(chan struct{}, numWorkers)
		
		for i := 0; i < len(nodes); i += chunkSize {
			end := i + chunkSize
			if end > len(nodes) {
				end = len(nodes)
			}
			
			semaphore <- struct{}{} // Acquire semaphore
			wg.Add(1)
			
			go func(start, end int) {
				defer func() {
					<-semaphore // Release semaphore
					wg.Done()
				}()
				
				// Check context before processing chunk
				select {
				case <-ctx.Done():
					return
				default:
				}
				
				// Process chunk of nodes with optimized loop
				for nodeIdx := start; nodeIdx < end; nodeIdx++ {
					n := &nodes[nodeIdx] // Use pointer to avoid copying
					histories := make([]HistoryResponse, len(n.Histories))
					
					// Optimized history processing
					for j := range n.Histories {
						h := &n.Histories[j] // Use pointer for better performance
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
					
					result[nodeIdx] = NodeResponse{
						ID:        n.ID,
						URL:       n.URL,
						Histories: histories,
					}
				}
				
				// Single atomic update per chunk
				atomic.AddInt64(&processedCount, int64(end-start))
			}(i, end)
		}

		// Wait for all workers
		wg.Wait()

		// Check if context was cancelled during processing
		if ctx.Err() != nil {
			return c.Status(408).JSON(BulkURLResponse{
				Code:    408,
				Msg:     "Request timeout during data processing",
				Success: false,
				Data:    nil,
			})
		}

		log.Printf("Last URLs: Processed %d nodes with %d total histories using %d workers in chunks of %d (ultra-optimized)", 
			atomic.LoadInt64(&processedCount), totalHistories, numWorkers, chunkSize)
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
