<template>
  <teleport to="body">
    <div class="dialog-backdrop" @click.self="$emit('cancel')">
      <div class="dialog-box">
        <div class="dialog-title">{{ title }}</div>
        <div class="dialog-message">{{ message }}</div>
        <div class="dialog-actions">
          <BaseButton v-if="cancelText" variant="secondary" @click="$emit('cancel')">
            {{ cancelText }}
          </BaseButton>
          <BaseButton v-if="confirmText" variant="primary" @click="$emit('confirm')">
            {{ confirmText }}
          </BaseButton>
        </div>
      </div>
    </div>
  </teleport>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import BaseButton from '@/components/base/BaseButton.vue'

export default defineComponent({
  name: 'BaseDialog',
  components: {
    BaseButton,
  },
  emits: ['confirm', 'cancel'],
  props: {
    title: { type: String, required: true },
    message: { type: String, required: true },
    type: {
      type: String as () => 'confirm' | 'error',
      default: 'confirm',
    },
    confirmText: { type: String, default: 'OK' },
    cancelText: { type: String, default: 'Cancel' },
  },
})
</script>

<style scoped>
.dialog-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background: rgba(0, 0, 0, 0.4);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dialog-box {
  background: white;
  border-radius: 8px;
  max-width: 400px;
  width: 90%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  animation: fadeIn 0.2s ease-out;
}

.dialog-title {
  background-color: #f5f5f5;
  padding: 1rem;
  font-weight: bold;
  font-size: 1.1rem;
  border-bottom: 1px solid #ddd;
}

.dialog-message {
  padding: 1rem;
  font-size: 1rem;
  color: #333;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  padding: 0.75rem 1rem;
  gap: 0.5rem;
  border-top: 1px solid #eee;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: scale(0.97);
  }

  to {
    opacity: 1;
    transform: scale(1);
  }
}
</style>
