package request

import user_command "github.com/poin4003/yourVibes_GoApi/internal/application/user/command"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (req *LoginRequest) ToLoginCommand() (*user_command.LoginCommand, error) {
	return &user_command.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	}, nil
}
