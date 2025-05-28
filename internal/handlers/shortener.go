package handlers

import (
	"net/http"
	"sync"

	"github.com/ImanaryPab/url-shortener/internal/core"
	"github.com/labstack/echo/v4"
)

var (
	shortener *core.URLShortener
	once      sync.Once
)

type ShortenerHandler struct{}

func NewShortenerHandler() *ShortenerHandler {
	once.Do(func() {
		shortener = core.NewURLShortener()
	})
	return &ShortenerHandler{}
}

func (h *ShortenerHandler) ShortenURL(c echo.Context) error {
	originalURL := c.FormValue("url")
	if originalURL == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "URL is required",
		})
	}

	shortCode := shortener.ShortenURL(originalURL)
	shortURL := "http://" + c.Request().Host + "/" + shortCode

	return c.JSON(http.StatusOK, map[string]string{
		"original_url": originalURL,
		"short_url":    shortURL,
	})
}

func (h *ShortenerHandler) Redirect(c echo.Context) error {
	shortCode := c.Param("code")
	originalURL, exists := shortener.GetOriginalURL(shortCode)

	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Short URL not found",
		})
	}

	return c.Redirect(http.StatusMovedPermanently, originalURL)
}
