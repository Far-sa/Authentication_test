package repository

import (
	"auth-svc/internal/entity"
	"database/sql"
	"log"
	"time"
)

type postgresDB struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) postgresDB {
	//TODO: for test purpose only-- add to config
	db, err := sql.Open("postgres", "postgresql://root:password@localhost:5432/authDB?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	return postgresDB{db: db}
}

func (db postgresDB) StoreToken(userID int, token string, expiration time.Time) error {
	_, err := db.db.Exec("INSERT INTO jwt_tokens (user_id, token, expiration) VALUES ($1, $2, $3)",
		userID, token, expiration)
	if err != nil {
		return err
	}
	return nil
}

func (db postgresDB) RetrieveToken(userID int) (*entity.Token, error) {
	var t entity.Token
	row := db.db.QueryRow("SELECT id, user_id, token, expiration FROM jwt_tokens WHERE user_id = $1", userID)
	err := row.Scan(&t.ID, &t.UserID, &t.TokenValue, &t.Expiration)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// AddRevokedToken adds the token ID to the list of revoked tokens in the repository
func (r postgresDB) AddRevokedToken(tokenID string) error {
	// Implement logic to add token ID to the list of revoked tokens
	return nil
}

// IsTokenRevoked checks if the token ID is in the list of revoked tokens
func (r postgresDB) IsTokenRevoked(tokenID string) bool {
	// Implement logic to check if token ID is in the list of revoked tokens
	return false
}

// TODO: implement save token
