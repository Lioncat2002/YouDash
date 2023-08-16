package controllers

import (
	"famtask/models"
	"famtask/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoData struct {
	Title    string `gorm:"size:255;"`
	Desc     string `gorm:"size:255;"`
	PubDate  time.Time
	ThumbUrl string
	Url      string
}

func GetAllVideo(c *gin.Context) {
	var videos []models.Video
	if err := services.DB.Find(&videos).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   videos,
	})
}

func CreateVideo(c *gin.Context) {
	var videoData VideoData
	if err := c.ShouldBindJSON(&videoData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	video := models.Video{
		Title:    videoData.Title,
		Desc:     videoData.Desc,
		PubDate:  videoData.PubDate,
		ThumbUrl: videoData.ThumbUrl,
		Url:      videoData.Url,
	}

	if err := services.DB.Create(&video).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
}
