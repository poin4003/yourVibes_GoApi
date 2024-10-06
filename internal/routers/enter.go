package routers

import (
	"github.com/poin4003/yourVibes_GoApi/internal/routers/admin"
	"github.com/poin4003/yourVibes_GoApi/internal/routers/user"
)

type RouterGroup struct {
	User  user.UserRouterGroup
	Admin admin.AdminRouterGroup
}

var RouterGroupApp = new(RouterGroup)
