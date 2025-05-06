<template>
  <div>
    <h3>Loan History</h3>
    <span>Current Status:
      <BaseBadge :label="getBookStatusLabel(status)" :variant="getBookStatusVariant(status)" />
    </span>

    <div class="history-section">
      <ul class="history-list" v-if="history.length > 0">
        <li v-for="entry in history" :key="entry.timestamp">
          <div class="history-entry">
            <span class="history-date">{{ formatDate(entry.timestamp) }}</span>
            <BaseBadge :label="getBookStatusLabel(entry.status)" :variant="getBookStatusVariant(entry.status)" />
          </div>
        </li>
      </ul>
      <div class="history-placeholder" v-else>
        No status changes recorded yet.
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import BaseBadge from '@/components/base/BaseBadge.vue'
import type { StatusChangeResponse } from '@/types/book'
import { getBookStatusLabel, getBookStatusVariant } from '@/utils/bookStatus'

export default defineComponent({
  name: 'BookHistory',
  components: {
    BaseBadge,
  },
  props: {
    status: {
      type: String,
      required: true,
    },
    history: {
      type: Array as () => StatusChangeResponse[],
      required: true,
    },
  },
  methods: {
    getBookStatusLabel,
    getBookStatusVariant,
    formatDate(iso: string) {
      return new Date(iso).toLocaleString()
    },
  },
})
</script>

<style scoped>
.history-section {
  margin-top: 2rem;
}

.history-title {
  margin-bottom: 0.5rem;
  font-size: 1.1rem;
  font-weight: bold;
  color: #333;
}

.history-list {
  list-style: none;
  padding-left: 0;
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.history-entry {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #f8f9fa;
  border: 1px solid #dee2e6;
  border-radius: 4px;
  padding: 0.5rem 0.75rem;
  font-size: 0.95rem;
}

.history-placeholder {
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

.history-date {
  color: #555;
}
</style>