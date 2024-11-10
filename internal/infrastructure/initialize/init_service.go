package initialize

import (
	comment_service "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services"
	comment_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/comment/services/implement"
	post_service "github.com/poin4003/yourVibes_GoApi/internal/application/post/services"
	post_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/post/services/implement"
	user_service "github.com/poin4003/yourVibes_GoApi/internal/application/user/services"
	user_service_impl "github.com/poin4003/yourVibes_GoApi/internal/application/user/services/implement"
	comment_repo "github.com/poin4003/yourVibes_GoApi/internal/domain/repositories"
	repo_impl2 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/comment/repo_impl"
	repo_impl3 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/post/repo_impl"
	repo_impl4 "github.com/poin4003/yourVibes_GoApi/internal/infrastructure/persistence/user/repo_impl"
	"gorm.io/gorm"
)

func InitServiceInterface(db *gorm.DB) {
	// 1. Initialize Repository
	userRepo := repo_impl4.NewUserRepositoryImplement(db)
	postRepo := repo_impl3.NewPostRepositoryImplement(db)
	postLikeRepo := repo_impl3.NewLikeUserPostRepositoryImplement(db)
	mediaRepo := repo_impl3.NewMediaRepositoryImplement(db)
	settingRepo := repo_impl4.NewSettingRepositoryImplement(db)
	commentRepo := repo_impl2.NewCommentRepositoryImplement(db)
	likeUserCommentRepo := repo_impl2.NewLikeUserCommentRepositoryImplement(db)
	notificationRepo := repo_impl4.NewNotificationRepositoryImplement(db)
	friendRepo := repo_impl4.NewFriendImplement(db)
	friendRequestRepo := repo_impl4.NewFriendRequestImplement(db)
	newFeedRepo := repo_impl3.NewNewFeedRepositoryImplement(db)

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
}
