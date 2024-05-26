package handlers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type fileService interface {
	WriteFile(fileHeader *multipart.FileHeader) (string, error)
}

type FileHandler struct {
	service fileService
}

func NewFileHandler(fs fileService) *FileHandler {
	return &FileHandler{
		service: fs,
	}
}

func (f FileHandler) UploadFile(c *gin.Context) {
	// read file content from request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	id, err := f.service.WriteFile(fileHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (f FileHandler) GetFileHeaders(c *gin.Context) {
}

func (f FileHandler) GetImage(c *gin.Context) {}
