package usecase

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tyarus/weather-app/internal/domain"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/internal/repository"
	"tyarus/weather-app/mocks"
)

func TestGetLocationsUsecase(t *testing.T) {
	t.Run("WHEN error occurred on get location from database, THEN should return error accordingly", func(t *testing.T) {
		mockRepo := mocks.NewLocationRepositoryInterface(t)
		usecase := NewLocationUsecase(mockRepo)
		ctx := context.Background()
		param := dto.GetLocationHandlerParam{
			Query:       "test",
			PageSize:    10,
			CurrentPage: 1,
			SortBy:      "name_ascend",
		}

		expectedError := errors.New("database error")
		mockRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			NameLike: "test",
			Limit:    10,
			Offset:   0,
			OrderBy:  "name_ascend",
		}).Return(nil, expectedError)

		_, err := usecase.GetLocationsUsecase(ctx, param)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("WHEN error occurred on get count from database, THEN should return error accordingly", func(t *testing.T) {
		mockRepo := mocks.NewLocationRepositoryInterface(t)
		usecase := NewLocationUsecase(mockRepo)
		ctx := context.Background()
		param := dto.GetLocationHandlerParam{
			Query:       "test",
			PageSize:    10,
			CurrentPage: 1,
			SortBy:      "name_ascend",
		}

		locations := []domain.Location{
			{
				ID:             1,
				Name:           "Test Location",
				Country:        "Test Country",
				Latitude:       1.0,
				Longitude:      1.0,
				CreatedAt:      time.Now(),
				LastModifiedAt: sql.NullTime{Time: time.Now()},
			},
		}

		expectedError := errors.New("count error")
		mockRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			NameLike: "test",
			Limit:    10,
			Offset:   0,
			OrderBy:  "name_ascend",
		}).Return(locations, nil)
		mockRepo.On("GetLocationsCount", ctx).Return(0, expectedError)

		_, err := usecase.GetLocationsUsecase(ctx, param)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("WHEN page size is zero or negative, THEN should use default page size", func(t *testing.T) {
		mockRepo := mocks.NewLocationRepositoryInterface(t)
		usecase := NewLocationUsecase(mockRepo)
		ctx := context.Background()
		param := dto.GetLocationHandlerParam{
			Query:       "test",
			PageSize:    0,
			CurrentPage: 1,
			SortBy:      "name_ascend",
		}

		locations := []domain.Location{}
		expectedCount := 0
		mockRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			NameLike: "test",
			Limit:    10,
			Offset:   0,
			OrderBy:  "name_ascend",
		}).Return(locations, nil)
		mockRepo.On("GetLocationsCount", ctx).Return(expectedCount, nil)

		result, err := usecase.GetLocationsUsecase(ctx, param)

		assert.NoError(t, err)
		assert.Equal(t, 10, result.Data.PageSize)
	})

	t.Run("WHEN current page is zero or negative, THEN should use default current page", func(t *testing.T) {
		mockRepo := mocks.NewLocationRepositoryInterface(t)
		usecase := NewLocationUsecase(mockRepo)
		ctx := context.Background()
		param := dto.GetLocationHandlerParam{
			Query:       "test",
			PageSize:    10,
			CurrentPage: 0,
			SortBy:      "name_ascend",
		}

		locations := []domain.Location{}
		expectedCount := 0
		mockRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			NameLike: "test",
			Limit:    10,
			Offset:   0,
			OrderBy:  "name_ascend",
		}).Return(locations, nil)
		mockRepo.On("GetLocationsCount", ctx).Return(expectedCount, nil)

		result, err := usecase.GetLocationsUsecase(ctx, param)

		assert.NoError(t, err)
		assert.Equal(t, 1, result.Data.CurrentPage)
	})

	t.Run("WHEN no locations found, THEN should return empty results with zero count", func(t *testing.T) {
		mockRepo := mocks.NewLocationRepositoryInterface(t)
		usecase := NewLocationUsecase(mockRepo)
		ctx := context.Background()
		param := dto.GetLocationHandlerParam{
			Query:       "test",
			PageSize:    10,
			CurrentPage: 1,
			SortBy:      "name_ascend",
		}

		locations := []domain.Location{}
		expectedCount := 0
		mockRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			NameLike: "test",
			Limit:    10,
			Offset:   0,
			OrderBy:  "name_ascend",
		}).Return(locations, nil)
		mockRepo.On("GetLocationsCount", ctx).Return(expectedCount, nil)

		result, err := usecase.GetLocationsUsecase(ctx, param)

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, result.Data.Total)
		assert.Empty(t, result.Data.Items)
	})

	t.Run("WHEN multiple locations returned, THEN should return all items with no error", func(t *testing.T) {
		mockRepo := mocks.NewLocationRepositoryInterface(t)
		usecase := NewLocationUsecase(mockRepo)
		ctx := context.Background()
		param := dto.GetLocationHandlerParam{
			PageSize:    10,
			CurrentPage: 1,
			SortBy:      "name_ascend",
		}

		now := time.Now()
		locations := []domain.Location{
			{
				ID:             1,
				Name:           "Location 1",
				Country:        "Country 1",
				Latitude:       1.0,
				Longitude:      1.0,
				CreatedAt:      now,
				LastModifiedAt: sql.NullTime{Time: now},
			},
			{
				ID:             2,
				Name:           "Location 2",
				Country:        "Country 2",
				Latitude:       2.0,
				Longitude:      2.0,
				CreatedAt:      now,
				LastModifiedAt: sql.NullTime{Time: now},
			},
		}

		expectedCount := 2
		mockRepo.On("GetLocations", ctx, repository.GetLocationsParam{
			Limit:   10,
			Offset:  0,
			OrderBy: "name_ascend",
		}).Return(locations, nil)
		mockRepo.On("GetLocationsCount", ctx).Return(expectedCount, nil)

		result, err := usecase.GetLocationsUsecase(ctx, param)

		assert.NoError(t, err)
		assert.Equal(t, expectedCount, result.Data.Total)
		assert.Len(t, result.Data.Items, 2)
		assert.Equal(t, locations[0].ID, result.Data.Items[0].ID)
		assert.Equal(t, locations[1].ID, result.Data.Items[1].ID)
	})
}

