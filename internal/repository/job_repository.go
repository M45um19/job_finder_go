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

func (j *JobRepository) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	query := "SELECT id, title, description, company, location, created_at, updated_at FROM jobs"

	rows, err := j.db.Query(ctx, query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []models.Job

	for rows.Next() {
		var job models.Job

		err := rows.Scan(
			&job.ID,
			&job.Title,
			&job.Description,
			&job.Company,
			&job.Location,
			&job.CreatedAt,
			&job.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (j *JobRepository) GetSingleJobDetails(ctx context.Context, jobId int64) (*models.Job, error) {
	var job models.Job

	query := "SELECT id, title, description, company, location, employerid, created_at, updated_at FROM jobs WHERE id=$1"

	err := j.db.QueryRow(ctx, query, jobId).Scan(
		&job.ID,
		&job.Title,
		&job.Description,
		&job.Company,
		&job.Location,
		&job.EmployerID,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (j *JobRepository) UpdateJob(ctx context.Context, job *models.Job) error {
	query := "UPDATE jobs SET title=$1, description=$2, company=$3, location=$4, updated_at=NOW() WHERE id=$5"

	_, err := j.db.Exec(ctx, query,
		job.Title,
		job.Description,
		job.Company,
		job.Location,
		job.ID,
	)

	return err
}
