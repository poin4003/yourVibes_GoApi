package consts

type PrivacyLevel string

const (
	PUBLIC      PrivacyLevel = "public"
	FRIEND_ONLY PrivacyLevel = "friend_only"
	PRIVATE     PrivacyLevel = "private"
)

var PrivacyLevels = []interface{}{
	PUBLIC,
	FRIEND_ONLY,
	PRIVATE,
}
