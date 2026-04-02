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
