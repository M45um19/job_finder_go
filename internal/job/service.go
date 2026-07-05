package job

import (
	"context"
	"errors"
	"log"

	"jobfinder/internal/platform/idgen"
)

type JobService struct {
	repo *JobRepository
}

func NewJobService(repo *JobRepository) *JobService {
	return &JobService{repo: repo}
}

func (j *JobService) CreateJob(ctx context.Context, title, description, company, location string, employerId int64) (*Job, error) {
	job := Job{
		ID:          idgen.NextID(),
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

func (j *JobService) GetAllJobs(ctx context.Context, search string) ([]Job, error) {
	return j.repo.GetAllJobs(ctx, search)
}

func (j *JobService) GetSingleJobDetails(ctx context.Context, jobId int64) (*Job, error) {
	return j.repo.GetSingleJobDetails(ctx, jobId)
}

func (j *JobService) UpdateJob(ctx context.Context, job *Job, userId int64) error {
	existing, err := j.repo.GetSingleJobDetails(ctx, job.ID)
	if err != nil {
		return err
	}

	if existing.EmployerID != userId {
		return errors.New("You are not job owner")
	}

	return j.repo.UpdateJob(ctx, job)
}

func (j *JobService) DeleteJob(ctx context.Context, userId int64, jobId int64) error {
	existing, err := j.repo.GetSingleJobDetails(ctx, jobId)
	if err != nil {
		return err
	}

	if existing.EmployerID != userId {
		return errors.New("You are not job owner")
	}

	return j.repo.DeleteJob(ctx, jobId)
}
