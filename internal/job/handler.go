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
	job, err := j.service.CreateJob(r.Context(), req.Title, req.Description, req.Company, req.Location, req.RequiredSkills, req.Salary, req.EmployerID)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, "Job created successfully", job)
}

func (j *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	} else if limit > 20 {
		limit = 20
	}

	jobs, totalItems, err := j.service.GetAllJobs(r.Context(), search, page, limit)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.PaginatedJSON(w, http.StatusOK, "Jobs retrieved successfully", jobs, totalItems, page, limit)
}

func (j *JobHandler) GetEmployerJobs(w http.ResponseWriter, r *http.Request) {
	employerID := r.Context().Value(auth.UserIdKey).(int64)

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	} else if limit > 20 {
		limit = 20
	}

	jobs, totalItems, err := j.service.GetJobsByEmployerID(r.Context(), employerID, page, limit)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.PaginatedJSON(w, http.StatusOK, "Employer jobs retrieved successfully", jobs, totalItems, page, limit)
}

func (j *JobHandler) GetSingleJobDetails(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	jobId, _ := strconv.ParseInt(idParam, 10, 64)

	job, err := j.service.GetSingleJobDetails(r.Context(), jobId)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "Job details retrieved successfully", job)
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
	utils.JSON(w, http.StatusOK, "Job updated successfully", job)
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

	utils.JSON(w, http.StatusOK, "Job deleted successfully", nil)
}
