package main

import (
	"net/http"
	"sql-blog/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}
func SetupRoutes(router *gin.Engine, database *gorm.DB) {
	protected := router.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	router.POST("/register", func(context *gin.Context) {
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
	router.POST("/users/login", func(context *gin.Context) {
		var loginData struct {
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}
		if err := context.ShouldBindJSON(&loginData); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var user User
		if err := database.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if loginData.Password != user.Password {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}
		token, err := middleware.GenerateJWT(user.ID.String())
		if err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"user":  user,
			"token": token})
	})

	protected.POST("/post", func(context *gin.Context) {
		var newPost Post
		userIDValue, ok := context.Get("userID")
		if !ok{
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		userID, ok := userIDValue.(uuid.UUID)
	
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}
		
		if err := context.ShouldBindJSON(&newPost); err != nil{
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newPost.UserID = userID
		if err := database.Create(&newPost).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Successfully created Post",
			"post":    newPost})
	})
	protected.GET("/users", func(context *gin.Context) {
		var users []User
		if err := database.Preload("Posts").Find(&users).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, users)
	})
	protected.GET("/users/:id", func(context *gin.Context) {
		var user User
		id := context.Param("id")
		userID, err := uuid.Parse(id)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
			return
		}
		if err := database.First(&user, userID).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, user)
	})
	protected.GET("/users/email/:email", func(context *gin.Context) {
		var user User
		email := context.Param("email")
		if err := database.Preload("Posts").Where("email = ?", email).First(&user).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, user)
	})

}
