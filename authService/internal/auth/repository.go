package auth

type AuthRepository interface {
	SaveToken(token string) error
}

type authRepository struct {
	// Database or any other storage mechanism can be injected here
}

func NewAuthRepository() AuthRepository {
	// Initialize and return a new instance of AuthRepository
	return &authRepository{}
}

func (repo *authRepository) SaveToken(token string) error {
	// Implementation to save the token in the database
	return nil
}
