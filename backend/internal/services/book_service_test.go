package services

import (
	"context"
	"github.com/santiago-buildit/code-challenge/backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/santiago-buildit/code-challenge/backend/internal/models"
	"github.com/stretchr/testify/mock"
)

// --- Mock definition ---

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) CreateBook(ctx context.Context, book *models.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *mockRepo) ListBooks(ctx context.Context, req models.ListBooksRequest) ([]models.Book, int, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.Book), args.Int(1), args.Error(2)
}

func (m *mockRepo) GetBookByID(ctx context.Context, id string) (*models.Book, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Book), args.Error(1)
}

func (m *mockRepo) UpdateBook(ctx context.Context, book *models.Book) error {
	args := m.Called(ctx, book)
	return args.Error(0)
}

func (m *mockRepo) DeleteBook(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepo) UpdateBookStatus(ctx context.Context, tx *sqlx.Tx, id string, status models.BookStatus, ts time.Time) error {
	args := m.Called(ctx, tx, id, status, ts)
	return args.Error(0)
}

func (m *mockRepo) AppendStatusChange(ctx context.Context, tx *sqlx.Tx, id string, status models.BookStatus, ts time.Time) error {
	args := m.Called(ctx, tx, id, status, ts)
	return args.Error(0)
}

func (m *mockRepo) GetBookWithHistory(ctx context.Context, id string) (*models.Book, []models.BookStatusChange, error) {
	args := m.Called(ctx, id)

	var book *models.Book
	if b := args.Get(0); b != nil {
		book = b.(*models.Book)
	}

	var history []models.BookStatusChange
	if h := args.Get(1); h != nil {
		history = h.([]models.BookStatusChange)
	}

	return book, history, args.Error(2)
}

// --- Test ---

func TestCreateBook_Success(t *testing.T) {
	ctx := context.Background()

	// Setup
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	// Request payload
	req := models.CreateBookRequest{
		ISBN:        "123456",
		Title:       "The Lord of the Rings",
		Author:      "J.R.R. Tolkien",
		Description: "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them",
	}

	// Capture the book received by CreateBook to validate
	var capturedBook *models.Book
	mockedRepo.On("CreateBook", ctx, mock.MatchedBy(func(b *models.Book) bool {
		capturedBook = b
		return true
	})).Return(nil)

	// Execute
	resp, err := service.CreateBook(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.ID)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, models.BookStatusAvailable, resp.Status)

	assert.WithinDuration(t, time.Now(), resp.CreatedAt, time.Second)
	assert.WithinDuration(t, resp.CreatedAt, resp.UpdatedAt, time.Second)

	mockedRepo.AssertExpectations(t)

	// Additional: validate that what was passed to the repo has the expected data
	assert.Equal(t, resp.ID, capturedBook.ID)
	assert.Equal(t, resp.Status, capturedBook.Status)
}

func TestCreateBook_RepositoryError(t *testing.T) {
	ctx := context.Background()

	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	req := models.CreateBookRequest{
		ISBN:        "123456",
		Title:       "Failing Book",
		Author:      "Unknown",
		Description: "This should fail",
	}

	mockedRepo.On("CreateBook", ctx, mock.AnythingOfType("*models.Book")).Return(assert.AnError)

	resp, err := service.CreateBook(ctx, req)

	assert.Nil(t, resp)
	assert.Equal(t, assert.AnError, err)
	mockedRepo.AssertExpectations(t)
}

func TestListBooks_Success(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	req := models.ListBooksRequest{
		Title:     "The Lord of the Rings",
		Page:      1,
		PageSize:  10,
		SortBy:    "title",
		SortOrder: "asc",
	}

	book := models.Book{
		ID:     "book-1",
		ISBN:   "123456",
		Title:  "The Lord of the Rings",
		Author: "J.R.R. Tolkien",
		Status: models.BookStatusAvailable,
	}
	books := []models.Book{book}
	total := 1

	mockedRepo.On("ListBooks", ctx, req).Return(books, total, nil)

	resp, err := service.ListBooks(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, total, resp.TotalItems)
	assert.Len(t, resp.Books, 1)
	assert.Equal(t, book.ID, resp.Books[0].ID)

	mockedRepo.AssertExpectations(t)
}

