package authenticate

type TokenRepository interface {
	AddRevokedToken(tokenID string) error
	IsTokenRevoked(tokenID string) bool
}

type TokenService interface {
	GenerateToken(userID string, roles []string) (string, error)
	RevokeToken(tokenID string) error
	ValidateToken(token string) bool
}
