package tokensrv

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"medods-api/core/domain"
	"medods-api/core/port"
)

type TokenService struct {
	configService port.ConfigService
	tokenRepo     port.TokenRepository
}

func New(configService port.ConfigService, tokenRepo port.TokenRepository) *TokenService {
	return &TokenService{
		tokenRepo:     tokenRepo,
		configService: configService,
	}
}

func (s *TokenService) GetToken(context context.Context, userId uuid.UUID) (*domain.Token, error) {
	token, _, refreshTkn, err := s.newToken(userId, nil)
	if err != nil {
		return nil, err
	}

	tokenData := port.TokenData{
		RefreshToken: token.Refresh,
		TokenFamily:  refreshTkn.TokenFamily.String(),
		UserId:       userId,
	}
	err = s.tokenRepo.WriteRefreshToken(context, &tokenData)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *TokenService) RotateToken(context context.Context, refreshToken string) (*domain.Token, error) {
	refreshTkn, err := parseRefreshToken(refreshToken, s.configService.GetRefreshTokenSignKey())
	if err != nil {
		return nil, err
	}

	isTokenSupposedToBeUsed, err := s.tokenRepo.CheckCorrectRefreshToken(context, refreshToken, refreshTkn.TokenFamily.String())
	if !isTokenSupposedToBeUsed {
		err = s.tokenRepo.RevokeTokenFamily(context, refreshToken)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("previous token usage detected")
	}

	userId, err := s.tokenRepo.GetUserByRefreshToken(context, refreshToken)
	if err != nil {
		return nil, err
	}

	return s.GetToken(context, userId)
}

func (s *TokenService) newToken(userId uuid.UUID, refreshTokenEncoded *string) (*domain.Token, *accessToken, *refreshToken, error) {
	refreshTkn := newRefreshToken()

	// Assign the family, if coming with an existing refreshToken
	if refreshTokenEncoded != nil {
		oldRefreshToken, err := parseRefreshToken(*refreshTokenEncoded, s.configService.GetRefreshTokenSignKey())
		if err != nil {
			return nil, nil, nil, err
		}

		refreshTkn.TokenFamily = oldRefreshToken.TokenFamily
	}

	refreshSigned, err := refreshTkn.Sign(s.configService.GetRefreshTokenSignKey())
	if err != nil {
		return nil, nil, nil, err
	}

	accessTkn := newAccessToken(userId)
	accessSigned, err := accessTkn.Sign(s.configService.GetAccessTokenSignKey())
	if err != nil {
		return nil, nil, nil, err
	}

	token := domain.Token{
		Refresh: refreshSigned,
		Access:  accessSigned,
	}

	return &token, accessTkn, refreshTkn, nil
}
