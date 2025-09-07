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

type BulkURLResponse struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

func GetBulkURL(c *fiber.Ctx) error {
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
	if err := database.DB.WithContext(ctx).Preload("Histories").Where("url IN ?", body.URLs).Order("id desc").Find(&nodes).Error; err != nil {
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

	result := make([]NodeResponse, len(nodes))
	
	// Optimize parallel processing with worker pool pattern for nested data
	if len(nodes) > 0 {
		// Calculate total work (nodes + their histories)
		totalHistories := 0
		for _, node := range nodes {
			totalHistories += len(node.Histories)
		}
		
		// Use worker pool for better resource management
		numWorkers := runtime.NumCPU()
		if len(nodes) < 100 {
			// For small datasets, use fewer workers
			numWorkers = utils.Min(numWorkers, utils.Max(1, len(nodes)/10))
		} else {
			// For large datasets, use more workers but cap appropriately
			numWorkers = utils.Min(numWorkers*2, 16)
		}
		
		if numWorkers == 0 {
			numWorkers = 1
		}

		// Create jobs channel for nodes
		jobs := make(chan int, len(nodes))
		var wg sync.WaitGroup
		var processedCount int64

		// Start workers
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				
				for {
					select {
					case <-ctx.Done():
						// Exit early if context is cancelled
						return
					case nodeIdx, ok := <-jobs:
						if !ok {
							return
						}
						
						// Process single node with all its histories
						n := nodes[nodeIdx]
						histories := make([]HistoryResponse, len(n.Histories))
						
						// Process histories for this node
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
						
						result[nodeIdx] = NodeResponse{
							ID:        n.ID,
							URL:       n.URL,
							Histories: histories,
						}
						
						atomic.AddInt64(&processedCount, 1)
					}
				}
			}()
		}

		// Send jobs to workers
		go func() {
			defer close(jobs)
			for i := 0; i < len(nodes); i++ {
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
			return c.Status(408).JSON(BulkURLResponse{
				Code:    408,
				Msg:     "Request timeout during data processing",
				Success: false,
				Data:    nil,
			})
		}

		log.Printf("Bulk URL: Processed %d nodes with %d total histories using %d workers", 
			atomic.LoadInt64(&processedCount), totalHistories, numWorkers)
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
