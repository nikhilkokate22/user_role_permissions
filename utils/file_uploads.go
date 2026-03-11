package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadFile(c *gin.Context, fieldName, traceID, folder string) (string, error) {

	file, err := c.FormFile(fieldName)
	if err != nil {
		return "", nil // optional file
	}

	if file.Size > 5*1024*1024 {
		return "", errors.New("file too large")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))

	allowed := map[string]bool{
		".pdf":  true,
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowed[ext] {
		return "", errors.New("invalid file type")
	}

	dir := "uploads/" + traceID + "/" + folder
	os.MkdirAll(dir, os.ModePerm)

	fileName := uuid.New().String() + ext
	fullPath := filepath.Join(dir, fileName)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}

	return fullPath, nil
}
