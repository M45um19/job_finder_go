package services

import (
	"context"
	"errors"
	"jobfinder/internal/models"
	"jobfinder/internal/repository"
	"log"
)

type JobService struct {
	repo *repository.JobRepository
}

func NewJobService(repo *repository.JobRepository) *JobService {
	return &JobService{repo: repo}
}

func (j *JobService) CreateJob(ctx context.Context, title, description, company, localtion string, employerId int64) (*models.Job, error) {

	job := models.Job{
		Title:       title,
		Description: description,
		Company:     company,
		Location:    localtion,
		EmployerID:  employerId,
	}

	err := j.repo.CreateJob(ctx, &job)

	if err != nil {
		log.Printf("ERROR: failed to create job in database: %v", err)
		return nil, errors.New("Job creation faild")
	}

	return &job, nil

}

func (j *JobService) GetAllJobs(ctx context.Context) ([]models.Job, error) {

	return j.repo.GetAllJobs(ctx)

}

func (j *JobService) GetSingleJobDetails(ctx context.Context, jobId int64) (*models.Job, error) {
	return j.repo.GetSingleJobDetails(ctx, jobId)
}
