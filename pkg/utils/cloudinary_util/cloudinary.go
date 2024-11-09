package cloudinary_util

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/poin4003/yourVibes_GoApi/global"
	"mime/multipart"
	"strings"
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

func DeleteMediaFromCloudinary(mediaUrl string) error {
	if mediaUrl == "" {
		return fmt.Errorf("MediaUrl is empty")
	}

	publicID, err := extractPublicID(mediaUrl)
	if err != nil {
		return err
	}

	cloudinaryClient := global.Cloudinary
	_, err = cloudinaryClient.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})

	if err != nil {
		return fmt.Errorf("failed to delete media from cloudinary: %w", err)
	}

	return nil
}

func extractPublicID(mediaUrl string) (string, error) {
	parts := strings.Split(mediaUrl, "/")

	if len(parts) < 2 {
		return "", fmt.Errorf("Invalid URL: %s", mediaUrl)
	}

	imageName := parts[len(parts)-1]

	publicID := strings.Split(imageName, ".")[0]

	return publicID, nil
}
