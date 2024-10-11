package cloudinary_util

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/poin4003/yourVibes_GoApi/global"
	"mime/multipart"
)

func UploadMediaToCloudinary(
	file multipart.File,
) (string, error) {
	if file == nil {
		return "", fmt.Errorf("file is nil")
	}

	cloudinaryClient := global.Cloudinary

	uploadParams, err := cloudinaryClient.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder: global.Config.CloudinarySetting.Folder,
	})

	if err != nil {
		return "", err
	}

	return uploadParams.SecureURL, nil
}
