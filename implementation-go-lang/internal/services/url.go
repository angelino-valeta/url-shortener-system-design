package services

import (
	"context"
	"fmt"
	"time"

	"github.com/angelino-valeta/url-shortener-system-design/internal/models"
	"github.com/angelino-valeta/url-shortener-system-design/internal/repository"
	"github.com/angelino-valeta/url-shortener-system-design/internal/utils"
	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"
)

// URLService business logic
type URLService struct {
	cassRepo  *repository.CassandraRepository
	redisRepo *repository.RedisRepository
	hashids   *hashids.HashID
	baseURL   string
	cacheTTL  int
	Logger    *zap.Logger
}

// NewURLService
func NewURLService(cassRepo *repository.CassandraRepository, redisRepo *repository.RedisRepository, salt string, minLength, cacheTTL int, baseURL string, logger *zap.Logger) *URLService {
	hd := hashids.NewData()
	hd.Salt = salt
	hd.MinLength = minLength
	hd.Alphabet = utils.Base62Alphabet
	h, _ := hashids.NewWithData(hd) // Ignore error for simplicity

	return &URLService{
		cassRepo:  cassRepo,
		redisRepo: redisRepo,
		hashids:   h,
		baseURL:   baseURL,
		cacheTTL:  cacheTTL,
		Logger:    logger,
	}
}

// ShortenURL logic
func (s *URLService) ShortenURL(ctx context.Context, longURL string) (string, error) {
	// Get unique ID from Redis
	id, err := s.redisRepo.IncrementCounter(ctx)
	if err != nil {
		return "", err
	}

	// Encode to shortcode
	shortcode, err := s.hashids.EncodeInt64([]int64{id})
	if err != nil {
		return "", err
	}

	// Check for collision (rare, but handle)
	existing, _ := s.cassRepo.GetURL(ctx, shortcode)
	if existing != nil {
		// Regenerate (simple retry, max 5)
		for i := 0; i < 5; i++ {
			id, err = s.redisRepo.IncrementCounter(ctx)
			if err != nil {
				return "", err
			}
			shortcode, err = s.hashids.EncodeInt64([]int64{id})
			if err != nil {
				return "", err
			}
			existing, _ = s.cassRepo.GetURL(ctx, shortcode)
			if existing == nil {
				break
			}
		}
		if existing != nil {
			return "", fmt.Errorf("collision retry limit exceeded")
		}
	}

	// Save to DB
	url := models.URL{
		Shortcode: shortcode,
		LongURL:   longURL,
		CreatedAt: time.Now(),
	}
	if err := s.cassRepo.SaveURL(ctx, url); err != nil {
		return "", err
	}

	return s.baseURL + "/" + shortcode, nil
}

// GetLongURL with cache
func (s *URLService) GetLongURL(ctx context.Context, shortcode string) (string, error) {
	// Check cache
	longURL, err := s.redisRepo.GetCachedURL(ctx, shortcode)
	if err == nil {
		return longURL, nil
	}

	// Fetch from DB
	url, err := s.cassRepo.GetURL(ctx, shortcode)
	if err != nil || url == nil {
		return "", err
	}

	// Set cache
	_ = s.redisRepo.SetCachedURL(ctx, shortcode, url.LongURL, s.cacheTTL)

	return url.LongURL, nil
}
