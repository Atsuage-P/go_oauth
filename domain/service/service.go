package service

import (
	"context"
	"errors"
	"fmt"
	"go_oauth/domain"
	"go_oauth/internal/cert"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
}

func NewAuthService() domain.AuthService {
	return &authService{}
}

func (as *authService) VerifyPassword(ctx context.Context, hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (as *authService) HashPassword(ctx context.Context, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// NOTE: golang-jwt(https://pkg.go.dev/github.com/golang-jwt/jwt/v5)
func (as *authService) CreateToken(ctx context.Context, userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		// TODO: 1を環境変数にする
		"exp": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	})
	tokenString, err := token.SignedString(cert.RawPrivKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (as *authService) InvalidateToken(ctx context.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return cert.RawPrivKey, nil
	})
	if token.Valid {
		fmt.Println("You look nice today")
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return fmt.Errorf("That's not even a token")
	} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		// Invalid signature
		return fmt.Errorf("Invalid signature")
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		// Token is either expired or not active yet
		return fmt.Errorf("Timing is everything")
	} else {
		return fmt.Errorf("Couldn't handle this token: %w", err)
	}
	return nil
}
