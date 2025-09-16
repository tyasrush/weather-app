package domain

import (
	"database/sql"
	"time"
)

type Location struct {
	ID             int64        `json:"id"`
	Name           string       `json:"name"`
	Region         string       `json:"region"`
	Country        string       `json:"country"`
	Latitude       float64      `json:"latitude"`
	Longitude      float64      `json:"longitude"`
	CreatedAt      time.Time    `json:"created_at"`
	LastModifiedAt sql.NullTime `json:"last_modified_at"`
	DeletedAt      sql.NullTime `json:"deleted_at"`
}
