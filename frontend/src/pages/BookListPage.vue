<template>
  <div class="list-page">
    <h2>Books</h2>

    <BookFilter v-model="filters" @search="fetchBooks" @add="goToNewBook" @reset="resetFilters" />

    <div class="overlay-container">
      <div class="overlay" v-if="isLoading">
        <div class="spinner"></div>
      </div>
      <BookTable :books="books" :sort-by="filters.sort_by!" :sort-order="filters.sort_order!" @view="viewBook"
        @edit="editBook" @delete="deleteBookById" @sort="sort" @checkout="handleCheckout" @checkin="handleCheckin" />

    </div>
    <div class="pagination-wrapper">
      <div class="pagination">
        <button class="pagination-button" :disabled="pagination.current_page === 1"
          @click="goToPage(pagination.current_page - 1)">
          ‹ Prev
        </button>

        <span class="pagination-info">
          Page {{ pagination.current_page }} of {{ pagination.total_pages }}
        </span>

        <button class="pagination-button" :disabled="pagination.current_page === pagination.total_pages"
          @click="goToPage(pagination.current_page + 1)">
          Next ›
        </button>
      </div>
    </div>

    <BaseDialog v-if="dialog.visible" :title="dialog.title" :message="dialog.message" :type="dialog.type"
      :confirmText="dialog.confirmText" :cancelText="dialog.cancelText" @confirm="dialog.onConfirm"
      @cancel="dialog.onCancel" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useDialog } from '@/composables/useDialog'
import BaseDialog from '@/components/base/BaseDialog.vue'
import BookFilter from '@/components/books/BookFilter.vue'
import BookTable from '@/components/books/BookTable.vue'
import { listBooks, checkoutBook, checkinBook, deleteBook } from '@/services/bookService'
import type { BookResponse, ListBooksRequest, ListBooksResponse } from '@/types/book'

