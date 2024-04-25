package ports

import (
	"context"
	"user-svc/internal/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (entity.User, error)
	GetUserByID(ctx context.Context, userID uint) (entity.User, error)
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
}
