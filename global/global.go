package global

import (
	"database/sql"
	"github.com/poin4003/yourVibes_GoApi/pkg/logger"
	"github.com/poin4003/yourVibes_GoApi/pkg/settings"
	"github.com/redis/go-redis/v9"
)

var (
	Config settings.Config
	Logger *logger.LoggerZap
	Rdb    *redis.Client
	Pdb    *sql.DB
)
