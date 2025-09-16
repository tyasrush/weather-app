package response

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

type PaginationData[T any] struct {
	Items       []T `json:"items"`
	Total       int `json:"total"`
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
}

func JSON[T any](w http.ResponseWriter, statusCode int, status, message string, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response[T]{
		Status:  status,
		Message: message,
		Data:    data,
	}

	_ = json.NewEncoder(w).Encode(resp)
}

func Error(w http.ResponseWriter, statusCode int, message string) {
	log.Println("error : ", message)
	JSON[any](w, statusCode, "error", message, nil)
}
