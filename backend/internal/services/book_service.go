package services

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/santiago-buildit/code-challenge/backend/internal/database"
	"github.com/santiago-buildit/code-challenge/backend/internal/models"
	"github.com/santiago-buildit/code-challenge/backend/internal/repositories"
)

// ProductService defines the interface for product-related operations
type BookService interface {

	// CRUD operations
	CreateBook(ctx context.Context, req models.CreateBookRequest) (*models.BookResponse, error)
	ListBooks(ctx context.Context, req models.ListBooksRequest) (*models.ListBooksResponse, error)
	GetBook(ctx context.Context, id string) (*models.BookResponse, error)
	UpdateBook(ctx context.Context, id string, req models.UpdateBookRequest) (*models.BookResponse, error)
	DeleteBook(ctx context.Context, id string) error

	// Status operations
	CheckoutBook(ctx context.Context, id string) error
	CheckinBook(ctx context.Context, id string) error

	// History
	GetBookWithHistory(ctx context.Context, id string) (*models.BookDetailResponse, error)
}

type bookServiceImpl struct {
	db   *sqlx.DB
	repo repositories.BookRepository
}

func NewBookService(db *sqlx.DB, repo repositories.BookRepository) BookService {
	return &bookServiceImpl{
		db:   db,
		repo: repo,
	}
}

func (s *bookServiceImpl) CreateBook(ctx context.Context, req models.CreateBookRequest) (*models.BookResponse, error) {

	now := time.Now()

	// Map request
	book := models.Book{
		ID:          uuid.New().String(), // Generate unique ID
		ISBN:        req.ISBN,
		Title:       req.Title,
		Author:      req.Author,
		Description: req.Description,
		Status:      models.BookStatusAvailable,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Create with repository
	err := s.repo.CreateBook(ctx, &book)
	if err != nil {
		return nil, err
	}

	// Map response
	return models.ToBookResponse(&book), nil
}

func (s *bookServiceImpl) ListBooks(ctx context.Context, req models.ListBooksRequest) (*models.ListBooksResponse, error) {

	// List with repository
	books, totalItems, err := s.repo.ListBooks(ctx, req)
	if err != nil {
		return nil, err
	}

	// Map response
	totalPages := int(math.Ceil(float64(totalItems) / float64(req.PageSize)))
	res := &models.ListBooksResponse{
		Books:       models.ToBookResponseList(books),
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
	}
	if res.TotalPages == 0 {
		res.TotalPages = 1 // 1 empty page
	}
	return res, nil
}

func (s *bookServiceImpl) GetBook(ctx context.Context, id string) (*models.BookResponse, error) {

	// Get with repository
	book, err := s.repo.GetBookByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Map response
	return models.ToBookResponse(book), nil
}

func (s *bookServiceImpl) UpdateBook(ctx context.Context, id string, req models.UpdateBookRequest) (*models.BookResponse, error) {

	// Get with repository
	book, err := s.repo.GetBookByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Map request
	book.ISBN = req.ISBN
	book.Title = req.Title
	book.Author = req.Author
	book.Description = req.Description

	// Update with repository
	err = s.repo.UpdateBook(ctx, book)
	if err != nil {
		return nil, err
	}

	// Map response
	return models.ToBookResponse(book), nil
}

func (s *bookServiceImpl) DeleteBook(ctx context.Context, id string) error {

	// Delete with repository
	return s.repo.DeleteBook(ctx, id)
}

func (s *bookServiceImpl) CheckoutBook(ctx context.Context, id string) error {

	// Change status to checked out
	return s.changeBookStatus(ctx, id, models.BookStatusCheckedOut)
}

func (s *bookServiceImpl) CheckinBook(ctx context.Context, id string) error {

	// Change status to available
	return s.changeBookStatus(ctx, id, models.BookStatusAvailable)
}

func (s *bookServiceImpl) GetBookWithHistory(ctx context.Context, id string) (*models.BookDetailResponse, error) {

	// Get with repository
	book, history, err := s.repo.GetBookWithHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	// Map response
	return &models.BookDetailResponse{
		Book:    *models.ToBookResponse(book),
		History: models.ToStatusChangeResponseList(history),
	}, nil
}

/* Helper functions */

func (s *bookServiceImpl) changeBookStatus(ctx context.Context, id string, status models.BookStatus) error {

	// Get with repository
	book, err := s.repo.GetBookByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if book is already in desired status (idempotence)
	if book.Status == status {
		return nil
	}

	now := time.Now() // Use same timestamp for book updated-at and status change timestamp

	// Transactional block
	return database.WithTransaction(ctx, s.db, func(tx *sqlx.Tx) error {

		// Update status with repository
		err = s.repo.UpdateBookStatus(ctx, tx, id, status, now)
		if err != nil {
			return err
		}

		// Append status change with repository
		return s.repo.AppendStatusChange(ctx, tx, id, status, now)
	})
}
