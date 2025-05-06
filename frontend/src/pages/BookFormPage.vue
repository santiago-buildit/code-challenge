<template>
  <div class="form-wrapper">
    <h2>{{ pageTitle }}</h2>
    <div v-if="mode !== 'new'" class="mode-toggle">
      <BaseButton variant="secondary" @click="toggleMode">
        Switch to {{ isViewMode ? 'Edit' : 'View' }}
      </BaseButton>
    </div>

    <div class="overlay-container">
      <div class="overlay" v-if="isLoading">
        <div class="spinner"></div>
      </div>
      <BookForm v-model="form" :mode="mode" @submit="handleSubmit" @cancel="goBack" />
    </div>

    <BookHistory v-if="isViewMode && form && form.status" :status="form.status" :history="history" />

    <BaseDialog v-if="dialog.visible" :title="dialog.title" :message="dialog.message" :type="dialog.type"
      :confirmText="dialog.confirmText" :cancelText="dialog.cancelText" @confirm="dialog.onConfirm"
      @cancel="dialog.onCancel" />
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import { useDialog } from '@/composables/useDialog'
import BaseDialog from '@/components/base/BaseDialog.vue'
import BookForm from '@/components/books/BookForm.vue'
import BookHistory from '@/components/books/BookHistory.vue'
import BaseButton from '@/components/base/BaseButton.vue'
import {
  createBook,
  updateBook,
  getBook,
  getBookWithHistory,
} from '@/services/bookService'
import type { BookResponse, StatusChangeResponse } from '@/types/book'

export default defineComponent({
  name: 'BookFormPage',
  setup() {
    const { dialog, showDialog } = useDialog()
    return { dialog, showDialog }
  },
  components: { BookForm, BookHistory, BaseButton, BaseDialog },
  data() {
    return {
      isLoading: false,
      form: {
        isbn: '',
        title: '',
        author: '',
        description: '',
      } as BookResponse,
      history: [] as StatusChangeResponse[],
    }
  },
  computed: {
    routeId(): string | null {
      return this.$route.params.id as string || null
    },
    mode(): 'new' | 'edit' | 'view' {
      const raw = this.$route.query.mode as string
      return raw === 'edit' ? 'edit' : raw === 'view' ? 'view' : 'new'
    },
    isViewMode(): boolean {
      return this.mode === 'view'
    },
    isEditMode(): boolean {
      return this.mode === 'edit'
    },
    pageTitle(): string {
      return this.isViewMode
        ? 'View Book'
        : this.isEditMode
          ? 'Edit Book'
          : 'Add Book'
    },
  },
  created() { // Lifecycle hook to load book data when the component is created
    if (this.isEditMode || this.isViewMode) {
      this.loadBook()
    }
  },
  methods: {
    async loadBook() { // Load book data from API
      this.isLoading = true
      try {
        if (!this.routeId) return
        if (this.isEditMode) {

          try {
            const book = await getBook(this.routeId) // Get Book from API
            this.form = { ...book }
          } catch (err) {
            const message = (err as any).response?.data?.error || 'An unexpected error occurred.'
            this.showDialog({
              title: 'An error occurred',
              message,
              type: 'error',
              confirmText: 'OK',
            })
          }
        } else if (this.isViewMode) {
          try {
            const res = await getBookWithHistory(this.routeId) // Get Book with History from API
            this.form = { ...res.book }
            this.history = res.history
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
      } finally {
        this.isLoading = false
      }
    },
    async handleSubmit() { // Handle form submission
      if (this.isEditMode && this.routeId) {

        try {
          await updateBook(this.routeId, this.form) // Update Book in API
          this.goBack()
        } catch (err) {
          const message = (err as any).response?.data?.error || 'An unexpected error occurred.'
          this.showDialog({
            title: 'An error occurred',
            message,
            type: 'error',
            confirmText: 'OK',
          })
        }

      } else {
        try {
          await createBook(this.form) // Create Book in API
          this.goBack()
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
    },
    toggleMode() { // Switch between view and edit mode
      const nextMode = this.isViewMode ? 'edit' : 'view'
      this
        .$router.replace({ path: `/book/${this.routeId}`, query: { mode: nextMode } })
    },
    goBack() { // Go back to the previous page
      if (window.history.length > 1) {
        this.$router.back() // Back to listing restoring filters from query params
      } else {
        this.$router.push({ path: '/' }) // Back to listing without filters
      }
    }
  },
})
</script>

<style scoped>
.form-wrapper {
  max-width: 500px;
  padding-bottom: 3rem;
}

.mode-toggle {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1rem;
}

.actions {
  margin-top: 2rem;
}
</style>
