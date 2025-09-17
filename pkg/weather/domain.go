package weather

type ForecastResponse struct {
	Location Location `json:"location"`
	Current  Hour     `json:"current"`
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

type Hour struct {
	ForecastTime string    `json:"time"`
	TempC        float64   `json:"temp_c"`
	TempF        float64   `json:"temp_f"`
	Condition    Condition `json:"condition"`
	Humidity     int       `json:"humidity"`
	WindKph      float64   `json:"wind_kph"`
}

type Forecast struct {
	Forecastday []ForecastDay `json:"forecastday"`
}

type ForecastDay struct {
	Date  string `json:"date"`
	Day   Day    `json:"day"`
	Hours []Hour `json:"hour"`
}

type Day struct {
	MaxtempC    float64   `json:"maxtemp_c"`
	MintempC    float64   `json:"mintemp_c"`
	AvgtempC    float64   `json:"avgtemp_c"`
	AvgtempF    float64   `json:"avgtemp_f"`
	AvgHumidity float64   `json:"avghumidity"`
	MaxWindKPH  float64   `json:"avgmaxwind_kph"`
	Condition   Condition `json:"condition"`
}

type Condition struct {
	Text string `json:"text"`
	Icon string `json:"icon"`
	Code int    `json:"code"`
}
