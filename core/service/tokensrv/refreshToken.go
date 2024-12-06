package tokensrv

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type refreshToken struct {
	TokenFamily uuid.UUID `json:"tokenFamily"`
	// The id to link the access and refresh tokens together
	LinkerId uuid.UUID `json:"linkerId"`
	jwt.RegisteredClaims
}

func (t *refreshToken) Sign(key []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, t)
	stringToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return stringToken, nil
}

func parseRefreshToken(tokenEncoded string, key []byte) (*refreshToken, error) {
	token, err := jwt.ParseWithClaims(tokenEncoded, &refreshToken{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	// At this point we know the token is valid
	if claims, ok := token.Claims.(*refreshToken); ok {
		return claims, nil
	}

	return nil, errors.New("not a refresh token")
}

func newRefreshToken(linkerId uuid.UUID) *refreshToken {
	token := refreshToken{
		TokenFamily: uuid.New(),
		LinkerId:    linkerId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// Expires after a month, arbitrary
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			ID:        uuid.New().String(),
		},
	}

	return &token
}
