package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, database *gorm.DB) {
	router.POST("/users", func(context *gin.Context) {
		var newUser User

		if err := context.ShouldBindJSON(&newUser); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := database.Create(&newUser).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusCreated, newUser)
	})
	router.POST("/users/multi", func(context *gin.Context) {
		var newUsers []User
		if err := context.ShouldBindJSON(&newUsers); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		if err := database.Create(&newUsers).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		context.IndentedJSON(http.StatusCreated, newUsers)
	})
	router.GET("/users", func(context *gin.Context) {
		var users []User
		if err := database.Preload("Posts").Find(&users).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		context.IndentedJSON(http.StatusOK, users)
	})
	router.GET("/users/:id", func(context *gin.Context) {
		var user User
		id := context.Param("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		if err := database.First(&user, idInt).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, user)
	})
	router.GET("/users/email/:email", func(context *gin.Context) {
		var user User
		email := context.Param("email")
		if err := database.Preload("Posts").Where("email = ?", email).First(&user).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, user)
	})
	router.POST("/users/email/:email/post", func(context *gin.Context) {
		var user User
		email := context.Param("email")
		if err := database.Where("email = ?", email).First(&user).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		var newPost Post
		newPost.UserID = user.ID
		if err := context.ShouldBindJSON(&newPost); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := database.Create(&newPost).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"ERROR": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusCreated, newPost)
	})
}
