package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin"
	super_admin_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth"
	auth_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/dto/request"
)

type AdminRouter struct{}

func (ar *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	AdminAuthController := admin_auth.NewAdminAuthController()
	SuperAdminController := admin_super_admin.NewSuperAdminController()

	// Public router
	adminRouterPublic := Router.Group("/admins")
	{
		// admin auth
		adminRouterPublic.POST("/login",
			helpers.ValidateJsonBody(&auth_request.AdminLoginRequest{}, auth_request.ValidateLoginRequest),
			AdminAuthController.Login,
		)

		// super admin
		adminRouterPublic.POST("/",
			helpers.ValidateJsonBody(&super_admin_request.CreateAdminRequest{}, super_admin_request.ValidateCreateAdminRequest),
			SuperAdminController.CreateAdmin,
		)
	}

	// Private router
}
