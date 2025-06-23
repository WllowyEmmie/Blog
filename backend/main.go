package main

import (
	// "database/sql"

	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning .env file not loaded")
		return
	}
}

func main() {
	db := InitializeDB()
	fmt.Println("Successfully connected to the database!")
	MigrateModels(db)
	err := db.AutoMigrate(&User{}, &Post{}, &Reaction{}, &Comment{})
	if err != nil {
		log.Fatalf("Migration failed %v", err)
	}
	fmt.Println("Migration Success")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:5500"},
		AllowMethods:     []string{"PUT", "GET", "PATCH", "DELETE", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge: 12 * time.Hour,
	}))
	SetupRoutes(router, db)
	router.Run("localhost:9090")
}
