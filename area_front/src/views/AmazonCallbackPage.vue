<template>
  <div class="callback-page">
    <div class="callback-container">
      <div class="loading-spinner"></div>
      <h2>Processing Amazon Authentication...</h2>
      <p>Please wait while we link your Amazon account.</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { linkAmazonAccount } = useAuth()

onMounted(async () => {
  const urlParams = new URLSearchParams(window.location.search)
  const code = urlParams.get('code')
  const error = urlParams.get('error')

  if (error) {
    router.push('/profile?amazon_error=' + encodeURIComponent(error))
    return
  }

  if (code) {
    try {
      await linkAmazonAccount(code)
      router.push('/profile?amazon_linked=true')
    } catch (err) {
      console.error('Amazon linking error:', err)
      router.push('/profile?amazon_error=' + encodeURIComponent('Failed to link Amazon account'))
    }
  } else {
    router.push('/profile?amazon_error=' + encodeURIComponent('No authorization code received'))
  }
})
</script>

<style scoped>
.callback-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--gradient-bg-primary);
}

.callback-container {
  text-align: center;
  color: white;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-top: 4px solid #ff9900;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1rem;
}
</style>
