<template>
  <div class="login-page">
    <div class="login-background">
      <div class="geometric-shape shape-1"></div>
      <div class="geometric-shape shape-2"></div>
      <div class="geometric-shape shape-3"></div>
    </div>

    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <div class="logo-container">
            <div class="logo-icon">
              <v-icon size="32" color="white">mdi-vector-square</v-icon>
            </div>
          </div>
          <h1 class="login-title">Welcome Back</h1>
          <p class="login-subtitle">Sign in to your AREA account</p>
        </div>

        <form class="login-form" @submit.prevent="handleLogin">
          <div class="form-group">
            <div class="input-container">
              <v-icon class="input-icon" size="20">mdi-email-outline</v-icon>
              <input
                v-model="form.email"
                type="email"
                required
                class="form-input"
                placeholder="Email address"
                :class="{ 'error': error && !form.email }"
              />
            </div>
          </div>

          <div class="form-group">
            <div class="input-container">
              <v-icon class="input-icon" size="20">mdi-lock-outline</v-icon>
              <input
                v-model="form.password"
                type="password"
                required
                class="form-input"
                placeholder="Password"
                :class="{ 'error': error && !form.password }"
              />
            </div>
          </div>

          <div v-if="error" class="error-message">
            <v-icon size="16" class="error-icon">mdi-alert-circle</v-icon>
            {{ error }}
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="login-button"
            :class="{ 'loading': loading }"
          >
            <div v-if="loading" class="loading-spinner"></div>
            <span>{{ loading ? 'Signing in...' : 'Sign In' }}</span>
          </button>

          <!-- Séparateur -->
          <div class="divider">
            <span class="divider-text">or</span>
          </div>

          <!-- Bouton Guest -->
          <button
            type="button"
            @click="continueAsGuest"
            class="guest-button"
          >
            <v-icon size="20" class="guest-icon">mdi-account-outline</v-icon>
            <span>Continue as Guest</span>
          </button>
        </form>

        <div class="login-footer">
          <p class="footer-text">
            Don't have an account?
            <router-link to="/register" class="footer-link">
              Create one here
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { authService, type LoginRequest } from '@/services/auth'

const router = useRouter()

const form = ref<LoginRequest>({
  email: '',
  password: ''
})

const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  if (!form.value.email || !form.value.password) {
    error.value = 'Please fill in all fields'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await authService.login(form.value)
    router.push('/')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Connection error'
  } finally {
    loading.value = false
  }
}

const continueAsGuest = () => {
  router.push('/')
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: var(--gradient-bg-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.login-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 1;
}

.geometric-shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  filter: blur(1px);
  animation: float 6s ease-in-out infinite;
}

.shape-1 {
  width: 200px;
  height: 200px;
  background: var(--gradient-accent);
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 150px;
  height: 150px;
  background: linear-gradient(135deg, var(--color-accent-secondary), var(--color-accent-tertiary));
  top: 60%;
  right: 15%;
  animation-delay: 2s;
}

.shape-3 {
  width: 100px;
  height: 100px;
  background: linear-gradient(135deg, var(--color-accent-tertiary), var(--color-accent-primary));
  bottom: 20%;
  left: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
}

.login-container {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 420px;
  padding: 2rem;
}

.login-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-2xl);
  padding: 3rem 2.5rem;
  backdrop-filter: blur(20px);
  box-shadow:
    0 20px 25px -5px rgba(0, 0, 0, 0.1),
    0 10px 10px -5px rgba(0, 0, 0, 0.04),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  animation: cardSlideIn 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

@keyframes cardSlideIn {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.login-header {
  text-align: center;
  margin-bottom: 2.5rem;
}

.logo-container {
  margin-bottom: 1.5rem;
}

.logo-icon {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-full);
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  box-shadow: var(--shadow-glow);
  animation: logoPulse 2s ease-in-out infinite;
}

@keyframes logoPulse {
  0%, 100% { transform: scale(1); }
  50% { transform: scale(1.05); }
}

.login-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
  background: var(--gradient-text);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.login-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 400;
}

