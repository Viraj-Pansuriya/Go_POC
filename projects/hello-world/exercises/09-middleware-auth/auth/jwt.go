package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// auth/jwt.go
type Claims struct {
	UserID uint64 `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var (
	secretKey = []byte("gbuhtegriesdjocklmsfvbgtnrhuiefjndkclmefjrgthiuefjndkrgthy")
)

func GenerateToken(userID uint64, email, role string, ttl time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userID,
			"email":  email,
			"role":   role,
			"exp":    time.Now().Add(ttl).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
