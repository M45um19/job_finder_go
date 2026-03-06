package handlers

import (
	"encoding/json"
	"jobfinder/internal/models"
	"jobfinder/internal/services"
	"jobfinder/internal/utils"
	"net/http"
)

type Jobhandler struct {
	service *services.JobService
}

func NewJobHandler(service *services.JobService) *Jobhandler {
	return &Jobhandler{service: service}
}

func (j *Jobhandler) CreateJob(w http.ResponseWriter, r http.Request) {
	userId := r.Context().Value("UserIdKey").(int64)

	req := models.Job{}

	json.NewDecoder(r.Body).Decode(&req)

	req.EmployerID = userId

	job, err := j.service.CreateJob(r.Context(), req.Title, req.Description, req.Company, req.Location, req.EmployerID)

	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, job)
}
