package initialize

import (
	adminServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/admin/services/implement"
	advertiseServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services/implement"
	commentProducer "github.com/poin4003/yourVibes_GoApi/internal/application/comment/producer"
	messageConsumer "github.com/poin4003/yourVibes_GoApi/internal/application/messages/consumer"
	messageProducer "github.com/poin4003/yourVibes_GoApi/internal/application/messages/producer"
	"github.com/poin4003/yourVibes_GoApi/internal/application/post/consumer"
	postProducer "github.com/poin4003/yourVibes_GoApi/internal/application/post/producer"
	reportProducer "github.com/poin4003/yourVibes_GoApi/internal/application/report/producer"
	statisticConsumer "github.com/poin4003/yourVibes_GoApi/internal/application/statistic/consumer"
	userProducer "github.com/poin4003/yourVibes_GoApi/internal/application/user/producer"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/grpc/comment_pb"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/rabbitmq"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/socket_hub"
	adminCacheImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/transient/admin"
	commentCacheImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/transient/comment"
	postCacheImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/transient/post"
	userCacheImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/transient/user"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/middlewares"
	"github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/routers"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"

	commentServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services/implement"

	mediaServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/media/services/implement"
	messageServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/messages/services/implement"
	postServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/post/services/implement"

	revenueServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/revenue/services/implement"

	userServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/user/services/implement"

	reportServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/report/services/implement"

	statisticServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/statistic/services/implement"

	notificationConsumer "github.com/poin4003/yourVibes_GoApi/internal/application/notification/consumer"
	notificationServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/notification/services/implement"

	adminRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/admin/repo_impl"
	advertiseRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/advertise/repo_impl"

	commentRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/repo_impl"
	messageRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/messages/repo_impl"
	notificationRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/notification/repo_impl"
	StatisticRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/statistic/repo_impl"

	postRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/repo_impl"

	voucherRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/voucher/repo_impl"

	userRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/repo_impl"

	reportRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/report/repo_impl"

	userAdvertiseControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_user/controller/impl"
	userAuthControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/user_auth/controller/impl"
	userCommentControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/comment/comment_user/controller/impl"
	userMessageControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/messages/message_user/controller/impl"
	userNotificationControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/notification/notification_user/controller/impl"
	userPostControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/post/post_user/controller/impl"
	userReportControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_user/controller/impl"
	userControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/user/user_user/controller/impl"

	mediaControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/media/controller/impl"

	adminControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_admin/controller/impl"
	superAdminControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/admin/admin_super_admin/controller/impl"
	adminAdvertiseControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/advertise/advertise_admin/controller/impl"
	adminAuthControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/auth/admin_auth/controller/impl"
	adminReportControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/report/report_admin/controller/impl"
	adminRevenueControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/revenue/revenue_admin/controller/impl"
	adminSystemControllerImpl "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/rest/system/system_admin/controller/impl"

	adminRouter "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/routers/admin"
	userRouter "github.com/poin4003/yourVibes_GoApi/internal/interfaces/api/routers/user"

	advertiseCronjob "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/cronjob"
	postCronjob "github.com/poin4003/yourVibes_GoApi/internal/application/post/cronjob"
)

