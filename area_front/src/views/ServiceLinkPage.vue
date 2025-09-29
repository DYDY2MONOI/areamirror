<template>
  <div class="service-link-page">
    <v-navigation-drawer class="sidebar-desktop text-white" color="#0d0d0d" elevation="0" permanent rail>
      <v-list class="text-white" density="comfortable" nav lines="false">
        <v-tooltip text="Home" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-home" class="text-white" rounded @click="goToHome"></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Search" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-magnify" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Create" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-plus" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Library" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-book-open-variant" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Profile" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-account-circle" class="text-white" rounded @click="goToProfile"></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Services" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-link-variant" class="text-white" rounded active></v-list-item>
          </template>
        </v-tooltip>

        <v-spacer></v-spacer>

        <v-tooltip text="Sign Out" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-logout" class="text-white" rounded @click="logout"></v-list-item>
          </template>
        </v-tooltip>
      </v-list>
    </v-navigation-drawer>

    <div class="content">
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

      <div class="page-header">
        <div class="header-content">
          <div class="header-text">
            <h1 class="page-title">Service Integration</h1>
            <p class="page-subtitle">Connect your favorite services to create powerful automations</p>
          </div>
          <div class="header-stats">
            <div class="stat-card">
              <div class="stat-number">{{ linkedServicesCount }}</div>
              <div class="stat-label">Connected</div>
            </div>
            <div class="stat-card">
              <div class="stat-number">{{ enabledServices.length }}</div>
              <div class="stat-label">Available</div>
            </div>
          </div>
        </div>
      </div>

      <div class="features-section">
        <div class="features-header">
          <h2>Why Connect Services?</h2>
          <p>Unlock the power of automation by connecting your favorite platforms</p>
        </div>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
              </svg>
            </div>
            <h3>Automation</h3>
            <p>Create powerful workflows that run automatically based on triggers</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z"/>
              </svg>
            </div>
            <h3>Integration</h3>
            <p>Seamlessly connect different platforms and services</p>
          </div>
          <div class="feature-card">
            <div class="feature-icon">
              <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M9 12l2 2 4-4"/>
                <path d="M21 12c-1 0-3-1-3-3s2-3 3-3 3 1 3 3-2 3-3 3"/>
                <path d="M3 12c1 0 3-1 3-3s-2-3-3-3-3 1-3 3 2 3 3 3"/>
                <path d="M12 3c0 1-1 3-3 3s-3-2-3-3 1-3 3-3 3 2 3 3"/>
                <path d="M12 21c0-1 1-3 3-3s3 2 3 3-1 3-3 3-3-2-3-3"/>
              </svg>
            </div>
            <h3>Efficiency</h3>
            <p>Save time and reduce manual work with smart automations</p>
          </div>
        </div>
      </div>

      <div class="services-section">
        <div class="services-header">
          <h2>Available Services</h2>
          <p>Connect your accounts to start building automations</p>
        </div>
        <div class="services-grid">
        <div
          v-for="service in enabledServices"
          :key="service.id"
          class="service-card"
          :class="{ 'linked': isServiceLinked(service.id) }"
        >
          <div class="service-header">
            <div class="service-icon" :style="{ backgroundColor: service.color }">
              <svg v-if="service.icon === 'github'" width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
              </svg>
              <svg v-else-if="service.icon === 'google'" width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
                <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              <svg v-else-if="service.icon === 'discord'" width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
                <path d="M20.317 4.37a19.791 19.791 0 0 0-4.885-1.515.074.074 0 0 0-.079.037c-.21.375-.444.864-.608 1.25a18.27 18.27 0 0 0-5.487 0 12.64 12.64 0 0 0-.617-1.25.077.077 0 0 0-.079-.037A19.736 19.736 0 0 0 3.677 4.37a.07.07 0 0 0-.032.027C.533 9.046-.32 13.58.099 18.057a.082.082 0 0 0 .031.057 19.9 19.9 0 0 0 5.993 3.03.078.078 0 0 0 .084-.028c.462-.63.874-1.295 1.226-1.994a.076.076 0 0 0-.041-.106 13.107 13.107 0 0 1-1.872-.892.077.077 0 0 1-.008-.128 10.2 10.2 0 0 0 .372-.292.074.074 0 0 1 .077-.01c3.928 1.793 8.18 1.793 12.062 0a.074.074 0 0 1 .078.01c.12.098.246.198.373.292a.077.077 0 0 1-.006.127 12.299 12.299 0 0 1-1.873.892.077.077 0 0 0-.041.107c.36.698.772 1.362 1.225 1.993a.076.076 0 0 0 .084.028 19.839 19.839 0 0 0 6.002-3.03.077.077 0 0 0 .032-.054c.5-5.177-.838-9.674-3.549-13.66a.061.061 0 0 0-.031-.03zM8.02 15.33c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.956-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.956 2.418-2.157 2.418zm7.975 0c-1.183 0-2.157-1.085-2.157-2.419 0-1.333.955-2.419 2.157-2.419 1.21 0 2.176 1.096 2.157 2.42 0 1.333-.946 2.418-2.157 2.418z"/>
              </svg>
              <svg v-else-if="service.icon === 'spotify'" width="32" height="32" viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 0C5.4 0 0 5.4 0 12s5.4 12 12 12 12-5.4 12-12S18.66 0 12 0zm5.521 17.34c-.24.359-.66.48-1.021.24-2.82-1.74-6.36-2.101-10.561-1.141-.418.122-.779-.179-.899-.539-.12-.421.18-.78.54-.9 4.56-1.021 8.52-.6 11.64 1.32.42.18.479.659.301 1.02zm1.44-3.3c-.301.42-.841.6-1.262.3-3.239-1.98-8.159-2.58-11.939-1.38-.479.12-1.02-.12-1.14-.6-.12-.48.12-1.021.6-1.141C9.6 9.9 15 10.561 18.72 12.84c.361.181.54.78.241 1.2zm.12-3.36C15.24 8.4 8.82 8.16 5.16 9.301c-.6.179-1.2-.181-1.38-.721-.18-.601.18-1.2.72-1.381 4.26-1.26 11.28-1.02 15.721 1.621.539.3.719 1.02.42 1.56-.299.421-1.02.599-1.559.3z"/>
              </svg>
            </div>
            <div class="service-info">
              <h3>{{ service.name }}</h3>
              <p>{{ service.description }}</p>
            </div>
          </div>

          <div class="service-status">
            <div v-if="isServiceLinked(service.id)" class="linked-status">
              <div class="status-icon success">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M18 6L6 18M6 6l12 12"/>
                </svg>
                Unlink
              </button>
            </div>
            <div v-else class="unlinked-status">
              <div class="status-icon">
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
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
                <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/>
                </svg>
                {{ isLoading ? 'Connecting...' : 'Connect' }}
              </button>
            </div>
          </div>

          <div v-if="errorMessages[service.id]" class="error-message">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="15" y1="9" x2="9" y2="15"/>
              <line x1="9" y1="9" x2="15" y2="15"/>
            </svg>
            {{ errorMessages[service.id] }}
          </div>

          <div v-if="successMessages[service.id]" class="success-message">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
              <polyline points="22,4 12,14.01 9,11.01"/>
            </svg>
            {{ successMessages[service.id] }}
          </div>
        </div>

        <div v-if="enabledServices.length === 0" class="no-services">
          <div class="no-services-icon">
            <svg width="64" height="64" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
          </div>
          <h3>No Services Available</h3>
          <p>There are currently no services available for integration.</p>
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
const { currentUser, linkGitHubAccount, unlinkGitHubAccount } = useAuth()

