package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tyarus/weather-app/internal/dto"
	"tyarus/weather-app/internal/usecase"
	"tyarus/weather-app/pkg/response"
)

type locationHandler struct {
	locationUc usecase.LocationUsecaseInterface
}

func NewLocationHandler(locaionUc usecase.LocationUsecaseInterface) locationHandler {
	return locationHandler{locationUc: locaionUc}
}

func (h *locationHandler) GetLocationHandler() http.HandlerFunc {
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

		param := dto.GetLocationHandlerParam{
			Query:       r.URL.Query().Get("query"),
			SortBy:      r.URL.Query().Get("sortBy"),
			PageSize:    pageSize,
			CurrentPage: currentPage,
		}

		err = param.Validate()
		if err != nil {
			response.Error(w, http.StatusBadRequest, err.Error())
			return
		}

		locations, err := h.locationUc.GetLocationsUsecase(ctx, param)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error occurred on fetch locations: "+err.Error())
			return
		}

		response.JSON(w, http.StatusOK, "success", "fetch locations successfully", locations)
	}
}

func (h *locationHandler) CreateLocationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req dto.PostLocationHandlerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body: "+err.Error())
			return
		}

		err := req.Validate()
		if err != nil {
			response.Error(w, http.StatusBadRequest, "invalid request body: "+err.Error())
			return
		}

		result, err := h.locationUc.CreateLocationUsecase(ctx, req)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "failed to insert location: "+err.Error())
			return
		}

		response.JSON(w, http.StatusOK, "success", "create location successfully", result)
	}
}
