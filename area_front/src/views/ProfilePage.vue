<template>
  <div class="profile-page">
    <div class="space-background">
      <div class="stars"></div>
      <div class="stars2"></div>
      <div class="stars3"></div>
      <div class="nebula"></div>
    </div>

    <div class="page-header">
      <button class="back-btn" @click="goBack">
        <v-icon size="20">mdi-arrow-left</v-icon>
        <span>Back to Dashboard</span>
      </button>
      <div class="header-content">
        <h1 class="page-title">My Profile</h1>
        <p class="page-subtitle">Manage your account information and preferences</p>
      </div>
    </div>

    <div class="content-container">
      <div class="profile-overview">
        <div class="profile-card">
          <div class="profile-header">
            <div class="avatar-section">
              <div class="profile-avatar" @click="handleImageUpload" :class="{ 'uploading': isUploading }">
                <img
                  v-if="profileImageUrl"
                  :src="profileImageUrl"
                  alt="Profile picture"
                  class="profile-image"
                />
                <div v-else class="default-avatar">
                  <v-icon size="48" color="white">mdi-account</v-icon>
                </div>
                <div v-if="isUploading" class="upload-overlay">
                  <v-progress-circular indeterminate size="32" color="white" width="3"></v-progress-circular>
                </div>
                <div v-else class="change-overlay">
                  <div class="change-content">
                    <v-icon size="20" color="white">mdi-camera-plus</v-icon>
                    <span class="change-text">Change</span>
                  </div>
                </div>
              </div>
              <input
                ref="fileInput"
                type="file"
                accept="image/*"
                @change="onFileSelected"
                style="display: none"
              />
            </div>
            <div class="profile-info">
              <h2 class="profile-name">{{ currentUser?.first_name || 'User' }} {{ currentUser?.last_name || 'Name' }}</h2>
              <p class="profile-email">{{ currentUser?.email || 'user@example.com' }}</p>
              <div class="profile-badges">
                <span class="badge premium">
                  <v-icon size="16">mdi-crown</v-icon>
                  Premium Member
                </span>
                <span class="badge verified">
                  <v-icon size="16">mdi-check-circle</v-icon>
                  Verified
                </span>
              </div>
            </div>
          </div>
          <div v-if="uploadError" class="error-message">
            <v-icon size="16">mdi-alert-circle</v-icon>
            {{ uploadError }}
          </div>
        </div>
      </div>

      <div class="info-section">
        <div class="info-card">
          <div class="card-header">
            <h3 class="section-title">
              <v-icon size="24" class="title-icon">mdi-account-details</v-icon>
              Personal Information
            </h3>
            <p class="section-description">Your account details and membership information</p>
          </div>
          <div class="card-content">
            <div class="info-grid">
              <div class="info-item">
                <div class="info-icon">
                  <v-icon size="20">mdi-account-outline</v-icon>
                </div>
                <div class="info-content">
                  <label class="info-label">First Name</label>
                  <div class="info-value">{{ currentUser?.first_name || 'Not provided' }}</div>
                </div>
              </div>
              <div class="info-item">
                <div class="info-icon">
                  <v-icon size="20">mdi-account-outline</v-icon>
                </div>
                <div class="info-content">
                  <label class="info-label">Last Name</label>
                  <div class="info-value">{{ currentUser?.last_name || 'Not provided' }}</div>
                </div>
              </div>
              <div class="info-item">
                <div class="info-icon">
                  <v-icon size="20">mdi-email-outline</v-icon>
                </div>
                <div class="info-content">
                  <label class="info-label">Email Address</label>
                  <div class="info-value">{{ currentUser?.email || 'Not provided' }}</div>
                </div>
              </div>
              <div class="info-item">
                <div class="info-icon">
                  <v-icon size="20">mdi-calendar-plus</v-icon>
                </div>
                <div class="info-content">
                  <label class="info-label">Member Since</label>
                  <div class="info-value">{{ formatDate(currentUser?.created_at) }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="linked-accounts-section">
        <div class="linked-accounts-card">
          <div class="card-header">
            <h3 class="section-title">
              <v-icon size="24" class="title-icon">mdi-link-variant</v-icon>
              Linked Accounts
            </h3>
            <p class="section-description">Connect your external services to AREA</p>
            <button class="add-account-btn" @click="requireAuth(() => {})">
              <v-icon size="16">mdi-plus</v-icon>
              <span>Add Account</span>
            </button>
          </div>
          <div class="card-content">
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
        </div>
      </div>

      <div class="stats-section">
        <div class="stats-card">
          <div class="card-header">
            <h3 class="section-title">
              <v-icon size="24" class="title-icon">mdi-chart-line</v-icon>
              Activity Statistics
            </h3>
            <p class="section-description">Your AREA automation performance metrics</p>
          </div>
          <div class="card-content">
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
                  <div class="stat-label">Total Executions</div>
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
        </div>
      </div>

      <div class="actions-section">
        <div class="actions-card">
          <div class="card-header">
            <h3 class="section-title">
              <v-icon size="24" class="title-icon">mdi-cog</v-icon>
              Account Management
            </h3>
            <p class="section-description">Manage your account settings and preferences</p>
          </div>
          <div class="card-content">
            <div class="actions-grid">
              <button class="action-btn primary" @click="editProfile">
                <v-icon size="20">mdi-pencil</v-icon>
                <span>Edit Profile</span>
              </button>
              <button class="action-btn secondary" @click="changePassword">
                <v-icon size="20">mdi-key</v-icon>
                <span>Change Password</span>
              </button>
              <button class="action-btn secondary" @click="manageNotifications">
                <v-icon size="20">mdi-bell</v-icon>
                <span>Notifications</span>
              </button>
              <button class="action-btn danger" @click="deleteAccount">
                <v-icon size="20">mdi-delete</v-icon>
                <span>Delete Account</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { authService } from '@/services/auth'
import { SERVICES_CONFIG, getEnabledServices, type ServiceConfig } from '@/config/services'

const router = useRouter()
const { currentUser, isAuthenticated, linkGitHubAccount, unlinkGitHubAccount, linkGoogleAccount, unlinkGoogleAccount, linkFacebookAccount, unlinkFacebookAccount, linkOneDriveAccount, unlinkOneDriveAccount, linkSpotifyAccount, unlinkSpotifyAccount, linkTwitterAccount, unlinkTwitterAccount, linkSlackAccount, unlinkSlackAccount, uploadProfileImage, getProfileImageUrl, refreshProfile } = useAuth()

const fileInput = ref<HTMLInputElement | null>(null)
const profileImageUrl = ref<string | null>(null)
const isUploading = ref(false)
const uploadError = ref<string | null>(null)
const errorMessages = ref<Record<string, string>>({})
const successMessages = ref<Record<string, string>>({})
const isLoading = ref(false)

// PKCE helper functions for Twitter OAuth 2.0
function generateCodeVerifier(): string {
  const array = new Uint8Array(32)
  crypto.getRandomValues(array)
  return base64URLEncode(array)
}

async function generateCodeChallenge(verifier: string): Promise<string> {
  const encoder = new TextEncoder()
  const data = encoder.encode(verifier)
  const hash = await crypto.subtle.digest('SHA-256', data)
  return base64URLEncode(new Uint8Array(hash))
}

function base64URLEncode(array: Uint8Array): string {
  return btoa(String.fromCharCode.apply(null, Array.from(array)))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '')
}

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

