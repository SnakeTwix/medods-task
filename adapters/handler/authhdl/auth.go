package authhdl

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"medods-api/core/port"
	"medods-api/util/base64"
	"net/http"
)

type AuthHandler struct {
	tokenService port.TokenService
}

func New(tokenService port.TokenService) *AuthHandler {
	return &AuthHandler{
		tokenService: tokenService,
	}
}

type getTokensRequest struct {
	UserId string `query:"userId"`
}

func (h *AuthHandler) GetTokens(c echo.Context) error {
	var requestData getTokensRequest

	// Default echo.Context.Bind() doesn't check for query strings in POST requests
	err := (&echo.DefaultBinder{}).BindQueryParams(c, &requestData)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	userId, err := uuid.Parse(requestData.UserId)
	if err != nil {
		return c.String(http.StatusBadRequest, "userId is not a guid")
	}

	token, err := h.tokenService.GetToken(c.Request().Context(), userId, c.RealIP())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Передаем в base64
	token.Refresh = base64.Encode(token.Refresh)
	return c.JSON(http.StatusOK, token)
}

type refreshTokenRequest struct {
	RefreshToken string `query:"refreshToken"`
	AccessToken  string `query:"accessToken"`
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var requestData refreshTokenRequest

	err := (&echo.DefaultBinder{}).BindQueryParams(c, &requestData)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	requestData.RefreshToken, err = base64.Decode(requestData.RefreshToken)
	if err != nil {
		return c.String(http.StatusBadRequest, "couldn't decode refresh token")
	}

	token, err := h.tokenService.RotateToken(c.Request().Context(), requestData.RefreshToken, requestData.AccessToken, c.RealIP())
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	token.Refresh = base64.Encode(token.Refresh)
	return c.JSON(http.StatusOK, token)
}
