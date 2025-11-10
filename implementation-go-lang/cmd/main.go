package main

import (
	"fmt"
	"os"

	"github.com/angelino-valeta/url-shortener-system-design/internal/config"
	"github.com/angelino-valeta/url-shortener-system-design/internal/handlers"
	"github.com/angelino-valeta/url-shortener-system-design/internal/repository"
	"github.com/angelino-valeta/url-shortener-system-design/internal/services"
	"github.com/angelino-valeta/url-shortener-system-design/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load env
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	// Init logger
	logger := utils.NewLogger(os.Getenv("LOG_LEVEL"))
	defer logger.Sync()

	// Load config
	cfg := config.LoadConfig(logger)

	// Init repositories
	cassRepo, err := repository.NewCassandraRepository(cfg.CassandraHosts, cfg.CassandraKeyspace, logger)
	if err != nil {
		logger.Fatal("Failed to init Cassandra", zap.Error(err))
	}
	redisRepo := repository.NewRedisRepository(cfg.RedisAddr, cfg.RedisPassword, logger)

	// Init services
	urlService := services.NewURLService(cassRepo, redisRepo, cfg.HashidsSalt, cfg.HashidsMinLength, cfg.CacheTTL, cfg.BaseURL, logger)

	// Init router
	r := gin.Default()
	r.POST("/api/v1/shorten", handlers.ShortenURL(urlService))
	r.GET("/:shortcode", handlers.RedirectURL(urlService))

	// Run server
	logger.Info("Starting server", zap.String("port", cfg.AppPort))
	if err := r.Run(":" + cfg.AppPort); err != nil {
		logger.Fatal("Server failed", zap.Error(err))
	}
}
