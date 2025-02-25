package handler

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type File struct {
	File multipart.FileHeader `form:"file" binding:"required"`
}

// uploadFile
// @Summary uploadFile
// @Description Upload a media file
// @Tags media
// @Accept multipart/form-data
// @Param file formData file true "UploadMediaForm"
// @Security BearerAuth
// @Success 201 {object} string
// @Failure 400 {object} error
// @Failure 500 {object} error
// @Router /minio/media [post]
func (h *Handler) Media(c *gin.Context) {
	// Faylni olish
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid file"})
		return
	}

	fileUrl := filepath.Join("./media", fileHeader.Filename)

	// Faylni saqlash
	err = c.SaveUploadedFile(fileHeader, fileUrl)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to save file"})
		return
	}

	// Faylning MIME turini aniqlash
	fileExt := filepath.Ext(fileHeader.Filename)
	contentType := "application/octet-stream"
	if fileExt == ".png" {
		contentType = "image/png"
	} else if fileExt == ".jpg" || fileExt == ".jpeg" {
		contentType = "image/jpeg"
	}

	// Yangi fayl nomini yaratish
	newFile := uuid.NewString() + fileExt

	// Faylni MinIO'ga yuklash
	info, err := h.MinIO.FPutObject(context.Background(), "photos", newFile, fileUrl, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": "Failed to upload file"})
		return
	}

	fmt.Println("Info Bucket:", info.Bucket)

	err = os.Remove(fileUrl)
	if err != nil {
		fmt.Println("Failed to delete local file:", err)
	}
	// URL yaratish
	madeUrl := fmt.Sprintf("https://aura.dilshodforever.uz/photos/%s", newFile)
	fmt.Println("url: ", madeUrl)
	c.JSON(201, gin.H{"url": madeUrl})
}

