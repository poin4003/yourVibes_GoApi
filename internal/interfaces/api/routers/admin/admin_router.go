package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin"
	adminRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/dto/request"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin"
	superAdminRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/request"
	superAdminQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/query"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth"
	authRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/dto/request"
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
			helpers.ValidateJsonBody(&authRequest.AdminLoginRequest{}, authRequest.ValidateLoginRequest),
			AdminAuthController.Login,
		)
	}

	// Private router
	adminRouterPrivate := Router.Group("/admins")
	//adminRouterPrivate.Use(middlewares.AdminAuthProtected())
	{
		// admin auth
		adminRouterPrivate.PATCH("/change_password",
			helpers.ValidateJsonBody(&authRequest.ChangeAdminPasswordRequest{}, authRequest.ValidateChangePasswordRequest),
			AdminAuthController.ChangeAdminPassword,
		)

		// admin info
		adminRouterPrivate.PATCH("/",
			helpers.ValidateJsonBody(&adminRequest.UpdateAdminInfoRequest{}, adminRequest.ValidateUpdateAdminInfoRequest),
			AdminController.UpdateAdminInfo,
		)

		// super admin
		adminRouterPrivate.POST("/super_admin",
			//middlewares.CheckSuperAdminRole(),
			helpers.ValidateJsonBody(&superAdminRequest.CreateAdminRequest{}, superAdminRequest.ValidateCreateAdminRequest),
			SuperAdminController.CreateAdmin,
		)

		// Change admin password
		adminRouterPrivate.POST("/super_admin/forgot_admin_password",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateJsonBody(&superAdminRequest.ForgotAdminPasswordRequest{}, superAdminRequest.ValidateForgotAdminPasswordRequest),
			SuperAdminController.ForgotAdminPassword,
		)

		adminRouterPrivate.GET("/:admin_id",
			middlewares.CheckSuperAdminRole(),
			SuperAdminController.GetAdminById,
		)

		adminRouterPrivate.GET("/",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateQuery(&superAdminQuery.AdminQueryObject{}, superAdminQuery.ValidateAdminQueryObject),
			SuperAdminController.GetManyAdmins,
		)

		adminRouterPrivate.PATCH("/super_admin",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateJsonBody(&superAdminRequest.UpdateAdminForSuperAdminRequest{}, superAdminRequest.ValidateUpdateAdminForSuperAdminRequest),
			SuperAdminController.UpdateAdmin,
		)
	}
}
