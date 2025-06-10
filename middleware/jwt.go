package middleware

import(
	"time"
	"github.com/golang-jwt/jwt/v5"
	"os"
)
var jwtKey = []byte(os.Getenv("JWT_SECRET"))
func GenerateJWT(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(), // Token valid for 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}