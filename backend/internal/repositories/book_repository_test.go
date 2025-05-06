package repositories_test

import (
	"context"
	"database/sql"
	"github.com/santiago-buildit/code-challenge/backend/internal/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/santiago-buildit/code-challenge/backend/internal/models"
	"github.com/santiago-buildit/code-challenge/backend/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestListBooks_FilterByTitle(t *testing.T) {
	// Setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	req := models.ListBooksRequest{
		Title:     "The Lord of the Rings",
		Page:      1,
		PageSize:  10,
		SortBy:    "title",
		SortOrder: "asc",
	}

	// Count query
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM books WHERE title ILIKE \$1 AND deleted = false`).
		WithArgs("%The Lord of the Rings%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Data query
	mock.ExpectQuery(`(?i)^SELECT \* FROM books WHERE title ILIKE .+ AND deleted = false ORDER BY title ASC LIMIT 10 OFFSET 0$`).
		WithArgs("%The Lord of the Rings%").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "isbn", "title", "author", "description", "status", "created_at", "updated_at", "deleted",
		}).AddRow(
			"1", "123456", "The Lord of the Rings", "J.R.R. Tolkien", "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them", "available",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			false,
		))

	// Invoke
	books, total, err := repo.ListBooks(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, 1, total)
	assert.Equal(t, "The Lord of the Rings", books[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListBooks_FilterByTitleAndStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	req := models.ListBooksRequest{
		Title:     "The Lord of the Rings",
		Status:    "available",
		Page:      1,
		PageSize:  10,
		SortBy:    "title",
		SortOrder: "asc",
	}

	mock.ExpectQuery(`(?i)^SELECT COUNT\(\*\) FROM books WHERE title ILIKE .+ AND status = .+ AND deleted = false$`).
		WithArgs("%The Lord of the Rings%", "available").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`(?i)^SELECT \* FROM books WHERE title ILIKE .+ AND status = .+ AND deleted = false ORDER BY title ASC LIMIT 10 OFFSET 0$`).
		WithArgs("%The Lord of the Rings%", "available").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "isbn", "title", "author", "description", "status", "created_at", "updated_at", "deleted",
		}).AddRow(
			"1", "123456", "The Lord of the Rings", "J.R.R. Tolkien", "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them", "available",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			false,
		))

	books, total, err := repo.ListBooks(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, 1, total)
	assert.Equal(t, "The Lord of the Rings", books[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListBooks_Pagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	req := models.ListBooksRequest{
		Page:      2,
		PageSize:  10,
		SortBy:    "title",
		SortOrder: "asc",
	}

	mock.ExpectQuery(`(?i)^SELECT COUNT\(\*\) FROM books WHERE deleted = false$`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(25))

	mock.ExpectQuery(`(?i)^SELECT \* FROM books WHERE deleted = false ORDER BY title ASC LIMIT 10 OFFSET 10$`).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "isbn", "title", "author", "description", "status", "created_at", "updated_at", "deleted",
		}).AddRow(
			"2", "789456", "The Lord of the Rings", "J.R.R. Tolkien", "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them", "checked_out",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			false,
		))

	books, total, err := repo.ListBooks(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, 25, total)
	assert.Equal(t, "The Lord of the Rings", books[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListBooks_FullTextSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	req := models.ListBooksRequest{
		Text:      "concurrency",
		Page:      1,
		PageSize:  10,
		SortBy:    "title",
		SortOrder: "asc",
	}

	mock.ExpectQuery(`(?i)^SELECT COUNT\(\*\) FROM books WHERE \(title ILIKE .+ OR description ILIKE .+\) AND deleted = false$`).
		WithArgs("%concurrency%", "%concurrency%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mock.ExpectQuery(`(?i)^SELECT \* FROM books WHERE \(title ILIKE .+ OR description ILIKE .+\) AND deleted = false ORDER BY title ASC LIMIT 10 OFFSET 0$`).
		WithArgs("%concurrency%", "%concurrency%").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "isbn", "title", "author", "description", "status", "created_at", "updated_at", "deleted",
		}).AddRow(
			"3", "000999", "The Lord of the Rings", "J.R.R. Tolkien", "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them", "available",
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			false,
		))

	books, total, err := repo.ListBooks(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, books, 1)
	assert.Equal(t, total, 1)
	assert.Equal(t, "The Lord of the Rings", books[0].Title)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookByID_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	bookID := "fac2b19c-e857-4d40-8233-8132b9759b55"

	mock.ExpectQuery(`(?i)^SELECT \* FROM books WHERE id = \$1 AND deleted = false$`).
		WithArgs(bookID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "isbn", "title", "author", "description", "status", "created_at", "updated_at", "deleted",
		}).AddRow(
			bookID, "987654321", "The Lord of the Rings", "J.R.R. Tolkien", "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them", "available",
			time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			false,
		))

	book, err := repo.GetBookByID(ctx, bookID)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, "The Lord of the Rings", book.Title)
	assert.Equal(t, "available", string(book.Status))
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetBookByID_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	bookID := "fac2b19c-e857-4d40-8233-8132b9759b55"

	// Simulate no rows for that ID
	mock.ExpectQuery(`(?i)^SELECT \* FROM books WHERE id = \$1 AND deleted = false$`).
		WithArgs(bookID).
		WillReturnError(sql.ErrNoRows)

	book, err := repo.GetBookByID(ctx, bookID)

	assert.Nil(t, book)
	assert.ErrorIs(t, err, utils.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBook_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	book := &models.Book{
		ID:          "fac2b19c-e857-4d40-8233-8132b9759b55",
		ISBN:        "1234567890",
		Title:       "The Lord of the Rings",
		Author:      "J.R.R. Tolkien",
		Description: "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them",
		Status:      models.BookStatusAvailable,
		UpdatedAt:   time.Now(),
	}

	mock.ExpectExec(`(?i)^UPDATE books SET`).
		WithArgs(book.ISBN, book.Title, book.Author, book.Description, book.Status, book.UpdatedAt, book.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repo.UpdateBook(ctx, book)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateBook_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	book := &models.Book{
		ID:          "fac2b19c-e857-4d40-8233-8132b9759b55",
		ISBN:        "000000",
		Title:       "Not found",
		Author:      "Ghost",
		Description: "Should not update",
		Status:      models.BookStatusAvailable,
		UpdatedAt:   time.Now(),
	}

	mock.ExpectExec(`(?i)^UPDATE books SET`).
		WithArgs(book.ISBN, book.Title, book.Author, book.Description, book.Status, book.UpdatedAt, book.ID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = repo.UpdateBook(ctx, book)

	assert.ErrorIs(t, err, utils.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBook_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	bookID := "fac2b19c-e857-4d40-8233-8132b9759b55"

	mock.ExpectExec(`(?i)^UPDATE books SET deleted = true WHERE id = \$1 AND deleted = false$`).
		WithArgs(bookID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err = repo.DeleteBook(ctx, bookID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteBook_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repositories.NewBookRepository(sqlxDB)

	ctx := context.Background()
	bookID := "fac2b19c-e857-4d40-8233-8132b9759b55"

	mock.ExpectExec(`(?i)^UPDATE books SET deleted = true WHERE id = \$1 AND deleted = false$`).
		WithArgs(bookID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err = repo.DeleteBook(ctx, bookID)

	assert.ErrorIs(t, err, utils.ErrNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}
