package service

import (
	"auth-svc/internal/entity"
	"auth-svc/internal/param"
	"auth-svc/internal/ports"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

// type Config struct {
// 	JwtSignKey                     string
// 	AccessTokenSubject             string
// 	RefreshTokenSubject            string
// 	AccessTokenExpirationDuration  time.Duration
// 	RefreshTokenExpirationDuration time.Duration
// }

const (
	JwtSignKey                     = "jwt-secret"
	AccessTokenSubject             = "at"
	RefreshTokenSubject            = "rt"
	AccessTokenExpirationDuration  = time.Hour * 24
	RefreshTokenExpirationDuration = time.Hour * 24 * 7
)

type authService struct {
	authRepo ports.AuthRepository
	event    ports.EventPublisher
}

// NewTokenHandler creates a new TokenHandler with the given authService
func NewAuthService(authRepo ports.AuthRepository, event ports.EventPublisher) authService {
	return authService{authRepo: authRepo}
}

func (s authService) Login(ctx context.Context, user param.LoginRequest) (param.LoginResponse, error) {

	//TODO: get user info from rabbitmq
	queue, _ := s.event.CreateQueue()
	_ := s.event.CreateBinding()
	messages, _ := s.event.Consume()

	if user.Password != getMD5Hash(req.Password) {
		return param.LoginResponse{}, fmt.Errorf("username/ password incorrect")
	}

	// create tokens
	accessToken, err := s.createAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	refreshToken, err := s.refreshAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error : %w", err)
	}

	//TODO: publish token and save to DB

	return param.LoginResponse{
		User:   param.UserInfo{ID: user.ID, Email: user.Email},
		Tokens: param.Tokens{AccessToken: accessToken, RefreshToken: refreshToken},
	}, nil
}

func (s authService) AddRevokedToken(tokenID string) error {
	panic("")
}

// func (s authService) IsRevokedToken(tokenID string) error {
// 	panic("")
// }

func (s authService) createAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, AccessTokenSubject, AccessTokenExpirationDuration)
}

func (s authService) refreshAccessToken(user entity.User) (string, error) {
	return s.createToken(user.ID, RefreshTokenSubject, RefreshTokenExpirationDuration)
}

func (s authService) VerifyToken(bearerToken string) (*Claims, error) {

	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSignKey), nil
	})

	if err != nil {
		return nil, err
	}

	// convert interface to conceret object
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil

	} else {
		return nil, err
	}
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}

func (s authService) createToken(userID uint, subject string, expiresDuration time.Duration) (string, error) {
	// create a signer for rsa 256
	//t := jwt.New(jwt.GetSigningMethod("RS256"))
	// TODO replace with rsa 256 RS256

	// set our claims
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: subject,
			// set the expire time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresDuration)),
		},
		UserID: userID,
	}

	// TODO add sign method to config
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := accessToken.SignedString([]byte(JwtSignKey))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// ! Middleware for token validation in Traefik
// func TokenValidationMiddleware(next http.Handler, authService authService) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Extract token from request
// 		token := ExtractTokenFromRequest(r)

// 		// Validate token
// 		if isValid := authService.ValidateToken(token); !isValid {
// 			// Token is invalid or revoked
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Token is valid, proceed to next handler
// 		next.ServeHTTP(w, r)
// 	})
// }

// // ExtractTokenFromRequest extracts JWT token from request
// func ExtractTokenFromRequest(r *http.Request) string {
// 	// Extract token from request headers, cookies, or query parameters
// 	// Example: Authorization: Bearer <token>
// 	token := r.Header.Get("Authorization")
// 	if token != "" {
// 		return strings.TrimPrefix(token, "Bearer ")
// 	}

// 	// Extract token from cookies or query parameters if needed

// 	return ""
// }
