package entity

import "time"

type Token struct {
	ID         uint
	UserID     uint
	Expiration time.Time
	TokenValue string
}
