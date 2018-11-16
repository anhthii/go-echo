package main

import (
	"log"
	"os"

	"github.com/anhthii/go-echo/db"
	"github.com/anhthii/go-echo/db/models"
	"github.com/anhthii/go-echo/router"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func main() {
	// Init database
	db.Init()
	db.Tables(&models.User{}, &models.Playlist{}, &models.Song{}, &models.Artist{})
	defer db.Close()

	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}

	MODE := os.Getenv("SERVER_MODE")
	router := router.SetupRouter(MODE)

	PORT := os.Getenv("SERVER_PORT")
	router.Run(":" + PORT)
}