const formatDate = (dateString?: string | null) => {
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
    case 'facebook':
      return !!currentUser.value.facebook_id
    case 'onedrive':
      return !!currentUser.value.onedrive_id
    case 'discord':
      return !!currentUser.value.discord_id
    case 'spotify':
      return !!currentUser.value.spotify_id
    case 'twitter':
      return !!currentUser.value.twitter_username
    case 'slack':
      return !!currentUser.value.slack_id
    default:
      return false
  }
}

const linkService = async (serviceId: string) => {
  isLoading.value = true
  errorMessages.value = { ...errorMessages.value, [serviceId]: '' }
  successMessages.value = { ...successMessages.value, [serviceId]: '' }

  try {
    const service = SERVICES_CONFIG.find(s => s.id === serviceId)
    if (!service) {
      throw new Error('Service not found')
    }

    if (serviceId === 'github') {
      const githubClientId = import.meta.env.VITE_GITHUB_CLIENT_ID

      if (!githubClientId || githubClientId === 'your_github_client_id') {
        errorMessages.value = { ...errorMessages.value, [serviceId]: 'GitHub OAuth not configured. Please set the VITE_GITHUB_CLIENT_ID environment variable.' }
        return
      }

      const redirectUri = encodeURIComponent(`${window.location.origin}${service.callbackPath}`)
      const githubAuthUrl = `${service.authUrl}?client_id=${githubClientId}&redirect_uri=${redirectUri}&scope=${service.scopes.join(',')}`

      window.location.href = githubAuthUrl
    } else if (serviceId === 'google') {
      const googleClientId = import.meta.env.VITE_GOOGLE_CLIENT_ID || 'your_google_client_id_here'

      if (!googleClientId || googleClientId === 'your_google_client_id_here') {
        errorMessages.value = { ...errorMessages.value, [serviceId]: 'Google OAuth not configured. Please set the VITE_GOOGLE_CLIENT_ID environment variable.' }
        return
      }

      const redirectUri = encodeURIComponent(`${window.location.origin}${service.callbackPath}`)
      const googleAuthUrl = `${service.authUrl}?client_id=${googleClientId}&redirect_uri=${redirectUri}&scope=${service.scopes.join(' ')}&response_type=code&access_type=offline&prompt=consent`

      window.location.href = googleAuthUrl
    } else if (serviceId === 'facebook') {
      const facebookClientId = import.meta.env.VITE_FACEBOOK_CLIENT_ID || 'your_facebook_client_id_here'

      if (!facebookClientId || facebookClientId === 'your_facebook_client_id_here') {
        errorMessages.value = { ...errorMessages.value, [serviceId]: 'Facebook OAuth not configured. Please set the VITE_FACEBOOK_CLIENT_ID environment variable.' }
        return
      }

      const redirectUri = encodeURIComponent(`${window.location.origin}${service.callbackPath}`)
      const facebookAuthUrl = `${service.authUrl}?client_id=${facebookClientId}&redirect_uri=${redirectUri}&scope=${service.scopes.join(',')}&response_type=code`

      window.location.href = facebookAuthUrl
    } else if (serviceId === 'onedrive') {
      const response = await fetch(service.authUrl!)
      const data = await response.json()

      if (data.authUrl) {
        window.location.href = data.authUrl
      } else {
        errorMessages.value = { ...errorMessages.value, [serviceId]: 'Failed to get OneDrive authorization URL' }
      }
    } else if (serviceId === 'spotify') {
      const spotifyClientId = import.meta.env.VITE_SPOTIFY_CLIENT_ID || 'your_spotify_client_id_here'

      if (!spotifyClientId || spotifyClientId === 'your_spotify_client_id_here') {
        errorMessages.value = { ...errorMessages.value, [serviceId]: 'Spotify OAuth not configured. Please set the VITE_SPOTIFY_CLIENT_ID environment variable.' }
        return
      }

      const overrideRedirect = import.meta.env.VITE_SPOTIFY_LINK_REDIRECT_URI || `${window.location.origin}${service.callbackPath}`
      const redirectUri = encodeURIComponent(overrideRedirect)
      const scopeParam = encodeURIComponent(service.scopes.join(' '))
      const spotifyAuthUrl = `${service.authUrl}?client_id=${spotifyClientId}&response_type=code&redirect_uri=${redirectUri}&scope=${scopeParam}&show_dialog=true&state=link`

      window.location.href = spotifyAuthUrl
    } else if (serviceId === 'twitter') {
      const twitterApiKey = import.meta.env.VITE_TWITTER_API_KEY || 'your_twitter_api_key_here'

      if (!twitterApiKey || twitterApiKey === 'your_twitter_api_key_here') {
        errorMessages.value = { ...errorMessages.value, [serviceId]: 'Twitter OAuth not configured. Please set the VITE_TWITTER_API_KEY environment variable.' }
        return
      }

      // Generate PKCE parameters for Twitter OAuth 2.0
      const codeVerifier = generateCodeVerifier()
      const codeChallenge = await generateCodeChallenge(codeVerifier)

      // Store code_verifier in sessionStorage for later use
      sessionStorage.setItem('twitter_code_verifier', codeVerifier)

      const redirectUri = encodeURIComponent(`${window.location.origin}${service.callbackPath}`)
      const twitterAuthUrl = `${service.authUrl}?client_id=${twitterApiKey}&response_type=code&redirect_uri=${redirectUri}&scope=${encodeURIComponent(service.scopes.join(' '))}&state=link&code_challenge=${codeChallenge}&code_challenge_method=S256`

      window.location.href = twitterAuthUrl
    } else if (serviceId === 'slack') {
      const response = await fetch(service.authUrl!)
      const data = await response.json()

      if (data.authUrl) {
        window.location.href = data.authUrl
      } else {
        errorMessages.value = { ...errorMessages.value, [serviceId]: 'Failed to get Slack authorization URL' }
      }
    } else {
      errorMessages.value = { ...errorMessages.value, [serviceId]: `${service.name} integration is not yet implemented.` }
    }
  } catch (error) {
    errorMessages.value = { ...errorMessages.value, [serviceId]: `Failed to initialize ${serviceId} authentication` }
    console.error(`${serviceId} auth error:`, error)
  } finally {
    isLoading.value = false
  }
}