func TestListBooks_RepoError(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	req := models.ListBooksRequest{
		Page: 1, PageSize: 10,
	}

	mockedRepo.On("ListBooks", ctx, req).Return([]models.Book(nil), 0, assert.AnError)

	resp, err := service.ListBooks(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	mockedRepo.AssertExpectations(t)
}

func TestGetBook_Success(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	expected := &models.Book{
		ID:     "book-1",
		ISBN:   "123",
		Title:  "The Lord of the Rings",
		Author: "J.R.R. Tolkien",
		Status: models.BookStatusAvailable,
	}

	mockedRepo.On("GetBookByID", ctx, "book-1").Return(expected, nil)

	resp, err := service.GetBook(ctx, "book-1")

	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
	assert.Equal(t, expected.Title, resp.Title)
	mockedRepo.AssertExpectations(t)
}

func TestGetBook_NotFound(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	mockedRepo.On("GetBookByID", ctx, "missing-id").Return((*models.Book)(nil), utils.ErrNotFound)

	resp, err := service.GetBook(ctx, "missing-id")

	assert.Nil(t, resp)
	assert.Equal(t, utils.ErrNotFound, err)
	mockedRepo.AssertExpectations(t)
}

func TestUpdateBook_Success(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	bookID := "book-1"
	existing := &models.Book{
		ID:     bookID,
		ISBN:   "111",
		Title:  "Old Title",
		Author: "Old Author",
		Status: models.BookStatusAvailable,
	}

	req := models.UpdateBookRequest{
		ISBN:        "222",
		Title:       "New Title",
		Author:      "New Author",
		Description: "Updated desc",
	}

	mockedRepo.On("GetBookByID", ctx, bookID).Return(existing, nil)
	mockedRepo.On("UpdateBook", ctx, mock.MatchedBy(func(b *models.Book) bool {
		return b.ISBN == req.ISBN && b.Title == req.Title && b.Author == req.Author && b.Description == req.Description
	})).Return(nil)

	result, err := service.UpdateBook(ctx, bookID, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Title, result.Title)
	mockedRepo.AssertExpectations(t)
}

func TestUpdateBook_NotFound(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	req := models.UpdateBookRequest{
		ISBN:        "222",
		Title:       "New Title",
		Author:      "New Author",
		Description: "Updated desc",
	}

	mockedRepo.On("GetBookByID", ctx, "missing-id").Return((*models.Book)(nil), utils.ErrNotFound)

	result, err := service.UpdateBook(ctx, "missing-id", req)

	assert.Nil(t, result)
	assert.Equal(t, utils.ErrNotFound, err)
	mockedRepo.AssertExpectations(t)
}

func TestDeleteBook_Success(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	bookID := "book-123"
	mockedRepo.On("DeleteBook", ctx, bookID).Return(nil)

	err := service.DeleteBook(ctx, bookID)

	assert.NoError(t, err)
	mockedRepo.AssertExpectations(t)
}

func TestDeleteBook_NotFound(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	bookID := "missing-book"
	mockedRepo.On("DeleteBook", ctx, bookID).Return(utils.ErrNotFound)

	err := service.DeleteBook(ctx, bookID)

	assert.Equal(t, utils.ErrNotFound, err)
	mockedRepo.AssertExpectations(t)
}

func TestGetBookWithHistory_Success(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	bookID := "book-1"
	book := &models.Book{
		ID:     bookID,
		ISBN:   "123",
		Title:  "The Lord of the Rings",
		Author: "J.R.R. Tolkien",
		Status: models.BookStatusCheckedOut,
	}
	history := []models.BookStatusChange{
		{Status: models.BookStatusAvailable, Timestamp: time.Now()},
		{Status: models.BookStatusCheckedOut, Timestamp: time.Now()},
	}

	mockedRepo.On("GetBookWithHistory", ctx, bookID).Return(book, history, nil)

	result, err := service.GetBookWithHistory(ctx, bookID)

	assert.NoError(t, err)
	assert.Equal(t, book.ID, result.Book.ID)
	assert.Len(t, result.History, 2)
	mockedRepo.AssertExpectations(t)
}

func TestGetBookWithHistory_NotFound(t *testing.T) {
	ctx := context.Background()
	mockedRepo := new(mockRepo)
	db := &sqlx.DB{}
	service := NewBookService(db, mockedRepo)

	bookID := "missing"
	mockedRepo.On("GetBookWithHistory", ctx, bookID).Return(nil, []models.BookStatusChange(nil), utils.ErrNotFound)

	result, err := service.GetBookWithHistory(ctx, bookID)

	assert.Nil(t, result)
	assert.Equal(t, utils.ErrNotFound, err)
	mockedRepo.AssertExpectations(t)
}
