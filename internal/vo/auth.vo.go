package vo

import "time"

type RegisterCredentials struct {
	FamilyName  string    `json:"family_name" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required,min=8"`
	PhoneNumber string    `json:"phone_number" binding:"required"`
	Birthday    time.Time `json:"birthday" binding:"required"`
	Otp         string    `json:"otp" binding:"required"`
}

type LoginCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type VerifyEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}
