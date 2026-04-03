package repository

import (
	"context"
	"jobfinder/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationRepository struct {
	db *pgxpool.Pool
}

func NewApplicationRepository(db *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (a *ApplicationRepository) CreateApplication(ctx context.Context, application *models.Application) error {
	query := "INSERT INTO applications (applicantUserId, jobId) VALUES ($1, $2) RETURNING id, created_at"

	return a.db.QueryRow(ctx, query, application.ApplicantUserId, application.JobId).Scan(&application.ID, &application.CreatedAt)
}

func (a *ApplicationRepository) GetApplicationByEmployeeId(ctx context.Context, employeeId int64) ([]models.Application, error) {
	query := "SELECT id, applicantUserId, jobId, created_at FROM applications WHERE applicantUserId=$1"

	rows, err := a.db.Query(ctx, query, employeeId)
	var applications []models.Application
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var application models.Application

		err = rows.Scan(
			&application.ID,
			&application.ApplicantUserId,
			&application.JobId,
			&application.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		applications = append(applications, application)
	}
	return applications, nil
}
