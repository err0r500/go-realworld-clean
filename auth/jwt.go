package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getEnvOrDefault("JWT_SECRET", "your-secret-key-change-in-production"))

// Claims represents the JWT claims
type Claims struct {
	KratosSessionToken string `json:"sid"` // Kratos session token
	UserID             string `json:"uid"` // Kratos identity ID
	jwt.RegisteredClaims
}

// GenerateToken generates a JWT token containing Kratos session information
func GenerateToken(sessionToken, userID string, expiresAt time.Time) (string, error) {
	claims := Claims{
		KratosSessionToken: sessionToken,
		UserID:             userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// validateToken validates a JWT token and returns the claims
func validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
