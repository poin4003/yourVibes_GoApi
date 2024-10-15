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
	mediaRepo := repository_implement.NewMediaRepositoryImplement(db)
	settingRepo := repository_implement.NewSettingRepositoryImplement(db)

	repository.InitUserRepository(userRepo)
	repository.InitPostRepository(postRepo)
	repository.InitMediaRepository(mediaRepo)

	// 2. Initialize Service
	userAuthService := service_implement.NewUserLoginImplement(userRepo, settingRepo)
	postUserService := service_implement.NewPostUserImplement(userRepo, postRepo, mediaRepo)
	userInfoService := service_implement.NewUserInfoImplement(userRepo, settingRepo)

	services.InitUserAuth(userAuthService)
	services.InitUserInfo(userInfoService)
	services.InitPostUser(postUserService)
}
