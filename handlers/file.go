package handlers

import (
	"dicomserver/service"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type fileService interface {
	WriteFile(fileHeader *multipart.FileHeader) (string, error)
	GetDicomHeaders(id string, query string) (*service.GetDicomHeaderAttributeResponse, error)
	SetImageEncoder(fileType string)
	GetDicomImage(id string) ([]byte, error)
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
			"error":   err.Error(),
			"message": "tag format should be in the format of (XXXX,XXXX)",
		})
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (f FileHandler) GetFileHeaders(c *gin.Context) {
	fileId := c.Param("id")
	if fileId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file id is required",
		})
	}
	tagQuery := c.Query("tag")
	if tagQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "tag is required",
		})
	}
	dicomHeader, err := f.service.GetDicomHeaders(fileId, tagQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, dicomHeader)
}

func (f FileHandler) GetImage(c *gin.Context) {
	fileId := c.Param("id")
	if fileId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file id is required",
		})
	}
	supportedFileTypes := map[string]bool{
		"png": true,
		"jpg": true,
	}
	fileType := c.Query("fileType")
	if _, ok := supportedFileTypes[fileType]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "file type is required and should be one of png, jpg",
		})
	}
	// set the content type based on the file type
	var contentType string
	switch fileType {
	case "png":
		contentType = "image/png"
	case "jpg":
		contentType = "image/jpeg"
	}

	// set the fileType encoder
	f.service.SetImageEncoder(fileType)

	image, err := f.service.GetDicomImage(fileId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.Data(http.StatusOK, contentType, image)
}
