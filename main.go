package main

import (
	"famtask/controllers"
	"famtask/services"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	services.ConnectDatabase()
	router := gin.Default()
	crn := cron.New()
	crn.AddFunc("*/15 * * * *", services.YouTubeCronJob)
	crn.Start()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	VideoRoute := router.Group("/api/video")
	VideoRoute.GET("/", controllers.GetAllVideo)
	VideoRoute.POST("/", controllers.CreateMultipleVideo)
	router.Run(":8080")
	crn.Stop()
}
