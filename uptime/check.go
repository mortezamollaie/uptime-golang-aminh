package uptime

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"
	"uptime/database"
	"uptime/models"
)

var SuspendedWords = []string{"suspended", "account suspended", "سایت مسدود است", "مسدود"}

const MaxWorkers = 50

func Check(nodes []models.Node) {
	var histories []models.History
	database.DB.Find(&histories)
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

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			req, _ := http.NewRequestWithContext(ctx, http.MethodGet, n.URL, nil)
			resp, err := http.DefaultClient.Do(req)
			delay := time.Since(start).Seconds()
			cancel()

			var status uint
			var up bool
			var suspended bool
			var exception *string

			if err != nil {
				exc := err.Error()
				exception = &exc
				status = 0
				up = false
				suspended = false
			} else {
				defer resp.Body.Close()
				status = uint(resp.StatusCode)
				up = status >= 200 && status < 300

				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				body := string(bodyBytes)
				for _, word := range SuspendedWords {
					if strings.Contains(strings.ToLower(body), strings.ToLower(word)) {
						suspended = true
						break
					}
				}
			}

			nodeLog := models.NodeLog{
				NodeID:    n.ID,
				Delay:     &delay,
				Status:    &status,
				Up:        up,
				Suspended: suspended,
				Exception: exception,
			}
			database.DB.Create(&nodeLog)

			if h, ok := historyMap[n.ID]; ok {
				h.Delay = &delay
				h.Status = &status
				h.Up = up
				h.Suspended = suspended
				h.Exception = exception
				database.DB.Save(h)
			} else {
				h := models.History{
					NodeID:    n.ID,
					Delay:     &delay,
					Status:    &status,
					Up:        up,
					Suspended: suspended,
					Exception: exception,
				}
				database.DB.Create(&h)
				historyMap[n.ID] = &h
			}

			fmt.Printf("Checked: %s | Status: %d | Up: %v | Suspended: %v | Delay: %.2fs\n",
				n.URL, status, up, suspended, delay)
		}
	}

	for i := 0; i < MaxWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	for _, n := range nodes {
		jobs <- n
	}
	close(jobs)

	wg.Wait()
}
