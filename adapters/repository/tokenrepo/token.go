package tokenrepo

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"medods-api/adapters/repository/model"
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
	hashedToken, err := r.hashToken(tokenData.RefreshTokenId)
	if err != nil {
		return err
	}

	dbToken := model.Token{
		HashedRefreshTokenId: hashedToken,
		TokenFamily:          tokenData.TokenFamily,
	}

	return r.db.WithContext(context).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "token_family"}},
		UpdateAll: true,
	}).Create(&dbToken).Error
}

func (r *TokenRepo) CheckCorrectGenerationRefreshToken(context context.Context, refreshTokenId string, tokenFamily string) (bool, error) {
	var dbToken model.Token

	err := r.db.WithContext(context).Take(&dbToken, "token_family = ?", tokenFamily).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, errors.New("no such family in db")
	}

	correctToken := bcrypt.CompareHashAndPassword([]byte(dbToken.HashedRefreshTokenId), []byte(refreshTokenId))
	if correctToken != nil {
		return false, nil
	}

	return true, nil
}

func (r *TokenRepo) RevokeTokenFamily(context context.Context, tokenFamily string) error {
	dbToken := model.Token{TokenFamily: tokenFamily}

	return r.db.WithContext(context).Delete(&dbToken).Error
}

func (r *TokenRepo) hashToken(token string) (string, error) {
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("failed token hashing")
	}

	return string(hashedToken), nil
}
