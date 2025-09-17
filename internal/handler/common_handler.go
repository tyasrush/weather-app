package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/internal/infra"
)

type commonHandler struct {
	db    *sql.DB
	cache infra.CacheInterface
}

func NewCommonHandler(db *sql.DB, cache infra.CacheInterface) commonHandler {
	return commonHandler{db: db, cache: cache}
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
		if err := h.cache.Ping(r.Context()); err != nil {
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
