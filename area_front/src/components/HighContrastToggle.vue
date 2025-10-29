<template>
  <button
    class="hc-toggle"
    :class="{ active: isHighContrast() }"
    @click="toggleHighContrast"
    :aria-pressed="isHighContrast() ? 'true' : 'false'"
    aria-label="Toggle high contrast mode"
    :title="isHighContrast() ? 'Désactiver le mode à haut contraste' : 'Activer le mode à haut contraste'"
  >
    <v-icon size="20">mdi-contrast</v-icon>
    <span class="hc-label">High Contrast</span>
  </button>
</template>

<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'
const { isHighContrast, toggleHighContrast } = useTheme()
</script>

<style scoped>
.hc-toggle {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 36px;
  padding: 0 12px;
  border-radius: 10px;
  border: 1px solid var(--border-primary);
  background: var(--bg-card);
  color: var(--text-primary);
  cursor: pointer;
  transition: var(--transition-colors);
}
.hc-toggle:hover { background: var(--overlay-hover); }
.hc-toggle.active {
  border-color: var(--border-accent);
  box-shadow: 0 0 0 2px var(--border-accent);
}
.hc-label { font-size: 0.85rem; font-weight: 600; }
</style>
