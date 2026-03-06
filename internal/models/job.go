package models

type Job struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Company     string `json:"company"`
	Location    string `json:"localtion"`
	EmployerID  int64  `json:"employerid"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `jsonL"updatedAt"`
}
