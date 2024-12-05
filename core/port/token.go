package port

import (
	"context"
	"github.com/google/uuid"
	"medods-api/core/domain"
)

type TokenService interface {
	GetToken(context context.Context, userId uuid.UUID) (*domain.Token, error)
	RotateToken(context context.Context, refreshToken string) (*domain.Token, error)
}

type TokenData struct {
	RefreshToken string
	TokenFamily  string
	UserId       uuid.UUID
}

type TokenRepository interface {
	WriteRefreshToken(context context.Context, tokenData *TokenData) error
	CheckCorrectRefreshToken(context context.Context, refreshToken string, tokenFamily string) (bool, error)
	RevokeTokenFamily(context context.Context, tokenFamily string) error
	GetUserByRefreshToken(context context.Context, refreshToken string) (uuid.UUID, error)
}
