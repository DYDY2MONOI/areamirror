<template>
  <button
    class="theme-toggle"
    @click="toggleTheme"
    :title="label"
    aria-label="Toggle theme"
  >
    <transition name="icon-fade" mode="out-in">
      <template v-if="currentTheme === 'dark'">
        <v-icon key="dark" size="22">mdi-weather-sunny</v-icon>
      </template>
      <template v-else-if="currentTheme === 'light'">
        <v-icon key="light" size="22">mdi-weather-night</v-icon>
      </template>
      <template v-else>
        <!-- When in high-contrast, show the base-theme icon that would be used when exiting HC -->
        <v-icon key="hc" size="22">{{ prevNonHCTheme === 'dark' ? 'mdi-weather-sunny' : 'mdi-weather-night' }}</v-icon>
      </template>
    </transition>
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useTheme } from '@/composables/useTheme'

const { currentTheme, prevNonHCTheme, toggleTheme } = useTheme()

const label = computed(() => {
  const t = currentTheme.value
  if (t === 'dark') return 'Passer en mode clair'
  if (t === 'light') return 'Passer en mode sombre'
  // In high-contrast, we still allow flipping the base theme setting
  return prevNonHCTheme.value === 'dark' ? 'Passer en mode clair (base)' : 'Passer en mode sombre (base)'
})
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







