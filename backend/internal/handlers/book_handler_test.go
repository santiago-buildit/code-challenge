package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/santiago-buildit/code-challenge/backend/internal/handlers"
	"github.com/santiago-buildit/code-challenge/backend/internal/models"
	"github.com/santiago-buildit/code-challenge/backend/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// MockBookService implements BookService for testing
type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(ctx context.Context, req models.CreateBookRequest) (*models.BookResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.BookResponse), args.Error(1)
}

func (m *MockBookService) ListBooks(ctx context.Context, req models.ListBooksRequest) (*models.ListBooksResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.ListBooksResponse), args.Error(1)
}
func (m *MockBookService) GetBook(ctx context.Context, id string) (*models.BookResponse, error) {
	args := m.Called(ctx, id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.BookResponse), args.Error(1)
}
func (m *MockBookService) UpdateBook(ctx context.Context, id string, req models.UpdateBookRequest) (*models.BookResponse, error) {
	args := m.Called(ctx, id, req)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.BookResponse), args.Error(1)
}
func (m *MockBookService) DeleteBook(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockBookService) CheckoutBook(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockBookService) CheckinBook(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockBookService) GetBookWithHistory(ctx context.Context, id string) (*models.BookDetailResponse, error) {
	args := m.Called(ctx, id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*models.BookDetailResponse), args.Error(1)
}

func TestCreateBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.POST("/books", handler.CreateBook)

	// Define and sanitize the payload (as the handler does)
	payload := models.CreateBookRequest{
		ISBN:        "123456",
		Title:       "The Lord of the Rings", // To test the sanitization
		Author:      "J.R.R. Tolkien",
		Description: "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them",
	}
	payload.Sanitize()

	// Mock response
	mockResp := &models.BookResponse{
		ID:          "book-1",
		ISBN:        payload.ISBN,
		Title:       payload.Title,
		Author:      payload.Author,
		Description: payload.Description,
		Status:      models.BookStatusAvailable,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Mock configuration
	mockSvc.On("CreateBook", mock.Anything, payload).Return(mockResp, nil)

	// Build request
	requestBody := models.CreateBookRequest{
		ISBN:        "123456",
		Title:       "The Lord of the Rings",
		Author:      "J.R.R. Tolkien",
		Description: "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them",
	}
	body, _ := json.Marshal(requestBody)

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var decoded models.BookResponse
	err := json.Unmarshal(resp.Body.Bytes(), &decoded)
	assert.NoError(t, err)
	assert.Equal(t, mockResp.Title, decoded.Title)

	mockSvc.AssertExpectations(t)
}

func TestCreateBook_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.POST("/books", handler.CreateBook)

	body := []byte(`{"isbn": "", "title": "", "author": ""}`) // Empty required fields

	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestListBooks_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.POST("/books/list", handler.ListBooks)

	reqBody := models.ListBooksRequest{
		Page:      1,
		PageSize:  5,
		SortBy:    "title",
		SortOrder: "asc",
		Title:     "The Lord of the Rings",
	}

	expectedResp := &models.ListBooksResponse{
		Books: []models.BookResponse{
			{
				ID:     "book-1",
				ISBN:   "001",
				Title:  "The Lord of the Rings",
				Author: "J.R.R. Tolkien",
				Status: models.BookStatusAvailable,
			},
		},
		TotalItems:  1,
		TotalPages:  1,
		CurrentPage: 1,
		PageSize:    5,
	}

	mockSvc.On("ListBooks", mock.Anything, reqBody).Return(expectedResp, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/books/list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var decoded models.ListBooksResponse
	err := json.Unmarshal(resp.Body.Bytes(), &decoded)
	assert.NoError(t, err)
	assert.Len(t, decoded.Books, 1)
	assert.Equal(t, "The Lord of the Rings", decoded.Books[0].Title)

	mockSvc.AssertExpectations(t)
}

