package repository

import (
	"context"
	"jobfinder/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(ctx context.Context, u *models.User) error {
	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	return ur.db.QueryRow(ctx, query, u.Name, u.Email, u.Password, u.Role).Scan(&u.Id, &u.CreatedAt)
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, password, role, created_at  FROM users WHERE email=$1`

	user := &models.User{}

	err := ur.db.QueryRow(ctx, query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role, &user.CreatedAt)

	return user, err
}
