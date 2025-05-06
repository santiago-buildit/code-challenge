<template>
  <div class="table-container">
    <table v-if="books && books.length > 0">
      <thead>
        <tr>
          <th class="col-isbn sortable hide-mobile" @click="sort('isbn')">
            ISBN
            <span v-if="sortBy === 'isbn'">{{ sortOrder === 'asc' ? '&#9650;' : '&#9660;' }}</span>
          </th>
          <th class="col-title sortable" @click="sort('title')">
            Title
            <span v-if="sortBy === 'title'">{{ sortOrder === 'asc' ? '&#9650;' : '&#9660;' }}</span>
          </th>
          <th class="col-author sortable hide-mobile" @click="sort('author')">
            Author
            <span v-if="sortBy === 'author'">{{ sortOrder === 'asc' ? '&#9650;' : '&#9660;' }}</span>
          </th>
          <th class="col-status sortable" hide-mobile @click="sort('status')">
            Status
            <span v-if="sortBy === 'status'">{{ sortOrder === 'asc' ? '&#9650;' : '&#9660;' }}</span>
          </th>
          <th class="col-loan">Loan</th>
          <th class="col-manage">Manage</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="book in books" :key="book.id">
          <td class="col-isbn hide-mobile">{{ book.isbn }}</td>
          <td class="col-title">{{ book.title }}</td>
          <td class="col-author hide-mobile">{{ book.author }}</td>
          <td class="col-status">
            <BaseBadge :label="getBookStatusLabel(book.status)" :variant="getBookStatusVariant(book.status)">
            </BaseBadge>
          </td>
          <td class="col-loan">
            <BaseButton v-if="book.status === 'available'" variant="secondary" @click="$emit('checkout', book.id)">
              <img src="@/assets/icons/checkout.png" alt="Checkout" class="button-icon" />
            </BaseButton>
            <BaseButton v-else-if="book.status === 'checked_out'" variant="secondary"
              @click="$emit('checkin', book.id)">
              <img src="@/assets/icons/checkin.png" alt="Checkin" class="button-icon" />
            </BaseButton>
          </td>
          <td class="col-manage">
            <BaseButton variant="secondary" @click="$emit('view', book.id)">
              <img src="@/assets/icons/view.png" alt="View" class="button-icon" />
            </BaseButton>
            <BaseButton variant="secondary" @click="$emit('edit', book.id)">
              <img src="@/assets/icons/edit.png" alt="Delete" class="button-icon" />
            </BaseButton>
            <BaseButton variant="secondary" @click="$emit('delete', book.id)">
              <img src="@/assets/icons/delete.png" alt="Delete" class="button-icon" />
            </BaseButton>
          </td>
        </tr>
      </tbody>
    </table>
    <div class="table-placeholder" v-if="!books || books.length === 0">
      No books match the current filters.
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import BaseBadge from '@/components/base/BaseBadge.vue'
import BaseButton from '@/components/base/BaseButton.vue'
import type { BookResponse } from '@/types/book'
import { getBookStatusLabel, getBookStatusVariant } from '@/utils/bookStatus'

export default defineComponent({
  name: 'BookTable',
  components: {
    BaseBadge,
    BaseButton
  },
  props: {
    books: {
      type: Array as () => BookResponse[],
      required: true,
    },
    sortBy: {
      type: String as () => 'title' | 'isbn' | 'author' | 'status',
      required: true,
    },
    sortOrder: {
      type: String as () => 'asc' | 'desc',
      required: true,
    },
  },
  emits: ['checkout', 'checkin', 'view', 'edit', 'delete', 'sort'],
  methods: {
    getBookStatusLabel,
    getBookStatusVariant,
    sort(field: 'title' | 'isbn' | 'author' | 'status') {
      this.$emit('sort', field)
    },
    statusBadgeLabel(status: string): string {
      switch (status) {
        case 'available':
          return 'Available'
        case 'checked_out':
          return 'Checked Out'
        default:
          return status.replace(/_/g, ' ').replace(/\b\w/g, c => c.toUpperCase()) // fallback
      }
    },
    statusBadgeVariant(status: string): 'success' | 'error' | 'neutral' {
      if (status === 'available') return 'success'
      if (status === 'checked_out') return 'error'
      return 'neutral'
    }
  },
})
</script>

<style scoped>
.table-container {
  width: 100%;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 1rem;
  font-size: 0.9rem;
}

thead {
  background-color: #f8f9fa;
}

th,
td {
  white-space: nowrap;
  padding: 0.40rem;
  text-align: left;
  border-bottom: 1px solid #dee2e6;
}

th.sortable {
  cursor: pointer;
  user-select: none;
}

.col-isbn {
  width: 8rem;
}

.col-title {
  width: auto;
}

.col-author {
  width: 12rem;
}

.col-status {
  width: 7rem;
  text-align: center;
}

.col-loan {
  width: 6rem;
  text-align: center;
}

.col-manage {
  width: 11rem;
  text-align: center;
}

tbody tr:hover {
  background-color: #f1f3f5;
}

td:last-child {
  white-space: nowrap;
}

tbody tr:nth-child(even) {
  background-color: #f9f9f9;
}

.button-icon {
  display: block;
  margin: 0 auto;
  width: 1.0rem;
  height: 1.25rem;
}

.table-placeholder {
  width: 100%;
  text-align: center;
  border: 1px solid #ccc;
  background-color: #f5f5f5;
  color: #666;
  padding: 1rem;
  border-radius: 4px;
  text-align: center;
  font-style: italic;
  font-size: 0.95rem;
  margin-top: 1rem;
}

@media (max-width: 600px) {
  .hide-mobile {
    display: none;
  }

  table {
    font-size: 0.8rem;
  }
}
</style>