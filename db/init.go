package db

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var _db *gorm.DB

// Init initialize postgres database instance
func Init(dotenvPath ...string) {
	var err error
	if len(dotenvPath) == 1 {
		err = godotenv.Load(dotenvPath[0])
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Panic("Error loading .env file")
	}

	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")

	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)
	fmt.Println("connecting to database: " + dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		log.Panic(err.Error())
	}

	_db = conn
}

func Tables(tables ...interface{}) {
	if _db == nil {
		log.Panic("you have to Init database first")
	}
	_db.AutoMigrate(tables...)
}

func Close() func() error {
	return _db.Close
}

//returns a handle to the DB object
func GetDB() *gorm.DB {
	return _db
}
