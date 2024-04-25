package service

import "order-svc/ports"

type Service struct {
	userRepository ports.UserRepository
	eventProducer  ports.EventProducer
	logger         ports.Logger // New: Logger
}
