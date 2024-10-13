package auth_dto

type VerifyEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}
