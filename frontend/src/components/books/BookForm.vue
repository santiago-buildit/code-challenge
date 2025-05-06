<template>
  <form @submit.prevent="submit" class="form">
    <div>
      <label for="isbn">ISBN *</label>
      <input type="text" id="isbn" name="isbn" v-model.trim="localForm.isbn" :disabled="isViewMode" required
        maxlength="20" />
    </div>
    <div>
      <label for="title">Title *</label>
      <input type="text" id="title" name="title" v-model.trim="localForm.title" :disabled="isViewMode" required
        maxlength="255" />
    </div>
    <div>
      <label for="author">Author *</label>
      <input type="text" id="author" name="author" v-model.trim="localForm.author" :disabled="isViewMode" required
        maxlength="255" />
    </div>
    <div>
      <label for="description">Description</label>
      <textarea id="description" name="description" v-model.trim="localForm.description" :disabled="isViewMode"
        maxlength="1000" rows="6"></textarea>
    </div>
    <p class="form-note">* Required fields</p>
    <div class="form-actions">
      <BaseButton v-if="!isViewMode" type="submit" variant="primary" :disabled="!isFormValid"
        :title="!isFormValid ? 'Please fill in all required fields' : ''">
        {{ submitLabel }}
      </BaseButton>

      <BaseButton type="button" variant="ghost" @click="$emit('cancel')">
        {{ mode === 'edit' ? 'Cancel' : 'Back to List' }}
      </BaseButton>
    </div>
  </form>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import BaseButton from '@/components/base/BaseButton.vue'
import type { BookResponse } from '@/types/book'

export default defineComponent({
  name: 'BookForm',
  components: {
    BaseButton,
  },
  props: {
    modelValue: {
      type: Object as () => BookResponse,
      required: true,
    },
    mode: {
      type: String as () => 'new' | 'edit' | 'view',
      default: 'new',
    },
  },
  emits: ['submit', 'cancel', 'update:modelValue'],
  computed: {
    localForm: {
      get(): BookResponse {
        return this.modelValue
      },
      set(val: BookResponse) {
        this.$emit('update:modelValue', val)
      },
    },
    isViewMode(): boolean {
      return this.mode === 'view'
    },
    submitLabel(): string {
      return this.mode === 'edit' ? 'Update' : 'Create'
    },
    isFormValid(): boolean {
      return (
        this.localForm.isbn.trim() !== '' &&
        this.localForm.title.trim() !== '' &&
        this.localForm.author.trim() !== ''
      )
    },
  },
  methods: {
    submit() {
      if (this.isViewMode) return // Prevents submit in view mode
      this.$emit('submit')
    },
  },
})
</script>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 0.5rem;
}

.form-note {
  font-size: 0.85rem;
  color: #666;
  margin-top: 0;
}

label {
  font-weight: bold;
}

input,
textarea {
  width: 100%;
  padding: 0.5rem;
  font-family: inherit;
  font-size: inherit;
}

input:focus,
textarea:focus {
  outline: none;
  border-color: #477777;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

button {
  padding: 0.5rem 1rem;
}
</style>