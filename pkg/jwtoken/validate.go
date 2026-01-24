package jwtoken

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
)

type TokenValidator interface {
	ValidateTokenWithClaims(tokenString string) (jwt.MapClaims, error)
	ValidateClaims(claims jwt.MapClaims) error
}

type tokenValidator struct{ secretKey []byte }

func NewTokenValidator(secret string) TokenValidator {
	return &tokenValidator{secretKey: []byte(secret)}
}

func (tv *tokenValidator) ValidateTokenWithClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method: %v", t.Header["alg"])
		}

		return tv.secretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

func (tv *tokenValidator) ValidateClaims(claims jwt.MapClaims) error {
	if _, ok := claims["userID"]; !ok {
		return errors.New("token hasn't userID")
	}

	if _, ok := claims["email"]; !ok {
		return errors.New("token hasn't email")
	}

	return nil
}
