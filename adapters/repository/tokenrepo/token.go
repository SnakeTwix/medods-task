package tokenrepo

import (
	"context"
	"gorm.io/gorm"
	"medods-api/core/port"
)

type TokenRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (r *TokenRepo) WriteRefreshToken(context context.Context, tokenData *port.TokenData) error {
	return nil
}

func (r *TokenRepo) CheckCorrectRefreshToken(context context.Context, refreshToken string, tokenFamily string) (bool, error) {
	return true, nil
}

func (r *TokenRepo) RevokeTokenFamily(context context.Context, tokenFamily string) error {
	return nil
}
