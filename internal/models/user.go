package models

import (
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	UpdatedAt time.Time `json:"updatedAt"`
}
