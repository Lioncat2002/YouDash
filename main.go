package main

import (
	"famtask/controllers"
	"famtask/services"

	"github.com/gin-gonic/gin"
)

func main() {
	services.ConnectDatabase()
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	VideoRoute := router.Group("/api/video")
	VideoRoute.GET("/", controllers.GetAllVideo)
	router.Run(":8080")
}
