package main

import (
	"log"

	"github.com/ImanaryPab/url-shortener/pkg/config"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	e := echo.New()

	e.GET("/:code", redirectHandler)
	e.POST("/shorten", shortenHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func redirectHandler(c echo.Context) error {
	return c.String(200, "Redirect logic here")
}

func shortenHandler(c echo.Context) error {
	return c.String(200, "Shorten logic here")
}
