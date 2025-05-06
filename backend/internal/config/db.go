package config

import (
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase(logger *zap.Logger) *sqlx.DB {

	// Build connection string
	connStr, safeStr := getConnectionString(logger)

	// Log connection (without password)
	logger.Info("Connecting to database", zap.String("dsn", safeStr))

	// Connect
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// Create tables and indexes (if not exists)
	createTables(db, logger)
	createIndexes(db, logger)

	return db
}

func getConnectionString(logger *zap.Logger) (full string, safe string) {

	// Get connection details from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")

	// Check variables are defined
	if host == "" || port == "" || name == "" || user == "" || pass == "" {
		logger.Fatal("Missing DB connection environment variables")
	}

	// Escape parameters that could break the connection string (e.g. symbol '#' in password)
	safeUser := url.QueryEscape(user)
	safePass := url.QueryEscape(pass)

	// Build connection string
	full = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", safeUser, safePass, host, port, name)
	safe = fmt.Sprintf("postgres://%s:****@%s:%s/%s", safeUser, host, port, name)

	return full, safe
}

func createTables(db *sqlx.DB, logger *zap.Logger) {

	// Define table creation queries
	queries := []string{
		`CREATE TABLE IF NOT EXISTS books (
			id UUID PRIMARY KEY,
			isbn TEXT NOT NULL,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			description TEXT,
			status TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			deleted BOOLEAN NOT NULL DEFAULT FALSE
		);`,
		`CREATE TABLE IF NOT EXISTS book_status_changes (
			id SERIAL PRIMARY KEY,
			book_id UUID REFERENCES books(id) ON DELETE CASCADE,
			status TEXT NOT NULL,
			timestamp TIMESTAMPTZ NOT NULL
		);`,
	}

	// Execute each query
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			logger.Fatal("failed to create table", zap.Error(err))
		}
	}
}

func createIndexes(db *sqlx.DB, logger *zap.Logger) {

	// Define index creation queries
	queries := []string{

		// Indexes for books table
		`CREATE INDEX IF NOT EXISTS idx_books_status ON books(status);`,
		`CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);`,
		`CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);`,
		`CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);`,

		// Index for book_status_changes lookup
		`CREATE INDEX IF NOT EXISTS idx_book_history_bookid_timestamp
			ON book_status_changes(book_id, timestamp DESC);`,
	}

	// Execute each query
	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			logger.Fatal("failed to create index", zap.Error(err))
		}
	}
}
