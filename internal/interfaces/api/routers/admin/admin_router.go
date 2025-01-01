package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin"
	admin_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin"
	super_admin_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/request"
	super_admin_query "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/query"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth"
	auth_request "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/dto/request"
)

type AdminRouter struct{}

func (ar *AdminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	AdminAuthController := admin_auth.NewAdminAuthController()
	SuperAdminController := admin_super_admin.NewSuperAdminController()
	AdminController := admin_admin.NewAdminController()

	// Public router
	adminRouterPublic := Router.Group("/admins")
	{
		// admin auth
		adminRouterPublic.POST("/login",
			helpers.ValidateJsonBody(&auth_request.AdminLoginRequest{}, auth_request.ValidateLoginRequest),
			AdminAuthController.Login,
		)
	}

	// Private router
	adminRouterPrivate := Router.Group("/admins")
	adminRouterPrivate.Use(middlewares.AdminAuthProtected())
	{
		// admin auth
		adminRouterPrivate.PATCH("/change_password",
			helpers.ValidateJsonBody(&auth_request.ChangeAdminPasswordRequest{}, auth_request.ValidateChangePasswordRequest),
			AdminAuthController.ChangeAdminPassword,
		)

		// admin info
		adminRouterPrivate.PATCH("/",
			helpers.ValidateJsonBody(&admin_request.UpdateAdminInfoRequest{}, admin_request.ValidateUpdateAdminInfoRequest),
			AdminController.UpdateAdminInfo,
		)

		// super admin
		adminRouterPrivate.POST("/super_admin",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateJsonBody(&super_admin_request.CreateAdminRequest{}, super_admin_request.ValidateCreateAdminRequest),
			SuperAdminController.CreateAdmin,
		)

		adminRouterPrivate.GET("/:admin_id",
			middlewares.CheckSuperAdminRole(),
			SuperAdminController.GetAdminById,
		)

		adminRouterPrivate.GET("/",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateQuery(&super_admin_query.AdminQueryObject{}, super_admin_query.ValidateAdminQueryObject),
			SuperAdminController.GetManyAdmins,
		)

		adminRouterPrivate.PATCH("/super_admin",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateJsonBody(&super_admin_request.UpdateAdminForSuperAdminRequest{}, super_admin_request.ValidateUpdateAdminForSuperAdminRequest),
			SuperAdminController.UpdateAdmin,
		)
	}
}
