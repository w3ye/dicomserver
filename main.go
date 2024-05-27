package main

import (
	"context"
	"dicomserver/db"
	"dicomserver/handlers"
	"dicomserver/repositories"
	"dicomserver/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	ctx := context.Background()

	address := os.Getenv("REDIS")
	if address == "" {
		address = "localhost:6379"
	}

	client, err := db.NewRedisClient(ctx, address)
	if err != nil {
		log.Println(err)
	}

	shouldExposeRedisEndpoints := true
	if err := client.Client.Ping(ctx).Err(); err != nil {
		shouldExposeRedisEndpoints = false
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fileRouter := r.Group("file")
	{
		fileRepo := repositories.NewLocalFileRepository("store")
		fileService := service.NewFileService(fileRepo)
		fileHandler := handlers.NewFileHandler(fileService)

		// GET /file/:id
		fileRouter.GET("/:id", fileHandler.GetFileHeaders)
		// GET /file/:id/image
		fileRouter.GET("/:id/image", fileHandler.GetImage)
		// POST /file/:id/upload
		fileRouter.POST("/upload", fileHandler.UploadFile)
	}

	if shouldExposeRedisEndpoints {
		redisFileRouter := r.Group("redisFile")
		{
			redisFileRepo := repositories.NewRedisFileRepository(client)
			redisFileService := service.NewFileService(redisFileRepo)
			redisFileHandler := handlers.NewFileHandler(redisFileService)

			// POST redisFile/upload
			redisFileRouter.POST("/upload", redisFileHandler.UploadFile)
			// GET redisFile/:id
			redisFileRouter.GET("/:id", redisFileHandler.GetFileHeaders)
			// GET redisFile/:id/image
			redisFileRouter.GET("/:id/image", redisFileHandler.GetImage)
		}
	}

	r.Run()
}
