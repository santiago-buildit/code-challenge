package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/santiago-buildit/code-challenge/backend/internal/models"
	"github.com/santiago-buildit/code-challenge/backend/internal/services"
	"github.com/santiago-buildit/code-challenge/backend/internal/utils"
	"go.uber.org/zap"
)

// Constants for pagination
const (
	defaultPageSize = 10
	maxPageSize     = 100
)

// Valid sort fields for pagination
var validSortFields = map[string]bool{
	"isbn":   true,
	"title":  true,
	"author": true,
	"status": true,
}

type BookHandler struct {
	service services.BookService
	logger  *zap.Logger
}

func NewBookHandler(service services.BookService, logger *zap.Logger) *BookHandler {
	return &BookHandler{
		service: service,
		logger:  logger,
	}
}

// CreateBook godoc
// @Summary Create a new book
// @Description Registers a new book with metadata
// @Tags books
// @Accept json
// @Produce json
// @Param request body models.CreateBookRequest true "Book data"
// @Success 201 {object} models.BookResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {

	h.logger.Info("Creating book")
	ctx := c.Request.Context()

	// Parse request body
	var req models.CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Sanitize input
	req.Sanitize()

	// Invoke service
	res, err := h.service.CreateBook(ctx, req)
	if err != nil {
		h.handleBookError(c, "", err, "create")
		return
	}
	h.logger.Info("Book created successfully",
		zap.String("id", res.ID),
		zap.String("isbn", res.ISBN),
		zap.String("title", res.Title),
	)
	c.JSON(http.StatusCreated, res)
}

// ListBooks godoc
// @Summary List books with filters, ordering, and pagination
// @Description Returns a paginated list of books. Supports filtering by ISBN, Title, Author, Status, and full-text search over Title/Description. Also supports ordering by field and direction.
// @Tags books
// @Accept json
// @Produce json
// @Param request body models.ListBooksRequest true "Filter and pagination parameters"
// @Success 200 {object} models.ListBooksResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/list [post]
func (h *BookHandler) ListBooks(c *gin.Context) {

	h.logger.Info("Listing books")
	ctx := c.Request.Context()

	// Parse request body
	var req models.ListBooksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Validate pagination parameters
	if req.PageSize <= 0 {
		req.PageSize = defaultPageSize
	}
	if req.PageSize > maxPageSize {
		req.PageSize = maxPageSize
	}
	if !validSortFields[req.SortBy] {
		req.SortBy = "title"
	}
	if req.SortOrder != "desc" {
		req.SortOrder = "asc"
	}

	// Invoke service
	res, err := h.service.ListBooks(ctx, req)
	if err != nil {
		h.logger.Error("Failed to list books", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to list books"})
		return
	}
	h.logger.Info("Books listed successfully", zap.Int("count", len(res.Books)))
	c.JSON(http.StatusOK, res)
}

// GetBook godoc
// @Summary Get a book by ID
// @Description Retrieves a book's metadata
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.BookResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBook(c *gin.Context) {

	h.logger.Info("Getting book")
	ctx := c.Request.Context()

	// Extract params
	id, ok := h.extractID(c)
	if !ok {
		return
	}

	// Invoke service
	res, err := h.service.GetBook(ctx, id)
	if err != nil {
		h.handleBookError(c, id, err, "get")
		return
	}
	h.logger.Info("Book retrieved successfully",
		zap.String("id", res.ID),
		zap.String("title", res.Title),
	)
	c.JSON(http.StatusOK, res)
}

// UpdateBook godoc
// @Summary Update a book by ID
// @Description Updates the metadata of a book
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param request body models.UpdateBookRequest true "Updated book data"
// @Success 200 {object} models.BookResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {

	h.logger.Info("Updating book")
	ctx := c.Request.Context()

	// Extract params
	id, ok := h.extractID(c)
	if !ok {
		return
	}

	// Parse request body
	var req models.UpdateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Sanitize input
	req.Sanitize()

	// Invoke service
	res, err := h.service.UpdateBook(ctx, id, req)
	if err != nil {
		h.handleBookError(c, id, err, "update")
		return
	}
	h.logger.Info("Book updated successfully",
		zap.String("id", res.ID),
		zap.String("isbn", res.ISBN),
		zap.String("title", res.Title),
	)
	c.JSON(http.StatusOK, res)
}

