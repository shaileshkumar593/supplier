package geoCoder

import (
	"context"
	"net/http"
	"time"

	httpipv4 "swallow-supplier/utils/client" // replace with your actual import path
)

// GoogleGeocoder implements the Geocoder interface using the Google Geocoding API.
type GoogleGeocoder struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

var domains = []string{
	".googleapis.com", "maps.googleapis.com",
}

// Geocoder defines the interface for a geocoding service.
type Geocoder interface {
	GetPlaceID(ctx context.Context, lat, lng float64) (string, error)
}

// NewGoogleGeoCoderWithClient lets you inject a custom http.Client (tests or overrides).
func NewGoogleGeoCoderWithClient(apiKey, baseURL string, httpClient *http.Client) (*GoogleGeocoder, error) {
	if httpClient == nil {
		httpClient = httpipv4.NewIPv4Client(8*time.Second, nil, domains)
	}
	return &GoogleGeocoder{
		APIKey:  apiKey,
		BaseURL: baseURL,
		HTTP:    httpClient,
	}, nil
}

func NewGoogleGeoCoder(apiKey, baseURL string) (*GoogleGeocoder, error) {
	client := httpipv4.NewIPv4Client(8*time.Second, nil, domains)
	return NewGoogleGeoCoderWithClient(apiKey, baseURL, client)
}
