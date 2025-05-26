package handlers

import (
    "net/http"
    
    "github.com/labstack/echo/v4"
)

type ShortenerHandler struct {
}

func NewShortenerHandler() *ShortenerHandler {
    return &ShortenerHandler{}
}

func (h *ShortenerHandler) ShortenURL(c echo.Context) error {
    url := c.FormValue("url")
    return c.JSON(http.StatusOK, map[string]string{
        "short_url": "http://localhost:8080/abc123",
        "original":  url,
    })
}

func (h *ShortenerHandler) Redirect(c echo.Context) error {
    code := c.Param("code")
    return c.Redirect(http.StatusMovedPermanently, "https://google.com")
}