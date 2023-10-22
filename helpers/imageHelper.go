package helpers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const uploadDirectory = "assets/"

func UploadImage(c *gin.Context, formKey string, prefix string) (string, error) {
	file, err := c.FormFile(formKey)
	if err != nil {
		return "", err
	}

	if !isImage(file.Filename) {
		return "", fmt.Errorf("file is not a valid image")
	}

	randomStr := generateRandomString(10)
	currentTime := time.Now().Format("20060102150405")
	// filePath := filepath.Join(uploadDirectory, prefix+"_"+file.Filename)
	// filePath := filepath.Join(uploadDirectory, prefix+"_"+randomStr+"_"+file.Filename)

	fileName := fmt.Sprintf("%s_%s", currentTime, randomStr)
	fileExtension := filepath.Ext(file.Filename)
	filePath := filepath.Join(uploadDirectory, fileName+fileExtension)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return filePath, nil
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp":
		return true
	default:
		return false
	}
}