func InitDependencyInjection(
	db *gorm.DB,
	rabbitmqConnection *rabbitmq.Connection,
	redis *redis.Client,
	notificationSocketHub *socket_hub.NotificationSocketHub,
	messageSocketHub *socket_hub.MessageSocketHub,
	commentCensorGrpcConn *grpc.ClientConn,
) *routers.RouterGroup {
	// 1. Initialize Repository
	userRepo := userRepoImpl.NewUserRepositoryImplement(db)
	postRepo := postRepoImpl.NewPostRepositoryImplement(db)
	postLikeRepo := postRepoImpl.NewLikeUserPostRepositoryImplement(db)
	mediaRepo := postRepoImpl.NewMediaRepositoryImplement(db)
	settingRepo := userRepoImpl.NewSettingRepositoryImplement(db)
	commentRepo := commentRepoImpl.NewCommentRepositoryImplement(db)
	likeUserCommentRepo := commentRepoImpl.NewLikeUserCommentRepositoryImplement(db)
	notificationRepo := notificationRepoImpl.NewNotificationRepositoryImplement(db)
	friendRepo := userRepoImpl.NewFriendImplement(db)
	friendRequestRepo := userRepoImpl.NewFriendRequestImplement(db)
	newFeedRepo := postRepoImpl.NewNewFeedRepositoryImplement(db)
	advertiseRepo := advertiseRepoImpl.NewAdvertiseRepositoryImplement(db)
	billRepo := advertiseRepoImpl.NewBillRepositoryImplement(db)
	adminRepo := adminRepoImpl.NewAdminRepositoryImplement(db)
	conversationRepo := messageRepoImpl.NewConversationRepositoryImplement(db)
	messageRepo := messageRepoImpl.NewMessageRepositoryImplement(db)
	conversationDetailRepo := messageRepoImpl.NewConversationDetailRepositoryImplement(db)
	reportRepo := reportRepoImpl.NewReportRepositoryImplement(db)
	voucherRepo := voucherRepoImpl.NewVoucherRepositoryImplement(db)
	statisticRepo := StatisticRepoImpl.NewStatisticRepository(db)

	// Init cache
	userAuthCache := userCacheImpl.NewUserAuthCache(redis)
	userCache := userCacheImpl.NewUserCache(redis)
	postCache := postCacheImpl.NewPostCacheImplement(redis)
	commentCache := commentCacheImpl.NewCommentCacheImplement(redis)
	adminCache := adminCacheImpl.NewAdminCache(redis)

	// Init publisher
	postEventPublisher := postProducer.NewPostEventPublisher(rabbitmqConnection)
	userNotificationPublisher := userProducer.NewNotificationPublisher(rabbitmqConnection)
	reportNotificationPublisher := reportProducer.NewNotificationPublisher(rabbitmqConnection)
	commentNotificationPublisher := commentProducer.NewNotificationPublisher(rabbitmqConnection)
	messagePublisher := messageProducer.NewMessagePublisher(rabbitmqConnection)

	// Init grpc
	commentCensorGrpcClient := comment_pb.NewCommentCensorServiceClient(commentCensorGrpcConn)

	// Initialize Service
	userAuthServiceInit := userServiceImpl.NewUserLoginImplement(userRepo, settingRepo, newFeedRepo, userAuthCache)
	userFriendServiceInit := userServiceImpl.NewUserFriendImplement(userRepo, friendRequestRepo, friendRepo, userCache, userNotificationPublisher)
	userNewFeedServiceInit := postServiceImpl.NewPostNewFeedImplement(userRepo, postRepo, postLikeRepo, newFeedRepo, postCache, postEventPublisher)
	userInfoServiceInit := userServiceImpl.NewUserInfoImplement(userRepo, settingRepo, friendRepo, friendRequestRepo, userCache)
	postUserServiceInit := postServiceImpl.NewPostUserImplement(userRepo, friendRepo, newFeedRepo, postRepo, mediaRepo, postLikeRepo, advertiseRepo, postCache, commentCache, postEventPublisher)
	postLikeServiceInit := postServiceImpl.NewPostLikeImplement(userRepo, postRepo, postLikeRepo, postCache, postEventPublisher)
	postShareServiceInit := postServiceImpl.NewPostShareImplement(userRepo, postRepo, mediaRepo, newFeedRepo, friendRepo, postCache, postEventPublisher)
	commentUserServiceInit := commentServiceImpl.NewCommentUserImplement(commentRepo, userRepo, postRepo, likeUserCommentRepo, commentCache, postCache, commentNotificationPublisher, commentCensorGrpcClient)
	likeCommentServiceInit := commentServiceImpl.NewCommentLikeImplement(userRepo, commentRepo, likeUserCommentRepo, commentCache, commentNotificationPublisher)
	advertiseServiceInit := advertiseServiceImpl.NewAdvertiseImplement(advertiseRepo, billRepo, voucherRepo, postCache)
	billServiceInit := advertiseServiceImpl.NewBillImplement(advertiseRepo, billRepo, postRepo, notificationRepo)
	adminAuthServiceInit := adminServiceImpl.NewAdminAuthImplement(adminRepo)
	adminInfoServiceInit := adminServiceImpl.NewAdminInfoImplement(adminRepo, adminCache)
	superAdminServiceInit := adminServiceImpl.NewSuperAdminImplement(adminRepo, adminCache)
	reportServiceInit := reportServiceImpl.NewReportFactoryImplment(reportRepo, voucherRepo, friendRepo, userCache, postCache, commentCache, reportNotificationPublisher)
	revenueServiceInit := revenueServiceImpl.NewRevenueImplement(billRepo, userRepo, postRepo)
	mediaServiceInit := mediaServiceImpl.NewMediaImplement()
	conversationServiceInit := messageServiceImpl.NewConversationImplement(conversationRepo, userCache)
	messageServiceInit := messageServiceImpl.NewMessageImplement(messageRepo, messagePublisher)
	messageMQServiceInit := messageServiceImpl.NewMessageMQImplement(conversationDetailRepo, messageSocketHub)
	conversationDetailServiceInit := messageServiceImpl.NewConversationDetailImplement(conversationRepo, messageRepo, conversationDetailRepo, userRepo, messagePublisher)
	notificationServiceInit := notificationServiceImpl.NewNotification(notificationRepo, notificationSocketHub)
	notificationUserInit := notificationServiceImpl.NewNotificationUserImplement(userRepo, notificationRepo)
	statisticServiceInit := statisticServiceImpl.NewStatisticImplement(statisticRepo)

	// Init Middleware
	userAuthProtectMiddleware := middlewares.NewUserAuthProtectedMiddleware(userInfoServiceInit)
	adminAuthProtectMiddleware := middlewares.NewAdminAuthProtectedMiddleware(adminInfoServiceInit)

	// Init controller
	userAuthControllerInit := userAuthControllerImpl.NewUserAuthController(userAuthServiceInit)
	userControllerInit := userControllerImpl.NewUserInfoController(userInfoServiceInit)
	userFriendControllerInit := userControllerImpl.NewUserFriendController(userFriendServiceInit)
	userReportControllerInit := userReportControllerImpl.NewReportController(reportServiceInit)
	userPostControllerInit := userPostControllerImpl.NewPostUserController(postUserServiceInit)
	userPostLikeControllerInit := userPostControllerImpl.NewPostLikeController(postLikeServiceInit)
	userPostShareControllerInit := userPostControllerImpl.NewPostShareController(postShareServiceInit)
	userPostNewFeedControllerInit := userPostControllerImpl.NewPostNewFeedController(userNewFeedServiceInit)
	userNotificationControllerInit := userNotificationControllerImpl.NewNotificationController(notificationUserInit, userInfoServiceInit, notificationSocketHub)
	userMessageControllerInit := userMessageControllerImpl.NewMessageController(messageServiceInit, messageSocketHub)
	userConversationControllerInit := userMessageControllerImpl.NewConversationController(conversationServiceInit)
	userConversationDetailControllerInit := userMessageControllerImpl.NewConversationDetailController(conversationDetailServiceInit)
	userCommentControllerInit := userCommentControllerImpl.NewCommentUserController(commentUserServiceInit)
	userCommentLikeControllerInit := userCommentControllerImpl.NewCommentLikeController(likeCommentServiceInit)
	userAdvertiseControllerInit := userAdvertiseControllerImpl.NewAdvertiseController(advertiseServiceInit, postUserServiceInit)
	userBillControllerInit := userAdvertiseControllerImpl.NewBillController(billServiceInit)

	mediaControllerInit := mediaControllerImpl.NewMediaController(mediaServiceInit)

	adminRevenueControllerInit := adminRevenueControllerImpl.NewRevenueAdminController(revenueServiceInit)
	adminReportControllerInit := adminReportControllerImpl.NewAdminReportController(reportServiceInit)
	adminAuthControllerInit := adminAuthControllerImpl.NewAdminAuthController(adminAuthServiceInit)
	adminControllerInit := adminControllerImpl.NewAdminController(adminInfoServiceInit)
	superAdminControllerInit := superAdminControllerImpl.NewSuperAdminController(superAdminServiceInit, adminAuthServiceInit)
	adminAdvertiseControllerInit := adminAdvertiseControllerImpl.NewAdvertiseAdminController(advertiseServiceInit)
	adminSystemCacheControllerInit := adminSystemControllerImpl.NewSystemAdminCacheController(userInfoServiceInit, postUserServiceInit, commentUserServiceInit)
	adminSystemPostControllerInit := adminSystemControllerImpl.NewSystemAdminPostController(userNewFeedServiceInit)

	// Init router
	userRouterInit := userRouter.NewUserRouter(userControllerInit, userFriendControllerInit, userAuthControllerInit, userAuthProtectMiddleware)
	userReportRouterInit := userRouter.NewReportRouter(userReportControllerInit, userAuthProtectMiddleware)
	userPostRouterInit := userRouter.NewPostRouter(userPostControllerInit, userPostLikeControllerInit, userPostShareControllerInit, userPostNewFeedControllerInit, userAuthProtectMiddleware)
	userNotificationRouterInit := userRouter.NewNotificationRouter(userNotificationControllerInit, userAuthProtectMiddleware)
	userMessageRouterInit := userRouter.NewMessagesRouter(userConversationControllerInit, userConversationDetailControllerInit, userMessageControllerInit, userAuthProtectMiddleware)
	mediaRouterInit := userRouter.NewMediaRouter(mediaControllerInit)
	userCommentRouterInit := userRouter.NewCommentRouter(userCommentControllerInit, userCommentLikeControllerInit, userAuthProtectMiddleware)
	userAdvertiseRouterInit := userRouter.NewAdvertiseRouter(userAdvertiseControllerInit, userBillControllerInit, userAuthProtectMiddleware)

	adminRevenueRouterInit := adminRouter.NewRevenueAdminRouter(adminRevenueControllerInit, adminAuthProtectMiddleware)
	adminReportRouterInit := adminRouter.NewAdminReportRouter(adminReportControllerInit, adminAuthProtectMiddleware)
	adminRouterInit := adminRouter.NewAdminRouter(adminAuthControllerInit, adminControllerInit, superAdminControllerInit, adminAuthProtectMiddleware)
	adminAdvertiesRouterInit := adminRouter.NewAdvertiseAdminRouter(adminAdvertiseControllerInit, adminAuthProtectMiddleware)
	adminSystemRouterInit := adminRouter.NewSystemRouter(adminSystemCacheControllerInit, adminSystemPostControllerInit, adminAuthProtectMiddleware)

	// Init router group
	userRouterGroup := userRouter.NewUserRouterGroup(
		*userRouterInit,
		*userPostRouterInit,
		*userCommentRouterInit,
		*userAdvertiseRouterInit,
		*mediaRouterInit,
		*userMessageRouterInit,
		*userReportRouterInit,
		*userNotificationRouterInit,
	)

	adminRouterGroup := adminRouter.NewAdminRouterGroup(
		*adminRouterInit,
		*adminAdvertiesRouterInit,
		*adminRevenueRouterInit,
		*adminReportRouterInit,
		*adminSystemRouterInit,
	)

	routerGroup := routers.NewRouterGroup(*userRouterGroup, *adminRouterGroup)

	// Init broker consumer
	notificationConsumer.InitNotificationConsumer(notificationServiceInit, rabbitmqConnection)
	messageConsumer.InitMessageConsumer(messageMQServiceInit, rabbitmqConnection)
	statisticConsumer.InitStatisticsConsumer(statisticServiceInit, rabbitmqConnection)
	consumer.InitPostModerationConsumer(postUserServiceInit, rabbitmqConnection)

	// Init cronjob
	advertiseCronjob.NewCheckExpiryCronJob(postRepo, newFeedRepo, postCache)
	advertiseCronjob.NewPushToNewFeedCronJob(newFeedRepo, postCache)
	postCronjob.NewPushFeaturePostToNewFeedCronJob(newFeedRepo, postCache)
	postCronjob.NewCheckExpiryFeaturePostCronJob(newFeedRepo, postCache)

	return routerGroup
}
