package repository

import (
	"context"

	"github.com/angelino-valeta/url-shortener-system-design/internal/models"
	"github.com/gocql/gocql"
	"go.uber.org/zap"
)

// CassandraRepository for DB ops
type CassandraRepository struct {
	session *gocql.Session
	logger  *zap.Logger
}

// NewCassandraRepository init
func NewCassandraRepository(hosts []string, keyspace string, logger *zap.Logger) (*CassandraRepository, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	// Create keyspace and table if not exists
	ctx := context.Background()
	if err := session.Query(`CREATE KEYSPACE IF NOT EXISTS ` + keyspace + ` WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 3}`).WithContext(ctx).Exec(); err != nil {
		logger.Warn("Keyspace creation failed", zap.Error(err))
	}
	if err := session.Query(`CREATE TABLE IF NOT EXISTS url (shortcode TEXT PRIMARY KEY, long_url TEXT, created_at TIMESTAMP)`).WithContext(ctx).Exec(); err != nil {
		return nil, err
	}

	return &CassandraRepository{session: session, logger: logger}, nil
}

// SaveURL
func (r *CassandraRepository) SaveURL(ctx context.Context, url models.URL) error {
	query := `INSERT INTO url (shortcode, long_url, created_at) VALUES (?, ?, ?)`
	return r.session.Query(query, url.Shortcode, url.LongURL, url.CreatedAt).WithContext(ctx).Exec()
}

// GetURL
func (r *CassandraRepository) GetURL(ctx context.Context, shortcode string) (*models.URL, error) {
	query := `SELECT shortcode, long_url, created_at FROM url WHERE shortcode = ?`
	var url models.URL
	if err := r.session.Query(query, shortcode).WithContext(ctx).Scan(&url.Shortcode, &url.LongURL, &url.CreatedAt); err != nil {
		if err == gocql.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &url, nil
}
