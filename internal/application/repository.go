package application

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ApplicationRepository struct {
	db *pgxpool.Pool
}

func NewApplicationRepository(db *pgxpool.Pool) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (a *ApplicationRepository) CreateApplication(ctx context.Context, application *Application) error {
	query := "INSERT INTO applications (id, applicantUserId, jobId) VALUES ($1, $2, $3) RETURNING created_at"

	return a.db.QueryRow(ctx, query, application.ID, application.ApplicantUserId, application.JobId).Scan(&application.CreatedAt)
}

func (a *ApplicationRepository) GetApplicationByEmployeeId(ctx context.Context, employeeId int64) ([]Application, error) {
	query := "SELECT id, applicantUserId, jobId, created_at FROM applications WHERE applicantUserId=$1"

	rows, err := a.db.Query(ctx, query, employeeId)
	var applications []Application
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var application Application
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

func (a *ApplicationRepository) GetApplicationByJobId(ctx context.Context, jobId int64) ([]Application, error) {
	query := "SELECT id, applicantUserId, jobId, created_at FROM applications WHERE jobId=$1"

	rows, err := a.db.Query(ctx, query, jobId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var applications []Application
	for rows.Next() {
		var application Application
		err := rows.Scan(
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
