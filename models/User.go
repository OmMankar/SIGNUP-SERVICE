package User

type User struct {
	EmailID  string `json:"emailid"`
	Name     string `json:"name"`
	Password string `json:"password"`
	OTP      int    `json:"otp"`
}
