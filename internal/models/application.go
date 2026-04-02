package models

import "time"

type Application struct {
	ID              int64     `json:"id"`
	ApplicantUserId int64     `json:"applicantUserId"`
	JobId           int64     `json:"jobId"`
	CreatedAt       time.Time `json:"createdAt"`
}