export default defineComponent({
  name: 'BookListPage',
  setup() {
    const { dialog, showDialog } = useDialog()
    return { dialog, showDialog }
  },
  components: {
    BaseDialog,
    BookFilter,
    BookTable,
  },
  data() {
    return {
      isLoading: false,
      books: [] as BookResponse[],
      filters: {
        page: 1,
        page_size: 10,
        sort_by: 'title',
        sort_order: 'asc',
        isbn: '',
        title: '',
        author: '',
        status: ''
      } as ListBooksRequest,
      pagination: {
        current_page: 1,
        total_pages: 1,
        total_items: 0,
        page_size: 10,
      },
    }
  },
  created() {
    const query = this.$route.query

    // Extract filters from query parameters (support navigating to the form and back to the list keeping the filters)
    this.filters = {
      page: parseInt(query.page as string) || 1,
      page_size: 10,
      sort_by: (query.sort_by as ListBooksRequest['sort_by']) || 'title',
      sort_order: (query.sort_order as 'asc' | 'desc') || 'asc',
      isbn: (query.isbn as string) || '',
      title: (query.title as string) || '',
      author: (query.author as string) || '',
      status: (query.status as string) || '',
    }

    // Fetch books when the component is created
    this.fetchBooks()
  },
  methods: {
    async fetchBooks() { // Fetch books from the API
      this.isLoading = true
      try {
        this.updateQueryParams() // Update query params based on filters

        const filtersToSend = { // Sanitize filters to send to the API
          ...this.filters,
          status: this.filters.status === 'all' ? '' : this.filters.status,
        }

        var result = {} as ListBooksResponse

        try {
          result = await listBooks(filtersToSend) // Call the API to get the list of books
        } catch (err) {
          const message = (err as any).response?.data?.error || 'An unexpected error occurred.'
          this.showDialog({
            title: 'An error occurred',
            message,
            type: 'error',
            confirmText: 'OK',
          })
        }

        // If the current page is greater than the total pages, set it to the last page
        // This usually happens when pagination is moved forward and a more restrictive filter is applied
        if (this.filters.page > result.total_pages && result.total_pages > 0) {
          this.filters.page = result.total_pages
          this.updateQueryParams()
          filtersToSend.page = this.filters.page

          try {
            result = await listBooks(filtersToSend) // Call the API again for the correct page
          } catch (err) {
            const message = (err as any).response?.data?.error || 'An unexpected error occurred.'
            this.showDialog({
              title: 'An error occurred',
              message,
              type: 'error',
              confirmText: 'OK',
            })
          }
        }

        this.books = result.books
        this.pagination.current_page = result.current_page
        this.pagination.total_pages = result.total_pages
        this.pagination.total_items = result.total_items
      } finally {
        this.isLoading = false
      }
    },
    async handleCheckout(id: string) { // Handle book checkout
      this.showDialog({ // Require confirmation
        title: 'Confirm Checkout',
        message: 'Do you want to check out this book?',
        onConfirm: async () => {
          try {
            await checkoutBook(id) // Call the API to check out the book
            await this.fetchBooks() // Refresh the book list
          } catch (err) {
            const message = (err as any).response?.data?.error || 'An unexpected error occurred.'
            this.showDialog({
              title: 'An error occurred',
              message,
              type: 'error',
              confirmText: 'OK',
            })
          }
        },
      })
    },
    async handleCheckin(id: string) { // Handle book checkin
      this.showDialog({ // Require confirmation
        title: 'Confirm Checkin',
        message: 'Do you want to check in this book?',
        onConfirm: async () => {
          try {
            await checkinBook(id) // Call the API to check in the book
            await this.fetchBooks() // Refresh the book list
          } catch (err) {
            const message = (err as any).response?.data?.error || 'An unexpected error occurred.'
            this.showDialog({
              title: 'An error occurred',
              message,
              type: 'error',
              confirmText: 'OK',
            })
          }
        },
      })
    },
    sort(field: 'title' | 'isbn' | 'author' | 'status') { // Handle sorting from table header
      this.updateQueryParams() // Update query params based on filters
      if (this.filters.sort_by === field) {
        this.filters.sort_order = this.filters.sort_order === 'asc' ? 'desc' : 'asc'
      } else {
        this.filters.sort_by = field
        this.filters.sort_order = 'asc'
      }
      this.fetchBooks() // Refresh the book list
    },
    goToPage(page: number) { // Handle pagination
      this.updateQueryParams()
      this.filters.page = page
      this.fetchBooks()
    },
    viewBook(id: string) { // Navigate to the form in view mode
      this.$router.push(`/book/${id}?mode=view`)
    },
    editBook(id: string) { // Navigate to the form in edit mode
      this.$router.push(`/book/${id}?mode=edit`)
    },
    goToNewBook() { // Navigate to the form in new mode
      this.$router.push('/book/new')
    },
    resetFilters() { // Reset filters to default values
      this.filters.page = 1
      this.filters.page_size = 10
      this.filters.sort_by = 'title'
      this.filters.sort_order = 'asc'
      this.filters.isbn = ''
      this.filters.title = ''
      this.filters.author = ''
      this.filters.status = ''
      this.updateQueryParams() // Update query params based on filters
      this.fetchBooks() // Refresh the book list
    },
    async deleteBookById(id: string) { // Handle book deletion
      this.showDialog({ // Require confirmation
        title: 'Confirm Deletion',
        message: 'Are you sure you want to delete this book?',
        confirmText: 'Delete',
        cancelText: 'Cancel',
        onConfirm: async () => {
          try {
            await deleteBook(id); // Call the API to delete the book
            await this.fetchBooks(); // Refresh the book list
          } catch (err) {
            const message = (err as any).response?.data?.error || 'An unexpected error occurred.'
            this.showDialog({
              title: 'An error occurred',
              message,
              type: 'error',
              confirmText: 'OK',
            })
          }
        },
      })
    },
    updateQueryParams() { // Update the query parameters in the URL based on the current filters
      const query: Record<string, string | undefined> = {
        page: String(this.filters.page),
        sort_by: this.filters.sort_by || undefined,
        sort_order: this.filters.sort_order || undefined,
        isbn: this.filters.isbn || undefined,
        title: this.filters.title || undefined,
        author: this.filters.author || undefined,
        status: this.filters.status || undefined,
      }

      this.$router.replace({ query })
    }
  },
})
</script>

<style scoped>
.list-page {
  max-width: 1000px;
}

.pagination-wrapper {
  margin-top: 0.5rem;
}

.pagination {
  display: flex;
  align-items: center;
  gap: 1rem;
  background-color: #f9f9f9;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.pagination-button {
  background-color: #477777;
  color: white;
  border: none;
  padding: 0.4rem 0.75rem;
  font-size: 0.95rem;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s ease;
}

.pagination-button:hover:not(:disabled) {
  background-color: #253e3e;
}

.pagination-button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}

.pagination-info {
  font-size: 0.95rem;
  color: #333;
}
</style>
