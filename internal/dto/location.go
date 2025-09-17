package dto

import (
	"errors"
	"time"
	"tyarus/weather-app/internal/domain"
	"tyarus/weather-app/pkg/response"
	"tyarus/weather-app/pkg/utils"
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

type GetLocationHandlerParam struct {
	Query       string
	PageSize    int
	CurrentPage int
	SortBy      string
}

func (p *GetLocationHandlerParam) Validate() error {
	if p.SortBy != "" && !utils.OrderByMaps[p.SortBy] {
		return errors.New("invalid sort_by parameter, only allow created_at_ascend, created_at_descend, name_ascend, name_descend")
	}

	return nil
}

type PostLocationHandlerRequest struct {
	Name      string  `json:"name"`
	Region    string  `json:"region"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
}

func (r *PostLocationHandlerRequest) Validate() error {
	if r.Name == "" {
		return errors.New("invalid name parameter, please check your parameter")
	}

	if r.Region == "" {
		return errors.New("invalid region parameter, please check your parameter")
	}

	if r.Country == "" {
		return errors.New("invalid country parameter, please check your parameter")
	}

	return nil
}

func (p *PostLocationHandlerRequest) PostLocationHandlerRequestToDomain() domain.Location {
	return domain.Location{
		Name:      p.Name,
		Region:    p.Region,
		Country:   p.Country,
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
	}
}
