package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	cAdminInfo "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/controller"
	adminRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/dto/request"
	cSuperAdmin "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/controller"
	superAdminRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/dto/request"
	superAdminQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/query"
	cAdminAuth "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/controller"
	authRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/dto/request"
)

type adminRouter struct {
	adminAuthController      cAdminAuth.IAdminAuthController
	adminController          cAdminInfo.IAdminInfoController
	superAdminController     cSuperAdmin.ISuperAdminController
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware
}

func NewAdminRouter(
	adminAuthController cAdminAuth.IAdminAuthController,
	adminController cAdminInfo.IAdminInfoController,
	superAdminController cSuperAdmin.ISuperAdminController,
	adminProtectedMiddleware middlewares.IAdminAuthProtectedMiddleware,
) *adminRouter {
	return &adminRouter{
		adminAuthController:      adminAuthController,
		adminController:          adminController,
		superAdminController:     superAdminController,
		adminProtectedMiddleware: adminProtectedMiddleware,
	}
}

func (r *adminRouter) InitAdminRouter(Router *gin.RouterGroup) {
	// Public router
	adminRouterPublic := Router.Group("/admins")
	{
		// admin auth
		adminRouterPublic.POST("/login",
			helpers.ValidateJsonBody(&authRequest.AdminLoginRequest{}, authRequest.ValidateLoginRequest),
			r.adminAuthController.Login,
		)
	}

	// Private router
	adminRouterPrivate := Router.Group("/admins")
	adminRouterPrivate.Use(r.adminProtectedMiddleware.AdminAuthProtected())
	{
		// admin auth
		adminRouterPrivate.PATCH("/change_password",
			helpers.ValidateJsonBody(&authRequest.ChangeAdminPasswordRequest{}, authRequest.ValidateChangePasswordRequest),
			r.adminAuthController.ChangeAdminPassword,
		)

		// admin info
		adminRouterPrivate.PATCH("/",
			helpers.ValidateJsonBody(&adminRequest.UpdateAdminInfoRequest{}, adminRequest.ValidateUpdateAdminInfoRequest),
			r.adminController.UpdateAdminInfo,
		)

		// super admin
		adminRouterPrivate.POST("/super_admin",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateJsonBody(&superAdminRequest.CreateAdminRequest{}, superAdminRequest.ValidateCreateAdminRequest),
			r.superAdminController.CreateAdmin,
		)

		// Change admin password
		adminRouterPrivate.POST("/super_admin/forgot_admin_password",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateJsonBody(&superAdminRequest.ForgotAdminPasswordRequest{}, superAdminRequest.ValidateForgotAdminPasswordRequest),
			r.superAdminController.ForgotAdminPassword,
		)

		adminRouterPrivate.GET("/:admin_id",
			middlewares.CheckSuperAdminRole(),
			r.superAdminController.GetAdminById,
		)

		adminRouterPrivate.GET("/",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateQuery(&superAdminQuery.AdminQueryObject{}, superAdminQuery.ValidateAdminQueryObject),
			r.superAdminController.GetManyAdmins,
		)

		adminRouterPrivate.PATCH("/super_admin",
			middlewares.CheckSuperAdminRole(),
			helpers.ValidateJsonBody(&superAdminRequest.UpdateAdminForSuperAdminRequest{}, superAdminRequest.ValidateUpdateAdminForSuperAdminRequest),
			r.superAdminController.UpdateAdmin,
		)
	}
}
