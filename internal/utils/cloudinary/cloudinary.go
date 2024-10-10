package cloudinary

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"os"
)

func uploadMediaToCloudinary(
	cld *cloudinary.Cloudinary,
	file *os.File,
	folder string,
) (mediaUrl string, err error) {
	params := uploader.UploadParams{
		Folder: folder,
	}

	result, err := cld.Upload.Upload(context.Background(), file, params)
	if err != nil {
		return "", fmt.Errorf("failed to upload cloudinary: %w", err)
	}

	return result.SecureURL, nil
}
