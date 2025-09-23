<template>
  <v-list-item rounded class="sidebar-button" @click="$emit('open')">
    <template #prepend>
      <v-icon color="#3b82f6" size="24">mdi-plus-circle</v-icon>
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
  background: rgba(59, 130, 246, 0.1);
}

.button-text {
  font-weight: 500;
  color: white;
}
</style>

