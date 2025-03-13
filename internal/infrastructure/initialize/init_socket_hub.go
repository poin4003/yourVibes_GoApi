package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/socket_hub"
)

func InitSocketHub() {
	global.SocketHub = socket_hub.NewWebSocketHub()
}
