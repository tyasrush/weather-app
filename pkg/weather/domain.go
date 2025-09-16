package weather

type ForecastResponse struct {
	Location Location `json:"location"`
	Current  Current  `json:"current"`
	Forecast Forecast `json:"forecast"`
}

type Location struct {
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Lat       float64 `json:"lat"`
	Lon       float64 `json:"lon"`
	TzID      string  `json:"tz_id"`
	Localtime string  `json:"localtime"`
}

type Current struct {
	LastUpdated string    `json:"last_updated"`
	TempC       float64   `json:"temp_c"`
	Condition   Condition `json:"condition"`
	Humidity    int       `json:"humidity"`
	WindKph     float64   `json:"wind_kph"`
}

type Forecast struct {
	Forecastday []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date string `json:"date"`
	Day  Day    `json:"day"`
}

type Day struct {
	MaxtempC  float64   `json:"maxtemp_c"`
	MintempC  float64   `json:"mintemp_c"`
	AvgtempC  float64   `json:"avgtemp_c"`
	Condition Condition `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}