const unlinkService = async (serviceId: string) => {
  isLoading.value = true
  errorMessages.value = { ...errorMessages.value, [serviceId]: '' }
  successMessages.value = { ...successMessages.value, [serviceId]: '' }

  try {
    if (serviceId === 'github') {
      await unlinkGitHubAccount()
      successMessages.value = { ...successMessages.value, [serviceId]: 'GitHub account unlinked successfully' }
    } else if (serviceId === 'google') {
      await unlinkGoogleAccount()
      successMessages.value = { ...successMessages.value, [serviceId]: 'Google account unlinked successfully' }
    } else if (serviceId === 'facebook') {
      await unlinkFacebookAccount()
      successMessages.value = { ...successMessages.value, [serviceId]: 'Facebook account unlinked successfully' }
    } else if (serviceId === 'onedrive') {
      await unlinkOneDriveAccount()
      successMessages.value = { ...successMessages.value, [serviceId]: 'OneDrive account unlinked successfully' }
    } else if (serviceId === 'spotify') {
      await unlinkSpotifyAccount()
      successMessages.value = { ...successMessages.value, [serviceId]: 'Spotify account unlinked successfully' }
    } else if (serviceId === 'twitter') {
      await unlinkTwitterAccount()
      successMessages.value = { ...successMessages.value, [serviceId]: 'Twitter account unlinked successfully' }
    } else if (serviceId === 'slack') {
      await unlinkSlackAccount()
      successMessages.value = { ...successMessages.value, [serviceId]: 'Slack workspace unlinked successfully' }
    } else {
      errorMessages.value = { ...errorMessages.value, [serviceId]: `${serviceId} unlinking is not yet implemented.` }
    }
  } catch (error) {
    errorMessages.value = { ...errorMessages.value, [serviceId]: `Failed to unlink ${serviceId} account` }
    console.error(`Unlink ${serviceId} error:`, error)
  } finally {
    isLoading.value = false
  }
}

