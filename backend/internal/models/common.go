package models

// MessageResponse is a generic API success response
type MessageResponse struct {
	Message string `json:"message" example:"A success message"`
}

// ErrorResponse is a generic API error response
type ErrorResponse struct {
	Error string `json:"error" example:"An error message"`
}
