package consts

import "time"

type RedisKey string

const (
	RK_PERSONAL_POST RedisKey = "personnal_posts"
	RK_USER_FEED     RedisKey = "userfeed"

	TTL_COMMON time.Duration = 355 * time.Minute
)
