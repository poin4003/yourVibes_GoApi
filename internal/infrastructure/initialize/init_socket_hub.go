package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/socket_hub"
)

func InitSocketHub() {
	global.NotificationSocketHub = socket_hub.NewNotificationSocketHub()
	global.MessageSocketHub = socket_hub.NewMessageSocketHub()
	global.Logger.Info("init socket hub success")
}
