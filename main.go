package main

import (
	// "database/sql"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func init(){
	godotenv.Load()
}	

func main() {
	db := InitializeDB()
	fmt.Println("Successfully connected to the database!")
	MigrateModels(db)

	router := gin.Default()
	SetupRoutes(router, db)
	router.Run("localhost:9090")
}
