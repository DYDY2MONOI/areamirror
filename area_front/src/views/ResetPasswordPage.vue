<template>
  <div class="reset-password-page">
    <div class="animated-background">
      <div class="floating-shapes">
        <div class="shape shape-1"></div>
        <div class="shape shape-2"></div>
        <div class="shape shape-3"></div>
        <div class="shape shape-4"></div>
        <div class="shape shape-5"></div>
      </div>
      <div class="gradient-overlay"></div>
    </div>

    <div class="content-container">
      <div class="reset-password-card">
        <div class="card-header">
          <div class="logo-container">
            <img src="/src/assets/logo.svg" alt="AREA Logo" class="logo" />
          </div>
          <h1 class="page-title">Nouveau mot de passe</h1>
          <p class="page-subtitle">Entrez votre nouveau mot de passe</p>
        </div>

        <div class="card-content">
          <form @submit.prevent="handleResetPassword" class="reset-password-form">
            <div class="form-group">
              <label for="newPassword" class="form-label">Nouveau mot de passe</label>
              <input
                id="newPassword"
                v-model="newPassword"
                type="password"
                class="form-input"
                :class="{ 'error': passwordError }"
                placeholder="Votre nouveau mot de passe"
                required
                :disabled="isLoading"
              />
              <div v-if="passwordError" class="error-message">
                <v-icon size="16" color="error">mdi-alert-circle</v-icon>
                {{ passwordError }}
              </div>
            </div>

            <div class="form-group">
              <label for="confirmPassword" class="form-label">Confirmer le mot de passe</label>
              <input
                id="confirmPassword"
                v-model="confirmPassword"
                type="password"
                class="form-input"
                :class="{ 'error': confirmPasswordError }"
                placeholder="Confirmez votre nouveau mot de passe"
                required
                :disabled="isLoading"
              />
              <div v-if="confirmPasswordError" class="error-message">
                <v-icon size="16" color="error">mdi-alert-circle</v-icon>
                {{ confirmPasswordError }}
              </div>
            </div>

            <button
              type="submit"
              class="submit-btn"
              :disabled="isLoading || !newPassword || !confirmPassword"
            >
              <div v-if="isLoading" class="loading-spinner"></div>
              <v-icon v-else size="20">mdi-key</v-icon>
              {{ isLoading ? 'Réinitialisation...' : 'Réinitialiser le mot de passe' }}
            </button>
          </form>

          <div v-if="successMessage" class="success-message">
            <v-icon size="20" color="success">mdi-check-circle</v-icon>
            {{ successMessage }}
          </div>

          <div v-if="errorMessage" class="error-message">
            <v-icon size="20" color="error">mdi-alert-circle</v-icon>
            {{ errorMessage }}
          </div>

          <div class="back-to-login">
            <p>Vous vous souvenez de votre mot de passe ?</p>
            <router-link to="/login" class="login-link">
              <v-icon size="16">mdi-arrow-left</v-icon>
              Retour à la connexion
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const route = useRoute()
const { resetPassword } = useAuth()

const newPassword = ref('')
const confirmPassword = ref('')
const passwordError = ref('')
const confirmPasswordError = ref('')
const successMessage = ref('')
const errorMessage = ref('')
const isLoading = ref(false)
const resetToken = ref('')

onMounted(() => {
  // Récupérer le token depuis l'URL
  const token = route.query.token as string
  if (!token) {
    errorMessage.value = 'Token de réinitialisation manquant ou invalide'
    return
  }
  resetToken.value = token
})

const validatePasswords = () => {
  passwordError.value = ''
  confirmPasswordError.value = ''

  if (newPassword.value.length < 6) {
    passwordError.value = 'Le mot de passe doit contenir au moins 6 caractères'
    return false
  }

  if (newPassword.value !== confirmPassword.value) {
    confirmPasswordError.value = 'Les mots de passe ne correspondent pas'
    return false
  }

  return true
}

