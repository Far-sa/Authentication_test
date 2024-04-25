package repository

type mysqlDB struct{}

func NewAuthRepository() mysqlDB {
	// Initialize and return a new TokenRepositoryImpl
	return mysqlDB{}
}

// AddRevokedToken adds the token ID to the list of revoked tokens in the repository
func (r mysqlDB) AddRevokedToken(tokenID string) error {
	// Implement logic to add token ID to the list of revoked tokens
	return nil
}

// IsTokenRevoked checks if the token ID is in the list of revoked tokens
func (r mysqlDB) IsTokenRevoked(tokenID string) bool {
	// Implement logic to check if token ID is in the list of revoked tokens
	return false
}

// TODO: implement save token
