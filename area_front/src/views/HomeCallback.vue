<template>
  <div class="home-callback">
    <div class="callback-container">
      <div class="callback-card">
        <div class="callback-header">
          <div class="logo-container">
            <v-icon size="48" color="primary">mdi-google</v-icon>
          </div>
          <h1 class="callback-title">Google Authentication</h1>
          <p class="callback-subtitle">{{ message }}</p>
        </div>

        <div class="loading-container" v-if="loading">
          <div class="loading-spinner"></div>
          <p class="loading-text">Processing authentication...</p>
        </div>

        <div class="error-container" v-if="error">
          <v-icon size="24" color="error">mdi-alert-circle</v-icon>
          <p class="error-text">{{ error }}</p>
          <v-btn
            @click="redirectToLogin"
            color="primary"
            variant="outlined"
            class="mt-4"
          >
            Return to Login
          </v-btn>
        </div>

        <div class="success-container" v-if="success">
          <v-icon size="24" color="success">mdi-check-circle</v-icon>
          <p class="success-text">Google account linked successfully!</p>
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
const success = ref(false)
const message = ref('Processing Google authentication...')

const redirectToLogin = () => {
  router.push('/login')
}

const redirectToProfile = () => {
  router.push('/profile')
}

onMounted(async () => {
  try {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const errorParam = urlParams.get('error')

    if (errorParam) {
      throw new Error('Google authentication was cancelled or failed')
    }

    if (!code) {
      router.push('/')
      return
    }

    message.value = 'Linking Google account...'
    
    await authService.linkGoogleAccount(code)
    
    success.value = true
    loading.value = false
    message.value = 'Google account linked successfully!'

    setTimeout(() => {
      redirectToProfile()
    }, 2000)

  } catch (err) {
    loading.value = false
    error.value = err instanceof Error ? err.message : 'Failed to link Google account'
    console.error('Google OAuth error:', err)
  }
})
</script>

<style scoped>
.home-callback {
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
  border: 3px solid rgba(var(--color-primary-rgb), 0.3);
  border-top: 3px solid var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
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
  text-align: center;
}

.success-text {
  color: var(--color-success);
  font-size: 0.875rem;
  margin: 0;
  text-align: center;
}

@media (max-width: 480px) {
  .callback-card {
    padding: 2rem 1.5rem;
  }

  .callback-title {
    font-size: 1.5rem;
  }
}
</style>
