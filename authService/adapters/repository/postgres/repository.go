package postgres

import (
	"auth-svc/internal/entity"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type postgresDB struct {
	db *sql.DB
}

func NewAuthRepository(dbPool *sql.DB) postgresDB {
	return postgresDB{db: dbPool}
}

func (db postgresDB) StoreToken(userID int, token string, expiration time.Time) error {
	result, err := db.db.Exec("INSERT INTO tokens (user_id, token, expiration) VALUES ($1, $2, $3)",
		userID, token, expiration)
	fmt.Println("Actual SQL:", result)

	if err != nil {
		return err
	}
	return nil
}

func (db postgresDB) RetrieveToken(userID int) (*entity.Token, error) {
	var t entity.Token
	row := db.db.QueryRow("SELECT id, user_id, token, expiration FROM tokens WHERE user_id = $1", userID)
	err := row.Scan(&t.ID, &t.UserID, &t.TokenValue, &t.Expiration)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// AddRevokedToken adds the token ID to the list of revoked tokens in the repository
// func (r postgresDB) AddRevokedToken(tokenID string) error {
// 	// Implement logic to add token ID to the list of revoked tokens
// 	return nil
// }

// IsTokenRevoked checks if the token ID is in the list of revoked tokens
// func (r postgresDB) IsTokenRevoked(tokenID string) bool {
// 	// Implement logic to check if token ID is in the list of revoked tokens
// 	return false
// }
