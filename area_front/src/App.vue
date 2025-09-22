<template>
  <div id="app">
    <ShootingStars />
    <AuthContainer />
  </div>
</template>

<script setup lang="ts">
import AuthContainer from './components/AuthContainer.vue'
import ShootingStars from './components/ShootingStars.vue'
</script>

<style>
#app { position: relative; }
</style>
