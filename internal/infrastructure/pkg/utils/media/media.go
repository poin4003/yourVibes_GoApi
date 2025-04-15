package media

import (
	"errors"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/IP"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

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

func AddUrlIntoMedia(fileName string) string {
	fileUrl := fmt.Sprintf("%s:%d/v1/2024/media/%s",
		global.Config.Server.ServerEndpoint,
		global.Config.Server.Port,
		fileName)

	return fileUrl
}

func GetUrlMedia() string {
	ipv4, err := IP.GetLocalIP()
	if err != nil {
		ipv4 = global.Config.Server.ServerEndpoint
	}

	fileUrl := fmt.Sprintf("http://%s:%d/v1/2024/media/",
		ipv4,
		global.Config.Server.Port,
	)

	return fileUrl
}

func saveMediaWithoutUrl(fileHeader *multipart.FileHeader) (string, error) {
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

	// 7. return filename
	return uniqueFileName, nil
}

func SaveManyMedia(media []multipart.FileHeader) ([]string, error) {
	if len(media) <= 0 {
		return nil, nil
	}

	var wg sync.WaitGroup
	urlChan := make(chan string, len(media))
	errChan := make(chan error, len(media))

	for _, file := range media {
		wg.Add(1)
		go func(file multipart.FileHeader) {
			defer wg.Done()
			mediaUrl, err := saveMediaWithoutUrl(&file)
			if err != nil {
				errChan <- response.NewServerFailedError(err.Error())
				return
			}
			urlChan <- mediaUrl
		}(file)
	}

	go func() {
		wg.Wait()
		close(urlChan)
		close(errChan)
	}()

	for err := range errChan {
		return nil, err
	}

	var mediaUrls []string
	for url := range urlChan {
		mediaUrls = append(mediaUrls, url)
	}

	return mediaUrls, nil
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

	time.Sleep(100 * time.Millisecond)

	// 4. Delete file
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func DeleteMediaByFilename(fileName string) error {
	// 1. Get media path
	filePath := filepath.Join(global.Config.Media.Folder, fileName)

	// 2. Check file exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	// 3. Delete file
	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func ParseRange(rangeHeader string, fileSize int64) (start, end int64, err error) {
	const prefix = "bytes="
	if !strings.HasPrefix(rangeHeader, prefix) {
		return 0, 0, errors.New("invalid range format")
	}

	rangeStr := strings.TrimPrefix(rangeHeader, prefix)
	rangeParts := strings.Split(rangeStr, "-")
	if len(rangeParts) != 2 {
		return 0, 0, errors.New("invalid range format")
	}

	start, err = strconv.ParseInt(rangeParts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	if rangeParts[1] == "" {
		end = fileSize - 1
	} else {
		end, err = strconv.ParseInt(rangeParts[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}

	if start < 0 || end < start || end >= fileSize {
		return 0, 0, errors.New("range out of bounds")
	}

	return start, end, nil
}