const handleServiceCallback = async (serviceId: string, code: string) => {
  isLoading.value = true
  errorMessages.value = { ...errorMessages.value, [serviceId]: '' }
  successMessages.value = { ...successMessages.value, [serviceId]: '' }

  try {
    if (serviceId === 'github') {
      const result = await linkGitHubAccount(code)
      successMessages.value = { ...successMessages.value, [serviceId]: 'GitHub account linked successfully!' }
    } else if (serviceId === 'google') {
      const result = await linkGoogleAccount(code)
      successMessages.value = { ...successMessages.value, [serviceId]: 'Google account linked successfully!' }
    } else if (serviceId === 'facebook') {
      const result = await linkFacebookAccount(code)
      successMessages.value = { ...successMessages.value, [serviceId]: 'Facebook account linked successfully!' }
    } else if (serviceId === 'onedrive') {
      const result = await linkOneDriveAccount(code)
      successMessages.value = { ...successMessages.value, [serviceId]: 'OneDrive account linked successfully!' }
    } else if (serviceId === 'spotify') {
      const result = await linkSpotifyAccount(code)
      successMessages.value = { ...successMessages.value, [serviceId]: 'Spotify account linked successfully!' }
    } else if (serviceId === 'twitter') {
      const codeVerifier = sessionStorage.getItem('twitter_code_verifier')
      if (!codeVerifier) {
        throw new Error('PKCE verification failed')
      }
      const result = await linkTwitterAccount(code, codeVerifier)
      sessionStorage.removeItem('twitter_code_verifier')
      successMessages.value = { ...successMessages.value, [serviceId]: 'Twitter account linked successfully!' }
    } else if (serviceId === 'slack') {
      const result = await linkSlackAccount(code)
      successMessages.value = { ...successMessages.value, [serviceId]: 'Slack workspace linked successfully!' }
    } else {
      errorMessages.value = { ...errorMessages.value, [serviceId]: `${serviceId} linking is not yet implemented.` }
    }

    await refreshProfile()
  } catch (error) {
    errorMessages.value = { ...errorMessages.value, [serviceId]: `Failed to link ${serviceId} account` }
    console.error(`Link ${serviceId} error:`, error)
  } finally {
    isLoading.value = false
  }
}

