package application

import (
	"context"
	"errors"

	"jobfinder/internal/platform/idgen"
)

type ApplicationService struct {
	repo *ApplicationRepository
}

func NewApplicationService(repo *ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: repo}
}

func (a *ApplicationService) CreateApplication(ctx context.Context, applicantUserId, jobId int64) (*Application, error) {
	application := Application{
		ID:              idgen.NextID(),
		ApplicantUserId: applicantUserId,
		JobId:           jobId,
	}
	err := a.repo.CreateApplication(ctx, &application)
	if err != nil {
		return nil, errors.New("Application creation faild")
	}

	return &application, nil
}

func (a *ApplicationService) GetApplicationByEmployeeId(ctx context.Context, employeeId int64) ([]Application, error) {
	return a.repo.GetApplicationByEmployeeId(ctx, employeeId)
}

func (a *ApplicationService) GetApplicationByJobId(ctx context.Context, jobId int64) ([]Application, error) {
	return a.repo.GetApplicationByJobId(ctx, jobId)
}
