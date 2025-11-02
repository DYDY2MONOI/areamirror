<template>
  <div class="settings-menu">
    <transition name="backdrop-fade">
      <div v-if="isOpen" class="settings-backdrop" @click="toggleMenu"></div>
    </transition>

    <button 
      class="settings-trigger" 
      @click="toggleMenu"
      :title="isOpen ? 'Close settings' : 'Open settings'"
    >
      <v-icon size="22">mdi-cog</v-icon>
    </button>

    <transition name="menu-fade">
      <div 
        v-if="isOpen" 
        class="settings-dropdown"
        :style="{ top: `${menuPosition.y}px`, left: `${menuPosition.x}px` }"
        @mousedown="startDrag"
      >
        <div class="settings-header" style="cursor: move;">
          <span class="settings-title">Settings</span>
          <button class="close-btn" @click="toggleMenu">
            <v-icon size="16">mdi-close</v-icon>
          </button>
        </div>
        
        <div class="settings-content">
          <div class="setting-item">
            <span class="setting-label">Theme</span>
            <ThemeToggle />
          </div>
          
          <div class="setting-item">
            <span class="setting-label">Accessibility</span>
            <HighContrastToggle @open="handleDaltonismOpen" />
          </div>
        </div>
      </div>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import ThemeToggle from './ThemeToggle.vue'
import HighContrastToggle from './HighContrastToggle.vue'

const isOpen = ref(false)
const menuPosition = ref({ x: 100, y: 100 })
const isDragging = ref(false)
const dragOffset = ref({ x: 0, y: 0 })

const emit = defineEmits<{
  openDaltonism: []
}>()

const toggleMenu = () => {
  isOpen.value = !isOpen.value
}

const handleDaltonismOpen = () => {
  emit('openDaltonism')
  isOpen.value = false
}

const startDrag = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.settings-header')) return
  
  isDragging.value = true
  dragOffset.value = {
    x: e.clientX - menuPosition.value.x,
    y: e.clientY - menuPosition.value.y
  }
  
  e.preventDefault()
}

const onDrag = (e: MouseEvent) => {
  if (!isDragging.value) return
  
  menuPosition.value = {
    x: e.clientX - dragOffset.value.x,
    y: e.clientY - dragOffset.value.y
  }
}

const stopDrag = () => {
  isDragging.value = false
}

onMounted(() => {
  menuPosition.value = {
    x: 100,
    y: 100
  }
  
  document.addEventListener('mousemove', onDrag)
  document.addEventListener('mouseup', stopDrag)
})

onUnmounted(() => {
  document.removeEventListener('mousemove', onDrag)
  document.removeEventListener('mouseup', stopDrag)
})
</script>

<style scoped>
.settings-menu {
  position: relative;
}

.settings-trigger {
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

.settings-trigger:hover {
  background: var(--color-hover-bg);
  transform: translateY(-1px);
}

.settings-dropdown {
  position: fixed;
  background: rgba(26, 26, 46, 0.95);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  padding: 16px;
  min-width: 240px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(20px);
  z-index: 1001;
  user-select: none;
}

[data-theme="light"] .settings-dropdown {
  background: rgba(255, 255, 255, 0.98);
  border: 2px solid rgba(0, 0, 0, 0.15);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(20px);
}

.settings-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--color-border-primary);
}

.settings-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.close-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  border: none;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: var(--transition-normal);
}

.close-btn:hover {
  background: var(--color-hover-bg);
  color: var(--color-text-primary);
}

.settings-content {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.setting-label {
  font-size: 0.8125rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  white-space: nowrap;
}

.settings-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.3);
  z-index: 1000;
}

[data-theme="light"] .settings-backdrop {
  background: rgba(0, 0, 0, 0.15);
}

.backdrop-fade-enter-active,
.backdrop-fade-leave-active {
  transition: opacity 0.25s ease, backdrop-filter 0.25s ease;
}

.backdrop-fade-enter-from,
.backdrop-fade-leave-to {
  opacity: 0;
}

.menu-fade-enter-active,
.menu-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.menu-fade-enter-from,
.menu-fade-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

.menu-fade-enter-to,
.menu-fade-leave-from {
  opacity: 1;
  transform: scale(1);
}
</style>

