package command

import "github.com/poin4003/yourVibes_GoApi/internal/application/post/common"

type LoginCommand struct {
	Email    string
	Password string
}

type LoginCommandResult struct {
	Admin       *common.AdminResult
	AccessToken string
}