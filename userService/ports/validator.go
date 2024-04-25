package ports

import "user-svc/internal/service/param"

type Validator interface {
	ValidateRegisterRequest(req param.RegisterRequest) (map[string]string, error)
	checkPhoneUniqueness(value interface{}) error
}
