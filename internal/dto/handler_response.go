package dto

import (
	"time"
	"tyarus/weather-app/internal/domain"
	"tyarus/weather-app/pkg/response"
)

type HealthResponse response.Response[any]

type GetLocationHandlerResponseItem struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Region         string    `json:"region"`
	Country        string    `json:"country"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	CreatedAt      time.Time `json:"createdAt"`
	LastModifiedAt time.Time `json:"lastModifiedAt"`
	DeletedAt      time.Time `json:"deletedAt"`
}

func ParseToGetLocationHandlerResponses(items []domain.Location) []GetLocationHandlerResponseItem {
	results := []GetLocationHandlerResponseItem{}
	for _, v := range items {
		results = append(results, GetLocationHandlerResponseItem{
			ID:             v.ID,
			Name:           v.Name,
			Region:         v.Region,
			Country:        v.Country,
			Latitude:       v.Latitude,
			Longitude:      v.Longitude,
			CreatedAt:      v.CreatedAt,
			LastModifiedAt: v.LastModifiedAt.Time,
			DeletedAt:      v.DeletedAt.Time,
		})
	}

	return results
}

func ParseToGetLocationHandlerResponse(item domain.Location) GetLocationHandlerResponseItem {
	result := GetLocationHandlerResponseItem{
		ID:             item.ID,
		Name:           item.Name,
		Region:         item.Region,
		Country:        item.Country,
		Latitude:       item.Latitude,
		Longitude:      item.Longitude,
		CreatedAt:      item.CreatedAt,
		LastModifiedAt: item.LastModifiedAt.Time,
		DeletedAt:      item.DeletedAt.Time,
	}

	return result
}
