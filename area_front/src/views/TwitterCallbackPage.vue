<template>
  <div class="callback-page">
    <div class="callback-container">
      <div class="callback-card">
        <div class="callback-header">
          <div class="logo-container">
            <v-icon size="48" color="#1DA1F2">mdi-twitter</v-icon>
          </div>
          <h1 class="callback-title">Twitter / X Authentication</h1>
          <p class="callback-subtitle">{{ message }}</p>
        </div>

        <div class="loading-container" v-if="loading">
          <div class="loading-spinner"></div>
          <p class="loading-text">Linking your Twitter/X account...</p>
        </div>

        <div class="error-container" v-else-if="error">
          <v-icon size="24" color="error">mdi-alert-circle</v-icon>
          <p class="error-text">{{ error }}</p>
          <v-btn
            @click="redirectToProfile"
            color="primary"
            variant="outlined"
            class="mt-4"
          >
            Return to Profile
          </v-btn>
        </div>

        <div class="success-container" v-else>
          <v-icon size="24" color="success">mdi-check-circle</v-icon>
          <p class="success-text">Twitter/X account linked successfully!</p>
          <v-btn
            @click="redirectToProfile"
            color="primary"
            class="mt-4"
          >
            Continue to Profile
          </v-btn>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authService } from '@/services/auth'

const router = useRouter()
const loading = ref(true)
const error = ref('')
const message = ref('Processing Twitter/X authentication...')

const redirectToProfile = () => {
  router.push('/profile')
}

onMounted(async () => {
  try {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const errorParam = urlParams.get('error')

    if (errorParam) {
      throw new Error('Twitter/X authentication was cancelled or failed')
    }

    if (!code) {
      throw new Error('No authorization code received from Twitter/X')
    }

    const codeVerifier = sessionStorage.getItem('twitter_code_verifier')
    if (!codeVerifier) {
      throw new Error('PKCE verification failed: code_verifier not found')
    }

    await authService.linkTwitterAccount(code, codeVerifier)
    
    sessionStorage.removeItem('twitter_code_verifier')

    loading.value = false
    message.value = 'Twitter/X account linked successfully!'

    setTimeout(() => {
      redirectToProfile()
    }, 1500)
  } catch (err) {
    loading.value = false
    error.value = err instanceof Error ? err.message : 'Failed to link Twitter/X account'
    console.error('Twitter/X linking error:', err)
  }
})
</script>

<style scoped>
.callback-page {
  min-height: 100vh;
  background: var(--gradient-bg-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.callback-container {
  width: 100%;
  max-width: 480px;
}

.callback-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-2xl);
  padding: 3rem 2.5rem;
  backdrop-filter: blur(20px);
  box-shadow:
    0 20px 25px -5px rgba(0, 0, 0, 0.1),
    0 10px 10px -5px rgba(0, 0, 0, 0.04),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  text-align: center;
}

.callback-header {
  margin-bottom: 2rem;
}

.logo-container {
  margin-bottom: 1.5rem;
}

.callback-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.callback-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 400;
}

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(29, 161, 242, 0.2);
  border-top: 3px solid #1DA1F2;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin: 0;
}

.error-container,
.success-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.error-text {
  color: var(--color-error);
  font-size: 0.875rem;
  margin: 0;
}

.success-text {
  color: var(--color-success);
  font-size: 0.975rem;
  margin: 0;
}
</style>

