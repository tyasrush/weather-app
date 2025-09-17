package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/internal/usecase"
	"tyarus/weather-app/pkg/response"
)

type weatherHandler struct {
	weatherUc usecase.WeatherUsecaseInterface
}

func NewWeatherHandler(weatherUc usecase.WeatherUsecaseInterface) weatherHandler {
	return weatherHandler{weatherUc: weatherUc}
}

func (h *weatherHandler) GetWeathersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		pageSize, currentPage := 0, 0
		var err error
		if r.URL.Query().Get("pageSize") != "" {
			pageSize, err = strconv.Atoi(r.URL.Query().Get("pageSize"))
			if err != nil {
				response.Error(w, http.StatusBadRequest, "invalid pageSize parameter, please check your parameter")
				return
			}
		}

		if r.URL.Query().Get("currentPage") != "" {
			currentPage, err = strconv.Atoi(r.URL.Query().Get("currentPage"))
			if err != nil {
				response.Error(w, http.StatusBadRequest, "invalid currentPage parameter, please check your parameter")
				return
			}
		}

		if r.URL.Query().Get("locationID") == "" {
			response.Error(w, http.StatusBadRequest, "locationID parameter is empty, please check your parameter")
			return
		}

		locationID, err := strconv.Atoi(r.URL.Query().Get("locationID"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid locationID parameter, please check your parameter")
			return
		}

		param := dto.GetWeathersParam{
			PageSize:    pageSize,
			CurrentPage: currentPage,
			LocationID:  locationID,
		}

		weathers, err := h.weatherUc.GetWeathersUsecase(ctx, param)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error occurred on fetch weathers: "+err.Error())
			return
		}

		response.JSON(w, http.StatusOK, "success", "fetch locations successfully", weathers)
	}
}

func (h *weatherHandler) SyncWeatherHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req dto.PostWeatherSyncUsecaseRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body: "+err.Error())
			return
		}

		err := h.weatherUc.SyncWeatherUsecase(ctx, req)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to sync weather: "+err.Error())
			return
		}

		response.JSON(w, http.StatusOK, "success", "sync weather successfully", req)
	}
}
