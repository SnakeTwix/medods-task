package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"medods-api/adapters/handler/authhdl"
	"medods-api/adapters/repository"
	"medods-api/adapters/repository/migration"
	"medods-api/adapters/repository/tokenrepo"
	"medods-api/core/service/configsrv"
	"medods-api/core/service/tokensrv"
	"medods-api/util/env"
)

func main() {
	db := repository.InitDB()
	err := migration.RunMigrations(db)
	if err != nil {
		log.Fatal("failed to run migrations", err)
	}

	// Repository initialization
	tokenRepo := tokenrepo.New(db)

	// Service initialization
	configService := configsrv.New()
	tokenService := tokensrv.New(configService, tokenRepo)

	// Handle initialization
	authHandler := authhdl.New(tokenService)

	e := echo.New()

	// Routes
	e.POST("/auth/login", authHandler.GetTokens)
	e.POST("/auth/refresh", authHandler.RefreshToken)

	apiAddress := fmt.Sprintf(":%s", env.Get(env.API_PORT))
	e.Logger.Fatal(e.Start(apiAddress))
}
