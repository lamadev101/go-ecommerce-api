package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lamadev101/ecommerce-api/database"
	"github.com/lamadev101/ecommerce-api/router"
)

func init() {
	// Load environment variables
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Loading the config from .env file")
		err = godotenv.Load(".env")

		if err != nil {
			log.Println("Error loading .env config file")
		}
		log.Println("Successfully loaded the config file")
	}
	database.ConnectDb()
}

func main() {
	router.ClientRoutes()
}
