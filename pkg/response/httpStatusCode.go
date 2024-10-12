package response

const (
	ErrCodeSuccess      = 20001 // Success
	ErrCodeParamInvalid = 20003 // Email is invalid
	ErrInvalidToken     = 30001 // Token is invalid
	ErrInvalidOTP       = 30002 // OTP is invalid
	ErrSendEmailOTP     = 30003

	// Register Code
	ErrCodeUserHasExists         = 50001 // user has already registered
	ErrCodeValidateParamRegister = 50002
	ErrCodeValidateParamEmail    = 50003
	ErrCodeValidate              = 50004

	// Err Login
	ErrCodeLoginFailed        = 60001
	ErrCodeValidateParamLogin = 60002
	ErrCodeOtpNotExists       = 60009

	// Err server failed
	ErrCreateUserFail   = 70001
	ErrHashPasswordFail = 70002
	ErrServerFailed     = 70003

	// Users Code
	NoUserID = 80001
	UserNotFound = 80002
	NoKeywordInFindUsers = 80003
	FoundUsersFailed = 80004
)

var msg = map[int]string{
	ErrCodeSuccess:      "Success",
	ErrCodeParamInvalid: "Email is invalid",
	ErrInvalidToken:     "Token is invalid",
	ErrInvalidOTP:       "OTP is invalid",
	ErrSendEmailOTP:     "Failed to send email otp",

	ErrCodeUserHasExists:         "User has already registered",
	ErrCodeValidateParamRegister: "Validate param register failed",
	ErrCodeValidateParamEmail:    "Validate param email failed",
	ErrCodeValidate:              "Validate param failed",

	ErrCodeLoginFailed:        "Account or Password is not correct",
	ErrCodeValidateParamLogin: "Validate param login",
	ErrCodeOtpNotExists:       "Otp exist but not registered",
	ErrCreateUserFail:         "Failed to create user",
	ErrHashPasswordFail:       "Failed to hash password",
	ErrServerFailed:           "Server failed",

	NoUserID: "User id not found",
	UserNotFound: "User not found",
	NoKeywordInFindUsers: "No keyword to find users",
	FoundUsersFailed: "Failed to find users",
}
