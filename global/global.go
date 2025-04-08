package global

import (
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/logger"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/settings"
)

var (
	Config settings.Config
	Logger *logger.LoggerZap
)
