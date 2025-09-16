package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"tyarus/weather-app/internal/dto"

	"github.com/redis/go-redis/v9"
)

type commonHandler struct {
	db          *sql.DB
	redisClient *redis.Client
}

func NewCommonHandler(db *sql.DB, redisClient *redis.Client) commonHandler {
	return commonHandler{db: db, redisClient: redisClient}
}

func (h *commonHandler) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		resp := dto.HealthResponse{Status: "success", Message: "app running!"}
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h *commonHandler) ReadyCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if err := h.db.Ping(); err != nil {
			http.Error(w, "MySQL not reachable", http.StatusInternalServerError)
			return
		}
		if _, err := h.redisClient.Ping(r.Context()).Result(); err != nil {
			http.Error(w, "Redis not reachable", http.StatusInternalServerError)
			return
		}

		resp := dto.HealthResponse{Status: "success", Message: "all resource running!"}
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
