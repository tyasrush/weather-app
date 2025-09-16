package usecase

import (
	"context"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/internal/repository"
	"tyarus/weather-app/pkg/response"
)

type LocationUsecaseInterface interface {
	GetLocationsUsecase(ctx context.Context, param dto.GetLocationHandlerParam) (response.Response[response.PaginationData[dto.GetLocationHandlerResponseItem]], error)
	CreateLocationUsecase(ctx context.Context, req dto.PostLocationHandlerRequest) (dto.GetLocationHandlerResponseItem, error)
}

type locationUsecase struct {
	locationRepo repository.LocationRepositoryInterface
}

func NewLocationUsecase(locationRepo repository.LocationRepositoryInterface) LocationUsecaseInterface {
	return &locationUsecase{locationRepo: locationRepo}
}

func (u *locationUsecase) GetLocationsUsecase(ctx context.Context, param dto.GetLocationHandlerParam) (response.Response[response.PaginationData[dto.GetLocationHandlerResponseItem]], error) {
	if param.PageSize <= 0 {
		param.PageSize = 10
	}
	if param.CurrentPage <= 0 {
		param.CurrentPage = 1
	}

	offset := (param.CurrentPage - 1) * param.PageSize
	resp := response.Response[response.PaginationData[dto.GetLocationHandlerResponseItem]]{}
	locations, err := u.locationRepo.GetLocations(ctx, repository.GetLocationsParam{
		NameLike: param.Query,
		Limit:    param.PageSize,
		Offset:   offset,
		OrderBy:  param.SortBy,
	})
	if err != nil {
		return resp, err
	}

	count, err := u.locationRepo.GetLocationsCount(ctx)
	if err != nil {
		return resp, err
	}

	resp.Data.Items = dto.ParseToGetLocationHandlerResponses(locations)
	resp.Data.Total = count
	resp.Data.CurrentPage = param.CurrentPage
	resp.Data.PageSize = param.PageSize

	return resp, nil

}

func (u *locationUsecase) CreateLocationUsecase(ctx context.Context, req dto.PostLocationHandlerRequest) (dto.GetLocationHandlerResponseItem, error) {
	location, err := u.locationRepo.InsertLocation(ctx, req.PostLocationHandlerRequestToDomain())
	if err != nil {
		return dto.GetLocationHandlerResponseItem{}, err
	}

	return dto.ParseToGetLocationHandlerResponse(location), nil
}
