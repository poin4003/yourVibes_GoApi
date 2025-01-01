package global

import (
	"github.com/poin4003/yourVibes_GoApi/pkg/logger"
	"github.com/poin4003/yourVibes_GoApi/pkg/settings"
	"github.com/poin4003/yourVibes_GoApi/pkg/socket_hub"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config    settings.Config
	Logger    *logger.LoggerZap
	Rdb       *redis.Client
	Pdb       *gorm.DB
	SocketHub *socket_hub.WebSocketHub
)
