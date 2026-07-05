package job

import "time"

type Job struct {
	ID          int64     `json:"id,string"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Company     string    `json:"company"`
	Location    string    `json:"location"`
	EmployerID  int64     `json:"employerid,string"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
