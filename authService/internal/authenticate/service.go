package authenticate

type authService struct {
	tokenRepo TokenRepository
}

// NewTokenHandler creates a new TokenHandler with the given authService
func NewTokenService(tokenRepo TokenRepository) *authService {
	return &authService{tokenRepo: tokenRepo}
}

func (s *authService) GenerateToken(userID string, roles []string) (string, error) {
	panic("")
	// Generate JWT token with user ID and roles as claims
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	//     "userID": userID,
	//     "roles":  roles,
	//     // Add other relevant claims as needed
	// })

	// // Sign the token with a secret key
	// tokenString, err := token.SignedString([]byte("secret"))
	// if err != nil {
	//     return "", err
	// }

	// return tokenString, nil
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
