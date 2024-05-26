package main

import (
	"dicomserver/handlers"
	"dicomserver/service"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fileRouter := r.Group("file")
	{
		fileService := service.NewFileService()
		fileHandler := handlers.NewFileHandler(fileService)

		// GET /file/:id
		fileRouter.GET("/:id", fileHandler.GetFileHeaders)
		// GET /file/:id/image
		fileRouter.GET("/:id/image", fileHandler.GetImage)
		// POST /file/:id/upload
		fileRouter.POST("/upload", fileHandler.UploadFile)
	}

	r.Run()
}
