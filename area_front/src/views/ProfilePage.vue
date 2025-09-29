<template>
  <div class="profile-page">
    <div class="profile-background">
      <div class="geometric-shape shape-1"></div>
      <div class="geometric-shape shape-2"></div>
      <div class="geometric-shape shape-3"></div>
    </div>

    <div class="profile-container">
      <!-- Header -->
      <div class="profile-header">
        <button class="back-button" @click="goBack">
          <v-icon size="20">mdi-arrow-left</v-icon>
          <span>Back</span>
        </button>
        <h1 class="profile-title">Profile</h1>
        <div class="header-spacer"></div>
      </div>

      <div class="profile-card">
        <div class="profile-section">
          <div class="avatar-section">
            <div class="profile-avatar">
              <v-icon size="48" color="white">mdi-account</v-icon>
            </div>
            <button class="edit-avatar-btn" @click="requireAuth(() => {})">
              <v-icon size="16">mdi-camera</v-icon>
            </button>
          </div>

          <div class="profile-info">
            <h2 class="profile-name">{{ currentUser?.first_name || 'User' }} {{ currentUser?.last_name || 'Name' }}</h2>
            <p class="profile-email">{{ currentUser?.email || 'user@example.com' }}</p>
            <div class="profile-badges">
              <span class="badge premium">Premium Member</span>
              <span class="badge verified">Verified</span>
            </div>
          </div>
        </div>

        <div class="profile-section">
          <h3 class="section-title">Personal Information</h3>
          <div class="info-grid">
            <div class="info-item">
              <label class="info-label">First Name</label>
              <div class="info-value">{{ currentUser?.first_name || 'Not provided' }}</div>
            </div>
            <div class="info-item">
              <label class="info-label">Last Name</label>
              <div class="info-value">{{ currentUser?.last_name || 'Not provided' }}</div>
            </div>
            <div class="info-item">
              <label class="info-label">Email Address</label>
              <div class="info-value">{{ currentUser?.email || 'Not provided' }}</div>
            </div>
            <div class="info-item">
              <label class="info-label">Member Since</label>
              <div class="info-value">{{ formatDate(currentUser?.created_at) }}</div>
            </div>
          </div>
        </div>

        <div class="profile-section">
          <h3 class="section-title">Statistics</h3>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-icon">
                <v-icon size="24" color="white">mdi-vector-square</v-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">12</div>
                <div class="stat-label">Active AREAs</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">
                <v-icon size="24" color="white">mdi-clock-outline</v-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">1,247</div>
                <div class="stat-label">Executions</div>
              </div>
            </div>
            <div class="stat-card">
              <div class="stat-icon">
                <v-icon size="24" color="white">mdi-calendar-check</v-icon>
              </div>
              <div class="stat-content">
                <div class="stat-number">89%</div>
                <div class="stat-label">Success Rate</div>
              </div>
            </div>
          </div>
        </div>

        <div class="profile-section">
          <div class="linked-accounts-header">
            <h3 class="section-title">Linked Accounts</h3>
            <button class="add-account-btn" @click="requireAuth(() => {})">
              <v-icon size="16">mdi-plus</v-icon>
              <span>Add Account</span>
            </button>
          </div>

          <div class="linked-accounts-grid">
            <div
              v-for="service in enabledServices"
              :key="service.id"
              class="account-card"
              :class="{ 'linked': isServiceLinked(service.id) }"
            >
              <div class="account-header">
                <div class="account-icon" :style="{ backgroundColor: service.color }">
                  <svg v-if="service.icon === 'github'" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
                  </svg>
                  <svg v-else-if="service.icon === 'google'" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                    <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                    <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                    <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                  </svg>
                  <svg v-else-if="service.icon === 'discord'" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"/>
                  </svg>
                  <svg v-else-if="service.icon === 'spotify'" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.42 1.56-.299.421-1.02.599-1.559.3z"/>
                  </svg>
                </div>
                <div class="account-info">
                  <h4>{{ service.name }}</h4>
                  <p>{{ service.description }}</p>
                </div>
              </div>

              <div class="account-status">
                <div v-if="isServiceLinked(service.id)" class="linked-status">
                  <div class="status-icon success">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                      <polyline points="22,4 12,14.01 9,11.01"/>
                    </svg>
                  </div>
                  <span class="status-text">Connected</span>
                  <button
                    @click="unlinkService(service.id)"
                    class="unlink-btn"
                    :disabled="isLoading"
                  >
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M18 6L6 18M6 6l12 12"/>
                    </svg>
                    Unlink
                  </button>
                </div>
                <div v-else class="unlinked-status">
                  <div class="status-icon">
                    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <circle cx="12" cy="12" r="10"/>
                      <line x1="15" y1="9" x2="9" y2="15"/>
                      <line x1="9" y1="9" x2="15" y2="15"/>
                    </svg>
                  </div>
                  <span class="status-text">Not Connected</span>
                  <button
                    @click="linkService(service.id)"
                    class="link-btn"
                    :disabled="isLoading"
                  >
                    <div v-if="isLoading" class="loading-spinner"></div>
                    <svg v-else width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/>
                    </svg>
                    {{ isLoading ? 'Connecting...' : 'Connect' }}
                  </button>
                </div>
              </div>

              <div v-if="errorMessages[service.id]" class="error-message">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <line x1="15" y1="9" x2="9" y2="15"/>
                  <line x1="9" y1="9" x2="15" y2="15"/>
                </svg>
                {{ errorMessages[service.id] }}
              </div>

              <div v-if="successMessages[service.id]" class="success-message">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
                  <polyline points="22,4 12,14.01 9,11.01"/>
                </svg>
                {{ successMessages[service.id] }}
              </div>
            </div>
          </div>
        </div>

        <div class="profile-section">
          <h3 class="section-title">Account Actions</h3>
          <div class="actions-grid">
            <button class="action-button primary" @click="requireAuth(() => {})">
              <v-icon size="20">mdi-pencil</v-icon>
              <span>Edit Profile</span>
            </button>
            <button class="action-button secondary" @click="requireAuth(() => {})">
              <v-icon size="20">mdi-key</v-icon>
              <span>Change Password</span>
            </button>
            <button class="action-button secondary" @click="requireAuth(() => {})">
              <v-icon size="20">mdi-bell</v-icon>
              <span>Notifications</span>
            </button>
            <button class="action-button danger" @click="requireAuth(() => {})">
              <v-icon size="20">mdi-delete</v-icon>
              <span>Delete Account</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { authService } from '@/services/auth'
