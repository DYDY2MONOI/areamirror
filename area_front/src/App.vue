<template>
  <v-app>
    <v-main>
      <router-view />
    </v-main>
  </v-app>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useTheme } from '@/composables/useTheme'

// Initialiser le thème
const { currentTheme } = useTheme()

onMounted(() => {
  // Appliquer le thème au chargement
  document.documentElement.setAttribute('data-theme', currentTheme.value)
})
</script>

<style>
#app { position: relative; }
</style>
