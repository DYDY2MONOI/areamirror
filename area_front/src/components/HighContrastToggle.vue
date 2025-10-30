<template>
  <div class="daltonism-control">
    <button class="hc-open-btn" @click="emitOpen" :title="`Daltonism: ${labelMap[daltonismMode]}`">
      <v-icon size="18">mdi-eye</v-icon>
      <span>Daltonism</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useTheme, type DaltonismMode } from '@/composables/useTheme'

const { daltonismMode: modeRef, setDaltonismMode } = useTheme()

const emit = defineEmits<{
  open: []
}>()

const options = [
  { title: 'No filter', value: 'none' as DaltonismMode },
  { title: 'Protanopia', value: 'protanopia' as DaltonismMode },
  { title: 'Deuteranopia', value: 'deuteranopia' as DaltonismMode },
  { title: 'Tritanopia', value: 'tritanopia' as DaltonismMode },
  { title: 'Monochrome', value: 'monochrome' as DaltonismMode },
]

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

// Expose for parent
defineExpose({
  setDaltonismMode,
  options,
  labelMap,
})
</script>

<style scoped>
.daltonism-control {
  display: flex;
  flex-direction: column;
}

.hc-open-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 36px;
  padding: 0 12px;
  border-radius: 10px;
  border: 1px solid var(--color-border-primary, var(--border-primary));
  background: var(--color-bg-card, var(--bg-card));
  color: var(--color-text-primary, var(--text-primary));
  cursor: pointer;
  transition: var(--transition-colors);
}
.hc-open-btn:hover { background: var(--color-hover-bg, var(--overlay-hover)); }
</style>
