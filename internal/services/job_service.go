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

func (j *JobService) CreateJob(ctx context.Context, title, description, company, location string, employerId int64) (*models.Job, error) {

	job := models.Job{
		Title:       title,
		Description: description,
		Company:     company,
		Location:    location,
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

func (j *JobService) UpdateJob(ctx context.Context, job *models.Job, userId int64) error {
	job, err := j.repo.GetSingleJobDetails(ctx, job.ID)

	if err != nil {
		return err
	}

	if job.EmployerID != userId {
		return errors.New("You are not job owner")
	}

	return j.repo.UpdateJob(ctx, job)
}
