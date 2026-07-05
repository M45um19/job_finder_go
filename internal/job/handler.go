package job

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"jobfinder/internal/auth"
	"jobfinder/internal/platform/utils"
)

type JobHandler struct {
	service *JobService
}

func NewJobHandler(service *JobService) *JobHandler {
	return &JobHandler{service: service}
}

func (j *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(auth.UserIdKey).(int64)

	req := Job{}
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
	search := r.URL.Query().Get("search")
	jobs, err := j.service.GetAllJobs(r.Context(), search)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, jobs)
}

func (j *JobHandler) GetSingleJobDetails(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	jobId, _ := strconv.ParseInt(idParam, 10, 64)

	job, err := j.service.GetSingleJobDetails(r.Context(), jobId)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, job)
}

func (j *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	jobId, _ := strconv.ParseInt(idParam, 10, 64)

	userId := r.Context().Value(auth.UserIdKey).(int64)

	job := Job{}
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
	userId := r.Context().Value(auth.UserIdKey).(int64)

	idParam := chi.URLParam(r, "id")
	jobId, _ := strconv.ParseInt(idParam, 10, 64)

	err := j.service.DeleteJob(r.Context(), userId, jobId)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, map[string]string{
		"message": "Job Deleted Successfully",
	})
}
