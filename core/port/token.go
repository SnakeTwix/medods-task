package port

import (
	"context"
	"github.com/google/uuid"
	"medods-api/core/domain"
)

type TokenService interface {
	GetToken(context context.Context, userId uuid.UUID, ip string) (*domain.Token, error)
	RotateToken(context context.Context, refreshToken string, accessToken string, ip string) (*domain.Token, error)
}

type TokenData struct {
	RefreshTokenId string
	TokenFamily    string
}

type TokenRepository interface {
	WriteRefreshToken(context context.Context, tokenData *TokenData) error
	CheckCorrectGenerationRefreshToken(context context.Context, refreshTokenId string, tokenFamily string) (bool, error)
	RevokeTokenFamily(context context.Context, tokenFamily string) error
}
