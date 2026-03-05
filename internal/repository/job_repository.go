package repository


type JobRepository struct {
	db *pgxpool.Pool
}

func NewJobRepository(db *pgxpool.Pool) *JobRepository{
	return &JobRepository{db: db}
}

func (*JobRepository) CreateJob(ctx context.Context, job *models.Job) error {
	
}