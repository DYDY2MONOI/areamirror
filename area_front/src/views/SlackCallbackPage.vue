<template>
  <div class="callback-page">
    <div class="callback-container">
      <div class="callback-card">
        <div class="callback-header">
          <div class="logo-container">
            <v-icon size="48" color="primary">mdi-slack</v-icon>
          </div>
          <h1 class="callback-title">Slack Authentication</h1>
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
            @click="redirectToProfile"
            color="primary"
            variant="outlined"
            class="mt-4"
          >
            Return to Profile
          </v-btn>
        </div>

        <div class="success-container" v-if="success">
          <v-icon size="24" color="success">mdi-check-circle</v-icon>
          <p class="success-text">Slack workspace linked successfully!</p>
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
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { linkSlackAccount, refreshProfile } = useAuth()

const loading = ref(true)
const error = ref('')
const success = ref(false)
const message = ref('Processing Slack authentication...')

const redirectToProfile = () => {
  router.push('/profile')
}

onMounted(async () => {
  try {
    const urlParams = new URLSearchParams(window.location.search)
    const code = urlParams.get('code')
    const errorParam = urlParams.get('error')

    if (errorParam) {
      throw new Error('Slack authentication was cancelled or failed')
    }

    if (!code) {
      throw new Error('No authorization code received from Slack')
    }

    message.value = 'Linking Slack workspace...'

    await linkSlackAccount(code)
    await refreshProfile()

    success.value = true
    loading.value = false
    message.value = 'Slack workspace linked successfully!'

    setTimeout(() => {
      redirectToProfile()
    }, 2000)

  } catch (err) {
    loading.value = false
    error.value = err instanceof Error ? err.message : 'Failed to link Slack workspace'
    console.error('Slack OAuth error:', err)
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
  border: 3px solid rgba(74, 21, 75, 0.2);
  border-top: 3px solid #4A154B;
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

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  color: var(--color-error);
}

.error-text {
  font-size: 0.875rem;
  margin: 0;
  text-align: center;
}

.success-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  color: var(--color-success);
}

.success-text {
  font-size: 0.875rem;
  margin: 0;
  text-align: center;
}

.mt-4 {
  margin-top: 1rem;
}
</style>
