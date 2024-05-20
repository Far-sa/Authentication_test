package service_test

import (
	"auth-svc/internal/param"
	"auth-svc/internal/service"
	"reflect"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestComparePassword(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	type testCase struct {
		name           string
		hashedPassword string
		reqPassword    string
		expectedResult bool
		expectedError  error
	}
	cases := []testCase{
		{
			name:           "Matching passwords",
			hashedPassword: string(hashedPassword),
			reqPassword:    "password123",
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name:           "Non-matching passwords",
			hashedPassword: string(hashedPassword),
			reqPassword:    "wrongpassword",
			expectedResult: false,
			expectedError:  nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := service.ComparePassword(tc.hashedPassword, tc.reqPassword)

			if (err != nil && tc.expectedError == nil) || (err == nil && tc.expectedError != nil) {
				t.Errorf("ComparePassword() error = %v, expected error = %v", err, tc.expectedError)
			}

			if result != tc.expectedResult {
				t.Errorf("ComparePassword() = %v, want %v", result, tc.expectedResult)
			}
		})
	}
}
func TestUnmarshalUser(t *testing.T) {

	type testCase struct {
		name     string
		input    []byte
		expected param.User
		wantErr  bool
	}
	cases := []testCase{
		{
			name:     "Valid JSON data",
			input:    []byte(`{"name": "Alice", "age": 30}`),
			expected: param.User{ID: 1, Email: "Teo", Password: "123456"},
			wantErr:  false,
		},
		{
			name:     "Invalid JSON data",
			input:    []byte(`invalid json`),
			expected: param.User{},
			wantErr:  true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			user, err := service.UnmarshalUser(c.input)
			if (err != nil) != c.wantErr {
				t.Errorf("unmarshalUser() error = %v, wantErr %v", err, c.wantErr)
				return
			}

			if !reflect.DeepEqual(user, c.expected) {
				t.Errorf("unmarshalUser() = %v, want %v", user, c.expected)
			}
		})
	}
}
