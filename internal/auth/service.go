package auth

import (
	"context"
	"errors"
	"io"

	"jobfinder/internal/platform/cloudinary"
	"jobfinder/internal/platform/idgen"
)

type AuthService struct {
	repo       *UserRepository
	jwtSecret  string
	cldService *cloudinary.Service
}

func NewAuthService(repo *UserRepository, secret string, cldService *cloudinary.Service) *AuthService {
	return &AuthService{repo: repo, jwtSecret: secret, cldService: cldService}
}

func (a *AuthService) Register(ctx context.Context, name, email, password, role string) (string, string, error) {
	if role != "employee" && role != "employer" {
		return "", "", errors.New("Invalid role")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return "", "", errors.New("Password hash failed")
	}

	userId := idgen.NextID()

	// Generate access & refresh tokens
	accessToken, err := GenerateAccessToken(userId, role, a.jwtSecret)
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := GenerateRefreshToken(userId, role, a.jwtSecret)
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	user := &User{
		Id:           userId,
		Name:         name,
		Email:        email,
		Password:     hash,
		Role:         role,
		RefreshToken: &refreshToken,
	}

	err = a.repo.Create(ctx, user)
	if err != nil {
		return "", "", errors.New("User can't be created")
	}

	return accessToken, refreshToken, nil
}

func (a *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := a.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", errors.New("User doesn't found")
	}

	err = ComparePassword(user.Password, password)
	if err != nil {
		return "", "", errors.New("Password not match")
	}

	accessToken, err := GenerateAccessToken(user.Id, user.Role, a.jwtSecret)
	if err != nil {
		return "", "", errors.New("access token generation failed")
	}

	refreshToken, err := GenerateRefreshToken(user.Id, user.Role, a.jwtSecret)
	if err != nil {
		return "", "", errors.New("refresh token generation failed")
	}

	// Update refresh token in DB on login
	err = a.repo.UpdateRefreshToken(ctx, user.Id, refreshToken)
	if err != nil {
		return "", "", errors.New("failed to save refresh token")
	}

	return accessToken, refreshToken, nil
}

func (a *AuthService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	// 1. Parse and validate the token string
	claims, err := ParseToken(refreshToken, a.jwtSecret)
	if err != nil {
		return "", errors.New("invalid or expired refresh token")
	}

	// 2. Fetch the user to verify this matches the active refresh token in database
	user, err := a.repo.GetUserByID(ctx, claims.UserId)
	if err != nil || user.RefreshToken == nil || *user.RefreshToken != refreshToken {
		return "", errors.New("refresh token revoked or not found")
	}

	// 3. Generate a new access token
	newAccessToken, err := GenerateAccessToken(user.Id, user.Role, a.jwtSecret)
	if err != nil {
		return "", errors.New("failed to generate access token")
	}

	return newAccessToken, nil
}

func (a *AuthService) GetProfile(ctx context.Context, userId int64) (*User, error) {
	user, err := a.repo.GetUserByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	// Sanitize sensitive info
	user.Password = ""
	user.RefreshToken = nil

	return user, nil
}

func (a *AuthService) UpdateProfilePhoto(ctx context.Context, userId int64, file io.Reader) (string, error) {
	if a.cldService == nil {
		return "", errors.New("Cloudinary service is not configured")
	}

	// 1. Retrieve current user info to check for existing photo
	user, err := a.repo.GetUserByID(ctx, userId)
	if err == nil && user.UserPhoto != nil && *user.UserPhoto != "" {
		// 2. Extract public ID and delete old file on Cloudinary
		publicID := cloudinary.ExtractPublicID(*user.UserPhoto)
		if publicID != "" {
			_ = a.cldService.DeleteFile(ctx, publicID) // ignore deletion error to avoid blocking the upload
		}
	}

	// 3. Upload new photo
	photoURL, err := a.cldService.UploadFile(ctx, file, "user_photos")
	if err != nil {
		return "", errors.New("failed to upload image to Cloudinary")
	}

	// 4. Update photo URL in database
	err = a.repo.UpdatePhoto(ctx, userId, photoURL)
	if err != nil {
		return "", errors.New("failed to update user photo in database")
	}

	return photoURL, nil
}

func (a *AuthService) UpdateEmployeeProfile(ctx context.Context, userId int64, name string, address, experience, skills, education, certification *string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	err := a.repo.UpdateEmployeeProfile(ctx, userId, name, address, experience, skills, education, certification)
	if err != nil {
		return errors.New("failed to update employee profile")
	}

	return nil
}

func (a *AuthService) UpdateEmployerProfile(ctx context.Context, userId int64, name string, companyName, companyDesc, address *string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}

	err := a.repo.UpdateEmployerProfile(ctx, userId, name, companyName, companyDesc, address)
	if err != nil {
		return errors.New("failed to update employer profile")
	}

	return nil
}
