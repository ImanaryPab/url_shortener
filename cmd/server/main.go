package main

import (
	"log"

	"github.com/ImanaryPab/url-shortener/internal/handlers"
	"github.com/ImanaryPab/url-shortener/pkg/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Инициализация обработчика
	handler := handlers.NewShortenerHandler()

	// Маршруты
	e.POST("/shorten", handler.ShortenURL)
	e.GET("/:code", handler.Redirect)

	// Запуск сервера
	port := cfg.ServerPort
	log.Printf("Starting server on :%d", port)
	e.Logger.Fatal(e.Start(":" + port))
}
