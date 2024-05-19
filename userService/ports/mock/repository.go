package mocks

import (
	"context"
	"reflect"
	"user-svc/internal/entity"

	"github.com/stretchr/testify/mock"
	// Replace with the path to your UserRepository interface and entity types
)

type MockUserRepository struct {
	mock.Mock
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user entity.User) (entity.User, error) {
	args := m.Called(ctx, user)
	var ret entity.User
	var err error
	if reflect.ValueOf(args.Get(0)).IsNil() {
		ret = entity.User{}
	} else {
		ret = args.Get(0).(entity.User)
	}
	if reflect.ValueOf(args.Get(1)).IsNil() {
		err = nil
	} else {
		err = args.Error(1)
	}
	return ret, err
}

func (m *MockUserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, userID uint) (entity.User, error) {
	args := m.Called(ctx, userID)
	var ret entity.User
	var err error
	if reflect.ValueOf(args.Get(0)).IsNil() {
		ret = entity.User{}
	} else {
		ret = args.Get(0).(entity.User)
	}
	if reflect.ValueOf(args.Get(1)).IsNil() {
		err = nil
	} else {
		err = args.Error(1)
	}
	return ret, err
}

func (m *MockUserRepository) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	args := m.Called(phoneNumber)
	var ret bool
	var err error
	if reflect.ValueOf(args.Get(0)).IsNil() {
		ret = false
	} else {
		ret = args.Get(0).(bool)
	}
	if reflect.ValueOf(args.Get(1)).IsNil() {
		err = nil
	} else {
		err = args.Error(1)
	}
	return ret, err
}

// func TestCreateUser(t *testing.T) {
// 	// Create a MockUserRepository instance
// 	repo := &mocks.MockUserRepository{}

// 	// Define the expected user and return value for CreateUser method
// 	expectedUser := entity.User{
// 		ID:          1,
// 		FirstName:   "John",
// 		LastName:    "Doe",
// 		PhoneNumber: "1234567890",
// 		Email:       "john@example.com",
// 	}
// 	repo.CreateUserFn = func(ctx context.Context, user entity.User) (entity.User, error) {
// 		return expectedUser, nil
// 	}

// 	// Call the function under test that uses CreateUser method
// 	createdUser, err := someFunctionThatCreatesUser(repo)

// 	// Assertions
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedUser, createdUser)
// }

// func TestGetUserByID(t *testing.T) {
// 	// Create a MockUserRepository instance
// 	repo := &mocks.MockUserRepository{}

// 	// Define the expected user and return value for GetUserByID method
// 	expectedUser := entity.User{
// 		ID:          1,
// 		FirstName:   "John",
// 		LastName:    "Doe",
// 		PhoneNumber: "1234567890",
// 		Email:       "john@example.com",
// 	}
// 	repo.GetUserByIDFn = func(ctx context.Context, userID uint) (entity.User, error) {
// 		return expectedUser, nil
// 	}

// 	// Call the function under test that uses GetUserByID method
// 	userByID, err := someFunctionThatGetsUserByID(repo, 1)

// 	// Assertions
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedUser, userByID)
// }

// func TestIsPhoneNumberUnique(t *testing.T) {
// 	// Create a MockUserRepository instance
// 	repo := &mocks.MockUserRepository{}

// 	// Define the expected return value for IsPhoneNumberUnique method
// 	phoneNumber := "1234567890"
// 	expectedUnique := true
// 	repo.IsPhoneNumberUniqueFn = func(phoneNumber string) (bool, error) {
// 		return expectedUnique, nil
// 	}

// 	// Call the function under test that checks phone number uniqueness
// 	isUnique, err := someFunctionThatChecksPhoneNumberUnique(repo, phoneNumber)

// 	// Assertions
// 	assert.NoError(t, err)
// 	assert.Equal(t, expectedUnique, isUnique)
// }
