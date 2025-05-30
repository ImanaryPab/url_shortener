package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ImanaryPab/url-shortener/internal/core"
	"github.com/ImanaryPab/url-shortener/internal/storage/postgres"
	"github.com/ImanaryPab/url-shortener/internal/storage/redis"
	"github.com/ImanaryPab/url-shortener/pkg/config"
	"github.com/labstack/echo/v4"
)

type ShortenerHandler struct {
	shortener *core.URLShortener
	cfg       *config.Config
}

func NewShortenerHandler(cfg *config.Config) (*ShortenerHandler, error) {
	dbStorage, err := postgres.NewPostgresStorage(cfg)
	if err != nil {
		return nil, err
	}

	redisCache := redis.NewRedisCache(cfg)

	shortener := core.NewURLShortener(dbStorage, redisCache)

	return &ShortenerHandler{
		shortener: shortener,
		cfg:       cfg,
	}, nil
}

func (h *ShortenerHandler) ShortenURL(c echo.Context) error {
	originalURL := c.FormValue("url")
	if originalURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL is required"})
	}

	// Парсим TTL (срок жизни ссылки)
	ttl := 0
	if ttlParam := c.FormValue("ttl"); ttlParam != "" {
		var err error
		ttl, err = strconv.Atoi(ttlParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid TTL value"})
		}
	}

	shortCode, err := h.shortener.ShortenURL(c.Request().Context(), originalURL, ttl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to shorten URL"})
	}

	shortURL := "http://" + c.Request().Host + "/" + shortCode

	return c.JSON(http.StatusOK, map[string]string{
		"original_url": originalURL,
		"short_url":    shortURL,
		"short_code":   shortCode,
		"expires_in":   fmt.Sprintf("%d hours", ttl),
	})
}

func (h *ShortenerHandler) GetStats(c echo.Context) error {
	shortCode := c.Param("code")
	stats, err := h.shortener.GetStats(c.Request().Context(), shortCode)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "URL not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, stats)
}

func (h *ShortenerHandler) Redirect(c echo.Context) error {
	shortCode := c.Param("code")
	originalURL, err := h.shortener.GetOriginalURL(c.Request().Context(), shortCode)
	if err != nil {
		if err == storage.ErrNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Short URL not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error",
		})
	}

	go func() {
		ctx := context.Background()
		_ = h.shortener.IncrementAccessCount(ctx, shortCode)
	}()

	return c.Redirect(http.StatusMovedPermanently, originalURL)
}
