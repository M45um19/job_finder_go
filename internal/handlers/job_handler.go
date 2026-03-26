package handlers

import (
	"encoding/json"
	"jobfinder/internal/middleware"
	"jobfinder/internal/models"
	"jobfinder/internal/services"
	"jobfinder/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type JobHandler struct {
	service *services.JobService
}

func NewJobHandler(service *services.JobService) *JobHandler {
	return &JobHandler{service: service}
}

func (j *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIdKey).(int64)

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

func (j *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {

	jobs, err := j.service.GetAllJobs(r.Context())

	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
	}

	utils.JSON(w, http.StatusOK, jobs)
}

func (j *JobHandler) GetSingleJobDetails(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	jobId, _ := strconv.ParseInt(idParam, 10, 64)

	job, err := j.service.GetSingleJobDetails(r.Context(), jobId)

	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
	}

	utils.JSON(w, http.StatusOK, job)
}

func (j *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	jobId, _ := strconv.ParseInt(idParam, 10, 64)

	userId := r.Context().Value(middleware.UserIdKey).(int64)

	job := models.Job{}

	json.NewDecoder(r.Body).Decode(&job)
	job.ID = jobId
	job.EmployerID = userId
	err := j.service.UpdateJob(r.Context(), &job, userId)

	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.JSON(w, http.StatusOK, job)
}

func (j *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middleware.UserIdKey).(int64)

	idParam := chi.URLParam(r, "id")

	jobId, _ := strconv.ParseInt(idParam, 10, 64)

	err := j.service.DeleteJob(r.Context(), userId, jobId)

	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
	}

	utils.JSON(w, http.StatusOK, map[string]string{
		"message": "Job Deleted Successfully",
	})
}
