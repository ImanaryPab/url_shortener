package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	handler, err := handlers.NewShortenerHandler(cfg)
	if err != nil {
		log.Fatal("Failed to create handler:", err)
	}

	e.POST("/shorten", handler.ShortenURL)
	e.GET("/:code", handler.Redirect)
	e.GET("/stats/:code", handler.GetStats)

	go func() {
		port := cfg.ServerPort
		if err := e.Start(":" + strconv.Itoa(port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
