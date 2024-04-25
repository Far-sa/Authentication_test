package ports

import "context"

type EventProducer interface {
	PublishUserRegisteredEvent(ctx context.Context, userID string) error
}
