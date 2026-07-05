package application

import "time"

type Application struct {
	ID              int64     `json:"id,string"`
	ApplicantUserId int64     `json:"applicantUserId,string"`
	JobId           int64     `json:"jobId,string"`
	CreatedAt       time.Time `json:"createdAt"`
}
