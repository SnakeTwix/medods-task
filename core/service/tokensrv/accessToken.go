package tokensrv

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type accessToken struct {
	UserGUID uuid.UUID
	jwt.RegisteredClaims
}

func (t *accessToken) Sign(key []byte) (string, error) {
	// SHA-512 signing
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, t)
	stringToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return stringToken, nil
}

func parseAccessToken(tokenEncoded string, key []byte) (*accessToken, error) {
	token, err := jwt.ParseWithClaims(tokenEncoded, accessToken{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	// At this point we know the token is valid
	if claims, ok := token.Claims.(*accessToken); ok {
		return claims, nil
	}

	return nil, errors.New("not an access token")
}

func newAccessToken(userId uuid.UUID) *accessToken {
	token := accessToken{
		UserGUID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	return &token
}
