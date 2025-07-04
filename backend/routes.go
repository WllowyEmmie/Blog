package main

import (
	"errors"
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
	//registering
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
	//login post
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
	// creating a new post
	protected.POST("/post", func(context *gin.Context) {
		var newPostInput struct {
			Title string `gorm:"not null" json:"title" binding:"required"`
			Body  string `gorm:"type:text;not null" json:"body"`
		}
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

		if err := context.ShouldBindJSON(&newPostInput); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		newPost := Post{
			Title:  newPostInput.Title,
			Body:   newPostInput.Body,
			UserID: userID,
		}
		if err := database.Create(&newPost).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Successfully created Post",
			"post":    newPost})
	})
	//Posting a comment on a post
	protected.POST("/post/comment/:postID", func(context *gin.Context) {
		var input struct {
			Body string `json:"body" binding:"required"`
		}
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

		if err := context.ShouldBindJSON(&input); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		comment := Comment{
			Body:   input.Body,
			UserID: userID,
			PostID: postID,
		}

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
	//Getting all users
	protected.GET("/users", func(context *gin.Context) {
		var users []User
		if err := database.Preload("Posts").Find(&users).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, users)
	})
	//getting all the post's comments
	protected.GET("/comments/:postID", func(context *gin.Context) {
		var comments []Comment
		var postIDStr = context.Param("postID")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := database.Preload("User").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message":  "Successfully retreived comments",
			"comments": comments,
		})
	})
	//Getting all blog posts
	protected.GET("/all-posts", func(context *gin.Context) {
		var posts []Post
		if err := database.Preload("Comments").Preload("Reactions").Preload("User").Find(&posts).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(posts) == 0 {
			context.IndentedJSON(http.StatusNotFound, gin.H{"message": "No posts found"})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Successfully retrieved posts",
			"posts":   posts,
		})
	})
	//Getting a user
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
	//Getting a posts's interaction
	protected.GET("/reactions/:postID", func(context *gin.Context) {
		postIDStr := context.Param("postID")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var reactions []Reaction
		if err := database.Where("post_id = ?", postID).Find(&reactions).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message":   "Retrieved Reactions",
			"reactions": reactions,
		})
	})
	//Getting a user's Post
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
	protected.PATCH("/posts/:postID", func(context *gin.Context) {

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
			Title string ` json:"title" binding:"required"`
			Body  string `json:"body"`
		}
		if err := context.ShouldBindJSON(&changedPost); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		post := Post{
			Title:  changedPost.Title,
			Body:   changedPost.Body,
			UserID: userID,
		}
		if err := database.Preload("Comments").Preload("Reactions").Where("user_id = ? AND id = ? ", userID, postID).First(&post).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			} else {
				context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		}

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
	protected.PATCH("/post/:postID/reactions", func(context *gin.Context) {

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
			Action string `json:"action"`
		}
		if err := context.ShouldBindJSON(&reactionData); err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var reaction Reaction
		if err := database.Where("post_id = ? AND user_id = ?", postID, userID).First(&reaction).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				reaction = Reaction{
					PostID:   postID,
					UserID:   userID,
					Likes:    0,
					Dislikes: 0,
				}
			} else {
				context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

		}
		switch reactionData.Action {
		case "like":
			reaction.Likes++
		case "dislike":
			reaction.Dislikes++
		case "not-like":
			reaction.Likes--
		case "not-dislike":
			reaction.Dislikes--
		default:
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid Action"})
			return
		}
		if err := database.Save(&reaction).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message":  "Reactions Updated",
			"reaction": reaction,
		})
	}) //patch ending
	protected.DELETE("/post/:postID/delete", func(context *gin.Context) {
		postIDStr := context.Param("postID")
		postID, err := uuid.Parse(postIDStr)
		if err != nil {
			context.IndentedJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		userIDValue, ok := context.Get("userID")
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "User ID is not in context"})
			return
		}
		userID, ok := userIDValue.(uuid.UUID)
		if !ok {
			context.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
			return
		}
		var post Post
		if err := database.Where("id = ? AND user_id = ?", postID, userID).First(&post).Error; err != nil {
			context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err := database.Delete(&post).Error; err != nil {
			context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		context.IndentedJSON(http.StatusOK, gin.H{
			"message": "Post successfully deleted",
		})
	})
}