func TestCreateLocationUsecase(t *testing.T) {
	t.Run("WHEN error occurred on insert location, THEN should return error", func(t *testing.T) {
		mockRepo := mocks.NewLocationRepositoryInterface(t)
		usecase := NewLocationUsecase(mockRepo)
		ctx := context.Background()
		req := dto.PostLocationHandlerRequest{
			Name:      "Test Location",
			Country:   "Test Country",
			Latitude:  1.0,
			Longitude: 1.0,
		}

		expectedError := errors.New("insert failed")
		mockRepo.On("InsertLocation", ctx, mock.Anything).Return(domain.Location{}, expectedError)

		_, err := usecase.CreateLocationUsecase(ctx, req)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("WHEN insert location succeeds, THEN should return no error", func(t *testing.T) {
		mockRepo := mocks.NewLocationRepositoryInterface(t)
		usecase := NewLocationUsecase(mockRepo)
		ctx := context.Background()
		req := dto.PostLocationHandlerRequest{
			Name:      "Test Location",
			Country:   "Test Country",
			Latitude:  1.0,
			Longitude: 1.0,
		}

		now := time.Now()
		expectedLocation := domain.Location{
			ID:             1,
			Name:           "Test Location",
			Country:        "Test Country",
			Latitude:       1.0,
			Longitude:      1.0,
			CreatedAt:      now,
			LastModifiedAt: sql.NullTime{Time: now},
		}

		mockRepo.On("InsertLocation", ctx, mock.Anything).Return(expectedLocation, nil)

		result, err := usecase.CreateLocationUsecase(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedLocation.ID, result.ID)
		assert.Equal(t, expectedLocation.Name, result.Name)
		assert.Equal(t, expectedLocation.Country, result.Country)
		assert.Equal(t, expectedLocation.Latitude, result.Latitude)
		assert.Equal(t, expectedLocation.Longitude, result.Longitude)
	})
}

