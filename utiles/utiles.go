package utiles

import (
	"database/sql"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func GenerateToken(Email string, Name string) (string, error) {
	godotenv.Load()
	expirationTime := time.Now().Add(10 * time.Minute)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"emailid":   Email,
		"name":      Name,
		"ExpiresAt": expirationTime.Unix(),
	})
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {

		return "", err
	}
	return tokenString, err
}
func InsertIntoDatabase(Db *sql.DB, Name string, Email string, Password string) error {
	newUser := `
	INSERT INTO users (emailid, name, password)
	VALUES (?, ?, ?)`
	_, err := Db.Exec(newUser, Email, Name, Password)
	return err
}
