package dto

type PostResponse struct {
	Status  string
	Message string
}

type HealthResponse PostResponse