import { SERVICES_CONFIG, getEnabledServices, type ServiceConfig } from '@/config/services'

const router = useRouter()
const { currentUser, isAuthenticated, linkGitHubAccount, unlinkGitHubAccount } = useAuth()

const isLoading = ref(false)
const errorMessages = ref<Record<string, string>>({})
const successMessages = ref<Record<string, string>>({})

const enabledServices = computed(() => getEnabledServices())

const goBack = () => {
  router.push('/')
}

const requireAuth = (action: () => void) => {
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  action()
}

const formatDate = (dateString?: string) => {
  if (!dateString) return 'Unknown'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const isServiceLinked = (serviceId: string): boolean => {
  if (!currentUser.value) return false

  switch (serviceId) {
    case 'github':
      return !!currentUser.value.github_username
    case 'google':
      return !!currentUser.value.google_id
    case 'discord':
      return !!currentUser.value.discord_id
    case 'spotify':
      return !!currentUser.value.spotify_id
    default:
      return false
  }
}

const linkService = async (serviceId: string) => {
  isLoading.value = true
  errorMessages.value[serviceId] = ''
  successMessages.value[serviceId] = ''

  try {
    const service = SERVICES_CONFIG.find(s => s.id === serviceId)
    if (!service) {
      throw new Error('Service not found')
    }

    if (serviceId === 'github') {
      const githubClientId = import.meta.env.VITE_GITHUB_CLIENT_ID

      if (!githubClientId || githubClientId === 'your_github_client_id') {
        errorMessages.value[serviceId] = 'GitHub OAuth not configured. Please set the VITE_GITHUB_CLIENT_ID environment variable.'
        return
      }

      const redirectUri = encodeURIComponent(`${window.location.origin}${service.callbackPath}`)
      const githubAuthUrl = `${service.authUrl}?client_id=${githubClientId}&redirect_uri=${redirectUri}&scope=${service.scopes.join(',')}`

      window.location.href = githubAuthUrl
    } else {
      errorMessages.value[serviceId] = `${service.name} integration is not yet implemented.`
    }
  } catch (error) {
    errorMessages.value[serviceId] = `Failed to initialize ${serviceId} authentication`
    console.error(`${serviceId} auth error:`, error)
  } finally {
    isLoading.value = false
  }
}

const unlinkService = async (serviceId: string) => {
  isLoading.value = true
  errorMessages.value[serviceId] = ''
  successMessages.value[serviceId] = ''

  try {
    if (serviceId === 'github') {
      await unlinkGitHubAccount()
      successMessages.value[serviceId] = 'GitHub account unlinked successfully'
    } else {
      errorMessages.value[serviceId] = `${serviceId} unlinking is not yet implemented.`
    }
  } catch (error) {
    errorMessages.value[serviceId] = `Failed to unlink ${serviceId} account`
    console.error(`Unlink ${serviceId} error:`, error)
  } finally {
    isLoading.value = false
  }
}

const handleServiceCallback = async (serviceId: string, code: string) => {
  isLoading.value = true
  errorMessages.value[serviceId] = ''
  successMessages.value[serviceId] = ''

  try {
    if (serviceId === 'github') {
      const result = await linkGitHubAccount(code)
      successMessages.value[serviceId] = 'GitHub account linked successfully!'
    } else {
      errorMessages.value[serviceId] = `${serviceId} linking is not yet implemented.`
    }
  } catch (error) {
    errorMessages.value[serviceId] = `Failed to link ${serviceId} account`
    console.error(`Link ${serviceId} error:`, error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  const urlParams = new URLSearchParams(window.location.search)
  const code = urlParams.get('code')
  const error = urlParams.get('error')
  const service = urlParams.get('service')

  if (code && service) {
    handleServiceCallback(service, code)
  } else if (error) {
    const serviceId = urlParams.get('service') || 'unknown'
    errorMessages.value[serviceId] = decodeURIComponent(error)
  }
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: var(--gradient-bg-primary);
  position: relative;
  overflow: hidden;
}

.profile-background {
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

.profile-container {
  position: relative;
  z-index: 2;
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

.profile-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2rem;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: transparent;
  border: 2px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: var(--transition-normal);
}

.back-button:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: translateY(-1px);
}

.profile-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.02em;
  background: var(--gradient-text);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-spacer {
  width: 120px;
}

.profile-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-2xl);
  padding: 3rem;
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

.profile-section {
  margin-bottom: 3rem;
}

.profile-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 1.5rem 0;
  letter-spacing: -0.01em;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.profile-avatar {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  box-shadow: var(--shadow-glow);
}

.edit-avatar-btn {
  position: absolute;
  bottom: -5px;
  right: -5px;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-full);
  background: var(--color-bg-card);
  border: 2px solid var(--color-border-primary);
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: var(--transition-normal);
}

.edit-avatar-btn:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: scale(1.1);
}

