<template>
  <button 
    class="theme-toggle" 
    @click="toggleTheme"
    :title="currentTheme === 'dark' ? 'Passer en mode clair' : 'Passer en mode sombre'"
  >
    <transition name="icon-fade" mode="out-in">
      <v-icon v-if="currentTheme === 'dark'" key="dark" size="22">mdi-weather-sunny</v-icon>
      <v-icon v-else key="light" size="22">mdi-weather-night</v-icon>
    </transition>
  </button>
</template>

<script setup lang="ts">
import { useTheme } from '@/composables/useTheme'

const { currentTheme, toggleTheme } = useTheme()
</script>

<style scoped>
.theme-toggle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: var(--radius-md);
  border: 1px solid var(--color-border-primary);
  background: var(--color-bg-card);
  color: var(--color-text-primary);
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(20px);
}

.theme-toggle:hover {
  background: var(--color-hover-bg);
  transform: translateY(-1px);
}

.icon-fade-enter-active,
.icon-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.icon-fade-enter-from {
  opacity: 0;
  transform: rotate(-90deg) scale(0.5);
}

.icon-fade-leave-to {
  opacity: 0;
  transform: rotate(90deg) scale(0.5);
}
</style>




