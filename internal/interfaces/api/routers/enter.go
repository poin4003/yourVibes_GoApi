package routers

import (
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/routers/admin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/routers/user"
)

type RouterGroup struct {
	User  user.UserRouterGroup
	Admin admin.AdminRouterGroup
}

var RouterGroupApp = new(RouterGroup)
