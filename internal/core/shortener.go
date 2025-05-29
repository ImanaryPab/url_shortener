package core

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/ImanaryPab/url-shortener/internal/storage"
)

const (
	cacheTTL = 24 * time.Hour
)

type URLShortener struct {
	storage storage.URLStorage
	cache   storage.URLCache
}

func NewURLShortener(storage storage.URLStorage, cache storage.URLCache) *URLShortener {
	rand.Seed(time.Now().UnixNano())
	return &URLShortener{
		storage: storage,
		cache:   cache,
	}
}

func (s *URLShortener) ShortenURL(ctx context.Context, originalURL string) (string, error) {
	shortCode := generateShortCode(6)

	for i := 0; i < 3; i++ {
		err := s.storage.Create(ctx, originalURL, shortCode)
		if err == nil {
			_ = s.cache.Set(ctx, shortCode, originalURL, cacheTTL)
			return shortCode, nil
		}
		shortCode = generateShortCode(6)
	}
	return "", errors.New("failed to generate unique short code")
}

func (s *URLShortener) GetOriginalURL(ctx context.Context, shortCode string) (string, error) {
	cachedURL, err := s.cache.Get(ctx, shortCode)
	if err == nil && cachedURL != "" {
		return cachedURL, nil
	}

	originalURL, err := s.storage.GetOriginalURL(ctx, shortCode)
	if err != nil {
		return "", err
	}

	_ = s.cache.Set(ctx, shortCode, originalURL, cacheTTL)

	return originalURL, nil
}

func (s *URLShortener) IncrementAccessCount(ctx context.Context, shortCode string) error {
	return s.storage.IncrementAccessCount(ctx, shortCode)
}

func generateShortCode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
