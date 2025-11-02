<template>
  <div class="oauth2-callback-page">
    <div class="callback-container">
      <div class="callback-card">
        <div class="callback-header">
          <div class="logo-container">
            <div class="logo-icon">🔐</div>
          </div>
          <h1 class="callback-title">Google Authentication</h1>
          <p class="callback-subtitle">Processing your Google login...</p>
        </div>

        <div v-if="loading" class="loading-container">
          <div class="loading-spinner"></div>
          <p class="loading-text">{{ message }}</p>
        </div>

        <div v-else-if="success" class="success-container">
          <div class="success-icon">✅</div>
          <h2 class="success-title">Login Successful!</h2>
          <p class="success-message">{{ message }}</p>
          <p class="redirect-text">Redirecting to dashboard...</p>
        </div>

        <div v-else-if="error" class="error-container">
          <div class="error-icon">❌</div>
          <h2 class="error-title">Authentication Failed</h2>
          <p class="error-message">{{ error }}</p>
          <button @click="retryLogin" class="retry-button">Try Again</button>
          <button @click="goToLogin" class="login-button">Go to Login</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { oauth2AuthService } from '@/services/oauth2-auth'
import { API_BASE_URL } from '@/config/api'

const router = useRouter()
const loading = ref(true)
const success = ref(false)
const error = ref('')
const message = ref('Connecting to Google...')

const redirectToDashboard = () => {
  router.push('/profile')
}

const retryLogin = () => {
  window.location.href = '/auth/google'
}

const goToLogin = () => {
  router.push('/login')
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
      throw new Error('No authorization code received from Google')
    }

    message.value = 'Authenticating with Google...'
    
    // Use the new OAuth2 direct login endpoint
    const response = await fetch(`${API_BASE_URL}/oauth2/google/callback?code=${code}`)
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Authentication failed' }))
      throw new Error(errorData.error || 'Google authentication failed')
    }

    const data = await response.json()
    
    // Handle the OAuth2 response
    oauth2AuthService.handleSuccessfulAuth(data)
    
    success.value = true
    loading.value = false
    message.value = 'Google login successful!'

    setTimeout(() => {
      redirectToDashboard()
    }, 2000)

  } catch (err) {
    loading.value = false
    error.value = err instanceof Error ? err.message : 'Failed to authenticate with Google'
    console.error('Google OAuth2 error:', err)
  }
})
</script>

<style scoped>
.oauth2-callback-page {
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

.logo-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
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
  border: 3px solid var(--color-border-primary);
  border-top: 3px solid var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.success-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.success-icon {
  font-size: 3rem;
  margin-bottom: 0.5rem;
}

.success-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-success);
  margin: 0;
}

.success-message {
  color: var(--color-text-secondary);
  margin: 0;
}

.redirect-text {
  color: var(--color-text-tertiary);
  font-size: 0.9rem;
  margin: 0;
}

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 0.5rem;
}

.error-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-error);
  margin: 0;
}

.error-message {
  color: var(--color-text-secondary);
  margin: 0;
  text-align: center;
}

.retry-button, .login-button {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: var(--radius-lg);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  margin: 0.25rem;
}

.retry-button {
  background: var(--color-primary);
  color: white;
}

.retry-button:hover {
  background: var(--color-primary-dark);
}

.login-button {
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border-primary);
}

.login-button:hover {
  background: var(--color-bg-tertiary);
}
</style>
