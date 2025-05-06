package models

import (
	"strings"
	"time"
)

/* Persistence */

type BookStatus string

const (
	BookStatusAvailable  BookStatus = "available"
	BookStatusCheckedOut BookStatus = "checked_out"
)

type Book struct {
	ID          string     `db:"id"` // Generated UUID
	ISBN        string     `db:"isbn"`
	Title       string     `db:"title"`
	Author      string     `db:"author"`
	Description string     `db:"description"`
	Status      BookStatus `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	Deleted     bool       `db:"deleted"` // Logical delete
}

type BookStatusChange struct {
	ID        string     `db:"id"`      // Generated UUID
	BookID    string     `db:"book_id"` // FK to Book.ID
	Status    BookStatus `db:"status"`
	Timestamp time.Time  `db:"timestamp"`
}

/* API */

// BookPayload is a common request for creating and updating books
type BookPayload struct {
	ISBN        string `json:"isbn" binding:"required,max=20"`
	Title       string `json:"title" binding:"required,max=255"`
	Author      string `json:"author" binding:"required,max=255"`
	Description string `json:"description" binding:"max=1000"`
}

type CreateBookRequest = BookPayload
type UpdateBookRequest = BookPayload

type ListBooksRequest struct {

	// Pagination
	Page      int    `json:"page" binding:"required,min=1"`      // 1-based index
	PageSize  int    `json:"page_size" binding:"required,min=1"` // items per page
	SortBy    string `json:"sort_by"`                            // isbn, title, author, status
	SortOrder string `json:"sort_order"`                         // asc / desc

	// Filters
	ISBN   string `json:"isbn" binding:"max=20"`
	Title  string `json:"title" binding:"max=255"`
	Author string `json:"author" binding:"max=255"`
	Status string `json:"status" binding:"max=20"`
	Text   string `json:"text" binding:"max=500"`
}

// BookResponse is a common request for creating, getting, updating, and listing books (within ListBooksResponse)
type BookResponse struct {
	ID          string     `json:"id"`
	ISBN        string     `json:"isbn"`
	Title       string     `json:"title"`
	Author      string     `json:"author"`
	Description string     `json:"description"`
	Status      BookStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ListBooksResponse struct {

	// Data
	Books []BookResponse `json:"books"`

	// Pagination
	TotalItems  int `json:"total_items"` // total items matching the filter
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

type StatusChangeResponse struct {
	Status    BookStatus `json:"status"`
	Timestamp time.Time  `json:"timestamp"`
}

type BookDetailResponse struct {
	Book    BookResponse           `json:"book"`
	History []StatusChangeResponse `json:"history"`
}

// Sanitize request fields
func (r *BookPayload) Sanitize() {
	r.ISBN = strings.TrimSpace(r.ISBN)
	r.Title = strings.TrimSpace(r.Title)
	r.Author = strings.TrimSpace(r.Author)
	r.Description = strings.TrimSpace(r.Description)
}
