package jwtoken

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type TokenValidator interface {
	ValidateToken(tokenString string) error
}

type tokenValidator struct{ secretKey []byte }

func NewTokenValidator(secret string) TokenValidator {
	return &tokenValidator{secretKey: []byte(secret)}
}

func (tv *tokenValidator) ValidateToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method: %v", t.Header["alg"])
		}

		return tv.secretKey, nil
	})
	if err != nil {
		return fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}