const handleResetPassword = async () => {
  if (!validatePasswords()) {
    return
  }

  if (!resetToken.value) {
    errorMessage.value = 'Token de réinitialisation manquant'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    await resetPassword(resetToken.value, newPassword.value)
    successMessage.value = 'Votre mot de passe a été réinitialisé avec succès !'
    
    // Rediriger vers la page de connexion après 3 secondes
    setTimeout(() => {
      router.push('/login')
    }, 3000)
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Une erreur est survenue'
  } finally {
    isLoading.value = false
  }
}
</script>

<style scoped>
.reset-password-page {
  min-height: 100vh;
  background: var(--gradient-bg-primary);
  position: relative;
  overflow-x: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
}

.animated-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
  pointer-events: none;
}

.floating-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: linear-gradient(45deg, rgba(87, 128, 232, 0.1), rgba(135, 81, 209, 0.1));
  filter: blur(2px);
  animation: float 8s ease-in-out infinite;
}

.shape-1 {
  width: 300px;
  height: 300px;
  top: 10%;
  left: 5%;
  animation-delay: 0s;
}

.shape-2 {
  width: 200px;
  height: 200px;
  top: 20%;
  right: 10%;
  animation-delay: 2s;
}

.shape-3 {
  width: 150px;
  height: 150px;
  bottom: 30%;
  left: 15%;
  animation-delay: 4s;
}

.shape-4 {
  width: 250px;
  height: 250px;
  bottom: 10%;
  right: 20%;
  animation-delay: 6s;
}

.shape-5 {
  width: 100px;
  height: 100px;
  top: 50%;
  left: 50%;
  animation-delay: 1s;
}

.gradient-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(87, 128, 232, 0.1) 0%, rgba(135, 81, 209, 0.1) 100%);
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  33% { transform: translateY(-20px) rotate(120deg); }
  66% { transform: translateY(10px) rotate(240deg); }
}

.content-container {
  position: relative;
  z-index: 2;
  width: 100%;
  max-width: 480px;
}

.reset-password-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-2xl);
  padding: 3rem 2.5rem;
  backdrop-filter: blur(20px);
  box-shadow:
    0 20px 25px -5px rgba(0, 0, 0, 0.1),
    0 10px 10px -5px rgba(0, 0, 0, 0.04),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
}

.card-header {
  text-align: center;
  margin-bottom: 2rem;
}

.logo-container {
  margin-bottom: 1.5rem;
}

.logo {
  height: 48px;
  width: auto;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.page-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 400;
}

.reset-password-form {
  margin-bottom: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 0.5rem;
}

.form-input {
  width: 100%;
  padding: 0.875rem 1rem;
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  background: var(--color-bg-input);
  color: var(--color-text-primary);
  font-size: 1rem;
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--color-accent-primary);
  box-shadow: 0 0 0 3px rgba(87, 128, 232, 0.1);
}

.form-input.error {
  border-color: var(--color-error);
}

.form-input:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: var(--color-error);
}

.submit-btn {
  width: 100%;
  padding: 0.875rem 1.5rem;
  background: var(--color-accent-primary);
  color: white;
  border: none;
  border-radius: var(--radius-lg);
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.submit-btn:hover:not(:disabled) {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
}

.submit-btn:disabled {
  opacity: 0.6;
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
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.success-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.2);
  border-radius: var(--radius-lg);
  margin-bottom: 1.5rem;
  font-size: 0.875rem;
  color: var(--color-success);
}

.back-to-login {
  text-align: center;
  padding-top: 1.5rem;
  border-top: 1px solid var(--color-border-primary);
}

.back-to-login p {
  margin: 0 0 1rem 0;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
}

.login-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--color-accent-primary);
  text-decoration: none;
  font-weight: 600;
  font-size: 0.875rem;
  transition: color 0.2s ease;
}

.login-link:hover {
  color: var(--color-accent-hover);
}
</style>