const isLoading = ref(false)
const errorMessages = ref<Record<string, string>>({})
const successMessages = ref<Record<string, string>>({})

const enabledServices = computed(() => getEnabledServices())

const linkedServicesCount = computed(() => {
  if (!currentUser.value) return 0
  let count = 0
  if (currentUser.value.github_username) count++
  if (currentUser.value.google_id) count++
  if (currentUser.value.discord_id) count++
  if (currentUser.value.spotify_id) count++
  return count
})

const goBack = () => {
  router.go(-1)
}

const goToHome = () => {
  router.push('/')
}

const goToProfile = () => {
  router.push('/profile')
}

const logout = async () => {
  try {
    await authService.logout()
    router.push('/login')
  } catch (error) {
    console.error('Logout error:', error)
  }
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
  // Handle callback parameters
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
.service-link-page {
  min-height: 100vh;
  position: relative;
  overflow: hidden;
  background: var(--gradient-bg-primary);
  display: flex;
}

.sidebar-desktop {
  position: fixed;
  left: 0;
  top: 0;
  height: 100vh;
  z-index: 1000;
}

.content {
  flex: 1;
  margin-left: 80px;
  position: relative;
  min-height: 100vh;
}

.animated-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
}

.floating-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: var(--gradient-accent);
  opacity: 0.1;
  animation: float 20s infinite ease-in-out;
}

.shape-1 {
  width: 80px;
  height: 80px;
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 120px;
  height: 120px;
  top: 60%;
  right: 15%;
  animation-delay: 5s;
}

.shape-3 {
  width: 60px;
  height: 60px;
  top: 40%;
  left: 80%;
  animation-delay: 10s;
}

.shape-4 {
  width: 100px;
  height: 100px;
  top: 80%;
  left: 20%;
  animation-delay: 15s;
}

