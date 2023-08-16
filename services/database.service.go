package services

import (
	"log"
	"os"
	"template/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}
	dburi := os.Getenv("DB_URI") //used a cockroachdb database but postgres is fine
	DB, err = gorm.Open(postgres.Open(dburi), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	} else {
		log.Println("connected to database")
	}
	//DB.Statement.Delete(&models.User{}, &models.Post{})
	//DB.Exec("DELETE FROM items")
	//DB.Exec("DELETE FROM users")
	//DB.Delete(&[]models.User{}, &[]models.Post{}).Where("1=1")
	DB.AutoMigrate(&models.Video{})
}
