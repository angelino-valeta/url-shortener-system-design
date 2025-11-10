package handlers

import (
	"net/http"
	"net/url"

	"github.com/angelino-valeta/url-shortener-system-design/internal/models"
	"github.com/angelino-valeta/url-shortener-system-design/internal/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/validator.v2"
)

// ShortenURL handler
func ShortenURL(svc *services.URLService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.ShortenRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			svc.Logger.Error("Invalid request", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		if errs := validator.Validate(req); errs != nil {
			svc.Logger.Error("Validation failed", zap.Error(errs))
			c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
			return
		}

		// Validate URL format
		if _, err := url.ParseRequestURI(req.URL); err != nil {
			svc.Logger.Error("Invalid URL", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
			return
		}

		shortURL, err := svc.ShortenURL(c.Request.Context(), req.URL)
		if err != nil {
			svc.Logger.Error("Shorten failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
			return
		}

		c.JSON(http.StatusCreated, models.ShortenResponse{ShortURL: shortURL})
	}
}

// RedirectURL handler
func RedirectURL(svc *services.URLService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortcode := c.Param("shortcode")
		if shortcode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing shortcode"})
			return
		}

		longURL, err := svc.GetLongURL(c.Request.Context(), shortcode)
		if err != nil {
			svc.Logger.Error("Redirect failed", zap.Error(err), zap.String("shortcode", shortcode))
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		c.Redirect(http.StatusMovedPermanently, longURL)
	}
}
