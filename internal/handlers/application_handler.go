package handlers

import (
	"jobfinder/internal/middleware"
	"jobfinder/internal/services"
	"jobfinder/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ApplicationHandler struct {
	service *services.ApplicationService
}

func NewApplicationHandler(service *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

func (a *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	applicantUserId := r.Context().Value(middleware.UserIdKey).(int64)

	jobIdParam := chi.URLParam(r, "id")

	jobId, _ := strconv.ParseInt(jobIdParam, 10, 64)

	application, err := a.service.CreateApplication(r.Context(), applicantUserId, jobId)

	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, application)
}
