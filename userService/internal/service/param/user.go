package param

type RegisterRequest struct {
	//Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	Email       string `json:"email"`
}

type UserInfo struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

type RegisterResponse struct {
	User UserInfo `json:"user"`
}

type LoginRequest struct {
	Email string `json:"email"`
	// PhoneNumber string `json:"phone_number"`
	Password string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	User      UserInfo `json:"user"`
	UserExist bool     `json:"userExist"`
	Error     string   `json:"error,omitempty"` // Optional field for error me
	// Tokens Tokens   `json:"tokens"`
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}