.profile-info {
  flex: 1;
}

.profile-name {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.01em;
}

.profile-email {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0 0 1rem 0;
}

.profile-badges {
  display: flex;
  gap: 0.75rem;
}

.badge {
  padding: 0.25rem 0.75rem;
  border-radius: var(--radius-md);
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.badge.premium {
  background: var(--gradient-accent);
  color: var(--color-text-primary);
}

.badge.verified {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.info-item {
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  padding: 1.25rem;
  transition: var(--transition-normal);
}

.info-item:hover {
  background: rgba(15, 23, 42, 0.6);
  border-color: var(--color-border-secondary);
}

.info-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: 0.5rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.info-value {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.stat-card {
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: var(--transition-normal);
}

.stat-card:hover {
  background: rgba(15, 23, 42, 0.6);
  border-color: var(--color-border-secondary);
  transform: translateY(-2px);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin-bottom: 0.25rem;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.action-button {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.25rem;
  border: 2px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-normal);
  text-align: left;
}

.action-button.primary {
  background: var(--gradient-accent);
  color: var(--color-text-primary);
  border-color: transparent;
}

.action-button.primary:hover {
  transform: translateY(-2px);
  box-shadow:
    var(--shadow-glow),
    0 10px 20px -5px rgba(6, 182, 212, 0.5);
}

.action-button.secondary {
  background: transparent;
  color: var(--color-text-primary);
}

.action-button.secondary:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: translateY(-1px);
}

.action-button.danger {
  background: transparent;
  color: #ef4444;
  border-color: rgba(239, 68, 68, 0.3);
}

.action-button.danger:hover {
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.5);
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .profile-container {
    padding: 1rem;
  }

  .profile-card {
    padding: 2rem 1.5rem;
  }

  .profile-title {
    font-size: 2rem;
  }

  .profile-header {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }

  .header-spacer {
    display: none;
  }

  .avatar-section {
    flex-direction: column;
    text-align: center;
  }

  .info-grid,
  .stats-grid,
  .actions-grid {
    grid-template-columns: 1fr;
  }
}

.linked-accounts-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.add-account-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--gradient-accent);
  color: var(--color-text-primary);
  border: none;
  border-radius: var(--radius-lg);
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-normal);
  box-shadow: var(--shadow-glow);
}

