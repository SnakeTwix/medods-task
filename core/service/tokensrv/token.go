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

func (s *TokenService) RotateToken(context context.Context, refreshToken string, accessToken string) (*domain.Token, error) {
	refreshTkn, err := parseRefreshToken(refreshToken, s.configService.GetRefreshTokenSignKey())
	if err != nil {
		return nil, err
	}

	accessTkn, err := parseAccessToken(accessToken, s.configService.GetAccessTokenSignKey())
	if err != nil {
		return nil, err
	}

	if accessTkn.linkerId != refreshTkn.linkerId {
		return nil, errors.New("access and refresh token pair do not match")
	}

	// Check if the token we got was the one we expected and not a previous one
	isTokenSupposedToBeUsed, err := s.tokenRepo.CheckCorrectRefreshToken(context, refreshToken, refreshTkn.TokenFamily.String())
	if !isTokenSupposedToBeUsed {
		err = s.tokenRepo.RevokeTokenFamily(context, refreshToken)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("previous token usage detected")
	}

	token, _, refreshTkn, err := s.newToken(accessTkn.UserGUID, refreshTkn)
	if err != nil {
		return nil, err
	}

	tokenData := port.TokenData{
		RefreshToken: token.Refresh,
		TokenFamily:  refreshTkn.TokenFamily.String(),
		UserId:       accessTkn.UserGUID,
	}

	// Update the token
	err = s.tokenRepo.WriteRefreshToken(context, &tokenData)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *TokenService) newToken(userId uuid.UUID, oldRefreshToken *refreshToken) (*domain.Token, *accessToken, *refreshToken, error) {
	linkerId := uuid.New()
	refreshTkn := newRefreshToken(linkerId)

	// Assign the family, if coming with an existing refreshToken
	if oldRefreshToken != nil {
		refreshTkn.TokenFamily = oldRefreshToken.TokenFamily
	}

	refreshSigned, err := refreshTkn.Sign(s.configService.GetRefreshTokenSignKey())
	if err != nil {
		return nil, nil, nil, err
	}

	accessTkn := newAccessToken(userId, linkerId)
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
