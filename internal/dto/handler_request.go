package dto

import (
	"errors"
	"tyarus/weather-app/internal/domain"
	"tyarus/weather-app/pkg/utils"
)

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

func (p *PostLocationHandlerRequest) PostLocationHandlerRequestToDomain() domain.Location {
	return domain.Location{
		Name:      p.Name,
		Region:    p.Region,
		Country:   p.Country,
		Latitude:  p.Latitude,
		Longitude: p.Longitude,
	}
}
