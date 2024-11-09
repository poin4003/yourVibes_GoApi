package request

type VerifyEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}
