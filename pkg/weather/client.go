package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type WeatherAPIClientInterface interface {
	GetForecast(location string, day int) (*ForecastResponse, error)
}

type WeatherAPIClient struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

func NewClient(apiKey, baseURL string) WeatherAPIClientInterface {
	return &WeatherAPIClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *WeatherAPIClient) GetForecast(location string, day int) (*ForecastResponse, error) {
	if day == 0 {
		day = 14
	}

	url := fmt.Sprintf("%s/forecast.json?key=%s&q=%s&days=%d",
		c.BaseURL, c.APIKey, location, day)

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
