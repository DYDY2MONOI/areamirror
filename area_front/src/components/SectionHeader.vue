<template>
  <div class="section-header">
    <div class="section-info">
      <h2 class="section-title">{{ title }}</h2>
      <p v-if="subtitle" class="section-subtitle">{{ subtitle }}</p>
    </div>
    <div v-if="hasAction" class="section-action">
      <slot name="action">
        <button
          v-if="actionButton"
          class="view-all-btn"
          @click="$emit('action-click')"
        >
          <span>{{ actionButton.text }}</span>
          <v-icon size="16">{{ actionButton.icon }}</v-icon>
        </button>
      </slot>
    </div>
  </div>
</template>

<script setup lang="ts">
interface ActionButton {
  text: string
  icon: string
}

interface Props {
  title: string
  subtitle?: string
  actionButton?: ActionButton
  hasAction?: boolean
}

interface Emits {
  (e: 'action-click'): void
}

withDefaults(defineProps<Props>(), {
  hasAction: false
})

defineEmits<Emits>()
</script>

<style scoped>
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2rem;
}

.section-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.section-title {
  font-weight: 800;
  font-size: 2rem;
  color: #ffffff;
  margin: 0;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff 0%, #e2e8f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  filter: drop-shadow(0 2px 8px rgba(255, 255, 255, 0.1));
}

.section-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 400;
}

.section-action {
  display: flex;
  align-items: center;
}

.view-all-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-lg);
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-weight: 500;
  font-size: 0.875rem;
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(20px);
}

.view-all-btn:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .section-action {
    align-self: flex-end;
  }
}
</style>
