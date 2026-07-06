package job

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository {
	return &JobRepository{db: db}
}

func (j *JobRepository) CreateJob(ctx context.Context, job *Job) error {
	query := "INSERT INTO jobs (id, title, description, company, location, employerid, required_skills, salary) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING created_at"

	return j.db.QueryRow(ctx, query, job.ID, job.Title, job.Description, job.Company, job.Location, job.EmployerID, job.RequiredSkills, job.Salary).Scan(&job.CreatedAt)
}

func (j *JobRepository) GetAllJobs(ctx context.Context, search string, page, limit int) ([]Job, int64, error) {
	// 1. Get count
	countQuery := "SELECT COUNT(*) FROM jobs"
	var countArgs []interface{}
	if search != "" {
		countQuery += " WHERE title ILIKE $1 OR description ILIKE $1 OR required_skills ILIKE $1"
		countArgs = append(countArgs, "%"+search+"%")
	}

	var totalItems int64
	err := j.db.QueryRow(ctx, countQuery, countArgs...).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// 2. Fetch records
	offset := (page - 1) * limit
	query := "SELECT id, title, description, company, location, employerid, required_skills, salary, created_at, updated_at FROM jobs"
	var args []interface{}

	if search != "" {
		query += " WHERE title ILIKE $1 OR description ILIKE $1 OR required_skills ILIKE $1"
		args = append(args, "%"+search+"%")
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := j.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
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
			&job.EmployerID,
			&job.RequiredSkills,
			&job.Salary,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}
	return jobs, totalItems, nil
}

func (j *JobRepository) GetSingleJobDetails(ctx context.Context, jobId int64) (*Job, error) {
	var job Job
	query := "SELECT id, title, description, company, location, employerid, required_skills, salary, created_at, updated_at FROM jobs WHERE id=$1"

	err := j.db.QueryRow(ctx, query, jobId).Scan(
		&job.ID,
		&job.Title,
		&job.Description,
		&job.Company,
		&job.Location,
		&job.EmployerID,
		&job.RequiredSkills,
		&job.Salary,
		&job.CreatedAt,
		&job.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (j *JobRepository) UpdateJob(ctx context.Context, job *Job) error {
	query := "UPDATE jobs SET title=$1, description=$2, company=$3, location=$4, required_skills=$5, salary=$6, updated_at=NOW() WHERE id=$7"

	_, err := j.db.Exec(ctx, query,
		job.Title,
		job.Description,
		job.Company,
		job.Location,
		job.RequiredSkills,
		job.Salary,
		job.ID,
	)
	return err
}

func (j *JobRepository) DeleteJob(ctx context.Context, jobId int64) error {
	query := "DELETE FROM jobs WHERE id=$1"

	_, err := j.db.Exec(ctx, query, jobId)
	return err
}

func (j *JobRepository) GetJobsByEmployerID(ctx context.Context, employerID int64, page, limit int) ([]Job, int64, error) {
	// 1. Get total count
	countQuery := "SELECT COUNT(*) FROM jobs WHERE employerid = $1"
	var totalItems int64
	err := j.db.QueryRow(ctx, countQuery, employerID).Scan(&totalItems)
	if err != nil {
		return nil, 0, err
	}

	// 2. Fetch records
	offset := (page - 1) * limit
	query := "SELECT id, title, description, company, location, employerid, required_skills, salary, created_at, updated_at FROM jobs WHERE employerid = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"

	rows, err := j.db.Query(ctx, query, employerID, limit, offset)
	if err != nil {
		return nil, 0, err
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
			&job.EmployerID,
			&job.RequiredSkills,
			&job.Salary,
			&job.CreatedAt,
			&job.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		jobs = append(jobs, job)
	}

	return jobs, totalItems, nil
}