.add-account-btn:hover {
  transform: translateY(-2px);
  box-shadow:
    var(--shadow-glow),
    0 10px 20px -5px rgba(6, 182, 212, 0.5);
}

.linked-accounts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.account-card {
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
  transition: var(--transition-normal);
}

.account-card:hover {
  background: rgba(15, 23, 42, 0.6);
  border-color: var(--color-border-secondary);
  transform: translateY(-2px);
}

.account-card.linked {
  border-color: rgba(16, 185, 129, 0.3);
  background: rgba(16, 185, 129, 0.05);
}

.account-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1rem;
}

.account-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.account-info h4 {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
}

.account-info p {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
  line-height: 1.4;
}

.account-status {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.linked-status,
.unlinked-status {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
}

.status-icon {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-text-secondary);
}

.status-icon.success {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.status-text {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  font-weight: 500;
  flex: 1;
}

.link-btn,
.unlink-btn {
  padding: 0.5rem 0.75rem;
  border-radius: var(--radius-md);
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: var(--transition-normal);
  display: flex;
  align-items: center;
  gap: 0.375rem;
  border: none;
}

.link-btn {
  background: var(--gradient-accent);
  color: white;
  box-shadow: var(--shadow-glow);
}

.link-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(87, 128, 232, 0.4);
}

.unlink-btn {
  background: transparent;
  border: 1px solid rgba(239, 68, 68, 0.5);
  color: #ef4444;
}

.unlink-btn:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.1);
  border-color: #ef4444;
}

.link-btn:disabled,
.unlink-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.error-message,
.success-message {
  margin-top: 0.75rem;
  padding: 0.5rem 0.75rem;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.75rem;
  font-weight: 500;
}

.error-message {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
}

.success-message {
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  color: #10b981;
}

.loading-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .linked-accounts-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .linked-accounts-grid {
    grid-template-columns: 1fr;
  }

  .account-header {
    flex-direction: column;
    text-align: center;
  }

  .account-icon {
    align-self: center;
  }
}

@media (prefers-reduced-motion: reduce) {
  .geometric-shape,
  .profile-card,
  .stat-card,
  .action-button,
  .account-card {
    animation: none !important;
  }

  .action-button:hover,
  .stat-card:hover,
  .back-button:hover,
  .account-card:hover {
    transform: none;
  }
}
</style>
