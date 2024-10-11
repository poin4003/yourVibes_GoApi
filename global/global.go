package global

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/poin4003/yourVibes_GoApi/pkg/logger"
	"github.com/poin4003/yourVibes_GoApi/pkg/settings"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	Config     settings.Config
	Logger     *logger.LoggerZap
	Rdb        *redis.Client
	Pdb        *gorm.DB
	Cloudinary *cloudinary.Cloudinary
)
