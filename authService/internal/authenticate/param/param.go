package param

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Tokens struct {
	// token
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type LoginResponse struct {
	User   UserInfo `json:"user"`
	Tokens Tokens   `json:"token"`
}
