package consts

// Define enum for privacy
type PrivacyLevel string

const (
	PUBLIC      PrivacyLevel = "public"
	FRIEND_ONLY PrivacyLevel = "friend_only"
	PRIVATE     PrivacyLevel = "private"
)
