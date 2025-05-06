package models

// Map Book to BookResponse
func ToBookResponse(book *Book) *BookResponse {
	return &BookResponse{
		ID:          book.ID,
		ISBN:        book.ISBN,
		Title:       book.Title,
		Author:      book.Author,
		Description: book.Description,
		Status:      book.Status,
		CreatedAt:   book.CreatedAt,
		UpdatedAt:   book.UpdatedAt,
	}
}

// Map Book[] to BookResponse[]
func ToBookResponseList(books []Book) []BookResponse {
	responses := make([]BookResponse, 0, len(books))
	for _, book := range books {
		responses = append(responses, *ToBookResponse(&book))
	}
	return responses
}

// Map BookStatusChange to StatusChangeResponse
func ToStatusChangeResponse(sc BookStatusChange) StatusChangeResponse {
	return StatusChangeResponse{
		Status:    sc.Status,
		Timestamp: sc.Timestamp,
	}
}

// Map BookStatusChange[] to StatusChangeResponse[]
func ToStatusChangeResponseList(changes []BookStatusChange) []StatusChangeResponse {
	result := make([]StatusChangeResponse, 0, len(changes))
	for _, sc := range changes {
		result = append(result, ToStatusChangeResponse(sc))
	}
	return result
}
