package media

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/poin4003/yourVibes_GoApi/global"

	"github.com/google/uuid"
)

func SaveMedia(fileHeader *multipart.FileHeader) (string, error) {
	// 1. Open file from file header
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 2. Create uuid for file name
	uniqueFileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))

	// 3. Define path to save file
	mediaFolder := global.Config.Media.Folder

	// 4. Ensure the directory exists
	err = os.MkdirAll(mediaFolder, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create media directory: %w", err)
	}

	filePath := filepath.Join(mediaFolder, uniqueFileName)

	// 5. Create file in directory
	outFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	// 6. Copy data into file
	_, err = outFile.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("failed to write data to file: %w", err)
	}

	// 7. Generate the file URL
	fileUrl := fmt.Sprintf("%s:%d/v1/2024/media/%s",
		global.Config.Server.ServerEndpoint,
		global.Config.Server.Port,
		uniqueFileName)

	return fileUrl, nil
}

func GetMedia(fileName string) (string, error) {
	// 1. Get path to file
	filePath := filepath.Join(global.Config.Media.Folder, fileName)

	// 2. Check file exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return filePath, nil
}

func DeleteMedia(mediaLink string) error {
	// 1. Get file name from link
	parts := strings.Split(mediaLink, "/")
	fileName := parts[len(parts)-1]

	// 2. Get media path
	filePath := filepath.Join(global.Config.Media.Folder, fileName)

	// 3. Check file exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	// 4. Delete file
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}
