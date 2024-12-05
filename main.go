package main

import (
	"fmt"
	"medods-api/util/env"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	apiAddress := fmt.Sprintf(":%s", env.Get(env.API_PORT))
	e.Logger.Fatal(e.Start(apiAddress))
}
