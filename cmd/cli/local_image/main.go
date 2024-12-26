package main

import (
	"github.com/gin-gonic/gin"
	"github.com/poin4003/yourVibes_GoApi/pkg/utils/media"
	"os"
	"path/filepath"
)

func uploadMedia(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"message": err,
		})
		return
	}

	fileUrl, err := media.SaveMedia(file)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err,
		})
		return
	}

	c.JSON(200, gin.H{"url": fileUrl})
}

func GetMedia(c *gin.Context) {
	fileName := c.Param("fileName")
	filePath := filepath.Join("./storages/media", fileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(400, gin.H{
			"error": "file not found",
		})
		return
	}

	c.File(filePath)
}

func main() {
	r := gin.Default()

	r.POST("/upload", uploadMedia)

	r.GET("/v1/2024/media/:fileName", GetMedia)

	r.Run(":8080")
}
