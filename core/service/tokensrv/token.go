package tokensrv

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"log"
	"medods-api/core/domain"
	"medods-api/core/port"
)

type TokenService struct {
	configService   port.ConfigService
	notifierService port.NotifierService

	tokenRepo port.TokenRepository
}

func New(configService port.ConfigService, tokenRepo port.TokenRepository, notifierService port.NotifierService) *TokenService {
	return &TokenService{
		configService:   configService,
		notifierService: notifierService,

		tokenRepo: tokenRepo,
	}
}

func (s *TokenService) GetToken(context context.Context, userId uuid.UUID, ip string) (*domain.Token, error) {
	token, _, refreshTkn, err := s.newToken(userId, ip, nil)
	if err != nil {
		return nil, err
	}

	tokenData := port.TokenData{
		RefreshTokenId: refreshTkn.ID,
		TokenFamily:    refreshTkn.TokenFamily.String(),
	}
	err = s.tokenRepo.WriteRefreshToken(context, &tokenData)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *TokenService) RotateToken(context context.Context, refreshToken string, accessToken string, ip string) (*domain.Token, error) {
	refreshTkn, err := parseRefreshToken(refreshToken, s.configService.GetRefreshTokenSignKey())
	if err != nil {
		return nil, err
	}

	accessTkn, err := parseAccessToken(accessToken, s.configService.GetAccessTokenSignKey())
	if err != nil {
		return nil, err
	}

	if accessTkn.LinkerId != refreshTkn.LinkerId {
		return nil, errors.New("access and refresh token pair do not match")
	}

	// If the ip changed, notify the user
	if accessTkn.UserIp != ip {
		err := s.notifierService.NotifyUserIPChange(context, accessTkn.UserGUID, ip)
		if err != nil {
			log.Println("failed to notify: ", err)
		}
	}

	// Check if the token we got was the one we expected and not a previous one
	isTokenSupposedToBeUsed, err := s.tokenRepo.CheckCorrectGenerationRefreshToken(context, refreshTkn.ID, refreshTkn.TokenFamily.String())
	if err != nil {
		return nil, err
	}

	if !isTokenSupposedToBeUsed {
		err = s.tokenRepo.RevokeTokenFamily(context, refreshTkn.TokenFamily.String())
		if err != nil {
			return nil, err
		}

		return nil, errors.New("previous token usage detected")
	}

	token, _, refreshTkn, err := s.newToken(accessTkn.UserGUID, ip, refreshTkn)
	if err != nil {
		return nil, err
	}

	tokenData := port.TokenData{
		RefreshTokenId: refreshTkn.ID,
		TokenFamily:    refreshTkn.TokenFamily.String(),
	}

	// Update the token
	err = s.tokenRepo.WriteRefreshToken(context, &tokenData)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *TokenService) newToken(userId uuid.UUID, userIP string, oldRefreshToken *refreshToken) (*domain.Token, *accessToken, *refreshToken, error) {
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

	accessTkn := newAccessToken(userId, userIP, linkerId)
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
