package cronjob

import (
	"bytes"
	"encoding/json"
	"net/http"
	"swallow-supplier/config"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// notifyTripAPI sends a request to /localhost/trip and retries until success.
func notifyTripAPI(logger log.Logger) {
	cf := config.Instance()
	tripURL := cf.TripSyncUrl
	payload := map[string]string{"message": "MonitorProductUpdates completed successfully"}
	jsonPayload, _ := json.Marshal(payload)

	client := &http.Client{Timeout: 10 * time.Second}
	retryInterval := 2 * time.Second // Initial retry interval
	maxRetryInterval := 30 * time.Second
	maxRetries := 5 // Maximum number of retries

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", tripURL, bytes.NewBuffer(jsonPayload))
		if err != nil {
			level.Error(logger).Log("error", "Failed to create request for trip API", "err", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			level.Error(logger).Log("error", "Failed to notify trip API", "attempt", attempt, "err", err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				level.Info(logger).Log("msg", "Successfully notified trip API", "attempt", attempt)
				return // Exit function on success
			}
			level.Error(logger).Log("msg", "Trip API responded with error", "status", resp.StatusCode, "attempt", attempt)
		}

		// If not successful, retry after a delay
		time.Sleep(retryInterval)
		if retryInterval < maxRetryInterval {
			retryInterval *= 2 // Exponential backoff
		}
	}

	level.Error(logger).Log("msg", "Failed to notify trip API after multiple attempts")
}
