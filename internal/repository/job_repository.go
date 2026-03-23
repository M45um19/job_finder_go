package repository

import (
	"context"
	"jobfinder/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{db: db}
}

func (j *JobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	query := "INSERT INTO jobs (title, description, company, location, employerid) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at"

	return j.db.QueryRow(ctx, query, job.Title, job.Description, job.Company, job.Location, job.EmployerID).Scan(&job.ID, &job.CreatedAt)

}
