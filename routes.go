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
func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
func (r *Reaction) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.New()
	return
}
func SetupRoutes(router *gin.Engine, database *gorm.DB) {
	protected := router.Group("/api")
	protected.Use(middleware.JWTMiddleware())
	//post beginning
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
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		userID, ok := userIDValue.(uuid.UUID)

		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}

		if err := context.ShouldBindJSON(&newPost); err != nil {
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
	protected.POST("/post/comment/:postID", func(context *gin.Context) {
		var comment Comment
		postIDStr := context.Param("postID")
		userIDValue, ok := context.Get("userID")
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User ID not in context"})
			return
		}
		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID format"})
			return
		}
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := context.ShouldBindJSON(&comment); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		comment.PostID = postID
		comment.UserID = userID

		if err := database.Create(&comment).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "successfully created comment",
			"comment": comment,
		})

	}) //post ending
	//get beginning
	protected.GET("/users", func(context *gin.Context) {
		var users []User
		if err := database.Preload("Posts").Find(&users).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, users)
	})
	protected.GET("/user", func(context *gin.Context) {
		var user User
		userIDValue, ok := context.Get("userID")
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}

		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
			return
		}
		if err := database.First(&user, "id = ?", userID).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, user)
	})
	protected.GET("/posts", func(context *gin.Context) {
		var user User
		userIDValue, ok := context.Get("userID")
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in context"})
			return
		}
		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID format"})
			return
		}
		if err := database.Preload("Posts").Where("id = ?", userID).First(&user).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Posts retrieved successfully",
			"posts":   user.Posts,
		})
	}) //get ending
	//patch beginning
	protected.PATCH("/posts/:postID/edit", func(context *gin.Context) {
		var post Post
		postIDStr := context.Param("postID")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userIDValue, ok := context.Get("userID")
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User ID not in context"})
			return
		}
		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID format"})
			return
		}
		var changedPost struct {
			Title string `gorm:"not null" json:"title" binding:"required"`
			Body  string `gorm:"type:text;not null" json:"body"`
		}
		if err := context.ShouldBindJSON(&changedPost); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := database.Preload("Comments", "Reactions").Where("user_id = ? AND id = ? ", userID, postID).First(&post).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		post.Body = changedPost.Body
		post.Title = changedPost.Title
		if err := database.Model(&post).Updates(map[string]interface{}{
			"title": changedPost.Title,
			"body":  changedPost.Body,
		}).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Post successfully updated",
			"post":    post,
		})

	})
	protected.PATCH("/post/postID/reactions", func (context *gin.Context){
	
		postIDStr := context.Param("postID")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userIDValue, ok := context.Get("userID")
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User ID not in context"})
			return
		}
		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID format"})
			return
		}
		var reactionData struct {
			Action string
		}
		if err := context.ShouldBindJSON(&reactionData); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var reaction Reaction
		if err := database.Where("post_id = ? AND user_id = ?", postID, userID).First(&reaction).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		switch reactionData.Action{
		case "like":
			reaction.Likes++
		case "dislike":
			reaction.Dislikes++
		default:
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid Action"})
			return
		}
		if err := database.Save(&reaction).Error; err != nil{
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Reactions Updated",
			"reaction": reaction,
		})
	})//patch ending
	protected.DELETE("/post/:postID/delete", func(context * gin.Context){
		postIDStr := context.Param("postID")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userIDValue, ok := context.Get("userID")
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error":"User ID is not in context"})
			return
		}
		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}
		var post Post
		if err := database.Where("id = ? AND user_id = ?", postID , userID).First(&post).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err := database.Delete(&post).Error; err != nil{
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Post successfully deleted",
		})
	})
}
