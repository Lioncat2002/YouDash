package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	Title    string `gorm:"size:255;"`
	Desc     string `gorm:"size:255;"`
	PubDate  time.Time
	ThumbUrl string
	Url      string
}
