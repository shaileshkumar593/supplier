package geoCoder

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type GoogleGeocodeResponse struct {
	Results []struct {
		PlaceID string `json:"place_id"`
	} `json:"results"`
	Status       string `json:"status"`
	ErrorMessage string `json:"error_message"`
}

// GetPlaceID logs DNS, connect, TLS, and remote addr to prove IPv4 usage.
func (g *GoogleGeocoder) GetPlaceID(ctx context.Context, lat, lng float64, logger log.Logger) (string, error) {
	if g.HTTP == nil {
		return "", fmt.Errorf("nil HTTP client; inject IPv4 client")
	}

	// deadline
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 8*time.Second)
		defer cancel()
	}

	// build URL
	u, err := url.Parse(g.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	q := u.Query()
	q.Set("latlng", fmt.Sprintf("%.8f,%.8f", lat, lng))
	q.Set("key", g.APIKey)
	u.RawQuery = q.Encode()

	// request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	// trace just enough to assert IPv4 without spam
	var usedIPv4 atomic.Bool

	trace := &httptrace.ClientTrace{
		GotConn: func(info httptrace.GotConnInfo) {
			if ta, ok := info.Conn.RemoteAddr().(*net.TCPAddr); ok {
				_ = ta.IP.String()
				if ta.IP.To4() != nil {
					usedIPv4.Store(true)
				}
			}
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(ctx, trace))

	resp, err := g.HTTP.Do(req)
	if err != nil {
		level.Error(logger).Log("op", "GetPlaceID", "err", err)
		return "", fmt.Errorf("http do: %w", err)
	}
	defer resp.Body.Close()

	if !usedIPv4.Load() {
		return "", fmt.Errorf("transport did not use IPv4; check client/transport and proxy bypass list")
	}

	// non-200
	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		level.Error(logger).Log("op", "GetPlaceID", "status", resp.StatusCode, "body", strings.TrimSpace(string(b)))
		return "", fmt.Errorf("geocode http %d", resp.StatusCode)
	}

	// decode
	var geocodeResp GoogleGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geocodeResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	// result
	switch geocodeResp.Status {
	case "OK":
		if len(geocodeResp.Results) == 0 {
			return "", fmt.Errorf("no place ID for %.8f,%.8f", lat, lng)
		}
		return geocodeResp.Results[0].PlaceID, nil
	case "OVER_QUERY_LIMIT", "REQUEST_DENIED", "INVALID_REQUEST":
		level.Error(logger).Log("status", geocodeResp.Status, "msg", geocodeResp.ErrorMessage)
		return "", fmt.Errorf("google %s: %s", geocodeResp.Status, geocodeResp.ErrorMessage)
	default:
		return "", fmt.Errorf("google status %s", geocodeResp.Status)
	}
}

func errString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

// GetPublicIP Signature unchanged.
func GetPublicIP(ctx context.Context, logger log.Logger) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://ifconfig.me/", nil)
	if err != nil {
		level.Error(logger).Log("error ", err)
		return "", fmt.Errorf("failed to create IP request: %w", err)
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		level.Error(logger).Log("error ", err)
		return "", fmt.Errorf("failed to fetch public IP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		level.Error(logger).Log("error ", err)
		return "", fmt.Errorf("failed to read IP response: %w", err)
	}

	return string(body), nil
}

/*
// for testing purpose

func main() {
	ctx := context.Background()
	// Replace with your actual API key and base URL.
	apiKey := "YOUR_API_KEY"
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"

	// Create a Geocoder interface instance using GoogleGeocoder.
	var geocoder Geocoder = &GoogleGeocoder{
		APIKey: apiKey,
		BaseURL: baseURL,
	}

	latitude := 37.7749
	longitude := -122.4194 // San Francisco, CA

	placeID, err := geocoder.GetPlaceID(ctx, latitude, longitude)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Google Place ID:", placeID)
}
*/
