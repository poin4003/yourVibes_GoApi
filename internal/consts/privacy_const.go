package consts

// Define enum for validator
type PrivacyLevel string

const (
	PUBLIC      PrivacyLevel = "public"
	FRIEND_ONLY PrivacyLevel = "friend_only"
	PRIVATE     PrivacyLevel = "private"
)

func IsValidPrivacyLevel(level PrivacyLevel) bool {
	switch level {
	case PUBLIC, FRIEND_ONLY, PRIVATE:
		return true
	default:
		return false
	}
}
