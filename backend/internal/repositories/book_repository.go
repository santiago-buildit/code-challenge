package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/santiago-buildit/code-challenge/backend/internal/models"
	"github.com/santiago-buildit/code-challenge/backend/internal/utils"
)

type BookRepository interface {

	// CRUD operations
	CreateBook(ctx context.Context, book *models.Book) error
	ListBooks(ctx context.Context, req models.ListBooksRequest) ([]models.Book, int /* total */, error)
	GetBookByID(ctx context.Context, id string) (*models.Book, error)
	UpdateBook(ctx context.Context, book *models.Book) error
	DeleteBook(ctx context.Context, id string) error

	// Status operations
	UpdateBookStatus(ctx context.Context, tx *sqlx.Tx, id string, status models.BookStatus, timestamp time.Time) error   // External TX
	AppendStatusChange(ctx context.Context, tx *sqlx.Tx, id string, status models.BookStatus, timestamp time.Time) error // External TX

	// History
	GetBookWithHistory(ctx context.Context, id string) (*models.Book, []models.BookStatusChange, error)
}

type bookRepositoryImpl struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) BookRepository {
	return &bookRepositoryImpl{
		db: db,
	}
}

func (r *bookRepositoryImpl) CreateBook(ctx context.Context, book *models.Book) error {

	// Execute insert
	_, err := r.db.NamedExecContext(ctx, `
		INSERT INTO books (
			id, isbn, title, author, description, status,
			created_at, updated_at, deleted
		) VALUES (
			:id, :isbn, :title, :author, :description, :status,
			:created_at, :updated_at, :deleted
		)
	`, book)
	return err
}

func (r *bookRepositoryImpl) ListBooks(ctx context.Context, req models.ListBooksRequest) ([]models.Book, int, error) {

	var (
		books      []models.Book
		args       []interface{}
		conditions []string
	)

	// Collect dynamic WHERE conditions (filters)
	if req.ISBN != "" {
		conditions = append(conditions, "isbn ILIKE ?")
		args = append(args, "%"+req.ISBN+"%")
	}
	if req.Title != "" {
		conditions = append(conditions, "title ILIKE ?")
		args = append(args, "%"+req.Title+"%")
	}
	if req.Author != "" {
		conditions = append(conditions, "author ILIKE ?")
		args = append(args, "%"+req.Author+"%")
	}
	if req.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, req.Status)
	}
	if req.Text != "" {
		conditions = append(conditions, "(title ILIKE ? OR description ILIKE ?)")
		args = append(args, "%"+req.Text+"%", "%"+req.Text+"%")
	}

	// Exclude deleted
	conditions = append(conditions, "deleted = false")

	// Build WHERE clause
	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Execute count query (for pagination)
	var total int
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM books %s`, where)
	countQuery = r.db.Rebind(countQuery) // Rebind converts '?' placeholders to PostgreSQL-style ($1, $2, ...)

	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, err
	}

	// Sanitize sort by
	sortBy := "title" // default
	if m := map[string]bool{"isbn": true, "title": true, "author": true, "status": true}; m[req.SortBy] {
		sortBy = req.SortBy
	}

	// Sanitize sort order
	sortOrder := "ASC"
	if strings.ToUpper(req.SortOrder) == "DESC" {
		sortOrder = "DESC"
	}

	// Pagination
	offset := (req.Page - 1) * req.PageSize

	// Build final query
	query := fmt.Sprintf(`
		SELECT * FROM books
		%s
		ORDER BY %s %s
		LIMIT %d OFFSET %d
	`, where, sortBy, sortOrder, req.PageSize, offset)
	query = r.db.Rebind(query) // Rebind converts '?' placeholders to PostgreSQL-style ($1, $2, ...)

	// Execute query
	if err := r.db.SelectContext(ctx, &books, query, args...); err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *bookRepositoryImpl) GetBookByID(ctx context.Context, id string) (*models.Book, error) {

	// Validate UUID format
	if err := validateUUIDOrNotFound(id); err != nil {
		return nil, err
	}

	// Execute query
	var book models.Book
	err := r.db.GetContext(ctx, &book, `
		SELECT * FROM books
		WHERE id = $1 AND deleted = false
	`, id)

	// Check for not found error
	if errors.Is(err, sql.ErrNoRows) {
		return nil, utils.ErrNotFound
	}
	return &book, err
}

func (r *bookRepositoryImpl) UpdateBook(ctx context.Context, book *models.Book) error {

	// Validate UUID format
	if err := validateUUIDOrNotFound(book.ID); err != nil {
		return err
	}

	// Execute update
	res, err := r.db.NamedExecContext(ctx, `
		UPDATE books SET
			isbn = :isbn,
			title = :title,
			author = :author,
			description = :description,
			status = :status,
			updated_at = :updated_at
		WHERE id = :id AND deleted = false
	`, book)
	if err != nil {
		return err
	}

	// Check for not found error
	return utils.CheckRowsAffected(res)
}

func (r *bookRepositoryImpl) DeleteBook(ctx context.Context, id string) error {

	// Validate UUID format
	if err := validateUUIDOrNotFound(id); err != nil {
		return err
	}

	// Execute update (logical delete)
	res, err := r.db.ExecContext(ctx, `
		UPDATE books SET deleted = true WHERE id = $1 AND deleted = false
	`, id)
	if err != nil {
		return err
	}

	// Check for not found error
	return utils.CheckRowsAffected(res)
}

func (r *bookRepositoryImpl) UpdateBookStatus(ctx context.Context, tx *sqlx.Tx, id string, status models.BookStatus, timestamp time.Time) error {

	// Validate UUID format
	if err := validateUUIDOrNotFound(id); err != nil {
		return err
	}

	// Execute update
	res, err := tx.ExecContext(ctx, `
		UPDATE books SET status = $1, updated_at = $2
		WHERE id = $3 AND deleted = false
	`, status, timestamp, id)
	if err != nil {
		return err
	}

	// Check for not found error
	return utils.CheckRowsAffected(res)
}

func (r *bookRepositoryImpl) AppendStatusChange(ctx context.Context, tx *sqlx.Tx, id string, status models.BookStatus, timestamp time.Time) error {

	// Validate UUID format
	if err := validateUUIDOrNotFound(id); err != nil {
		return err
	}

	// Execute insert
	_, err := tx.ExecContext(ctx, `
		INSERT INTO book_status_changes (book_id, status, timestamp)
		VALUES ($1, $2, $3)
	`, id, status, time.Now())
	return err
}

func (r *bookRepositoryImpl) GetBookWithHistory(ctx context.Context, id string) (*models.Book, []models.BookStatusChange, error) {

	// Validate UUID format
	if err := validateUUIDOrNotFound(id); err != nil {
		return nil, nil, err
	}

	// Execute query (book)
	book, err := r.GetBookByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}

	// Execute query (book status changes)
	var history []models.BookStatusChange
	err = r.db.SelectContext(ctx, &history, `
		SELECT status, timestamp
		FROM book_status_changes
		WHERE book_id = $1
		ORDER BY timestamp DESC
	`, id)
	if err != nil {
		return nil, nil, err
	}

	return book, history, nil
}

// validateUUIDOrNotFound checks if the given ID is a valid UUID.
func validateUUIDOrNotFound(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return utils.ErrNotFound
	}
	return nil
}
