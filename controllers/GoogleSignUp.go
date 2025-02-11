package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"signUp/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type GoogleSignUp struct {
	Db *sql.DB
}

func (ptr *GoogleSignUp) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	googleConfig := config.SetupConfig()
	url := googleConfig.AuthCodeURL("randomState")

	//redirect to google login page
	http.Redirect(w, r, url, http.StatusSeeOther)
}
func (ptr *GoogleSignUp) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	//inside google console i mentioned this to be our redirect url

	//verifying the state which we mentioned during GoogleLogin
	state := r.URL.Query()["state"][0]
	if state != "randomState" {
		fmt.Println("state does not exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//we will get a code which we will exchange with google server for token
	code := r.URL.Query()["code"][0]
	fmt.Println("code", code)
	//we create google configuration
	googleConfig := config.SetupConfig()

	//exchange the code for token and use token to get user details from google server
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Fprintln(w, "Code-Token Exchange Failed")
	}
	fmt.Println("token", token)
	fmt.Println("Access Token:", token.AccessToken)

	//use gogle api to get user info
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	if err != nil {
		fmt.Println("Error in Fetching User Info from google server : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//parse response
	type user struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	userDetails := user{}
	err = json.NewDecoder(resp.Body).Decode(&userDetails)
	if err != nil {
		fmt.Println("Error pasing JSON : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	fmt.Println(userDetails)

	//add email id password to the data base

	newUser := `
	INSERT INTO users (emailid, name, password)
	VALUES (?, ?, ?)`
	_, err = ptr.Db.Exec(newUser, userDetails.Email, userDetails.Name, "")
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	//now create a jwt token out of this user data presently not implementing refresh token feature
	expirationTime := time.Now().Add(10 * time.Minute)
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"emailid":   userDetails.Email,
		"name":      userDetails.Name,
		"ExpiresAt": expirationTime.Unix(),
	})
	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, "Error signing token", http.StatusInternalServerError)
		return
	}
	fmt.Println(tokenString)

	// Set the cookie in the response
	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenString,
	}
	http.SetCookie(w, cookie)
	//return success or error
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully", "success": "true", "token": tokenString})

}
