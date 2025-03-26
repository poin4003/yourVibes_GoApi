package response

import "net/http"

type CustomResponse struct {
	Code           int
	Message        string
	HttpStatusCode int
}

var codeRegistry = map[int]CustomResponse{}

func registerCode(code int, message string, httpStatusCode int) {
	codeRegistry[code] = CustomResponse{
		Code:           code,
		Message:        message,
		HttpStatusCode: httpStatusCode,
	}
}

func GetCustomCode(code int) (CustomResponse, bool) {
	customCode, exists := codeRegistry[code]
	return customCode, exists
}

const (
	ErrCodeSuccess                    = 20001
	ErrCodeParamInvalid               = 20003
	ErrInvalidToken                   = 30001
	ErrInvalidOTP                     = 30002
	ErrSendEmailOTP                   = 30003
	ErrCodeAccountBlockedByAdmin      = 30004
	ErrCodeAccountBlockedBySuperAdmin = 30005
	ErrCodeEmailOrPasswordIsWrong     = 30006
	ErrCodeInvalidLocalAuthType       = 30007
	ErrCodeInvalidGoogleAuthType      = 30008
	ErrCodeOldPasswordIsWrong         = 30009

	// Register Code
	ErrCodeUserHasExists                = 50001
	ErrCodeValidateParamRegister        = 50002
	ErrCodeValidateParamEmail           = 50003
	ErrCodeValidate                     = 50004
	ErrDataNotFound                     = 50005
	ErrFriendRequestHasAlreadyExists    = 50006
	ErrFriendHasAlreadyExists           = 50007
	ErrFriendRequestNotExists           = 50008
	ErrFriendNotExist                   = 50009
	ErrMakeFriendWithYourSelf           = 50010
	ErrAdsExpired                       = 50011
	ErrPostFriendAccess                 = 50012
	ErrPostPrivateAccess                = 50013
	ErrAdMustBePublic                   = 50014
	ErrUserFriendAccess                 = 50015
	ErrUserPrivateAccess                = 50016
	ErrCodeAdminHasExist                = 50017
	ErrCodeUserReportHasAlreadyExist    = 50018
	ErrCodePostReportHasAlreadyExist    = 50019
	ErrCodeCommentReportHasAlreadyExist = 50020
	ErrCodeReportIsAlreadyHandled       = 50021
	ErrCodeUserIsAlreadyActivated       = 50022
	ErrCodePostIsAlreadyActivated       = 50023
	ErrCodeCommentIsAlreadyActivated    = 50024
	ErrCodeGoogleAuth                   = 50025
	ErrDataHasAlreadyExist              = 50026
	ErrVoucherExpired                   = 50027
	ErrConversationAlreadyExist         = 50028

	// Err Decentralization
	ErrCodeLoginFailed        = 60001
	ErrCodeValidateParamLogin = 60002
	ErrCodeOtpNotExists       = 60009
	ErrSuperAdminPermission   = 60010

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

func InitCustomCode() {
	registerCode(ErrCodeSuccess, "Success", http.StatusOK)
	registerCode(ErrCodeParamInvalid, "Some param is invalid", http.StatusBadRequest)
	registerCode(ErrInvalidToken, "Token is invalid", http.StatusUnauthorized)
	registerCode(ErrInvalidOTP, "OTP is invalid", http.StatusBadRequest)
	registerCode(ErrSendEmailOTP, "Failed to send email otp", http.StatusInternalServerError)
	registerCode(ErrCodeAccountBlockedByAdmin, "Account blocked", http.StatusBadRequest)
	registerCode(ErrCodeAccountBlockedBySuperAdmin, "Account blocked by super admin, contact to activate", http.StatusBadRequest)
	registerCode(ErrCodeEmailOrPasswordIsWrong, "Email or password is wrong", http.StatusBadRequest)
	registerCode(ErrCodeInvalidLocalAuthType, "Invalid local auth type, you must use form to login", http.StatusBadRequest)
	registerCode(ErrCodeInvalidGoogleAuthType, "Invalid google auth type, you must use google to login", http.StatusBadRequest)
	registerCode(ErrCodeOldPasswordIsWrong, "Old password is wrong", http.StatusBadRequest)

	registerCode(ErrCodeUserHasExists, "User has already registered", http.StatusBadRequest)
	registerCode(ErrCodeValidateParamRegister, "Validate param register failed", http.StatusBadRequest)
	registerCode(ErrCodeValidateParamEmail, "Validate param email failed", http.StatusBadRequest)
	registerCode(ErrCodeValidate, "Validate param failed", http.StatusBadRequest)
	registerCode(ErrDataNotFound, "Data not found", http.StatusBadRequest)
	registerCode(ErrFriendRequestHasAlreadyExists, "Friend request has already exist", http.StatusBadRequest)
	registerCode(ErrFriendHasAlreadyExists, "Friend has already exist", http.StatusBadRequest)
	registerCode(ErrFriendRequestNotExists, "Friend request has not exist", http.StatusBadRequest)
	registerCode(ErrFriendNotExist, "Friend has not exist", http.StatusBadRequest)
	registerCode(ErrMakeFriendWithYourSelf, "You can't make friend with yourself", http.StatusBadRequest)
	registerCode(ErrAdsExpired, "Previous ads have not expired yet, you can't promote 2 advertise at a same time", http.StatusBadRequest)
	registerCode(ErrPostFriendAccess, "You must be friend to get this post", http.StatusBadRequest)
	registerCode(ErrPostPrivateAccess, "You can't get this post because it's private", http.StatusBadRequest)
	registerCode(ErrAdMustBePublic, "You must update privacy of post to PUBLIC before create advertise", http.StatusBadRequest)
	registerCode(ErrUserFriendAccess, "You must be friend to get full info", http.StatusOK)
	registerCode(ErrUserPrivateAccess, "You can't get this private info", http.StatusOK)
	registerCode(ErrCodeAdminHasExist, "Admin has already exist", http.StatusBadRequest)
	registerCode(ErrCodeUserReportHasAlreadyExist, "You already report this user!", http.StatusBadRequest)
	registerCode(ErrCodePostReportHasAlreadyExist, "You already report this post!", http.StatusBadRequest)
	registerCode(ErrCodeCommentReportHasAlreadyExist, "You already report this comment!", http.StatusBadRequest)
	registerCode(ErrCodeReportIsAlreadyHandled, "Report is already handle", http.StatusBadRequest)
	registerCode(ErrCodeUserIsAlreadyActivated, "User account is already activated", http.StatusBadRequest)
	registerCode(ErrCodePostIsAlreadyActivated, "Post is already activated", http.StatusBadRequest)
	registerCode(ErrCodeCommentIsAlreadyActivated, "Comment is already activated", http.StatusBadRequest)
	registerCode(ErrCodeGoogleAuth, "Failed to login with Google", http.StatusBadRequest)
	registerCode(ErrDataHasAlreadyExist, "Data has already exist", http.StatusBadRequest)
	registerCode(ErrVoucherExpired, "Voucher has expired", http.StatusBadRequest)
	registerCode(ErrConversationAlreadyExist, "Conversation has already exist", http.StatusBadRequest)

	registerCode(ErrCodeLoginFailed, "Account or Password is not correct", http.StatusBadRequest)
	registerCode(ErrCodeValidateParamLogin, "Validate param login", http.StatusBadRequest)
	registerCode(ErrCodeOtpNotExists, "Otp exist but not registered", http.StatusBadRequest)
	registerCode(ErrSuperAdminPermission, "You must be super admin to access this function", http.StatusForbidden)

	registerCode(ErrCreateUserFail, "Failed to create user", http.StatusBadRequest)
	registerCode(ErrHashPasswordFail, "Failed to hash password", http.StatusBadRequest)
	registerCode(ErrServerFailed, "Server failed", http.StatusInternalServerError)

	registerCode(NoUserID, "User id not found", http.StatusBadRequest)
	registerCode(UserNotFound, "User not found", http.StatusBadRequest)
	registerCode(NoKeywordInFindUsers, "No keyword to find users", http.StatusBadRequest)
	registerCode(FoundUsersFailed, "Failed to find users", http.StatusBadRequest)
}
