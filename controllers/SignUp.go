package controller

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	User "signUp/models"
	"text/template"
	"time"

	"github.com/ReneKroon/ttlcache"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
)

type SignUpController struct {
	Db    *sql.DB
	Cache *ttlcache.Cache
}

func validateUserInput(user User.User) bool {
	if user.EmailID == "" || user.Name == "" || user.Password == "" {
		fmt.Println("User Details Incomplete")
		return false
	}

	return true
}
func generateOneTimePassword() int {
	return rand.Intn(90000) + 100000
}
func (s *SignUpController) SendEmail(to string, one_time_password int) error {
	//sender data.
	godotenv.Load(".env")
	from := os.Getenv("SENDER_EMAIL")
	password := os.Getenv("SENDER_PASSWORD")

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message
	t, err := template.ParseFiles("./template.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	var body bytes.Buffer
	// this we need to write when using an html file.
	mimeHeaders := "MIME-version: 1.0;\nContent-Type:text/html;charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Forgot Password Request \n%s \n\n", mimeHeaders)))
	t.Execute(&body, struct {
		OTP int
	}{OTP: one_time_password})

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, body.Bytes())
	if err != nil {

		return err
	}
	//catching the otp into the cache
	s.Cache.SetWithTTL(to, one_time_password, 10*time.Minute)
	fmt.Println("Email Sent Successfully!")
	return nil
}

func (s *SignUpController) emailValidation(emailid string) bool {

	/**1.emailid is valid or not**/
	_, err := mail.ParseAddress(emailid)
	if err != nil {
		fmt.Println("Invalid email ID", err)
		return false
	}

	/**2.send an verification email to the user**/
	s.SendEmail(emailid, generateOneTimePassword())

	return true
}
func (s *SignUpController) UserDetails(w http.ResponseWriter, r *http.Request) {

	//emailid , name , password , from request body
	var user User.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//validate user input
	if !validateUserInput(user) {
		http.Error(w, "Invalid user input", http.StatusBadRequest)
		return
	}

	// Email Verification
	if !s.emailValidation(user.EmailID) {
		http.Error(w, "Invalid email ID", http.StatusBadRequest)
		return
	}

	//return success or error
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "OTP sent successfully", "success": "true"})
}
func (s *SignUpController) VerifyEmail(w http.ResponseWriter, r *http.Request) {

	//emailid , otp,password,name from request body
	var user User.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	//otp verification
	otp, ok := s.Cache.Get(user.EmailID)
	if !ok {
		http.Error(w, "Invalid OTP", http.StatusBadRequest)
		return
	}
	if user.OTP != otp {
		http.Error(w, "Invalid OTP", http.StatusBadRequest)
		return
	}
	//check if emailid already exists
	var exists bool
	// Query for a value based on a single row.
	err := s.Db.QueryRow("SELECT EXISTS (SELECT 1 FROM users WHERE emailid = ?)", user.EmailID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			//no worries
		}
		fmt.Println("Error checking email ID", err)
		http.Error(w, "Error checking email ID", http.StatusInternalServerError)
		return
	}

	// Password Hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	newUser := `
	INSERT INTO users (emailid, name, password)
	VALUES (?, ?, ?)`
	_, err = s.Db.Exec(newUser, user.EmailID, user.Name, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	fmt.Println("User Stored Successfully into Database")

	//creating a jwt token for the user presently not implementing refresh token feature

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"emailid": user.EmailID,
		"name":    user.Name,
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}

	// Set the cookie in the response
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
	// return success or error
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully", "success": "true", "token": tokenString})

}
