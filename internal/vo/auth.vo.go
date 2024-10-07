package vo

import "time"

type RegisterCredentials struct {
	FamilyName  string    `json:"family_name" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	Password    string    `json:"password" validate:"required,min=8"`
	PhoneNumber string    `json:"phone_number" validate:"required"`
	Birthday    time.Time `json:"birthday" validate:"required"`
	Otp         string    `json:"otp" validate:"required"`
}

type LoginCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type VerifyEmailInput struct {
	Email string `json:"email" validate:"required,email"`
}
