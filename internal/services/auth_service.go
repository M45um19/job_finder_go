package services

import (
	"context"
	"errors"
	"jobfinder/internal/models"
	"jobfinder/internal/repository"
	"jobfinder/internal/utils"
)

type AuthService struct {
	repo      *repository.UserRepository
	jwtSecret string
}

func NewAuthService(repo *repository.UserRepository, secret string) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret}
}

func (a *AuthService) Register(ctx context.Context, name, email, password, role string) (*models.User, error) {
	if role != "employee" && role != "employer" {
		return nil, errors.New("Invalid role")
	}

	hash, err := utils.HashPassword(password)

	if err != nil {
		return nil, errors.New("Password hash faild")

	}

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     role,
	}

	err = a.repo.Create(ctx, user)

	if err != nil {
		return nil, errors.New("User can't be created")
	}

	return user, nil

}

func (a *AuthService) Login(ctx context.Context, email, password string) (string, error) {

	user, err := a.repo.GetUserByEmail(ctx, email)

	if err != nil {
		return "", errors.New("User doesn't found")
	}

	err = utils.ComparePassword(user.Password, password)

	if err != nil {
		return "", errors.New("Password not match")
	}

	token, err := utils.GenerateToken(user.Id, user.Role, a.jwtSecret)

	if err != nil {
		return "", errors.New("token generation faild")
	}

	return token, nil
}
