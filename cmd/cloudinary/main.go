package main

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	r := gin.Default()

	var cld, err = cloudinary.NewFromParams("dkf51e57t", "238396968984682", "WRS5sVfgYzYwHeFDG0IyQZE9ZSE")

	if err != nil {
		log.Fatal(err)
	}

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Received file: %s, Size: %bytes\n", file.Filename, file.Size)

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		log.Printf("file:::", src)

		log.Printf("File metadata:\nFilename: %s\nSize: %d bytes\n", file.Filename, file.Size)

		params := uploader.UploadParams{
			Folder: "yourVibes",
		}

		result, err := cld.Upload.Upload(context.Background(), src, params)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": result.SecureURL})
	})

	r.Run(":8080")
}
