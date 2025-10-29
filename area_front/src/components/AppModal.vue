<template>
  <div
    v-if="modelValue"
    class="modal-overlay"
    @click="handleOverlayClick"
  >
    <div class="modal-content" @click.stop>
      <div v-if="showHeader" class="modal-header">
        <div v-if="icon" class="modal-icon-container">
          <div class="modal-icon-bg">
            <v-icon :size="iconSize" :color="iconColor">{{ icon }}</v-icon>
          </div>
        </div>
        <h2 class="modal-title">{{ title }}</h2>
        <p v-if="subtitle" class="modal-subtitle">{{ subtitle }}</p>
        <button v-if="showCloseButton" class="modal-close-btn" @click="close">
          <v-icon size="20">mdi-close</v-icon>
        </button>
      </div>

      <div class="modal-body">
        <slot />
      </div>

      <div v-if="hasActions" class="modal-actions">
        <slot name="actions">
          <button
            v-if="cancelButton"
            class="modal-btn modal-cancel-btn"
            @click="handleCancel"
          >
            {{ cancelButton.text }}
          </button>
          <button
            v-if="confirmButton"
            class="modal-btn modal-confirm-btn"
            :disabled="confirmButton.disabled"
            @click="handleConfirm"
          >
            {{ confirmButton.text }}
          </button>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, nextTick } from 'vue'

interface ModalButton {
  text: string
  disabled?: boolean
}

interface Props {
  modelValue: boolean
  title?: string
  subtitle?: string
  icon?: string
  iconSize?: number | string
  iconColor?: string
  showHeader?: boolean
  showCloseButton?: boolean
  hasActions?: boolean
  cancelButton?: ModalButton
  confirmButton?: ModalButton
  closeOnOverlay?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
  (e: 'close'): void
}

const props = withDefaults(defineProps<Props>(), {
  showHeader: true,
  showCloseButton: true,
  hasActions: true,
  closeOnOverlay: true,
  iconSize: 24,
  iconColor: '#ff3b30'
})

const emit = defineEmits<Emits>()

const close = () => {
  emit('update:modelValue', false)
  emit('close')
}

const handleOverlayClick = () => {
  if (props.closeOnOverlay) {
    close()
  }
}

const handleCancel = () => {
  emit('cancel')
  close()
}

const handleConfirm = () => {
  emit('confirm')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease-out;
  padding: var(--spacing-xl);
}

.modal-content {
  max-width: 600px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(40px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow:
    0 20px 40px rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.modal-header {
  padding: 32px 24px 24px 24px;
  text-align: center;
  position: relative;
}

.modal-icon-container {
  margin-bottom: 20px;
}

.modal-icon-bg {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ff3b30, #ff6b6b);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  box-shadow: 0 8px 20px rgba(255, 59, 48, 0.3);
}

.modal-title {
  font-size: 24px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0 0 8px 0;
  letter-spacing: -0.02em;
}

.modal-subtitle {
  font-size: 16px;
  color: #86868b;
  margin: 0;
  line-height: 1.4;
  font-weight: 400;
}

.modal-close-btn {
  position: absolute;
  top: 16px;
  right: 16px;
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(142, 142, 147, 0.12);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.modal-close-btn:hover {
  background: rgba(142, 142, 147, 0.18);
  transform: scale(1.1);
}

.modal-body {
  padding: 0 24px;
  max-height: 60vh;
  overflow-y: auto;
}

.modal-actions {
  padding: 24px;
  display: flex;
  gap: 12px;
  justify-content: center;
}

.modal-btn {
  flex: 1;
  padding: 14px 20px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  letter-spacing: -0.01em;
}

.modal-cancel-btn {
  background: rgba(142, 142, 147, 0.12);
  color: #1d1d1f;
}

.modal-cancel-btn:hover {
  background: rgba(142, 142, 147, 0.18);
  transform: translateY(-1px);
}

.modal-confirm-btn {
  background: linear-gradient(135deg, #ff3b30, #ff6b6b);
  color: white;
  box-shadow: 0 4px 12px rgba(255, 59, 48, 0.3);
}

.modal-confirm-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #ff2d55, #ff5252);
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(255, 59, 48, 0.4);
}

.modal-confirm-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

@media (prefers-color-scheme: dark) {
  .modal-content {
    background: rgba(28, 28, 30, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .modal-title {
    color: #f2f2f7;
  }

  .modal-subtitle {
    color: #8e8e93;
  }

  .modal-cancel-btn {
    background: rgba(142, 142, 147, 0.2);
    color: #f2f2f7;
  }

  .modal-cancel-btn:hover {
    background: rgba(142, 142, 147, 0.3);
  }
}

@media (max-width: 480px) {
  .modal-overlay {
    padding: 1rem;
  }

  .modal-content {
    margin: 1rem;
    max-width: calc(100vw - 2rem);
  }

  .modal-header {
    padding: 24px 20px 20px 20px;
  }

  .modal-icon-bg {
    width: 56px;
    height: 56px;
  }

  .modal-title {
    font-size: 20px;
  }

  .modal-subtitle {
    font-size: 14px;
  }

  .modal-body {
    padding: 0 20px;
  }

  .modal-actions {
    padding: 20px;
    flex-direction: column;
  }

  .modal-btn {
    width: 100%;
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}
</style>