.shape-5 {
  width: 140px;
  height: 140px;
  top: 10%;
  right: 30%;
  animation-delay: 8s;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

.gradient-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, rgba(0,0,0,0.3) 0%, rgba(0,0,0,0.1) 100%);
}

.page-header {
  position: relative;
  z-index: 10;
  padding: 2rem;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 2rem;
}

.header-text {
  flex: 1;
}

.page-title {
  font-size: 3rem;
  font-weight: 700;
  background: var(--gradient-accent);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin: 0 0 0.5rem 0;
}

.page-subtitle {
  font-size: 1.25rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  font-weight: 400;
}

.header-stats {
  display: flex;
  gap: 1rem;
}

.stat-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1rem;
  padding: 1.5rem;
  text-align: center;
  backdrop-filter: blur(10px);
  min-width: 120px;
}

.stat-number {
  font-size: 2rem;
  font-weight: 700;
  color: var(--color-accent-primary);
  margin-bottom: 0.25rem;
}

.stat-label {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
  font-weight: 500;
}

.features-section {
  position: relative;
  z-index: 10;
  padding: 0 2rem 4rem;
}

.features-header {
  text-align: center;
  margin-bottom: 3rem;
}

.features-header h2 {
  font-size: 2.5rem;
  font-weight: 700;
  color: white;
  margin: 0 0 1rem 0;
}

.features-header p {
  font-size: 1.25rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.feature-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1.5rem;
  padding: 2rem;
  text-align: center;
  backdrop-filter: blur(20px);
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.4);
  border-color: rgba(87, 128, 232, 0.3);
}

.feature-icon {
  width: 64px;
  height: 64px;
  background: var(--gradient-accent);
  border-radius: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1.5rem;
  color: white;
}

.feature-card h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: white;
  margin: 0 0 1rem 0;
}

.feature-card p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1rem;
  line-height: 1.6;
  margin: 0;
}

.services-section {
  position: relative;
  z-index: 10;
  padding: 0 2rem 4rem;
}

.services-header {
  text-align: center;
  margin-bottom: 3rem;
}

.services-header h2 {
  font-size: 2.5rem;
  font-weight: 700;
  color: white;
  margin: 0 0 1rem 0;
}

.services-header p {
  font-size: 1.25rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.content-container {
  position: relative;
  z-index: 10;
  padding: 0 2rem 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 2rem;
}

.service-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 1.5rem;
  padding: 2rem;
  backdrop-filter: blur(20px);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  transition: all 0.3s ease;
}

.service-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.4);
}

.service-card.linked {
  border-color: rgba(16, 185, 129, 0.3);
  background: rgba(16, 185, 129, 0.05);
}

.service-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.service-icon {
  width: 64px;
  height: 64px;
  border-radius: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.service-info h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: white;
  margin: 0 0 0.5rem 0;
}

.service-info p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.875rem;
  line-height: 1.5;
  margin: 0;
}

.service-status {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.linked-status,
.unlinked-status {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  width: 100%;
}

.status-icon {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.6);
}

.status-icon.success {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.status-text {
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.875rem;
  font-weight: 500;
  flex: 1;
}

.link-btn,
.unlink-btn {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  border: none;
}

.link-btn {
  background: var(--gradient-accent);
  color: white;
  box-shadow: var(--shadow-glow);
}

.link-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(87, 128, 232, 0.4);
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
  margin-top: 1rem;
  padding: 0.75rem;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
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
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.no-services {
  text-align: center;
  padding: 4rem 2rem;
  color: rgba(255, 255, 255, 0.7);
}

.no-services-icon {
  color: rgba(255, 255, 255, 0.4);
  margin-bottom: 1rem;
}

.no-services h3 {
  font-size: 1.5rem;
  font-weight: 600;
  color: white;
  margin-bottom: 0.5rem;
}

.no-services p {
  font-size: 1rem;
  margin: 0;
}

@media (max-width: 768px) {
  .content {
    margin-left: 0;
  }

  .sidebar-desktop {
    display: none;
  }

  .page-header {
    padding: 1rem;
  }

  .header-content {
    flex-direction: column;
    gap: 1rem;
  }

  .page-title {
    font-size: 2rem;
  }

  .page-subtitle {
    font-size: 1rem;
  }

  .header-stats {
    justify-content: center;
  }

  .features-section,
  .services-section {
    padding: 0 1rem 2rem;
  }

  .features-header h2,
  .services-header h2 {
    font-size: 2rem;
  }

  .features-header p,
  .services-header p {
    font-size: 1rem;
  }

  .features-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .content-container {
    padding: 0 1rem 1rem;
  }

  .services-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .service-card {
    padding: 1.5rem;
  }

  .service-header {
    flex-direction: column;
    text-align: center;
  }

  .service-icon {
    align-self: center;
  }
}
</style>
