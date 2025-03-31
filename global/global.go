package global

import (
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/logger"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/settings"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/socket_hub"
)

var (
	Config                settings.Config
	Logger                *logger.LoggerZap
	NotificationSocketHub *socket_hub.NotificationSocketHub
	MessageSocketHub      *socket_hub.MessageSocketHub
)
