package services

import (
	"context"
	"errors"
	"jobfinder/internal/models"
	"jobfinder/internal/repository"
)

type ApplicationService struct {
	repo *repository.ApplicationRepository
}

func NewApplicationService(repo *repository.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (a *ApplicationService) CreateApplication(ctx context.Context, applicantUserId, jobId int64) (*models.Application, error) {
	application := models.Application{
		ApplicantUserId: applicantUserId,
		JobId:           jobId,
	}
	err := a.repo.CreateApplication(ctx, &application)

	if err != nil {
		return nil, errors.New("Application creation faild")
	}

	return &application, nil
}

func (a *ApplicationService) GetApplicationByEmployeeId(ctx context.Context, employeeId int64) ([]models.Application, error) {
	return a.repo.GetApplicationByEmployeeId(ctx, employeeId)
}

func (a *ApplicationService) GetApplicationByJobId(ctx context.Context, jobId int64) ([]models.Application, error) {
	return a.repo.GetApplicationByJobId(ctx, jobId)
}
