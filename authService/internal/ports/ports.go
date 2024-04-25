package ports

import (
	"auth-svc/internal/param"
	"context"
)

type AuthRepository interface {
	AddRevokedToken(tokenID string) error
	IsTokenRevoked(tokenID string) bool
}

type AuthService interface {
	RevokeToken(ctx context.Context, tokenID string) error
	Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error)
}
