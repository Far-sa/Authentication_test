package ports

import (
	"context"
	"user-svc/internal/service/param"
)

type Service interface {
	Register(ctx context.Context, req param.RegisterRequest) (param.RegisterResponse, error)
}
