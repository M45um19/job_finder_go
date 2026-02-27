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

func (r *UserRepository) Create(u *models.User) error {
	query := `INSERT INTO users (name, email, password, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	return r.db.QueryRow(context.Background(), query, u.Name, u.Email, u.Password, u.Role).Scan(&u.Id, &u.CreatedAt)
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT id, name, email, password, role, created_at  FROM users WHERE email=$1`

	user := &models.User{}

	err := r.db.QueryRow(context.Background(), query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.CreatedAt)

	return user, err
}
