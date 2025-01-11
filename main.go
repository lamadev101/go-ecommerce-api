package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lamadev101/ecommerce-api/constant"
	"github.com/lamadev101/ecommerce-api/database"
	"github.com/lamadev101/ecommerce-api/router"
	"github.com/lamadev101/ecommerce-api/types"
	"github.com/lamadev101/ecommerce-api/utils"
)

func init() {
	// Load environment variables
	if _, err := os.Stat(".env"); err == nil {
		log.Println("Loading the config from .env file")
		err = godotenv.Load(".env")

		if err != nil {
			log.Println("Error loading .env config file")
		}
	}
	database.ConnectDb()

	// creating system admin
	hashPassword := utils.GenerateHashPassword("root@123")
	user := types.User{
		Name:     "Admin",
		Email:    "rootuser@gmail.com",
		Password: hashPassword,
		UserType: constant.ADMIN_ROLE,
	}
	if u := database.Mgr.GetSingleRecordByEmailForUser(user.Email, constant.USERS_COLLECTION); u.Email == "" {
		_, err := database.Mgr.Insert(user, constant.USERS_COLLECTION)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	router.ClientRoutes()
}
