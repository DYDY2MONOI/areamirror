<template>
  <div class="callback-page">
    <div class="callback-container">
      <div class="callback-card">
        <div class="callback-header">
          <div class="logo-container">
            <svg width="48" height="48" viewBox="0 0 24 24" fill="#1db954">
              <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.42 1.56-.299.421-1.02.599-1.559.3z"/>
            </svg>
          </div>
          <h1 class="callback-title">Spotify Authentication</h1>
          <p class="callback-subtitle">{{ message }}</p>
        </div>

        <div class="loading-container" v-if="loading">
          <div class="loading-spinner"></div>
          <p class="loading-text">Linking your Spotify account...</p>
        </div>

        <div class="success-container" v-else-if="success">
          <v-icon size="24" color="success">mdi-check-circle</v-icon>
          <p class="success-text">Spotify account linked successfully!</p>
          <v-btn @click="redirectToProfile" color="success" class="mt-4">
            Continue to Profile
          </v-btn>
        </div>

        <div class="error-container" v-else>
          <v-icon size="24" color="error">mdi-alert-circle</v-icon>
          <p class="error-text">{{ error }}</p>
          <v-btn @click="redirectToProfile" color="primary" variant="outlined" class="mt-4">
            Back to Profile
          </v-btn>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { linkSpotifyAccount } = useAuth()

const loading = ref(true)
const success = ref(false)
const error = ref('An unexpected error occurred while linking your Spotify account.')
const message = ref('Processing Spotify authentication...')

const redirectToProfile = () => {
  router.push('/profile')
}

onMounted(async () => {
  try {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const errorParam = urlParams.get('error')

    if (errorParam) {
      throw new Error('Spotify authentication was cancelled or failed.')
    }

    if (!code) {
      throw new Error('No authorization code received from Spotify.')
    }

    message.value = 'Finalizing Spotify account linking...'

    await linkSpotifyAccount(code)

    success.value = true
    error.value = ''
    loading.value = false
    message.value = 'Spotify account linked successfully!'

    setTimeout(() => {
      redirectToProfile()
    }, 2000)
  } catch (err) {
    loading.value = false
    success.value = false
    error.value = err instanceof Error ? err.message : 'Failed to link Spotify account.'
    console.error('Spotify OAuth error:', err)
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
  border: 3px solid rgba(29, 185, 84, 0.2);
  border-top: 3px solid #1db954;
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

.success-container,
.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.success-text,
.error-text {
  font-size: 0.95rem;
  margin: 0;
  color: var(--color-text-primary);
  text-align: center;
}

.error-text {
  color: var(--color-error);
}
</style>
