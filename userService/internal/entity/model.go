package entity

import "time"

// User represents a user entity
type User struct {
	ID          uint      `json:"id"`
	PhoneNumber string    `json:"phone_number"`
	Email       string    `json:"email"`
	Password    string    `json:"password"`
	CreatedAt   time.Time `db:"created_at"`
}
