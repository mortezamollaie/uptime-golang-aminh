package uptime

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
	"uptime/config"
	"uptime/database"
	"uptime/models"
)

var SuspendedWords = []string{"suspended", "Suspended", "account suspended", "سایت مسدود است", "مسدود"}

// readBodyWithTimeout reads the response body with a separate timeout context
func readBodyWithTimeout(body io.ReadCloser, ctx context.Context) ([]byte, error) {
	done := make(chan struct {
		data []byte
		err  error
	}, 1)
	
	go func() {
		data, err := io.ReadAll(body)
		done <- struct {
			data []byte
			err  error
		}{data, err}
	}()
	
	select {
	case result := <-done:
		return result.data, result.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func Check(nodes []models.Node) {
	maxWorkers := config.AppConfig.UptimeChecker.MaxWorkers
	requestTimeout := config.AppConfig.UptimeChecker.RequestTimeout
	var histories []models.History
	if err := database.DB.Find(&histories).Error; err != nil {
		log.Printf("Error fetching histories: %v", err)
		return
	}
	
	historyMap := make(map[uint]*models.History)
	for i := range histories {
		h := &histories[i]
		historyMap[h.NodeID] = h
	}

	jobs := make(chan models.Node, len(nodes))
	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
		for n := range jobs {
			start := time.Now()

			ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, n.URL, nil)
			if err != nil {
				cancel()
				log.Printf("Error creating request for %s: %v", n.URL, err)
				continue
			}
			
			resp, err := http.DefaultClient.Do(req)
			delay := time.Since(start).Seconds()
			
			var status uint
			var up bool
			var suspended bool
			var exception *string

			if err != nil {
				// Handle timeout and other errors
				exc := err.Error()
				exception = &exc
				status = 0
				up = false
				suspended = false
				
				// For timeout errors, set delay to the actual timeout duration
				if strings.Contains(err.Error(), "context deadline exceeded") {
					delay = requestTimeout.Seconds()
				}
			} else {
				defer resp.Body.Close()
				status = uint(resp.StatusCode)
				up = status >= 200 && status < 300

				// Read body with a separate timeout context to avoid cancellation issues
				bodyCtx, bodyCancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer bodyCancel()
				
				bodyBytes, err := readBodyWithTimeout(resp.Body, bodyCtx)
				if err != nil {
					log.Printf("Error reading response body for %s: %v", n.URL, err)
				} else {
					body := string(bodyBytes)
					for _, word := range SuspendedWords {
						if strings.Contains(strings.ToLower(body), strings.ToLower(word)) {
							suspended = true
							break
						}
					}
				}
			}
			
			// Always cancel the request context after we're done
			cancel()

			nodeLog := models.NodeLog{
				NodeID:    n.ID,
				Delay:     &delay,
				Status:    &status,
				Up:        up,
				Suspended: suspended,
				Exception: exception,
			}
			if err := database.DB.Create(&nodeLog).Error; err != nil {
				log.Printf("Error creating node log for %s: %v", n.URL, err)
			}

			if h, ok := historyMap[n.ID]; ok {
				h.Delay = &delay
				h.Status = &status
				h.Up = up
				h.Suspended = suspended
				h.Exception = exception
				if err := database.DB.Save(h).Error; err != nil {
					log.Printf("Error updating history for %s: %v", n.URL, err)
				}
			} else {
				h := models.History{
					NodeID:    n.ID,
					Delay:     &delay,
					Status:    &status,
					Up:        up,
					Suspended: suspended,
					Exception: exception,
				}
				if err := database.DB.Create(&h).Error; err != nil {
					log.Printf("Error creating history for %s: %v", n.URL, err)
				} else {
					historyMap[n.ID] = &h
				}
			}

			fmt.Printf("Checked: %s | Status: %d | Up: %v | Suspended: %v | Delay: %.2fs\n",
				n.URL, status, up, suspended, delay)
		}
	}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	for _, n := range nodes {
		jobs <- n
	}
	close(jobs)

	wg.Wait()
}
