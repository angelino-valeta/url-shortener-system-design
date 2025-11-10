package config

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

// Config holds application configuration
type Config struct {
	AppPort           string
	CassandraHosts    []string
	CassandraKeyspace string
	RedisAddr         string
	RedisPassword     string
	HashidsSalt       string
	HashidsMinLength  int
	BaseURL           string
	CacheTTL          int // seconds
}

// LoadConfig loads from env
func LoadConfig(logger *zap.Logger) Config {
	hosts := []string{os.Getenv("CASSANDRA_HOSTS")}
	if hosts[0] == "" {
		hosts = []string{"localhost"}
	}

	minLength, _ := strconv.Atoi(os.Getenv("HASHIDS_MIN_LENGTH"))
	if minLength == 0 {
		minLength = 7
	}

	ttl, _ := strconv.Atoi(os.Getenv("CACHE_TTL"))
	if ttl == 0 {
		ttl = 86400
	}

	cfg := Config{
		AppPort:           os.Getenv("APP_PORT"),
		CassandraHosts:    hosts,
		CassandraKeyspace: os.Getenv("CASSANDRA_KEYSPACE"),
		RedisAddr:         os.Getenv("REDIS_ADDR"),
		RedisPassword:     os.Getenv("REDIS_PASSWORD"),
		HashidsSalt:       os.Getenv("HASHIDS_SALT"),
		HashidsMinLength:  minLength,
		BaseURL:           os.Getenv("BASE_URL"),
		CacheTTL:          ttl,
	}

	if cfg.AppPort == "" {
		cfg.AppPort = "8080"
	}
	if cfg.CassandraKeyspace == "" {
		cfg.CassandraKeyspace = "url_shortener"
	}
	if cfg.RedisAddr == "" {
		cfg.RedisAddr = "localhost:6379"
	}
	if cfg.HashidsSalt == "" {
		logger.Warn("HASHIDS_SALT not set, using default (insecure for prod)")
		cfg.HashidsSalt = "default_salt"
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://bit.ly"
	}

	return cfg
}
