package main

import (
	// "database/sql"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func init(){
	err:= godotenv.Load()
	if err != nil {
		fmt.Println("Warning .env file not loaded")
	}
}	

func main() {
	db := InitializeDB()
	fmt.Println("Successfully connected to the database!")
	MigrateModels(db)
	err := db.AutoMigrate(&User{}, &Post{}, &Reaction{}, &Comment{})
	if err != nil{
		log.Fatalf("Migration failed %v", err)
	}
	fmt.Println("Migration Success")
	router := gin.Default()
	SetupRoutes(router, db)
	router.Run("localhost:9090")
}
