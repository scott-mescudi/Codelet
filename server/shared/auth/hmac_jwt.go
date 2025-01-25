package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var HMACSecretKey = []byte("apeirbvpijebvejbfpvibfevqepirjvb")
var Issuer = "codelet"

type Claims struct {
	UserID int
	TokenType int8
	jwt.RegisteredClaims
}

func GenerateHMac(userID int, tokenType int8, timeframe time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(timeframe),
		},
	})

	tkstring, err := token.SignedString(HMACSecretKey)
	if err != nil {
		panic(err)
	}

	return tkstring
}
