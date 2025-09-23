<template>
  <v-list-item rounded class="sidebar-button" @click="$emit('open')">
    <template #prepend>
      <v-icon color="var(--color-accent-primary)" size="24">mdi-plus-circle</v-icon>
    </template>
    <v-list-item-title class="button-text">Create Area</v-list-item-title>
  </v-list-item>
</template>

<script setup lang="ts">
defineEmits<{ (e: 'open'): void }>()
</script>

<style scoped>
.sidebar-button {
  margin: 0.25rem 0;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.sidebar-button:hover {
  background: var(--color-hover-bg);
}

.button-text {
  font-weight: 500;
  color: var(--color-text-primary);
}
</style>

