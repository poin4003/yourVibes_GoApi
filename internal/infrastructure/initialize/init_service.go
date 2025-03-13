package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	adminService "github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	adminServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/admin/services/implement"
	advertiseService "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	advertiseServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services/implement"

	commentService "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	commentServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services/implement"

	mediaService "github.com/poin4003/yourVibes_GoApi/internal/application/media/services"
	mediaServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/media/services/implement"
	messageService "github.com/poin4003/yourVibes_GoApi/internal/application/messages/services"
	messageServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/messages/services/implement"
	postService "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	postServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/post/services/implement"

	revenueService "github.com/poin4003/yourVibes_GoApi/internal/application/revenue/services"
	revenueServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/revenue/services/implement"

	userService "github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	userServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/user/services/implement"

	reportService "github.com/poin4003/yourVibes_GoApi/internal/application/report/services"
	reportServiceImpl "github.com/poin4003/yourVibes_GoApi/internal/application/report/services/implement"

	repository "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"

	adminRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/admin/repo_impl"
	advertiseRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/advertise/repo_impl"

	commentRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/repo_impl"
	messageRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/messages/repo_impl"
	notificationRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/notification/repo_impl"

	postRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/repo_impl"

	userRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/repo_impl"

	reportRepoImpl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/report/repo_impl"
)

func InitServiceInterface() {
	db := global.Pdb

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

	repository.InitUserRepository(userRepo)
	repository.InitPostRepository(postRepo)
	repository.InitLikeUserPostRepository(postLikeRepo)
	repository.InitMediaRepository(mediaRepo)
	repository.InitSettingRepository(settingRepo)
	repository.InitCommentRepository(commentRepo)
	repository.InitLikeUserCommentRepository(likeUserCommentRepo)
	repository.InitNotificationRepository(notificationRepo)
	repository.InitFriendRepository(friendRepo)
	repository.InitFriendRequestRepository(friendRequestRepo)
	repository.InitNewFeedRepository(newFeedRepo)
	repository.InitAdvertiseRepository(advertiseRepo)
	repository.InitBillRepository(billRepo)
	repository.InitAdminRepository(adminRepo)
	repository.InitConversationRepository(conversationRepo)
	repository.InitMessageRepository(messageRepo)
	repository.InitConversationDetailRepository(conversationDetailRepo)
	repository.InitReportRepository(reportRepo)

	// 2. Initialize Service
	userAuthServiceInit := userServiceImpl.NewUserLoginImplement(userRepo, settingRepo)
	userNotificationInit := userServiceImpl.NewUserNotificationImplement(userRepo, notificationRepo)
	userFriendServiceInit := userServiceImpl.NewUserFriendImplement(userRepo, friendRequestRepo, friendRepo, notificationRepo)
	userNewFeedServiceInit := postServiceImpl.NewPostNewFeedImplement(userRepo, postRepo, postLikeRepo, newFeedRepo)
	userInfoServiceInit := userServiceImpl.NewUserInfoImplement(userRepo, settingRepo, friendRepo, friendRequestRepo)
	postUserServiceInit := postServiceImpl.NewPostUserImplement(userRepo, friendRepo, newFeedRepo, postRepo, mediaRepo, postLikeRepo, notificationRepo, advertiseRepo)
	postLikeServiceInit := postServiceImpl.NewPostLikeImplement(userRepo, postRepo, postLikeRepo, notificationRepo)
	postShareServiceInit := postServiceImpl.NewPostShareImplement(userRepo, postRepo, mediaRepo)
	commentUserServiceInit := commentServiceImpl.NewCommentUserImplement(commentRepo, userRepo, postRepo, likeUserCommentRepo)
	likeCommentServiceInit := commentServiceImpl.NewCommentLikeImplement(userRepo, commentRepo, likeUserCommentRepo)
	advertiseServiceInit := advertiseServiceImpl.NewAdvertiseImplement(advertiseRepo, billRepo, notificationRepo)
	billServiceInit := advertiseServiceImpl.NewBillImplement(advertiseRepo, billRepo, postRepo, notificationRepo)
	adminAuthServiceInit := adminServiceImpl.NewAdminAuthImplement(adminRepo)
	adminInfoServiceInit := adminServiceImpl.NewAdminInfoImplement(adminRepo)
	superAdminServiceInit := adminServiceImpl.NewSuperAdminImplement(adminRepo)
	reportServiceInit := reportServiceImpl.NewReportFactoryImplment(reportRepo)
	revenueServiceInit := revenueServiceImpl.NewRevenueImplement(billRepo, userRepo, postRepo)
	mediaServiceInit := mediaServiceImpl.NewMediaImplement()
	conversationServiceInit := messageServiceImpl.NewConversationImplement(conversationRepo)
	messagerServiceInit := messageServiceImpl.NewMessageImplement(conversationRepo, messageRepo)
	conversationDetailServiceInit := messageServiceImpl.NewConversationDetailImplement(conversationRepo, messageRepo, conversationDetailRepo)

	userService.InitUserAuth(userAuthServiceInit)
	userService.InitUserInfo(userInfoServiceInit)
	userService.InitUserNotification(userNotificationInit)
	userService.InitUserFriend(userFriendServiceInit)
	postService.InitPostNewFeed(userNewFeedServiceInit)
	postService.InitLikeUserPost(postLikeServiceInit)
	postService.InitPostUser(postUserServiceInit)
	postService.InitPostShare(postShareServiceInit)
	commentService.InitCommentUser(commentUserServiceInit)
	commentService.InitCommentLike(likeCommentServiceInit)
	advertiseService.InitAdvertise(advertiseServiceInit)
	advertiseService.InitBill(billServiceInit)
	adminService.InitAdminAuth(adminAuthServiceInit)
	adminService.InitAdminInfo(adminInfoServiceInit)
	adminService.InitSuperAdmin(superAdminServiceInit)
	reportService.InitReport(reportServiceInit)
	revenueService.InitRevenue(revenueServiceInit)
	mediaService.InitMedia(mediaServiceInit)
	messageService.InitConversation(conversationServiceInit)
	messageService.InitMessage(messagerServiceInit)
	messageService.InitConversationDetail(conversationDetailServiceInit)

}