const handleImageUpload = () => {
  requireAuth(() => {
    fileInput.value?.click()
  })
}

const onFileSelected = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (!file) return

  if (!file.type.startsWith('image/')) {
    uploadError.value = 'Please select a valid image file'
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    uploadError.value = 'Image must not exceed 5MB'
    return
  }

  try {
    isUploading.value = true
    uploadError.value = null

    await uploadProfileImage(file)

    profileImageUrl.value = getProfileImageUrl()

    if (target) {
      target.value = ''
    }
  } catch (error) {
    console.error('Upload error:', error)
    uploadError.value = error instanceof Error ? error.message : 'Error uploading image'
  } finally {
    isUploading.value = false
  }
}

const editProfile = () => {
  requireAuth(() => {
    router.push('/profile/edit')
  })
}

const changePassword = () => {
  requireAuth(() => {
    console.log('Change password')
  })
}

const manageNotifications = () => {
  requireAuth(() => {
    console.log('Manage notifications')
  })
}

const deleteAccount = () => {
  requireAuth(() => {
    console.log('Delete account')
  })
}

onMounted(async () => {
  const urlParams = new URLSearchParams(window.location.search)
  const code = urlParams.get('code')
  const error = urlParams.get('error')
  const service = urlParams.get('service')

  if (code && service) {
    handleServiceCallback(service, code)
  } else if (error) {
    const serviceId = urlParams.get('service') || 'unknown'
    errorMessages.value = { ...errorMessages.value, [serviceId]: decodeURIComponent(error) }
  }

  if (isAuthenticated.value) {
    await refreshProfile()
    profileImageUrl.value = getProfileImageUrl()
  }
})
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background: var(--gradient-bg-primary);
  position: relative;
  overflow-x: hidden;
}

/* Space Background */
.space-background {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  background:
    radial-gradient(ellipse at center,
      rgba(87, 128, 232, 0.05) 0%,
      rgba(135, 81, 209, 0.03) 50%,
      rgba(0, 0, 0, 0.8) 100%),
    linear-gradient(135deg,
      rgba(0, 0, 0, 0.9) 0%,
      rgba(26, 26, 51, 0.8) 25%,
      rgba(0, 0, 0, 0.9) 50%,
      rgba(26, 26, 51, 0.8) 75%,
      rgba(0, 0, 0, 0.9) 100%);
  z-index: 0;
  pointer-events: none;
}

