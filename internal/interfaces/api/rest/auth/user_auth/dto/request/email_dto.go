package request

type VerifyEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}
