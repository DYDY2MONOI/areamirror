<template>
  <div class="oauth2-callback-page">
    <div class="callback-container">
      <div class="callback-card">
        <div class="callback-header">
          <div class="logo-container">
            <div class="logo-icon"></div>
          </div>
          <h1 class="callback-title">Twitter / X Authentication</h1>
          <p class="callback-subtitle">{{ subtitle }}</p>
        </div>

        <div v-if="loading" class="loading-container">
          <div class="loading-spinner"></div>
          <p class="loading-text">{{ message }}</p>
        </div>

        <div v-else-if="success" class="success-container">
          <div class="success-icon"></div>
          <h2 class="success-title">{{ successTitle }}</h2>
          <p class="success-message">{{ message }}</p>
          <p class="redirect-text">Redirecting to dashboard...</p>
        </div>

        <div v-else class="error-container">
          <div class="error-icon"></div>
          <h2 class="error-title">Authentication Failed</h2>
          <p class="error-message">{{ error }}</p>
          <button @click="retryLogin" class="retry-button">Try Again</button>
          <button @click="goToLogin" class="login-button">Go to Login</button>
        </div>
      </div>
    </div>
  </div>
  <div v-if="Object.keys(callbackParams).length" class="debug-section">
    <h2 class="debug-title">Callback Parameters</h2>
    <p class="debug-description">Copy these values when checking the Twitter redirect.</p>
    <ul class="debug-list">
      <li v-for="(value, key) in callbackParams" :key="key">
        <span class="debug-key">{{ key }}</span>
        <span class="debug-value">{{ value || '—' }}</span>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { oauth2AuthService } from '@/services/oauth2-auth'
import { API_BASE_URL } from '@/config/api'

const router = useRouter()
const loading = ref(true)
const success = ref(false)
const error = ref('')
const message = ref('Connecting to Twitter/X...')
const subtitle = ref('Processing your Twitter/X login...')
const successTitle = ref('Login Successful!')
const callbackParams = ref<Record<string, string>>({})

const { linkTwitterAccount } = useAuth()

const redirectToDashboard = () => {
  router.push('/profile')
}

const retryLogin = () => {
  window.location.href = '/login'
}

const goToLogin = () => {
  router.push('/login')
}

onMounted(async () => {
  try {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const errorParam = urlParams.get('error')
    const stateParam = urlParams.get('state')
    callbackParams.value = Object.fromEntries(urlParams.entries())

    if (errorParam) {
      throw new Error('Twitter/X authentication was cancelled or failed')
    }

    if (!code) {
      throw new Error('No authorization code received from Twitter/X')
    }

    if (stateParam === 'link') {
      subtitle.value = 'Linking your Twitter/X account...'
      message.value = 'Linking your Twitter/X account...'

      const codeVerifier = sessionStorage.getItem('twitter_code_verifier')
      if (!codeVerifier) {
        throw new Error('PKCE verification failed: code_verifier not found')
      }

      const result = await linkTwitterAccount(code, codeVerifier)
      sessionStorage.removeItem('twitter_code_verifier')

      successTitle.value = 'Twitter/X Account Linked!'
      message.value = result?.twitter_username
        ? `Twitter/X account linked (@${result.twitter_username})`
        : 'Twitter/X account linked successfully!'
      success.value = true
      loading.value = false

      setTimeout(() => {
        redirectToDashboard()
      }, 2000)
      return
    }

    message.value = 'Authenticating with Twitter/X...'

    const response = await fetch(`${API_BASE_URL}/oauth2/twitter/callback?code=${code}`)

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: 'Authentication failed' }))
      throw new Error(errorData.error || 'Twitter/X authentication failed')
    }

    const data = await response.json()

    oauth2AuthService.handleSuccessfulAuth(data)

    success.value = true
    loading.value = false
    message.value = 'Twitter/X login successful!'

    setTimeout(() => {
      redirectToDashboard()
    }, 2000)
  } catch (err) {
    loading.value = false
    error.value = err instanceof Error ? err.message : 'Failed to authenticate with Twitter/X'
    console.error('Twitter/X OAuth2 error:', err)
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
  font-size: 0.9rem;
  color: var(--color-text-secondary);
  margin-top: 0.5rem;
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
}

.retry-button,
.login-button {
  width: 100%;
  padding: 0.875rem;
  border-radius: var(--radius-lg);
  border: none;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-normal);
}

.retry-button {
  background: linear-gradient(135deg, #1DA1F2, #1991DB);
  color: white;
  border: none;
}

.retry-button:hover {
  background: linear-gradient(135deg, #1991DB, #1781C2);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(29, 161, 242, 0.3);
}

.login-button {
  background: transparent;
  border: 2px solid var(--color-border-primary);
  color: var(--color-text-primary);
}

.login-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.debug-section {
  margin: 1.5rem auto;
  max-width: 480px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-xl);
  padding: 1.25rem 1.5rem;
  color: var(--color-text-secondary);
}

.debug-title {
  margin: 0 0 0.5rem 0;
  font-size: 1.1rem;
  color: var(--color-text-primary);
}

.debug-description {
  margin: 0 0 1rem 0;
  font-size: 0.85rem;
}

.debug-list {
  list-style: none;
  margin: 0;
  padding: 0;
  display: grid;
  gap: 0.5rem;
}

.debug-key {
  font-weight: 600;
  margin-right: 0.5rem;
  color: var(--color-text-primary);
}

.debug-value {
  font-family: 'Fira Code', 'Source Code Pro', monospace;
  word-break: break-all;
}

@media (max-width: 640px) {
  .callback-card {
    padding: 2rem 1.5rem;
  }

  .callback-title {
    font-size: 1.5rem;
  }
}
</style>
