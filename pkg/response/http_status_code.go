package response

const (
	ErrCodeSuccess      = 20001 // Success
	ErrCodeParamInvalid = 20003 // Email is invalid
	ErrInvalidToken     = 30001 // Token is invalid
	ErrInvalidOTP       = 30002 // OTP is invalid
	ErrSendEmailOTP     = 30003

	// Register Code
	ErrCodeUserHasExists             = 50001 // user has already registered
	ErrCodeValidateParamRegister     = 50002
	ErrCodeValidateParamEmail        = 50003
	ErrCodeValidate                  = 50004
	ErrDataNotFound                  = 50005
	ErrFriendRequestHasAlreadyExists = 50006
	ErrFriendHasAlreadyExists        = 50007
	ErrFriendRequestNotExists        = 50008
	ErrFriendNotExist                = 50009
	ErrMakeFriendWithYourSelf        = 50010
	ErrAdsExpired                    = 50011
	ErrPostFriendAccess              = 50012
	ErrPostPrivateAccess             = 50013
	ErrAdMustBePublic                = 50014
	ErrUserFriendAccess              = 50015
	ErrUserPrivateAccess             = 50016

	// Err Login
	ErrCodeLoginFailed        = 60001
	ErrCodeValidateParamLogin = 60002
	ErrCodeOtpNotExists       = 60009

	// Err server failed
	ErrCreateUserFail   = 70001
	ErrHashPasswordFail = 70002
	ErrServerFailed     = 70003

	// Users Code
	NoUserID             = 80001
	UserNotFound         = 80002
	NoKeywordInFindUsers = 80003
	FoundUsersFailed     = 80004
)

var msg = map[int]string{
	ErrCodeSuccess:      "Success",
	ErrCodeParamInvalid: "Email is invalid",
	ErrInvalidToken:     "Token is invalid",
	ErrInvalidOTP:       "OTP is invalid",
	ErrSendEmailOTP:     "Failed to send email otp",

	ErrCodeUserHasExists:             "User has already registered",
	ErrCodeValidateParamRegister:     "Validate param register failed",
	ErrCodeValidateParamEmail:        "Validate param email failed",
	ErrCodeValidate:                  "Validate param failed",
	ErrDataNotFound:                  "Data not found",
	ErrFriendRequestHasAlreadyExists: "Friend request has already exist",
	ErrFriendHasAlreadyExists:        "Friend has already exist",
	ErrFriendRequestNotExists:        "Friend request has not exist",
	ErrFriendNotExist:                "Friend has not exist",
	ErrMakeFriendWithYourSelf:        "You can't make friend with yourself",
	ErrAdsExpired:                    "Previous ads have not expired yet, you can't promote 2 advertise at a same time",
	ErrPostFriendAccess:              "You must be friend to get this post",
	ErrPostPrivateAccess:             "You can't get this post because it's private",
	ErrAdMustBePublic:                "You must update privacy of post to PUBLIC before create advertise",
	ErrUserFriendAccess:              "You must be friend to get full info",
	ErrUserPrivateAccess:             "You can't get this private info",

	ErrCodeLoginFailed:        "Account or Password is not correct",
	ErrCodeValidateParamLogin: "Validate param login",
	ErrCodeOtpNotExists:       "Otp exist but not registered",
	ErrCreateUserFail:         "Failed to create user",
	ErrHashPasswordFail:       "Failed to hash password",
	ErrServerFailed:           "Server failed",

	NoUserID:             "User id not found",
	UserNotFound:         "User not found",
	NoKeywordInFindUsers: "No keyword to find users",
	FoundUsersFailed:     "Failed to find users",
}
