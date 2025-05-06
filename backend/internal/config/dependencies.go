package config

import (
	"github.com/santiago-buildit/code-challenge/backend/internal/handlers"
	"github.com/santiago-buildit/code-challenge/backend/internal/repositories"
	"github.com/santiago-buildit/code-challenge/backend/internal/services"
)

// Dependencies holds all application dependencies
type Dependencies struct {
	BookHandler *handlers.BookHandler
}

// InitDependencies initializes and returns all dependencies
func InitDependencies() *Dependencies {

	// Initialize logger (ZAP)
	logger := NewLogger()

	// Initialize database client
	db := NewDatabase(logger)

	// Initialize repositories
	bookRepo := repositories.NewBookRepository(db)

	// Initialize services
	bookService := services.NewBookService(db, bookRepo)

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(bookService, logger)

	// Build dependencies holder
	return &Dependencies{
		BookHandler: bookHandler,
	}
}
