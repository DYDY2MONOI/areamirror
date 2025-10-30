<template>
  <div class="daltonism-control">
    <button class="hc-open-btn" @click="emitOpen" :title="`Daltonism: ${labelMap[daltonismMode]}`">
      <v-icon size="22">mdi-eye</v-icon>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useTheme, type DaltonismMode } from '@/composables/useTheme'

const { daltonismMode: modeRef } = useTheme()

const emit = defineEmits<{
  open: []
}>()

const daltonismMode = computed(() => modeRef.value)

const labelMap: Record<DaltonismMode, string> = {
  none: 'No filter',
  protanopia: 'Protanopia',
  deuteranopia: 'Deuteranopia',
  tritanopia: 'Tritanopia',
  monochrome: 'Monochrome',
}

const emitOpen = () => {
  emit('open')
}
</script>

<style scoped>
.daltonism-control {
  display: flex;
  flex-direction: column;
}

.hc-open-btn {
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

.hc-open-btn:hover {
  background: var(--color-hover-bg);
  transform: translateY(-1px);
}
</style>
