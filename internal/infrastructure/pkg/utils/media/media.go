package media

import (
	"bytes"
	"fmt"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/response"
	"github.com/poin4003/yourVibes_GoApi/internal/infrastructure/pkg/utils/IP"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/poin4003/yourVibes_GoApi/global"

	"github.com/google/uuid"
)

func SaveMedia(fileHeader *multipart.FileHeader) (string, error) {
	// 1. Open file from file header
	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// 2. Create unique file name
	rawFileName := fmt.Sprintf("raw_%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))
	processedFileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))

	mediaFolder := global.Config.Media.Folder

	err = os.MkdirAll(mediaFolder, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create media directory: %w", err)
	}

	rawPath := filepath.Join(mediaFolder, rawFileName)
	finalPath := filepath.Join(mediaFolder, processedFileName)

	// 3. Save raw file temporarily
	rawOut, err := os.Create(rawPath)
	if err != nil {
		return "", fmt.Errorf("failed to create raw file: %w", err)
	}
	defer rawOut.Close()

	_, err = io.Copy(rawOut, srcFile)
	if err != nil {
		return "", fmt.Errorf("failed to save raw file: %w", err)
	}

	// 4. Run ffmpeg to move moov atom (faststart)
	cmd := exec.Command("ffmpeg", "-i", rawPath, "-movflags", "+faststart", "-c", "copy", finalPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v - %s", err, stderr.String())
	}

	// 5. (Optional) Remove raw file
	os.Remove(rawPath)

	// 6. Return URL
	fileUrl := fmt.Sprintf("%s:%d/v1/2024/media/%s",
		global.Config.Server.ServerEndpoint,
		global.Config.Server.Port,
		processedFileName)

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
	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	rawFileName := fmt.Sprintf("raw_%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))
	finalFileName := fmt.Sprintf("%s%s", uuid.New().String(), filepath.Ext(fileHeader.Filename))

	mediaFolder := global.Config.Media.Folder

	err = os.MkdirAll(mediaFolder, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create media directory: %w", err)
	}

	rawPath := filepath.Join(mediaFolder, rawFileName)
	finalPath := filepath.Join(mediaFolder, finalFileName)

	rawOut, err := os.Create(rawPath)
	if err != nil {
		return "", fmt.Errorf("failed to create raw file: %w", err)
	}
	defer rawOut.Close()

	_, err = io.Copy(rawOut, srcFile)
	if err != nil {
		return "", fmt.Errorf("failed to save raw file: %w", err)
	}

	cmd := exec.Command("ffmpeg", "-i", rawPath, "-movflags", "+faststart", "-c", "copy", finalPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ffmpeg error: %v - %s", err, stderr.String())
	}

	_ = os.Remove(rawPath)

	return finalFileName, nil
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
