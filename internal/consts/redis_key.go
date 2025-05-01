package consts

import "time"

type RedisKey string

const (
	RK_PERSONAL_POST RedisKey = "personal_posts"
	RK_USER_FEED     RedisKey = "userfeed"
	RK_TRENDING      RedisKey = "trending"

	TTL_COMMON time.Duration = 355 * time.Minute
)
