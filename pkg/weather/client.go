package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ClientInterface interface {
	GetForecast(location string, days int) (*ForecastResponse, error)
}

type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

func NewClient(apiKey, baseURL string) ClientInterface {
	return &Client{
		APIKey:  apiKey,
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetForecast(location string, days int) (*ForecastResponse, error) {
	url := fmt.Sprintf("%s/forecast.json?key=%s&q=%s&days=%d&aqi=no&alerts=no",
		c.BaseURL, c.APIKey, location, days)

	resp, err := c.HTTP.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to request weather api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather api returned status %d", resp.StatusCode)
	}

	var forecast ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &forecast, nil
}
