package consts

type AuthType string

const (
	LOCAL_AUTH  AuthType = "local"
	GOOGLE_AUTH AuthType = "google"
)

var AuthTypes = []interface{}{
	LOCAL_AUTH,
	GOOGLE_AUTH,
}
