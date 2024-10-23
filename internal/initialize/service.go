package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/internal/repository/repository_implement"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/internal/services/service_implement"
	"gorm.io/gorm"
)

func InitServiceInterface(db *gorm.DB) {
	// 1. Initialize Repository
	userRepo := repository_implement.NewUserRepositoryImplement(db)
	postRepo := repository_implement.NewPostRepositoryImplement(db)
	postLikeRepo := repository_implement.NewLikeUserPostRepositoryImplement(db)
	mediaRepo := repository_implement.NewMediaRepositoryImplement(db)
	settingRepo := repository_implement.NewSettingRepositoryImplement(db)
	commentRepo := repository_implement.NewCommentRepositoryImplement(db)

	repository.InitUserRepository(userRepo)
	repository.InitPostRepository(postRepo)
	repository.InitLikeUserPostRepository(postLikeRepo)
	repository.InitMediaRepository(mediaRepo)
	repository.InitSettingRepository(settingRepo)
	repository.InitCommentRepository(commentRepo)

	// 2. Initialize Service
	userAuthService := service_implement.NewUserLoginImplement(userRepo, settingRepo)
	postUserService := service_implement.NewPostUserImplement(userRepo, postRepo, mediaRepo)
	postLikeService := service_implement.NewPostLikeImplement(userRepo, postRepo, postLikeRepo)
	postShareService := service_implement.NewPostShareImplement(userRepo, postRepo, mediaRepo)
	userInfoService := service_implement.NewUserInfoImplement(userRepo, settingRepo)
	commentUserService := service_implement.NewCommentUserImplement(commentRepo, userRepo, postRepo)

	services.InitUserAuth(userAuthService)
	services.InitUserInfo(userInfoService)
	services.InitLikeUserPost(postLikeService)
	services.InitPostUser(postUserService)
	services.InitPostShare(postShareService)
	services.InitCommentUser(commentUserService)
}
