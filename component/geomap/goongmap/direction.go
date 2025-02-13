package goongmap

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/spf13/viper"
)

const (
	baseURL           = "https://rsapi.goong.io"
	defaultTimeout    = 10 * time.Second
	defaultMaxRetries = 3
)

// Config represents the configuration for the Goong client
type Config struct {
	APIKey     string
	HTTPClient *http.Client
	Timeout    time.Duration
	MaxRetries int
	BaseURL    string
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() *Config {
	goongAPIKey := viper.GetString("GOONG_API_KEY")
	return &Config{
		HTTPClient: &http.Client{Timeout: defaultTimeout},
		Timeout:    defaultTimeout,
		MaxRetries: defaultMaxRetries,
		BaseURL:    baseURL,
		APIKey:     goongAPIKey,
	}
}

// GoongClient represents a Goong Directions API client
type GoongClient struct {
	config *Config
}

// Location represents a geographical point
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// DirectionsRequest represents the parameters for a directions request
type DirectionsRequest struct {
	Origin      Location // Required
	Destination Location // Required
	Vehicle     string   // Required: car, bike, taxi, truck, hd
}

// Distance represents distance information
type Distance struct {
	Text  string `json:"text"`  // e.g., "1155.66 km"
	Value int    `json:"value"` // Distance in meters
}

// Duration represents duration information
type Duration struct {
	Text  string `json:"text"`  // e.g., "22 giờ 10 phút"
	Value int    `json:"value"` // Duration in seconds
}

// GeocodedWaypoint represents a geocoding result
type GeocodedWaypoint struct {
	GeocoderStatus string `json:"geocoder_status"`
	PlaceID        string `json:"place_id"`
}

// Leg represents a leg of the journey
type Leg struct {
	Distance      Distance `json:"distance"`
	Duration      Duration `json:"duration"`
	EndAddress    string   `json:"end_address"`
	EndLocation   Location `json:"end_location"`
	StartAddress  string   `json:"start_address"`
	StartLocation Location `json:"start_location"`
}

// Route represents a route
type Route struct {
	Bounds struct {
		Northeast Location `json:"northeast"`
		Southwest Location `json:"southwest"`
	} `json:"bounds"`
	Legs             []Leg `json:"legs"`
	OverviewPolyline struct {
		Points string `json:"points"` // Encoded polyline string
	} `json:"overview_polyline"`
	Warnings      []string `json:"warnings"`
	WaypointOrder []int    `json:"waypoint_order"`
}

// DirectionsResponse represents the API response
type DirectionsResponse struct {
	GeocodedWaypoints []GeocodedWaypoint `json:"geocoded_waypoints"`
	Routes            []Route            `json:"routes"`
}

// NewClient creates a new Goong Directions API client with the provided configuration
func NewClient(config *Config) (*GoongClient, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{Timeout: config.Timeout}
	}

	if config.BaseURL == "" {
		config.BaseURL = baseURL
	}

	return &GoongClient{
		config: config,
	}, nil
}

// GetDirections gets directions between two points
func (c *GoongClient) GetDirections(req DirectionsRequest) (*DirectionsResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Add("origin", fmt.Sprintf("%f,%f", req.Origin.Lat, req.Origin.Lng))
	params.Add("destination", fmt.Sprintf("%f,%f", req.Destination.Lat, req.Destination.Lng))
	params.Add("vehicle", req.Vehicle)
	params.Add("api_key", c.config.APIKey)

	// Build full URL
	reqURL := fmt.Sprintf("%s/direction?%s", c.config.BaseURL, params.Encode())

	// Make HTTP request with retries
	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		resp, err = c.config.HTTPClient.Get(reqURL)
		if err == nil {
			break
		}
		if attempt == c.config.MaxRetries {
			return nil, fmt.Errorf("error making request after %d attempts: %v", c.config.MaxRetries, err)
		}
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	// Parse response
	var result DirectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
