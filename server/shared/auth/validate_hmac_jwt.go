package auth

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateHmac(tokenString string) (int, int8, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return HMACSecretKey, nil
	})

	if err != nil {
		return -1, -1, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return -1, -1, fmt.Errorf("invalid token")
	}

	if claims.Issuer != Issuer {
		return -1, -1, fmt.Errorf("invalid issuer")
	}

	return claims.UserID, claims.TokenType, nil
}
