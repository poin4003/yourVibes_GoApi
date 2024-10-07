package initialize

import (
	"github.com/poin4003/yourVibes_GoApi/internal/repository"
	"github.com/poin4003/yourVibes_GoApi/internal/repository/repository_implement"
	"github.com/poin4003/yourVibes_GoApi/internal/services"
	"github.com/poin4003/yourVibes_GoApi/internal/services/service_implement"
	"gorm.io/gorm"
)

func InitServiceInterface(db *gorm.DB) {
	// 1. Khởi tạo UserRepository
	userRepo := repository_implement.NewUserRepositoryImplement(db)
	repository.InitUserRepository(userRepo)

	// 2. Khởi tạo UserAuthService với UserRepository
	userAuthService := service_implement.NewUserLoginImplement(userRepo)
	services.InitUserAuth(userAuthService)
}
