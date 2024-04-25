package ports

import (
	"auth-svc/internal/authenticate/param"
	"context"
)

type AuthRepository interface {
	AddRevokedToken(tokenID string) error
	IsTokenRevoked(tokenID string) bool
	SaveTokens()
}

type AuthService interface {
	AddRevokeToken(tokenID string) error
	ValidateToken(token string) bool

	Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error)
}
