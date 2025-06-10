package main
import(
	"github.com/google/uuid"
)
type User struct {
	ID     uuid.UUID  `gorm:"primaryKey;autoIncrement" json:"id"`
	Name  string `gorm:"not null" json:"name" binding:"required"`
	Email string `gorm:"unique;not null" json:"email" binding:"required,email"`
	Posts []Post `gorm:"foreignKey:UserID" json:"posts"`
	Password string `gorm:"not null" json:"password" binding:"required"`
}
type Post struct {
	ID           uuid.UUID         `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string        `gorm:"not null" json:"title" binding:"required"`
	Body         string        `gorm:"type:text;not null" json:"body"`
	Interactions []Interaction `gorm:"foreignKey:PostID" json:"interactions"`
	UserID       uuid.UUID          `json:"user_id"`
}
type Interaction struct {
	ID       uuid.UUID   `gorm:"primaryKey;autoIncrement" json:"id"`
	Comment  string `gorm:"type:text;not null" json:"comment"`
	Likes    int    `gorm:"default:0" json:"likes"`
	Dislikes int    `gorm:"default:0" json:"dislikes"`
	PostID   uuid.UUID   `json:"post_id"`
}
