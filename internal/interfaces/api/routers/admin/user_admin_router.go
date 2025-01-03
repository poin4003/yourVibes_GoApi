package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_admin"
	adminUserReportQuery "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_admin/query"
)

type UserAdminRouter struct{}

func (uar *UserAdminRouter) InitUserAdminRouter(Router *gin.RouterGroup) {
	adminUserReportController := user_admin.NewAdminUserReportController()

	// Private router
	adminRouterPrivate := Router.Group("/users")
	adminRouterPrivate.Use(middlewares.AdminAuthProtected())
	{
		// user report
		adminRouterPrivate.GET("/report/:user_id/:reported_user_id",
			adminUserReportController.GetUserReport,
		)

		adminRouterPrivate.GET("/report",
			helpers.ValidateQuery(&adminUserReportQuery.UserReportQueryObject{}, adminUserReportQuery.ValidateUserReportQueryObject),
			adminUserReportController.GetManyUserReports,
		)

		adminRouterPrivate.PATCH("/report/:user_id/:reported_user_id",
			adminUserReportController.HandleUserReport,
		)

		adminRouterPrivate.DELETE("/report/:user_id/:reported_user_id",
			adminUserReportController.DeleteUserReport,
		)

		adminRouterPrivate.PATCH("/report/activate/:user_id",
			adminUserReportController.ActivateUserAccount,
		)
	}
}
