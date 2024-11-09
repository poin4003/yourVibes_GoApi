package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	"github.com/poin4003/yourVibes_GoApi/pkg/socket_hub"
)

func InitSocketHub() {
	global.SocketHub = socket_hub.NewWebSocketHub()
}
