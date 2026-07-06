package auth

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"jobfinder/internal/platform/utils"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := a.service.Register(r.Context(), req.Name, req.Email, req.Password, req.Role)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusCreated, "User registered successfully", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := a.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "Login successfully", map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (a *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	newAccessToken, err := a.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		utils.Error(w, http.StatusUnauthorized, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "Token refreshed successfully", map[string]string{
		"access_token": newAccessToken,
	})
}

func (a *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	userId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := a.service.GetProfile(r.Context(), userId)
	if err != nil {
		utils.Error(w, http.StatusNotFound, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "User profile retrieved successfully", user)
}

func (a *AuthHandler) UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserIdKey).(int64)

	// Parse multipart form (max 5MB file size)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to parse form: file too large")
		return
	}

	file, _, err := r.FormFile("photo")
	if err != nil {
		utils.Error(w, http.StatusBadRequest, "photo file is required")
		return
	}
	defer file.Close()

	photoURL, err := a.service.UpdateProfilePhoto(r.Context(), userId, file)
	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "Profile photo updated successfully", map[string]string{"photo_url": photoURL})
}

func (a *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(UserIdKey).(int64)
	role := r.Context().Value(RoleKey).(string)

	var req struct {
		Name          string  `json:"name"`
		Address       *string `json:"address"`
		Experience    *string `json:"experience"`
		Skills        *string `json:"skills"`
		Education     *string `json:"education"`
		Certification *string `json:"certification"`
		CompanyName   *string `json:"company_name"`
		CompanyDesc   *string `json:"company_desc"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var err error
	if role == "employee" {
		err = a.service.UpdateEmployeeProfile(r.Context(), userId, req.Name, req.Address, req.Experience, req.Skills, req.Education, req.Certification)
	} else if role == "employer" {
		err = a.service.UpdateEmployerProfile(r.Context(), userId, req.Name, req.CompanyName, req.CompanyDesc, req.Address)
	} else {
		utils.Error(w, http.StatusForbidden, "Invalid user role")
		return
	}

	if err != nil {
		utils.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSON(w, http.StatusOK, "Profile updated successfully", nil)
}
