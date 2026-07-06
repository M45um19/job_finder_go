package utils

import (
	"encoding/json"
	"net/http"
)

type PaginationMeta struct {
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
	Limit       int   `json:"limit"`
}

type PaginatedResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

type StandardResponse struct {
	Success    bool            `json:"success"`
	StatusCode int             `json:"statusCode"`
	Message    string          `json:"message"`
	Data       interface{}     `json:"data"`
	Meta       *PaginationMeta `json:"meta,omitempty"`
}

func JSON(w http.ResponseWriter, status int, message string, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	isSuccess := status < 400

	if sr, ok := payload.(StandardResponse); ok {
		json.NewEncoder(w).Encode(sr)
		return
	}

	wrapped := StandardResponse{
		Success:    isSuccess,
		StatusCode: status,
		Message:    message,
		Data:       payload,
	}
	json.NewEncoder(w).Encode(wrapped)
}

func PaginatedJSON(w http.ResponseWriter, status int, message string, data interface{}, totalItems int64, page int, limit int) {
	totalPages := int((totalItems + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	payload := StandardResponse{
		Success:    status < 400,
		StatusCode: status,
		Message:    message,
		Data:       data,
		Meta: &PaginationMeta{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: page,
			Limit:       limit,
		},
	}
	json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, message, nil)
}
