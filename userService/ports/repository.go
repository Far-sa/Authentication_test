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

// MigrationManager interface defines methods for managing migrations
type MigrationManager interface {
	Up() error   // Apply pending migrations
	Down() error // Revert migrations (optional)
}
