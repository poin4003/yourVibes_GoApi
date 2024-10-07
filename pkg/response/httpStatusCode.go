package response

const (
	ErrCodeSuccess      = 20001 // Success
	ErrCodeParamInvalid = 20003 // Email is invalid
	ErrInvalidToken     = 30001 // Token is invalid
	ErrInvalidOTP       = 30002 // OTP is invalid
	ErrSendEmailOTP     = 30003

	// Register Code
	ErrCodeUserHasExists = 50001 // user has already registered

	// Err Login
	ErrCodeOtpNotExists = 60009

	// Err server failed
	ErrCreateUserFail   = 70001
	ErrHashPasswordFail = 70002
)

var msg = map[int]string{
	ErrCodeSuccess:      "Success",
	ErrCodeParamInvalid: "Email is invalid",
	ErrInvalidToken:     "Token is invalid",
	ErrInvalidOTP:       "OTP is invalid",
	ErrSendEmailOTP:     "Failed to send email otp",

	ErrCodeUserHasExists: "User has already registered",

	ErrCodeOtpNotExists: "Otp exist but not registered",
	ErrCreateUserFail:   "Failed to create user",
	ErrHashPasswordFail: "Failed to hash password",
}
