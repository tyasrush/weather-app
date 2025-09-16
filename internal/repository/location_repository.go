package repository

import (
	"context"
	"database/sql"
	"fmt"
	"tyarus/weather-app/internal/domain"
	"tyarus/weather-app/pkg/utils"
)

type GetLocationsParam struct {
	NameLike string
	Limit    int
	Offset   int
	OrderBy  string
}

type LocationRepositoryInterface interface {
	GetLocations(ctx context.Context, param GetLocationsParam) ([]domain.Location, error)
	InsertLocation(ctx context.Context, location domain.Location) (domain.Location, error)
	GetLocationsCount(ctx context.Context) (int, error)
}

type locationRepository struct {
	db *sql.DB
}

func NewLocationRepository(db *sql.DB) LocationRepositoryInterface {
	return &locationRepository{db: db}
}

func (r *locationRepository) GetLocations(ctx context.Context, param GetLocationsParam) ([]domain.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultDBTimeout)
	defer cancel()

	query := `SELECT * FROM locations `
	statement := ""
	params := []interface{}{}
	if param.NameLike != "" {
		query = query + "WHERE name ILIKE %?%"
		params = append(params, param.NameLike)
	}

	if param.OrderBy != "" {
		ascStatement := "ORDER BY %s ASC"
		descStatement := "ORDER BY %s DESC"
		switch param.OrderBy {
		case utils.OrderByCreatedAtAsc:
			statement = statement + fmt.Sprintf(ascStatement, "created_at")
		case utils.OrderByCreatedAtDesc:
			statement = statement + fmt.Sprintf(descStatement, "created_at")
		case utils.OrderByNameAsc:
			statement = statement + fmt.Sprintf(ascStatement, "name")
		case utils.OrderByNameDesc:
			statement = statement + fmt.Sprintf(descStatement, "name")
		default:
			statement = statement + fmt.Sprintf(ascStatement, "created_at")
		}
	}

	if param.Limit > 0 {
		statement = statement + "LIMIT ?"
		params = append(params, param.Limit)
	}

	if param.Offset > 0 {
		statement = statement + "OFFSET ?"
		params = append(params, param.Offset)
	}

	query = query + statement
	rows, err := r.db.QueryContext(ctx, query, params...)
	if err != nil {
		return nil, fmt.Errorf("execute query locations: %w", err)
	}
	defer rows.Close()

	var locations []domain.Location
	for rows.Next() {
		var item domain.Location
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Region,
			&item.Country,
			&item.Latitude,
			&item.Longitude,
			&item.CreatedAt,
			&item.LastModifiedAt,
			&item.DeletedAt); err != nil {
			return nil, fmt.Errorf("scan location: %w", err)
		}
		locations = append(locations, item)
	}

	return locations, nil
}

func (r *locationRepository) InsertLocation(ctx context.Context, location domain.Location) (domain.Location, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultDBTimeout)
	defer cancel()

	res, err := r.db.ExecContext(ctx, `INSERT INTO locations (name, region, country, latitude, longitude) VALUES (?, ?, ?, ?, ?)`,
		location.Name, location.Region, location.Country, location.Latitude, location.Longitude,
	)
	if err != nil {
		return location, fmt.Errorf("insert location: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return location, fmt.Errorf("get last insert id: %w", err)
	}

	location.ID = id
	return location, nil
}

func (r *locationRepository) GetLocationsCount(ctx context.Context) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, utils.DefaultDBTimeout)
	defer cancel()

	var count int
	query := "SELECT COUNT(*) FROM locations"
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("get locations count: %w", err)
	}

	return count, err
}