.login-form {
  margin-bottom: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 1rem;
  color: var(--color-text-secondary);
  z-index: 1;
  transition: var(--transition-normal);
}

.form-input {
  width: 100%;
  padding: 1rem 1rem 1rem 3rem;
  background: rgba(15, 23, 42, 0.6);
  border: 2px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-size: 1rem;
  font-weight: 400;
  transition: var(--transition-normal);
  outline: none;
}

.form-input::placeholder {
  color: var(--color-text-secondary);
}

.form-input:focus {
  border-color: var(--color-border-focus);
  background: rgba(15, 23, 42, 0.8);
  box-shadow: 0 0 0 4px var(--color-focus-ring);
}

.form-input:focus + .input-icon,
.input-container:focus-within .input-icon {
  color: var(--color-accent-primary);
  transform: scale(1.1);
}

.form-input.error {
  border-color: var(--color-accent-tertiary);
  background: rgba(249, 115, 22, 0.1);
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(249, 115, 22, 0.1);
  border: 1px solid rgba(249, 115, 22, 0.3);
  border-radius: var(--radius-md);
  color: var(--color-accent-tertiary);
  font-size: 0.875rem;
  font-weight: 500;
  margin-bottom: 1rem;
  animation: errorSlideIn 0.3s ease-out;
}

.error-icon {
  color: var(--color-accent-tertiary);
}

@keyframes errorSlideIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.login-button {
  width: 100%;
  padding: 1rem;
  background: var(--gradient-accent);
  border: none;
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-size: 1rem;
  font-weight: 600;
  letter-spacing: 0.02em;
  cursor: pointer;
  transition: var(--transition-normal);
  position: relative;
  overflow: hidden;
  box-shadow: var(--shadow-glow);
}

.login-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow:
    var(--shadow-glow),
    0 10px 20px -5px rgba(6, 182, 212, 0.5);
}

.login-button:active:not(:disabled) {
  transform: translateY(0);
}

.login-button:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-right: 0.5rem;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* Séparateur */
.divider {
  display: flex;
  align-items: center;
  margin: 1.5rem 0;
  position: relative;
}

.divider::before {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--color-border-primary);
}

.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--color-border-primary);
}

.divider-text {
  padding: 0 1rem;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  background: var(--color-bg-card);
}

/* Bouton Guest */
.guest-button {
  width: 100%;
  padding: 1rem;
  background: transparent;
  border: 2px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-size: 1rem;
  font-weight: 500;
  letter-spacing: 0.02em;
  cursor: pointer;
  transition: var(--transition-normal);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.guest-button:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: translateY(-1px);
}

.guest-button:active {
  transform: translateY(0);
}

.guest-icon {
  color: var(--color-text-secondary);
  transition: var(--transition-normal);
}

.guest-button:hover .guest-icon {
  color: var(--color-accent-primary);
  transform: scale(1.1);
}

.login-footer {
  text-align: center;
}

.footer-text {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin: 0;
}

.footer-link {
  color: var(--color-accent-primary);
  text-decoration: none;
  font-weight: 600;
  transition: var(--transition-normal);
  position: relative;
}

.footer-link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: var(--gradient-accent);
  transition: width 0.3s ease;
}

.footer-link:hover {
  color: var(--color-accent-secondary);
}

.footer-link:hover::after {
  width: 100%;
}

@media (max-width: 480px) {
  .login-container {
    padding: 1rem;
  }

  .login-card {
    padding: 2rem 1.5rem;
  }

  .login-title {
    font-size: 1.75rem;
  }

  .logo-icon {
    width: 56px;
    height: 56px;
  }

  .form-input {
    padding: 0.875rem 0.875rem 0.875rem 2.75rem;
  }

  .input-icon {
    left: 0.875rem;
  }
}

@media (prefers-reduced-motion: reduce) {
  .geometric-shape,
  .logo-icon,
  .login-card,
  .error-message {
    animation: none !important;
  }

  .login-button:hover:not(:disabled) {
    transform: none;
  }
}
</style>
