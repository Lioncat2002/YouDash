package controllers

import (
	"famtask/models"
	"famtask/services"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// the video data for the POST requests
type VideoData struct {
	Title    string    `json:"title" binding:"required"`
	Desc     string    `json:"desc" binding:"required"`
	PubDate  time.Time `json:"pub_date" binding:"required"`
	ThumbUrl string    `json:"thumb_url" binding:"required"`
	VideoId  string    `json:"video_id" binding:"required"`
}

/**
* Gets all the video data stored in the db
 */
func GetAllVideo(c *gin.Context) {
	var videos []models.Video
	if err := services.DB.Find(&videos).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	model := services.DB.Model(&models.Video{})
	c.JSON(http.StatusOK, services.PG.With(model).Request(c.Request).Response(&[]models.Video{}))
}

/**
* Gets the Video based on the video id
 */
func GetVideoById(c *gin.Context) {
	id := c.Param("id")
	log.Println("id", id)
	video := models.Video{}
	if err := services.DB.Where("video_id = ?", id).First(&video).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"status": "success",
		"data":   video,
	})
}

/**
* Performs multiple inserts into the db
 */
func CreateMultipleVideo(c *gin.Context) {
	var videoDatas []VideoData
	if err := c.ShouldBindJSON(&videoDatas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	videos := []models.Video{}
	for _, m := range videoDatas {
		video := models.Video{
			Title:    m.Title,
			Desc:     m.Desc,
			PubDate:  m.PubDate,
			ThumbUrl: m.ThumbUrl,
			Url:      "https://youtube.com/watch?v=" + m.VideoId,
			VideoId:  m.VideoId,
		}
		videos = append(videos, video)
	}

	if err := services.DB.Create(&videos).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
}

/**
* Performs a single insert into the db
 */
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
		Url:      "https://youtube.com/watch?v=" + videoData.VideoId, //we just get the videoID from youtube data api
		VideoId:  videoData.VideoId,
	}

	if err := services.DB.Create(&video).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}
}
