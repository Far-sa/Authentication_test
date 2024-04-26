package mocks

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/mock"
)

type MockEventPublisher struct {
	mock.Mock
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{}
}

func (m *MockEventPublisher) Publish(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {
	args := m.Called("Publish", ctx, exchange, routingKey, options)
	return args.Error(1)
}

func (m *MockEventPublisher) CreateQueue(queueName string, durable, autodelete bool) (amqp.Queue, error) {
	args := m.Called("CreateQueue", queueName, durable, autodelete)
	return args.Get(0).(amqp.Queue), args.Error(1)
}

func (m *MockEventPublisher) CreateBinding(name, binding, exchange string) error {
	args := m.Called("CreateBinding", name, binding, exchange)
	return args.Error(0)
}
