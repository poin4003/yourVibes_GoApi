package consts

// Define enum for validator
type PrivacyLevel string

const (
	PUBLIC      PrivacyLevel = "public"
	FRIEND_ONLY PrivacyLevel = "friend_only"
	PRIVATE     PrivacyLevel = "private"
)
