package auth

// internal/auth/auth.go

type AuthRepository interface {
	SaveToken(token string) error
}

type AuthService interface {
	CreateAccessToken(userID string) (string, error)
	CreateRefreshToken(userID string) (string, error)
}

type authRepository struct {
	// Database or any other storage mechanism can be injected here
}

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{
		repo: repo,
	}
}
