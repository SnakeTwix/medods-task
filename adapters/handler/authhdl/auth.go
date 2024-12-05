package authhdl

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"medods-api/core/port"
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

	err := c.Bind(&requestData)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	userId, err := uuid.Parse(requestData.UserId)
	if err != nil {
		return c.String(http.StatusBadRequest, "userId is not a guid")
	}

	token, err := h.tokenService.GetToken(c.Request().Context(), userId)
	if err != nil {
		return c.String(http.StatusInternalServerError, "something went wrong")
	}

	return c.JSON(http.StatusOK, token)
}

type refreshTokenRequest struct {
	RefreshToken string `query:"refreshToken"`
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var requestData refreshTokenRequest

	err := c.Bind(&requestData)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	token, err := h.tokenService.RotateToken(c.Request().Context(), requestData.RefreshToken)
	if err != nil {
		return c.String(http.StatusInternalServerError, "something went wrong")
	}

	return c.JSON(http.StatusOK, token)
}
