package ports

import "context"

type EventPublisher interface {
	//PublishUserRegisteredEvent(ctx context.Context, userID string) error
	PublishUserRegisteredEvent(ctx context.Context, data []byte) error
}
