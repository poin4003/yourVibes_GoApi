package utils

import "fmt"

func GetUserKey(hashKey string) string {
	return fmt.Sprint("u:%s:otp", hashKey)
}

func GetOtpForgotPasswordUser(hashKey string) string {
	return fmt.Sprint("cpu:%s:otp", hashKey)
}

func GetOtpForgotPasswordAdmin(hashKey string) string {
	return fmt.Sprint("cpa:%s:otp", hashKey)
}
