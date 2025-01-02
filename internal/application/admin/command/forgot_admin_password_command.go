package command

type ForgotAdminPasswordCommand struct {
	Email       string
	NewPassword string
}

type ForgotAdminPasswordCommandResult struct {
	ResultCode     int
	HttpStatusCode int
}