.stars {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    radial-gradient(2px 2px at 20px 30px, #fff, transparent),
    radial-gradient(2px 2px at 40px 70px, rgba(255,255,255,0.8), transparent),
    radial-gradient(1px 1px at 90px 40px, #fff, transparent),
    radial-gradient(1px 1px at 130px 80px, rgba(255,255,255,0.6), transparent),
    radial-gradient(2px 2px at 160px 30px, #fff, transparent);
  background-repeat: repeat;
  background-size: 200px 100px;
  animation: twinkle 4s ease-in-out infinite alternate;
}

.stars2 {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    radial-gradient(1px 1px at 30px 50px, rgba(133, 206, 235, 0.8), transparent),
    radial-gradient(2px 2px at 60px 20px, rgba(135, 81, 209, 0.6), transparent),
    radial-gradient(1px 1px at 100px 60px, rgba(133, 206, 235, 0.7), transparent),
    radial-gradient(2px 2px at 140px 40px, rgba(135, 81, 209, 0.5), transparent),
    radial-gradient(1px 1px at 170px 90px, rgba(133, 206, 235, 0.8), transparent);
  background-repeat: repeat;
  background-size: 180px 120px;
  animation: twinkle 6s ease-in-out infinite alternate;
}

.stars3 {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image:
    radial-gradient(1px 1px at 25px 80px, rgba(135, 81, 209, 0.7), transparent),
    radial-gradient(2px 2px at 55px 30px, rgba(133, 206, 235, 0.5), transparent),
    radial-gradient(1px 1px at 85px 70px, rgba(135, 81, 209, 0.6), transparent),
    radial-gradient(2px 2px at 125px 50px, rgba(133, 206, 235, 0.4), transparent),
    radial-gradient(1px 1px at 155px 25px, rgba(135, 81, 209, 0.8), transparent);
  background-repeat: repeat;
  background-size: 160px 140px;
  animation: twinkle 8s ease-in-out infinite alternate;
}

@keyframes twinkle {
  0% {
    opacity: 0.8;
    transform: translateY(0px) scale(1);
  }
  50% {
    opacity: 1;
    transform: translateY(-5px) scale(1.1);
  }
  100% {
    opacity: 0.6;
    transform: translateY(0px) scale(0.9);
  }
}

/* Nebula Effect */
.nebula {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background:
    radial-gradient(ellipse at 20% 20%,
      rgba(87, 128, 232, 0.1) 0%,
      transparent 50%),
    radial-gradient(ellipse at 80% 80%,
      rgba(135, 81, 209, 0.08) 0%,
      transparent 50%),
    radial-gradient(ellipse at 60% 30%,
      rgba(133, 206, 235, 0.05) 0%,
      transparent 50%);
  animation: nebula-drift 30s ease-in-out infinite;
}

@keyframes nebula-drift {
  0%, 100% {
    transform: translate(0, 0) rotate(0deg);
    opacity: 0.6;
  }
  33% {
    transform: translate(10px, -15px) rotate(1deg);
    opacity: 0.8;
  }
  66% {
    transform: translate(-5px, 10px) rotate(-1deg);
    opacity: 1;
  }
}

.page-header {
  position: relative;
  z-index: 10;
  padding: 2rem 2rem 1rem;
  max-width: 1200px;
  margin: 0 auto;
}

.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: white;
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.header-content {
  text-align: center;
  margin-top: 2rem;
}

.page-title {
  font-size: 3rem;
  font-weight: 800;
  color: white;
  margin: 0 0 0.5rem 0;
  background: var(--gradient-accent);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -0.02em;
}

[data-theme="light"] .page-title {
  background: linear-gradient(135deg, #1a1a1a 0%, #3b82f6 50%, #8b5cf6 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  font-weight: 400;
}

[data-theme="light"] .page-subtitle {
  color: #4b5563;
}

.content-container {
  position: relative;
  z-index: 10;
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem 4rem;
  display: grid;
  gap: 2rem;
}

.profile-card,
.info-card,
.stats-card,
.actions-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: 24px;
  padding: 2rem;
  backdrop-filter: blur(20px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

[data-theme="light"] .profile-card,
[data-theme="light"] .info-card,
[data-theme="light"] .stats-card,
[data-theme="light"] .actions-card {
  background: #e5e7eb;
  border: 2px solid rgba(0, 0, 0, 0.15);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.profile-card:hover,
.info-card:hover,
.stats-card:hover,
.actions-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 30px 60px rgba(0, 0, 0, 0.15);
}

[data-theme="light"] .profile-card:hover,
[data-theme="light"] .info-card:hover,
[data-theme="light"] .stats-card:hover,
[data-theme="light"] .actions-card:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 2rem;
  margin-bottom: 1rem;
}

.avatar-section {
  position: relative;
}

.profile-avatar {
  width: 100px;
  height: 100px;
  border-radius: 50%;
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: var(--shadow-glow);
}

.profile-avatar:hover:not(.uploading) {
  transform: scale(1.1);
  box-shadow: 0 20px 40px rgba(87, 128, 232, 0.4);
}

.profile-avatar.uploading {
  cursor: not-allowed;
}

.profile-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.default-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
}

.change-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  opacity: 0;
  transition: all 0.3s ease;
}

.profile-avatar:hover .change-overlay {
  opacity: 1;
}

.change-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.change-text {
  font-size: 0.75rem;
  font-weight: 600;
  color: white;
}

.profile-info {
  flex: 1;
}

.profile-name {
  font-size: 2rem;
  font-weight: 700;
  color: white;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.01em;
}

[data-theme="light"] .profile-name {
  color: #1a1a1a;
}

.profile-email {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 1rem 0;
}

[data-theme="light"] .profile-email {
  color: rgba(0, 0, 0, 0.6);
}

.profile-badges {
  display: flex;
  gap: 0.75rem;
}

.badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
}

.badge.premium {
  background: var(--gradient-accent);
  color: white;
}

.badge.verified {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.card-header {
  margin-bottom: 2rem;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.5rem;
  font-weight: 700;
  color: white;
  margin: 0 0 0.5rem 0;
}

[data-theme="light"] .section-title {
  color: #1a1a1a;
}

.title-icon {
  color: var(--color-accent-primary);
}

.section-description {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
  margin: 0;
}

[data-theme="light"] .section-description {
  color: rgba(0, 0, 0, 0.6);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  transition: all 0.3s ease;
}

[data-theme="light"] .info-item {
  background: #d1d5db;
  border: 2px solid rgba(0, 0, 0, 0.15);
}

.info-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

[data-theme="light"] .info-item:hover {
  background: #d1d5db;
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.info-icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.info-content {
  flex: 1;
}

.info-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  margin-bottom: 0.25rem;
}

[data-theme="light"] .info-label {
  color: rgba(0, 0, 0, 0.5);
}

.info-value {
  font-size: 1rem;
  font-weight: 600;
  color: white;
}

[data-theme="light"] .info-value {
  color: #1a1a1a;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  transition: all 0.3s ease;
}

[data-theme="light"] .stat-card {
  background: #d1d5db;
  border: 2px solid rgba(0, 0, 0, 0.15);
}

.stat-card:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-3px);
}

