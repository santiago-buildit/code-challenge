<template>
  <form @submit.prevent="onSearch" class="filters">
    <input type="text" id="isbn" name="isbn" v-model="localFilters.isbn" placeholder="ISBN" maxlength="20" />
    <input type="text" id="title" name="title" v-model="localFilters.title" placeholder="Title" maxlength="255" />
    <input type="text" id="author" name="author" v-model="localFilters.author" placeholder="Author" maxlength="255" />
    <select id="status" name="status" v-model="localFilters.status"
      :class="{ 'select-placeholder': !localFilters.status }">
      <option v-if="!localFilters.status" disabled value="">Status</option>
      <option value="all">All</option>
      <option value="available">Available</option>
      <option value="checked_out">Checked Out</option>
    </select>
    <BaseButton type="submit" variant="primary">Search</BaseButton>
    <BaseButton type="button" variant="secondary" @click="$emit('add')">Add Book</BaseButton>
    <BaseButton type="button" variant="secondary" @click="$emit('reset')">Reset Filters</BaseButton>
  </form>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import BaseButton from '@/components/base/BaseButton.vue'
import type { ListBooksRequest } from '@/types/book'

export default defineComponent({
  name: 'BookFilter',
  components: {
    BaseButton,
  },
  props: {
    modelValue: {
      type: Object as () => ListBooksRequest,
      required: true,
    },
  },
  emits: ['update:modelValue', 'search', 'add', 'reset'],
  computed: {
    localFilters: {
      get(): ListBooksRequest {
        return this.modelValue
      },
      set(val: ListBooksRequest) {
        this.$emit('update:modelValue', val)
      },
    },
  },
  methods: {
    onSearch() {
      this.$emit('search')
    },
  },
})
</script>

<style scoped>
.filters {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
  align-items: stretch;
}

.filters input,
.filters select {
  flex: 1 1 150px;
  min-width: 120px;
  height: 2.2rem;
  padding: 0 0.5rem;
}

.filters .button-group {
  display: flex;
  gap: 0.5rem;
  flex: 1 1 100%;
  justify-content: flex-end;
}

.filters .button-group button {
  height: 2.2rem;
  min-width: 100px;
}
</style>