package initialize

import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/poin4003/yourVibes_GoApi/global"
)

func InitCloudinary() {
	cloudinaryClient, err := cloudinary.NewFromParams(
		global.Config.CloudinarySetting.CloudName,
		global.Config.CloudinarySetting.ApiKey,
		global.Config.CloudinarySetting.ApiSecretKey,
	)

	if err != nil {
		panic("Failed to initialize Cloudinary client: " + err.Error())
	}

	global.Cloudinary = cloudinaryClient
}