[data-theme="light"] .stat-card:hover {
  background: #d1d5db;
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 16px;
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
  font-size: 1.75rem;
  font-weight: 700;
  color: white;
  margin-bottom: 0.25rem;
}

[data-theme="light"] .stat-number {
  color: #1a1a1a;
}

.stat-label {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
  font-weight: 500;
}

[data-theme="light"] .stat-label {
  color: rgba(0, 0, 0, 0.6);
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  border-radius: 16px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  text-align: left;
}

.action-btn.primary {
  background: var(--gradient-accent);
  color: white;
  box-shadow: 0 10px 30px rgba(87, 128, 232, 0.3);
}

.action-btn.primary:hover {
  transform: translateY(-3px);
  box-shadow: 0 20px 40px rgba(87, 128, 232, 0.4);
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
}

[data-theme="light"] .action-btn.secondary {
  background: #ffffff;
  border: 2px solid rgba(0, 0, 0, 0.2);
  color: #1a1a1a;
}

.action-btn.secondary:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

[data-theme="light"] .action-btn.secondary:hover {
  background: #f9fafb;
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.action-btn.danger {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
}

.action-btn.danger:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: rgba(239, 68, 68, 0.5);
  transform: translateY(-2px);
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
  padding: 0.75rem 1rem;
  border-radius: 12px;
  font-size: 0.875rem;
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .page-header {
    padding: 1rem;
  }

  .page-title {
    font-size: 2rem;
  }

  .content-container {
    padding: 0 1rem 2rem;
    gap: 1.5rem;
  }

  .profile-card,
  .info-card,
  .stats-card,
  .actions-card {
    padding: 1.5rem;
  }

  .profile-header {
    flex-direction: column;
    text-align: center;
    gap: 1rem;
  }

  .info-grid,
  .stats-grid,
  .actions-grid {
    grid-template-columns: 1fr;
  }
}

