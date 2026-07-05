package application

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"jobfinder/internal/auth"
	"jobfinder/internal/platform/utils"
)

type ApplicationHandler struct {
	service *ApplicationService
}

func NewApplicationHandler(service *ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: service}
}

func (a *ApplicationHandler) CreateApplication(w http.ResponseWriter, r *http.Request) {
	applicantUserId := r.Context().Value(auth.UserIdKey).(int64)

	jobIdParam := chi.URLParam(r, "id")
	jobId, _ := strconv.ParseInt(jobIdParam, 10, 64)

	application, err := a.service.CreateApplication(r.Context(), applicantUserId, jobId)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, application)
}

func (a *ApplicationHandler) GetApplicationByEmployeeId(w http.ResponseWriter, r *http.Request) {
	employeeId := r.Context().Value(auth.UserIdKey).(int64)

	applications, err := a.service.GetApplicationByEmployeeId(r.Context(), employeeId)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, applications)
}

func (a *ApplicationHandler) GetApplicationByJobId(w http.ResponseWriter, r *http.Request) {
	jobIdParam := chi.URLParam(r, "id")
	jobId, _ := strconv.ParseInt(jobIdParam, 10, 64)

	applications, err := a.service.GetApplicationByJobId(r.Context(), jobId)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, applications)
}
