import axios from 'axios'
import type {
  BookPayload,
  BookResponse,
  ListBooksRequest,
  ListBooksResponse,
  BookDetailResponse,
  MessageResponse
} from '@/types/book'

/* This file contains the API calls related to books */

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000', // Backend URL property set in /.env file
})

// List
export async function listBooks(payload: ListBooksRequest): Promise<ListBooksResponse> {
  const res = await api.post('/books/list', payload)
  return res.data
}

// Create
export async function createBook(payload: BookPayload): Promise<BookResponse> {
  const res = await api.post('/books', payload)
  return res.data
}

// Get
export async function getBook(id: string): Promise<BookResponse> {
  const res = await api.get(`/books/${id}`)
  return res.data
}

// Update
export async function updateBook(id: string, payload: BookPayload): Promise<BookResponse> {
  const res = await api.put(`/books/${id}`, payload)
  return res.data
}

// Delete
export async function deleteBook(id: string): Promise<MessageResponse> {
  const res = await api.delete(`/books/${id}`)
  return res.data
}

// Checkout
export async function checkoutBook(id: string): Promise<MessageResponse> {
  const res = await api.put(`/books/${id}/checkout`)
  return res.data
}

// Checkin
export async function checkinBook(id: string): Promise<MessageResponse> {
  const res = await api.put(`/books/${id}/checkin`)
  return res.data
}

// Get with history
export async function getBookWithHistory(id: string): Promise<BookDetailResponse> {
  const res = await api.get(`/books/${id}/details`)
  return res.data
}
