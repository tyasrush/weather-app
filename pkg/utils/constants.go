package utils

import "time"

const (
	DefaultDBTimeout   time.Duration = 5 * time.Second
	DateFormat         string        = "2006-01-02"
	DateFormatWithHour string        = "2006-01-02 15:04"
	WeatherLocationKey string        = "weather:location:%d"
)

const (
	OrderByCreatedAtAsc  = "created_at_ascend"
	OrderByCreatedAtDesc = "created_at_descend"
	OrderByNameAsc       = "name_ascend"
	OrderByNameDesc      = "name_descend"
)

var OrderByMaps = map[string]bool{
	OrderByCreatedAtAsc:  true,
	OrderByCreatedAtDesc: true,
	OrderByNameAsc:       true,
	OrderByNameDesc:      true,
}
