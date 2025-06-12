package main

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Name     string    `gorm:"not null" json:"name" binding:"required"`
	Email    string    `gorm:"unique;not null" json:"email" binding:"required,email"`
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts"`
	Password string    `gorm:"not null" json:"-" binding:"required"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments"`
}
type Post struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;" json:"id"`
	Title        string     `gorm:"not null" json:"title" binding:"required"`
	Body         string     `gorm:"type:text;not null" json:"body"`
	Reactions []Reaction `gorm:"foreignKey:PostID" json:"interactions"`
	UserID       uuid.UUID  `gorm:"type:uuid;" json:"user_id"`
	Comments     []Comment  `gorm:"foreignKey:PostID" json:"comments"`
}
type Reaction struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Likes    int       `gorm:"default:0" json:"likes"`
	Dislikes int       `gorm:"default:0" json:"dislikes"`
	PostID   uuid.UUID `gorm:"type:uuid;" json:"post_id"`
	UserID   uuid.UUID `gorm:"type:uuid;" json:"user_id"`
}
type Comment struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	Body   string    `gorm:"type:text;not null" json:"body"`
	UserID uuid.UUID `gorm:"type:uuid;" json:"user_id"`
	PostID uuid.UUID `gorm:"type:uuid;" json:"post_id"`
}
