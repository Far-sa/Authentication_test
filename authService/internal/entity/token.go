package entity

import "time"

type Token struct {
	ID         string
	UserID     string
	Expiration time.Time
	TokenValue string
}
