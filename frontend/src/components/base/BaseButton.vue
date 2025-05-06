<template>
  <button :class="['base-button', variantClass]" v-bind="$attrs">
    <slot></slot>
  </button>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'BaseButton',
  props: {
    variant: {
      type: String as () => 'primary' | 'secondary' | 'danger' | 'ghost',
      default: 'primary',
    },
  },
  computed: {
    variantClass(): string {
      return `variant-${this.variant}`
    },
  },
})
</script>

<style scoped>
.base-button {
  padding: 0.5rem 1rem;
  border: none;
  font-size: 1rem;
  cursor: pointer;
  border-radius: 4px;
  transition: background-color 0.2s ease;
}

/* Variants */
.variant-primary {
  background-color: #477777;
  color: white;
}

.variant-primary:hover {
  background-color: #253e3e;
}

.variant-secondary {
  background-color: #e2e6ea;
  color: #333;
}

.variant-secondary:hover {
  background-color: #d6d8db;
}

.variant-danger {
  background-color: #dc3545;
  color: white;
}

.variant-danger:hover {
  background-color: #c82333;
}

.variant-ghost {
  background-color: transparent;
  color: #333;
  border: 1px solid #ccc;
}

.variant-ghost:hover {
  background-color: rgba(0, 0, 0, 0.05);
}

.base-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  background-color: #ccc !important;
  color: #666 !important;
  border: none !important;
}
</style>