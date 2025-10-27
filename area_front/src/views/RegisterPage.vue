<template>
  <div class="register-page">
    <div class="register-background">
      <div class="geometric-shape shape-1"></div>
      <div class="geometric-shape shape-2"></div>
      <div class="geometric-shape shape-3"></div>
    </div>

    <div class="register-container">
      <div class="register-card">
        <div class="register-header">
          <div class="logo-container">
            <div class="logo-icon">
              <v-icon size="32" color="white">mdi-account-plus</v-icon>
            </div>
          </div>
          <h1 class="register-title">Create Account</h1>
          <p class="register-subtitle">Join AREA and start automating</p>
        </div>

        <form class="register-form" @submit.prevent="handleRegister">
          <div class="form-row">
            <div class="form-group">
              <div class="input-container">
                <v-icon class="input-icon" size="20">mdi-account-outline</v-icon>
                <input
                  v-model="form.first_name"
                  type="text"
                  class="form-input"
                  placeholder="First name"
                />
              </div>
            </div>
            <div class="form-group">
              <div class="input-container">
                <v-icon class="input-icon" size="20">mdi-account-outline</v-icon>
                <input
                  v-model="form.last_name"
                  type="text"
                  class="form-input"
                  placeholder="Last name"
                />
              </div>
            </div>
          </div>

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

          <div class="form-group">
            <div class="input-container">
              <v-icon class="input-icon" size="20">mdi-lock-check-outline</v-icon>
              <input
                v-model="confirmPassword"
                type="password"
                required
                class="form-input"
                placeholder="Confirm password"
                :class="{ 'error': error && form.password !== confirmPassword }"
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
            class="register-button"
            :class="{ 'loading': loading }"
          >
            <div v-if="loading" class="loading-spinner"></div>
            <span>{{ loading ? 'Creating account...' : 'Create Account' }}</span>
          </button>
        </form>

        <div class="register-footer">
          <p class="footer-text">
            Already have an account?
            <router-link to="/login" class="footer-link">
              Sign in here
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
import { authService, type RegisterRequest } from '@/services/auth'

const router = useRouter()

const form = ref<RegisterRequest>({
  email: '',
  password: '',
  first_name: '',
  last_name: ''
})

const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')

const handleRegister = async () => {
  console.log('🚀 Début de la création de compte')

  if (!form.value.email || !form.value.password) {
    error.value = 'Please fill in all required fields'
    return
  }

  if (form.value.password !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  if (form.value.password.length < 6) {
    error.value = 'Password must contain at least 6 characters'
    return
  }

  console.log('📝 Données valides, démarrage de l\'enregistrement')
  loading.value = true
  error.value = ''

  try {
    console.log('🔄 Appel du service d\'authentification...')
    await authService.register(form.value)
    console.log('✅ Compte créé avec succès')
    localStorage.setItem('area_new_user', 'true')
    router.push('/')
  } catch (err) {
    console.error('❌ Erreur lors de la création:', err)
    error.value = err instanceof Error ? err.message : 'Registration error'
  } finally {
    console.log('🏁 Fin du processus de création')
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  background: var(--gradient-bg-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.register-background {
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

.register-container {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 480px;
  padding: 2rem;
}

.register-card {
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

.register-header {
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

.register-title {
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

.register-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 400;
}

.register-form {
  margin-bottom: 2rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1.5rem;
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

.register-button {
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

.register-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow:
    var(--shadow-glow),
    0 10px 20px -5px rgba(6, 182, 212, 0.5);
}

.register-button:active:not(:disabled) {
  transform: translateY(0);
}

.register-button:disabled {
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

.register-footer {
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
  .register-container {
    padding: 1rem;
  }

  .register-card {
    padding: 2rem 1.5rem;
  }

  .register-title {
    font-size: 1.75rem;
  }

  .logo-icon {
    width: 56px;
    height: 56px;
  }

  .form-row {
    grid-template-columns: 1fr;
    gap: 0;
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
  .register-card,
  .error-message {
    animation: none !important;
  }

  .register-button:hover:not(:disabled) {
    transform: none;
  }
}
</style>








