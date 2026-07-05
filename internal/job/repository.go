package job

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{db: db}
}

func (j *JobRepository) CreateJob(ctx context.Context, job *Job) error {
	query := "INSERT INTO jobs (id, title, description, company, location, employerid) VALUES ($1, $2, $3, $4, $5, $6) RETURNING created_at"

	return j.db.QueryRow(ctx, query, job.ID, job.Title, job.Description, job.Company, job.Location, job.EmployerID).Scan(&job.CreatedAt)
}

func (j *JobRepository) GetAllJobs(ctx context.Context, search string) ([]Job, error) {
	query := "SELECT id, title, description, company, location, created_at, updated_at FROM jobs"
	var args []interface{}

	if search != "" {
		query += " WHERE title ILIKE $1 OR description ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	rows, err := j.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
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

func (j *JobRepository) GetSingleJobDetails(ctx context.Context, jobId int64) (*Job, error) {
	var job Job
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

func (j *JobRepository) UpdateJob(ctx context.Context, job *Job) error {
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

func (j *JobRepository) DeleteJob(ctx context.Context, jobId int64) error {
	query := "DELETE FROM jobs WHERE id=$1"

	_, err := j.db.Exec(ctx, query, jobId)
	return err
}
