package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/global"
	advertise_service "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services"
	advertise_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/advertise/services/implement"
	comment_service "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	comment_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services/implement"
	post_service "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	post_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/post/services/implement"
	user_service "github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	user_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/user/services/implement"
	advertise_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	comment_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
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
	postRepo := post_repo_impl.NewPostRepositoryImplement(db)
	postLikeRepo := post_repo_impl.NewLikeUserPostRepositoryImplement(db)
	mediaRepo := post_repo_impl.NewMediaRepositoryImplement(db)
	settingRepo := user_repo_impl.NewSettingRepositoryImplement(db)
	commentRepo := comment_repo_impl.NewCommentRepositoryImplement(db)
	likeUserCommentRepo := comment_repo_impl.NewLikeUserCommentRepositoryImplement(db)
	notificationRepo := notification_repo_impl.NewNotificationRepositoryImplement(db)
	friendRepo := user_repo_impl.NewFriendImplement(db)
	friendRequestRepo := user_repo_impl.NewFriendRequestImplement(db)
	newFeedRepo := post_repo_impl.NewNewFeedRepositoryImplement(db)
	advertiseRepo := advertise_repo_impl.NewAdvertiseRepositoryImplement(db)
	billRepo := advertise_repo_impl.NewBillRepositoryImplement(db)

	comment_repo.InitUserRepository(userRepo)
	comment_repo.InitPostRepository(postRepo)
	comment_repo.InitLikeUserPostRepository(postLikeRepo)
	comment_repo.InitMediaRepository(mediaRepo)
	comment_repo.InitSettingRepository(settingRepo)
	comment_repo.InitCommentRepository(commentRepo)
	comment_repo.InitLikeUserCommentRepository(likeUserCommentRepo)
	comment_repo.InitNotificationRepository(notificationRepo)
	comment_repo.InitFriendRepository(friendRepo)
	comment_repo.InitFriendRequestRepository(friendRequestRepo)
	comment_repo.InitNewFeedRepository(newFeedRepo)
	advertise_repo.InitAdvertiseRepository(advertiseRepo)
	advertise_repo.InitBillRepository(billRepo)

	// 2. Initialize Service
	userAuthService := user_service_impl.NewUserLoginImplement(userRepo, settingRepo)
	userNotification := user_service_impl.NewUserNotificationImplement(userRepo, notificationRepo)
	userFriendService := user_service_impl.NewUserFriendImplement(userRepo, friendRequestRepo, friendRepo, notificationRepo)
	userNewFeedService := post_service_impl.NewPostNewFeedImplement(userRepo, postRepo, postLikeRepo, newFeedRepo)
	userInfoService := user_service_impl.NewUserInfoImplement(userRepo, settingRepo, friendRepo, friendRequestRepo)
	postUserService := post_service_impl.NewPostUserImplement(userRepo, friendRepo, newFeedRepo, postRepo, mediaRepo, postLikeRepo, notificationRepo)
	postLikeService := post_service_impl.NewPostLikeImplement(userRepo, postRepo, postLikeRepo, notificationRepo)
	postShareService := post_service_impl.NewPostShareImplement(userRepo, postRepo, mediaRepo)
	commentUserService := comment_service_impl.NewCommentUserImplement(commentRepo, userRepo, postRepo, likeUserCommentRepo)
	likeCommentService := comment_service_impl.NewCommentLikeImplement(userRepo, commentRepo, likeUserCommentRepo)
	advertiseService := advertise_service_impl.NewAdvertiseImplement(advertiseRepo, billRepo, notificationRepo)
	billSerivce := advertise_service_impl.NewBillImplement(advertiseRepo, billRepo, postRepo, notificationRepo)

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
	advertise_service.InitBill(billSerivce)
}
