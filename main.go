package main

import (
	"famtask/controllers"
	"famtask/services"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	services.ConnectDatabase()
	services.InitPagination()
	router := gin.Default()
	//Cron job
	crn := cron.New()
	//Add the Youtube Cron JOB and set to run every 1 min
	crn.AddFunc("*/1 * * * *", services.YouTubeCronJob)
	//Start CronJob
	crn.Start()
	// for root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})
	//Video Route
	VideoRoute := router.Group("/api/video")
	VideoRoute.GET("/", controllers.GetAllVideo)
	VideoRoute.GET("/:id", controllers.GetVideoById)
	VideoRoute.POST("/", controllers.CreateMultipleVideo)
	//Run Gin Server
	router.Run(":8080")
	//Stop the cron job after the execution completes
	crn.Stop()
}
