package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	fileRouter := r.Group("file")
	{
		// GET /file/:id
		fileRouter.GET("/:id")
		// GET /file/:id/image
		fileRouter.GET("/:id/image")
		// POST /file/:id/upload
		fileRouter.POST("/:id/upload")
	}

	r.Run()
}
