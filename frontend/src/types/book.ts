
// Generic success response
export interface MessageResponse {
  message: string
}

// Generic error response
export interface ErrorResponse {
  error: string
}

// Common request for CreateBook and UpdateBook
export interface BookPayload {
  isbn: string
  title: string
  author: string
  description?: string
}

// Common response for CreateBook, GetBook, UpdateBook, and ListBooks (within ListBooksResponse)
export interface BookResponse extends BookPayload {
  id: string
  status: 'available' | 'checked_out'
  created_at: string
  updated_at: string
}

export interface ListBooksRequest {
  page: number
  page_size: number
  sort_by?: 'isbn' | 'title' | 'author' | 'status'
  sort_order?: 'asc' | 'desc'
  isbn?: string
  title?: string
  author?: string
  status?: string
  text?: string
}

export interface ListBooksResponse {
  books: BookResponse[]
  total_items: number
  total_pages: number
  current_page: number
  page_size: number
}

export interface StatusChangeResponse {
  status: 'available' | 'checked_out'
  timestamp: string
}

export interface BookDetailResponse {
  book: BookResponse
  history: StatusChangeResponse[]
}