func TestListBooks_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.POST("/books/list", handler.ListBooks)

	// Page = 0 (invalid by binding:"min=1")
	body := []byte(`{"page":0, "page_size":5}`)

	req := httptest.NewRequest(http.MethodPost, "/books/list", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestGetBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.GET("/books/:id", handler.GetBook)

	bookID := "book-1"
	expected := &models.BookResponse{
		ID:     bookID,
		ISBN:   "111",
		Title:  "The Lord of the Rings",
		Author: "J.R.R. Tolkien",
		Status: models.BookStatusAvailable,
	}

	mockSvc.On("GetBook", mock.Anything, bookID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/books/"+bookID, nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var got models.BookResponse
	err := json.Unmarshal(resp.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, got.ID)
	assert.Equal(t, expected.Title, got.Title)

	mockSvc.AssertExpectations(t)
}

func TestGetBook_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.GET("/books/:id", handler.GetBook)

	bookID := "missing-book"
	mockSvc.On("GetBook", mock.Anything, bookID).Return(nil, utils.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/books/"+bookID, nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.PUT("/books/:id", handler.UpdateBook)

	bookID := "book-1"
	reqBody := models.UpdateBookRequest{
		ISBN:        "999",
		Title:       "The Lord of the Rings",
		Author:      "J.R.R. Tolkien",
		Description: "One Ring to rule them all, One Ring to find them, One Ring to bring them all and in the darkness bind them",
	}

	expected := &models.BookResponse{
		ID:          bookID,
		ISBN:        reqBody.ISBN,
		Title:       reqBody.Title,
		Author:      reqBody.Author,
		Description: reqBody.Description,
		Status:      models.BookStatusAvailable,
	}

	mockSvc.On("UpdateBook", mock.Anything, bookID, reqBody).Return(expected, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/books/"+bookID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var got models.BookResponse
	err := json.Unmarshal(resp.Body.Bytes(), &got)
	assert.NoError(t, err)
	assert.Equal(t, expected.Title, got.Title)

	mockSvc.AssertExpectations(t)
}

func TestUpdateBook_InvalidPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.PUT("/books/:id", handler.UpdateBook)

	body := []byte(`{"title": "", "author": ""}`) // Missing ISBN, blank fields

	req := httptest.NewRequest(http.MethodPut, "/books/book-1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUpdateBook_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.PUT("/books/:id", handler.UpdateBook)

	bookID := "missing-book"
	reqBody := models.UpdateBookRequest{
		ISBN:        "000",
		Title:       "Ghost Book",
		Author:      "Nobody",
		Description: "Does not exist",
	}

	mockSvc.On("UpdateBook", mock.Anything, bookID, reqBody).Return(nil, utils.ErrNotFound)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/books/"+bookID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.DELETE("/books/:id", handler.DeleteBook)

	bookID := "book-123"
	mockSvc.On("DeleteBook", mock.Anything, bookID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/books/"+bookID, nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteBook_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.DELETE("/books/:id", handler.DeleteBook)

	bookID := "not-found"
	mockSvc.On("DeleteBook", mock.Anything, bookID).Return(utils.ErrNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/books/"+bookID, nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestCheckoutBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.PUT("/books/:id/checkout", handler.CheckoutBook)

	bookID := "book-1"
	mockSvc.On("CheckoutBook", mock.Anything, bookID).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/books/"+bookID+"/checkout", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestCheckoutBook_NotFound(t *testing.T) {
	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.PUT("/books/:id/checkout", handler.CheckoutBook)

	bookID := "not-found"
	mockSvc.On("CheckoutBook", mock.Anything, bookID).Return(utils.ErrNotFound)

	req := httptest.NewRequest(http.MethodPut, "/books/"+bookID+"/checkout", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestCheckinBook_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.PUT("/books/:id/checkin", handler.CheckinBook)

	bookID := "book-1"
	mockSvc.On("CheckinBook", mock.Anything, bookID).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/books/"+bookID+"/checkin", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestCheckinBook_NotFound(t *testing.T) {
	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.PUT("/books/:id/checkin", handler.CheckinBook)

	bookID := "not-found"
	mockSvc.On("CheckinBook", mock.Anything, bookID).Return(utils.ErrNotFound)

	req := httptest.NewRequest(http.MethodPut, "/books/"+bookID+"/checkin", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetBookWithHistory_Success(t *testing.T) {
	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.GET("/books/:id/details", handler.GetBookWithHistory)

	bookID := "book-1"
	expected := &models.BookDetailResponse{
		Book: models.BookResponse{
			ID:     bookID,
			ISBN:   "123",
			Title:  "The Lord of the Rings",
			Author: "J.R.R. Tolkien",
			Status: models.BookStatusCheckedOut,
		},
		History: []models.StatusChangeResponse{
			{Status: models.BookStatusAvailable, Timestamp: time.Now()},
		},
	}

	mockSvc.On("GetBookWithHistory", mock.Anything, bookID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/books/"+bookID+"/details", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetBookWithHistory_NotFound(t *testing.T) {
	mockSvc := new(MockBookService)
	logger := zaptest.NewLogger(t)
	handler := handlers.NewBookHandler(mockSvc, logger)

	r := gin.New()
	r.GET("/books/:id/details", handler.GetBookWithHistory)

	bookID := "not-found"
	mockSvc.On("GetBookWithHistory", mock.Anything, bookID).Return(nil, utils.ErrNotFound)

	req := httptest.NewRequest(http.MethodGet, "/books/"+bookID+"/details", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	mockSvc.AssertExpectations(t)
}
