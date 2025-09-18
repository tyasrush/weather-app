package repository

import (
	"context"
	"database/sql"
	"fmt"
	"tyarus/weather-app/internal/domain"
	"tyarus/weather-app/pkg/utils"
)

type GetWeathersParam struct {
	LocationID int64
	Limit      int
	Offset     int
	OrderBy    string
}

type WeatherRepositoryInterface interface {
	GetWeathers(ctx context.Context, param GetWeathersParam) ([]domain.Weather, error)
	BulkUpsertWeather(ctx context.Context, weathers []domain.Weather) ([]domain.Weather, error)
	GetWeathersCount(ctx context.Context) (int, error)
}

type weatherRepository struct {
	db *sql.DB
}

func NewWeatherRepository(db *sql.DB) WeatherRepositoryInterface {
	return &weatherRepository{db: db}
}

func (r *weatherRepository) GetWeathers(ctx context.Context, param GetWeathersParam) ([]domain.Weather, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultDBTimeout)
	defer cancel()

	query := `SELECT * FROM weathers WHERE deleted_at IS NULL`
	params := []interface{}{}
	if param.LocationID != 0 {
		query += " AND location_id = ?"
		params = append(params, param.LocationID)
	}

	if param.OrderBy != "" {
		query += " ORDER BY " + param.OrderBy
	} else {
		query += " ORDER BY created_at DESC"
	}

	if param.Limit > 0 {
		query += " LIMIT ?"
		params = append(params, param.Limit)
	}

	if param.Offset > 0 {
		query += " OFFSET ?"
		params = append(params, param.Offset)
	}

	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to query weathers: %w", err)
	}
	defer rows.Close()

	var weathers []domain.Weather
	for rows.Next() {
		var w domain.Weather
		err := rows.Scan(
			&w.ID,
			&w.LocationID,
			&w.TemperatureCelcius,
			&w.TemperatureFahrenheit,
			&w.Humidity,
			&w.WindSpeed,
			&w.ConditionStatus,
			&w.ConditionIconURL,
			&w.ForecastTime,
			&w.ForecastType,
			&w.CreatedAt,
			&w.LastModifiedAt,
			&w.DeletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan weather: %w", err)
		}
		weathers = append(weathers, w)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return weathers, nil
}

func (r *weatherRepository) BulkUpsertWeather(ctx context.Context, weathers []domain.Weather) ([]domain.Weather, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultDBTimeout)
	defer cancel()

	if len(weathers) == 0 {
		return weathers, nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO weathers (location_id, temperature_celcius, temperature_fahrenheit, humidity, wind_speed, condition_status, condition_icon_url, forecast_time, forecast_type) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			  ON DUPLICATE KEY UPDATE 
			  temperature_celcius = VALUES(temperature_celcius),
			  temperature_fahrenheit = VALUES(temperature_fahrenheit),
			  humidity = VALUES(humidity),
			  wind_speed = VALUES(wind_speed),
			  condition_status = VALUES(condition_status),
			  condition_icon_url = VALUES(condition_icon_url)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	upsertedWeathers := make([]domain.Weather, len(weathers))
	for _, weather := range weathers {
		_, err := stmt.ExecContext(
			ctx,
			weather.LocationID,
			weather.TemperatureCelcius,
			weather.TemperatureFahrenheit,
			weather.Humidity,
			weather.WindSpeed,
			weather.ConditionStatus,
			weather.ConditionIconURL,
			weather.ForecastTime,
			weather.ForecastType,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to upsert weather: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return upsertedWeathers, nil
}

func (r *weatherRepository) GetWeathersCount(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultDBTimeout)
	defer cancel()

	query := `SELECT COUNT(*) FROM weathers WHERE deleted_at IS NULL`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count weathers: %w", err)
	}
	return count, nil
}
