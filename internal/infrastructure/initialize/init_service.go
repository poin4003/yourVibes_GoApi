package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	admin_service "github.com/poin4003/yourVibes_GoApi/internal/application/admin/services"
	admin_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/admin/services/implement"
	advertise_service "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	advertise_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services/implement"
	comment_service "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	comment_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services/implement"
	media_service "github.com/poin4003/yourVibes_GoApi/internal/application/media/services"
	media_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/media/services/implement"
	post_service "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	post_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/post/services/implement"
	revenue_service "github.com/poin4003/yourVibes_GoApi/internal/application/revenue/services"
	revenue_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/revenue/services/implement"
	user_service "github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	user_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/user/services/implement"
	repository "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	admin_repo_impl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/admin/repo_impl"
	advertise_repo_impl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/advertise/repo_impl"
	comment_repo_impl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/repo_impl"
	notification_repo_impl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/notification/repo_impl"
	post_repo_impl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/repo_impl"
	user_repo_impl "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/repo_impl"
)

func InitServiceInterface() {
	db := global.Pdb

	// 1. Initialize Repository
	userRepo := user_repo_impl.NewUserRepositoryImplement(db)
	userReportRepo := user_repo_impl.NewUserReportRepositoryImplement(db)
	postRepo := post_repo_impl.NewPostRepositoryImplement(db)
	postReportRepo := post_repo_impl.NewPostReportRepositoryImplement(db)
	postLikeRepo := post_repo_impl.NewLikeUserPostRepositoryImplement(db)
	mediaRepo := post_repo_impl.NewMediaRepositoryImplement(db)
	settingRepo := user_repo_impl.NewSettingRepositoryImplement(db)
	commentRepo := comment_repo_impl.NewCommentRepositoryImplement(db)
	commentReportRepo := comment_repo_impl.NewCommentReportRepositoryImplement(db)
	likeUserCommentRepo := comment_repo_impl.NewLikeUserCommentRepositoryImplement(db)
	notificationRepo := notification_repo_impl.NewNotificationRepositoryImplement(db)
	friendRepo := user_repo_impl.NewFriendImplement(db)
	friendRequestRepo := user_repo_impl.NewFriendRequestImplement(db)
	newFeedRepo := post_repo_impl.NewNewFeedRepositoryImplement(db)
	advertiseRepo := advertise_repo_impl.NewAdvertiseRepositoryImplement(db)
	billRepo := advertise_repo_impl.NewBillRepositoryImplement(db)
	adminRepo := admin_repo_impl.NewAdminRepositoryImplement(db)

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
	repository.InitUserReportRepository(userReportRepo)
	repository.InitPostReportRepository(postReportRepo)
	repository.InitCommentReportRepository(commentReportRepo)

	// 2. Initialize Service
	userAuthService := user_service_impl.NewUserLoginImplement(userRepo, settingRepo)
	userNotification := user_service_impl.NewUserNotificationImplement(userRepo, notificationRepo)
	userFriendService := user_service_impl.NewUserFriendImplement(userRepo, friendRequestRepo, friendRepo, notificationRepo)
	userNewFeedService := post_service_impl.NewPostNewFeedImplement(userRepo, postRepo, postLikeRepo, newFeedRepo)
	userInfoService := user_service_impl.NewUserInfoImplement(userRepo, settingRepo, friendRepo, friendRequestRepo)
	postUserService := post_service_impl.NewPostUserImplement(userRepo, friendRepo, newFeedRepo, postRepo, mediaRepo, postLikeRepo, notificationRepo, advertiseRepo)
	postLikeService := post_service_impl.NewPostLikeImplement(userRepo, postRepo, postLikeRepo, notificationRepo)
	postShareService := post_service_impl.NewPostShareImplement(userRepo, postRepo, mediaRepo)
	commentUserService := comment_service_impl.NewCommentUserImplement(commentRepo, userRepo, postRepo, likeUserCommentRepo)
	likeCommentService := comment_service_impl.NewCommentLikeImplement(userRepo, commentRepo, likeUserCommentRepo)
	advertiseService := advertise_service_impl.NewAdvertiseImplement(advertiseRepo, billRepo, notificationRepo)
	billService := advertise_service_impl.NewBillImplement(advertiseRepo, billRepo, postRepo, notificationRepo)
	adminAuthService := admin_service_impl.NewAdminAuthImplement(adminRepo)
	adminInfoService := admin_service_impl.NewAdminInfoImplement(adminRepo)
	superAdminService := admin_service_impl.NewSuperAdminImplement(adminRepo)
	userReportService := user_service_impl.NewUserReportImplement(userReportRepo, userRepo, postRepo, commentRepo)
	postReportService := post_service_impl.NewPostReportImplement(postReportRepo, postRepo)
	commentReportService := comment_service_impl.NewCommentReportImplement(commentReportRepo, commentRepo)
	revenueService := revenue_service_impl.NewRevenueImplement(billRepo, userRepo, postRepo)
	mediaService := media_service_impl.NewMediaImplement()

	user_service.InitUserAuth(userAuthService)
	user_service.InitUserInfo(userInfoService)
	user_service.InitUserNotification(userNotification)
	user_service.InitUserFriend(userFriendService)
	post_service.InitPostNewFeed(userNewFeedService)
	post_service.InitLikeUserPost(postLikeService)
	post_service.InitPostUser(postUserService)
	post_service.InitPostShare(postShareService)
	comment_service.InitCommentUser(commentUserService)
	comment_service.InitCommentLike(likeCommentService)
	advertise_service.InitAdvertise(advertiseService)
	advertise_service.InitBill(billService)
	admin_service.InitAdminAuth(adminAuthService)
	admin_service.InitAdminInfo(adminInfoService)
	admin_service.InitSuperAdmin(superAdminService)
	user_service.InitUserReport(userReportService)
	post_service.InitPostReport(postReportService)
	comment_service.InitCommentReport(commentReportService)
	revenue_service.InitRevenue(revenueService)
	media_service.InitMedia(mediaService)
}
