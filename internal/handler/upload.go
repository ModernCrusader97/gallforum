package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "파일을 선택하세요"})
		return
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "이미지 파일만 업로드 가능합니다"})
		return
	}
	if file.Size > 10<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "파일 크기는 10MB 이하여야 합니다"})
		return
	}
	os.MkdirAll("uploads", 0755)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join("uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"url": "/uploads/" + filename})
}
