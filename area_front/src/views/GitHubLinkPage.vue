<template>
  <div class="github-link-page">
    <div class="link-container">
      <div class="link-card">
        <div class="link-header">
          <v-icon size="48" color="white" class="header-icon">mdi-github</v-icon>
          <h1 class="link-title">GitHub Account Linking</h1>
        </div>

        <div v-if="loading" class="loading-container">
          <div class="loading-spinner"></div>
          <p class="loading-text">Processing your GitHub account...</p>
        </div>

        <div v-else-if="success" class="success-container">
          <div class="success-icon">✅</div>
          <h2 class="success-title">Success!</h2>
          <p class="success-message">Your GitHub account has been linked successfully!</p>
          <p class="redirect-text">Redirecting to your profile...</p>
        </div>

        <div v-else-if="error" class="error-container">
          <div class="error-icon">❌</div>
          <h2 class="error-title">Linking Failed</h2>
          <p class="error-message">{{ error }}</p>
          <div class="error-actions">
            <button @click="goToProfile" class="action-button primary">
              Go to Profile
            </button>
            <button @click="retryLink" class="action-button secondary">
              Try Again
            </button>
          </div>
        </div>

        <div v-else class="info-container">
          <div class="info-icon">ℹ️</div>
          <h2 class="info-title">No Link Request</h2>
          <p class="info-message">No GitHub link request is currently active.</p>
          <button @click="goToProfile" class="action-button primary">
            Go to Profile
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const success = ref(false)
const error = ref('')

const goToProfile = () => {
  router.push('/profile')
}

const retryLink = () => {
  router.push('/profile')
}

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  const successParam = urlParams.get('success')
  const errorParam = urlParams.get('error')

  if (successParam === 'true') {
    success.value = true
    setTimeout(() => {
      goToProfile()
    }, 2000)
  } else if (errorParam) {
    error.value = decodeURIComponent(errorParam)
  }
})
</script>

<style scoped>
.github-link-page {
  min-height: 100vh;
  background: var(--gradient-bg-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.link-container {
  width: 100%;
  max-width: 500px;
}

.link-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-xl);
  padding: 2rem;
  text-align: center;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
}

.link-header {
  margin-bottom: 2rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.header-icon {
  color: #24292e;
}

.link-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
}

.loading-container,
.success-container,
.error-container,
.info-container {
  animation: fadeIn 0.3s ease-in;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid rgba(255, 255, 255, 0.2);
  border-top: 4px solid var(--color-accent-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1.5rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  color: var(--color-text-secondary);
  font-size: 1rem;
  margin: 0;
}

.success-icon,
.error-icon,
.info-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.success-title,
.error-title,
.info-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 0.75rem;
  color: var(--color-text-primary);
}

.success-title {
  color: #10b981;
}

.error-title {
  color: #ef4444;
}

.success-message,
.error-message,
.info-message {
  color: var(--color-text-secondary);
  font-size: 1rem;
  margin-bottom: 1.5rem;
  line-height: 1.6;
}

.redirect-text {
  color: var(--color-text-secondary);
  font-size: 0.9rem;
  margin-top: 1rem;
}

.error-actions {
  display: flex;
  gap: 1rem;
  flex-direction: column;
}

.action-button {
  padding: 0.75rem 1.5rem;
  border-radius: var(--radius-lg);
  font-weight: 600;
  font-size: 1rem;
  border: none;
  cursor: pointer;
  transition: var(--transition-normal);
  text-transform: none;
}

.action-button.primary {
  background: var(--gradient-accent);
  color: white;
}

.action-button.primary:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-glow);
}

.action-button.secondary {
  background: transparent;
  color: var(--color-accent-primary);
  border: 2px solid var(--color-accent-primary);
}

.action-button.secondary:hover {
  background: rgba(59, 130, 246, 0.1);
  transform: translateY(-2px);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .link-card {
    padding: 1.5rem;
  }

  .link-title {
    font-size: 1.25rem;
  }

  .success-title,
  .error-title,
  .info-title {
    font-size: 1.25rem;
  }
}
</style>