// DeleteBook godoc
// @Summary Delete a book by ID
// @Description Performs a logical delete on a book
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {

	h.logger.Info("Deleting book")
	ctx := c.Request.Context()

	// Extract params
	id, ok := h.extractID(c)
	if !ok {
		return
	}

	// Invoke service
	err := h.service.DeleteBook(ctx, id)
	if err != nil {
		h.handleBookError(c, id, err, "delete")
		return
	}
	h.logger.Info("Book deleted successfully", zap.String("id", id))
	c.JSON(http.StatusOK, models.MessageResponse{Message: "Book deleted"})
}

// CheckoutBook godoc
// @Summary Checkout a book by ID
// @Description Marks the book as checked out and updates history
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id}/checkout [put]
func (h *BookHandler) CheckoutBook(c *gin.Context) {

	h.logger.Info("Checking out book")
	ctx := c.Request.Context()

	// Extract params
	id, ok := h.extractID(c)
	if !ok {
		return
	}

	// Invoke service
	err := h.service.CheckoutBook(ctx, id)
	if err != nil {
		h.handleBookError(c, id, err, "checkout")
		return
	}
	h.logger.Info("Book checked out successfully", zap.String("id", id))
	c.JSON(http.StatusOK, models.MessageResponse{Message: "Book checked out"})
}

// CheckinBook godoc
// @Summary Checkin a book by ID
// @Description Marks the book as available and updates history
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.MessageResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id}/checkin [put]
func (h *BookHandler) CheckinBook(c *gin.Context) {

	h.logger.Info("Checking in book")
	ctx := c.Request.Context()

	// Extract params
	id, ok := h.extractID(c)
	if !ok {
		return
	}

	// Invoke service
	err := h.service.CheckinBook(ctx, id)
	if err != nil {
		h.handleBookError(c, id, err, "checkin")
		return
	}
	h.logger.Info("Book checked in successfully", zap.String("id", id))
	c.JSON(http.StatusOK, models.MessageResponse{Message: "Book checked in"})
}

// GetBookWithHistory godoc
// @Summary Get a book by ID with status change history
// @Description Retrieves book metadata and its full status change history
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} models.BookDetailResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/{id}/details [get]
func (h *BookHandler) GetBookWithHistory(c *gin.Context) {

	h.logger.Info("Getting book with history")
	ctx := c.Request.Context()

	// Extract params
	id, ok := h.extractID(c)
	if !ok {
		return
	}

	// Invoke service
	res, err := h.service.GetBookWithHistory(ctx, id)
	if err != nil {
		h.handleBookError(c, id, err, "get (with history)")
		return
	}
	h.logger.Info("Book and history retrieved", zap.String("id", id))
	c.JSON(http.StatusOK, res)
}

/* Helper functions */

func (h *BookHandler) extractID(c *gin.Context) (string, bool) {

	id := c.Param("id")
	if id == "" {
		h.logger.Warn("Missing parameter ID")
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Missing parameter ID"})
		return "", false
	}
	return id, true
}

func (h *BookHandler) handleBookError(c *gin.Context, id string, err error, action string) {

	// Handle specific errors
	if errors.Is(err, utils.ErrNotFound) { // Not found error
		h.logger.Warn("Book not found", zap.String("id", id))
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Book not found"})
	} else { // Generic error
		h.logger.Error("Failed to "+action+" book",
			zap.String("id", id),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to " + action + " book"})
	}
}
