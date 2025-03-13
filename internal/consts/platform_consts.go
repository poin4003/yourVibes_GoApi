package consts

type Platform string

const (
	WEB     Platform = "web"
	ANDROID Platform = "android"
	IOS     Platform = "ios"
)

var Platforms = []interface{}{
	WEB,
	ANDROID,
	IOS,
}
