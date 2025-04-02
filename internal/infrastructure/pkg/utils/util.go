package utils

import "fmt"

func GetUserKey(hashKey string) string {
	return fmt.Sprintf("u:%s:otp", hashKey)
}

func GetOtpForgotPasswordUser(hashKey string) string {
	return fmt.Sprintf("cpu:%s:otp", hashKey)
}
