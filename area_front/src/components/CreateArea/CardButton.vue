<template>
  <div class="card-button-container">
    <div class="create-card" @click="$emit('open')">
      <div class="card-header">
        <div class="card-icon">
          <v-icon size="32" color="white">mdi-plus-circle</v-icon>
        </div>
        <div class="card-badge">
          <span>New</span>
        </div>
      </div>

      <div class="card-content">
        <h3 class="card-title">Create New Area</h3>
        <p class="card-subtitle">Build custom automation</p>
        <p class="card-description">Connect your favorite services and create intelligent workflows</p>
      </div>

      <div class="card-footer">
        <div class="card-action">
          <span class="action-text">Get Started</span>
          <v-icon size="16" color="var(--color-accent-primary)">mdi-arrow-right</v-icon>
        </div>
      </div>

      <div class="card-overlay"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
defineEmits<{ (e: 'open'): void }>()
</script>

<style scoped>

.card-button-container {
  width: 100%;
  height: 100%;
  max-width: 320px;
  box-sizing: border-box;
}

.create-card {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 280px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
  border-radius: 20px;
  padding: 1.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.create-card:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-4px);
  box-shadow:
    0 12px 32px rgba(0, 0, 0, 0.3),
    0 0 0 1px rgba(255, 255, 255, 0.1);
}

.create-card:active {
  transform: translateY(-2px);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
}

.card-icon {
  width: 56px;
  height: 56px;
  background: linear-gradient(135deg, var(--color-accent-primary), var(--color-accent-secondary));
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(255, 255, 255, 0.1);
}

.card-badge {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 20px;
  padding: 0.25rem 0.75rem;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--color-accent-primary);
  letter-spacing: 0.02em;
}

.card-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.02em;
  line-height: 1.3;
}

.card-subtitle {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin: 0;
  letter-spacing: -0.01em;
}

.card-description {
  font-size: 0.8125rem;
  color: var(--color-text-secondary);
  margin: 0;
  line-height: 1.5;
  opacity: 0.8;
}

.card-footer {
  margin-top: auto;
  padding-top: 1rem;
}

.card-action {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  opacity: 0.7;
  transition: all 0.2s ease;
}

.create-card:hover .card-action {
  opacity: 1;
  transform: translateX(4px);
}

.action-text {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-accent-primary);
  letter-spacing: 0.01em;
}

.card-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(135deg,
    rgba(255, 255, 255, 0.05) 0%,
    transparent 50%,
    rgba(255, 255, 255, 0.02) 100%);
  opacity: 0;
  transition: opacity 0.3s ease;
  border-radius: 20px;
  pointer-events: none;
}

.create-card:hover .card-overlay {
  opacity: 1;
}

@media (max-width: 768px) {
  .create-card {
    min-height: 240px;
    padding: 1.25rem;
    border-radius: 16px;
  }

  .card-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
  }

  .card-title {
    font-size: 1.125rem;
  }

  .card-subtitle {
    font-size: 0.8125rem;
  }

  .card-description {
    font-size: 0.75rem;
  }

  .action-text {
    font-size: 0.8125rem;
  }
}

@media (max-width: 480px) {
  .create-card {
    min-height: 200px;
    padding: 1rem;
  }

  .card-header {
    margin-bottom: 1rem;
  }

  .card-content {
    gap: 0.5rem;
  }

  .card-footer {
    padding-top: 0.75rem;
  }
}

.create-card {
  animation: cardSlideIn 0.4s ease-out;
}

@keyframes cardSlideIn {
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

