package auth

import (
	"time"
)

type SeekerProfile struct {
	Address       *string `json:"address"`
	Experience    *string `json:"experience"`
	Skills        *string `json:"skills"`
	Education     *string `json:"education"`
	Certification *string `json:"certification"`
}

type EmployerProfile struct {
	CompanyName *string `json:"company_name"`
	CompanyDesc *string `json:"company_desc"`
	Address     *string `json:"address"`
}

type User struct {
	Id              int64            `json:"id,string"`
	Name            string           `json:"name"`
	Email           string           `json:"email"`
	Password        string           `json:"-"`
	Role            string           `json:"role"`
	UserPhoto       *string          `json:"userphoto"`
	RefreshToken    *string          `json:"refresh_token,omitempty"`
	SeekerProfile   *SeekerProfile   `json:"seeker_profile,omitempty"`
	EmployerProfile *EmployerProfile `json:"employer_profile,omitempty"`
	CreatedAt       time.Time        `json:"createdAt"`
}