.linked-accounts-card {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: 24px;
  padding: 2rem;
  backdrop-filter: blur(20px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

[data-theme="light"] .linked-accounts-card {
  background: #e5e7eb;
  border: 2px solid rgba(0, 0, 0, 0.15);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.linked-accounts-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 30px 60px rgba(0, 0, 0, 0.15);
}

[data-theme="light"] .linked-accounts-card:hover {
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
}

.add-account-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: var(--gradient-accent);
  color: white;
  border: none;
  border-radius: 12px;
  padding: 0.75rem 1rem;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px rgba(87, 128, 232, 0.3);
  margin-top: 1rem;
}

.add-account-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 20px 40px rgba(87, 128, 232, 0.4);
}

.linked-accounts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.account-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  padding: 1.5rem;
  transition: all 0.3s ease;
}

[data-theme="light"] .account-card {
  background: #e5e7eb;
  border: 2px solid rgba(0, 0, 0, 0.15);
}

.account-card:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
}

[data-theme="light"] .account-card:hover {
  background: #d1d5db;
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
}

.account-card.linked {
  border-color: rgba(16, 185, 129, 0.3);
  background: rgba(16, 185, 129, 0.05);
}

[data-theme="light"] .account-card.linked {
  border-color: rgba(16, 185, 129, 0.4);
  background: rgba(16, 185, 129, 0.08);
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
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  flex-shrink: 0;
}

.account-info h4 {
  font-size: 1.125rem;
  font-weight: 600;
  color: white;
  margin: 0 0 0.25rem 0;
}

[data-theme="light"] .account-info h4 {
  color: #1a1a1a;
}

.account-info p {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  line-height: 1.4;
}

[data-theme="light"] .account-info p {
  color: rgba(0, 0, 0, 0.7);
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
  color: rgba(255, 255, 255, 0.6);
}

.status-icon.success {
  background: rgba(16, 185, 129, 0.2);
  color: #10b981;
}

.status-text {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
  font-weight: 500;
  flex: 1;
}

[data-theme="light"] .status-text {
  color: rgba(0, 0, 0, 0.7);
}

.link-btn,
.unlink-btn {
  padding: 0.5rem 0.75rem;
  border-radius: 8px;
  font-size: 0.75rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  gap: 0.375rem;
  border: none;
}

.link-btn {
  background: var(--gradient-accent);
  color: white;
  box-shadow: 0 10px 30px rgba(87, 128, 232, 0.3);
}

[data-theme="light"] .link-btn {
  box-shadow: 0 4px 12px rgba(87, 128, 232, 0.4);
}

.link-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 20px rgba(87, 128, 232, 0.4);
}

[data-theme="light"] .link-btn:hover:not(:disabled) {
  box-shadow: 0 6px 20px rgba(87, 128, 232, 0.5);
}

.unlink-btn {
  background: transparent;
  border: 1px solid rgba(239, 68, 68, 0.5);
  color: #ef4444;
}

[data-theme="light"] .unlink-btn {
  border: 2px solid rgba(239, 68, 68, 0.7);
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
  border-radius: 8px;
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

@media (max-width: 480px) {
  .page-title {
    font-size: 1.75rem;
  }

  .profile-avatar {
    width: 80px;
    height: 80px;
  }

  .section-title {
    font-size: 1.25rem;
  }
}

@media (prefers-reduced-motion: reduce) {
  .stars,
  .stars2,
  .stars3,
  .nebula,
  .profile-card,
  .stat-card,
  .action-btn,
  .account-card {
    animation: none !important;
  }

  .action-btn:hover,
  .stat-card:hover,
  .back-btn:hover,
  .account-card:hover {
    transform: none;
  }
}
</style>






