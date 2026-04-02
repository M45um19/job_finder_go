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
