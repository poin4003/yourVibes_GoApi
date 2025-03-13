package user

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/helpers"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	reportUser "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_user"
	reportRequest "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_user/dto/request"
)

type ReportRouter struct{}

func (rr *ReportRouter) InitReportRouter(Router *gin.RouterGroup) {
	// Public router
	reportUserController := reportUser.NewReportController()

	// Private router
	reportRouterPrivate := Router.Group("/report")
	reportRouterPrivate.Use(middlewares.UserAuthProtected())
	{
		reportRouterPrivate.POST("/",
			helpers.ValidateJsonBody(&reportRequest.ReportRequest{}, reportRequest.ValidateReportRequest),
			reportUserController.Report,
		)
	}
}
