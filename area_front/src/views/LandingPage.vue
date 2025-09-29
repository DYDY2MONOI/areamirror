<template>
  <div class="landing-dark">
    <v-navigation-drawer class="sidebar-desktop text-white" color="#0d0d0d" elevation="0" permanent rail>
      <!-- Section utilisateur dans la sidebar -->
      <div class="sidebar-user-section" v-if="isAuthenticated">
        <v-avatar size="32" class="sidebar-avatar">
          <img
            v-if="getProfileImageUrl()"
            :src="getProfileImageUrl() || ''"
            alt="Profile picture"
            class="sidebar-profile-image"
          />
          <v-icon v-else color="white" size="20">mdi-account</v-icon>
        </v-avatar>
        <div class="sidebar-user-info">
          <div class="sidebar-user-name">{{ currentUser?.first_name || 'User' }}</div>
        </div>
      </div>

      <v-list class="text-white" density="comfortable" nav lines="false">
        <v-tooltip text="Home" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-home" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Search" location="end">
          <template #activator="{ props }">
            <v-list-item v-bind="props" prepend-icon="mdi-magnify" class="text-white" rounded></v-list-item>
          </template>
        </v-tooltip>
        <SidebarButton tooltip="Create" @open="() => requireAuth(() => showCreateModal = true)" />
        <v-tooltip text="Library" location="end">
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-book-open-variant"
              class="text-white"
              rounded
              @click="requireAuth(() => {})"
            ></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="Profile" location="end">
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-account-circle"
              class="text-white"
              rounded
              @click="requireAuth(() => router.push('/profile'))"
            ></v-list-item>
          </template>
        </v-tooltip>
        <v-tooltip text="GitHub Link" location="end">
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-github"
              class="text-white"
              rounded
              @click="requireAuth(() => router.push('/github-link'))"
            ></v-list-item>
          </template>
        </v-tooltip>

        <v-spacer></v-spacer>

        <v-tooltip text="Connexion" location="end" v-if="!isAuthenticated">
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-login"
              class="text-white"
              rounded
              @click="goToLogin"
            ></v-list-item>
          </template>
        </v-tooltip>

        <v-tooltip text="Sign Out" location="end" v-if="isAuthenticated">
          <template #activator="{ props }">
            <v-list-item
              v-bind="props"
              prepend-icon="mdi-logout"
              class="text-white"
              rounded
              @click="showLogoutDialog = true"
            ></v-list-item>
          </template>
        </v-tooltip>
      </v-list>
    </v-navigation-drawer>

    <div class="content">
    <div class="search-section">
      <div class="search-container">
        <div class="search-header">
          <h1 class="search-title">Find Your Perfect Automation</h1>
          <p class="search-subtitle">Discover templates, browse services, or create something new</p>
        </div>
        <div class="search-bar">
          <div class="search-input-container">
            <v-icon size="20" color="#9ca3af" class="search-icon">mdi-magnify</v-icon>
            <input type="text" placeholder="Search automations, services, or templates..." class="search-input">
            <button class="search-filter-btn">
              <v-icon size="16" color="#9ca3af">mdi-tune</v-icon>
            </button>
          </div>
          <div class="search-suggestions">
            <span class="suggestion-label">Popular:</span>
            <button class="suggestion-chip" @click="requireAuth(() => {})">Gmail</button>
            <button class="suggestion-chip" @click="requireAuth(() => {})">Discord</button>
            <button class="suggestion-chip" @click="requireAuth(() => {})">Spotify</button>
            <button class="suggestion-chip" @click="requireAuth(() => {})">GitHub</button>
          </div>
        </div>
      </div>
    </div>

    <v-container class="pt-6 pb-4">
      <div class="d-flex align-center justify-space-between">
        <div class="user-section" v-if="isAuthenticated">
          <v-avatar size="48" class="gradient-avatar">
            <img
              v-if="getProfileImageUrl()"
              :src="getProfileImageUrl() || ''"
              alt="Profile picture"
              class="profile-image"
            />
            <v-icon v-else color="white">mdi-account</v-icon>
          </v-avatar>
          <div class="user-info">
            <span class="user-name">{{ currentUser?.first_name || 'User' }} {{ currentUser?.last_name || 'Name' }}</span>
            <span class="user-status">Premium Member</span>
          </div>
        </div>

        <div class="guest-section" v-else>
          <div class="guest-content">
            <div class="guest-icon">
              <v-icon size="32" color="white">mdi-account-plus</v-icon>
            </div>
            <div class="guest-text">
              <h3 class="guest-title">Join AREA Today</h3>
              <p class="guest-subtitle">Start automating your workflow</p>
            </div>
          </div>
          <div class="guest-actions">
            <button class="guest-btn primary" @click="goToLogin">
              <v-icon size="16">mdi-login</v-icon>
              <span>Sign In</span>
            </button>
            <button class="guest-btn secondary" @click="router.push('/register')">
              <v-icon size="16">mdi-account-plus</v-icon>
              <span>Join Us</span>
            </button>
          </div>
        </div>
        <div class="filter-tabs">
          <button class="filter-tab active" @click="requireAuth(() => {})">All</button>
          <button class="filter-tab" @click="requireAuth(() => {})">My AREAs</button>
          <button class="filter-tab" @click="requireAuth(() => {})">Popular</button>
          <button class="filter-tab" @click="requireAuth(() => {})">Templates</button>
        </div>
        <div class="action-buttons">
          <button class="action-btn-icon" @click="requireAuth(() => {})">
            <v-icon size="20">mdi-magnify</v-icon>
          </button>
          <button class="action-btn-icon" @click="requireAuth(() => {})">
            <v-icon size="20">mdi-bell-outline</v-icon>
          </button>
        </div>
      </div>
    </v-container>

    <v-container>
      <div class="section-header">
        <div class="section-info">
          <h2 class="section-title">Popular AREAs</h2>
          <p class="section-subtitle">Most used automation templates</p>
        </div>
        <button class="view-all-btn" @click="requireAuth(() => {})">
          <span>View All</span>
          <v-icon size="16">mdi-arrow-right</v-icon>
        </button>
      </div>
      <div class="cards-grid">
        <div class="card-col">
          <v-sheet class="area-card gradient-red" rounded="xl">
            <v-icon size="64" color="white">mdi-email-outline</v-icon>
          </v-sheet>
          <div class="card-title">Gmail → Discord</div>
          <div class="card-subtitle">Auto notification</div>
          <div class="card-description">Send Discord message when you receive important emails</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-green" rounded="xl">
            <v-icon size="64" color="white">mdi-music-note</v-icon>
          </v-sheet>
          <div class="card-title">Spotify → Twitter</div>
          <div class="card-subtitle">Auto sharing</div>
          <div class="card-description">Automatically tweet your favorite tracks</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-indigo" rounded="xl">
            <v-icon size="64" color="white">mdi-github</v-icon>
          </v-sheet>
          <div class="card-title">GitHub → Slack</div>
          <div class="card-subtitle">Team alerts</div>
          <div class="card-description">Notify your Slack channel when issues are opened</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-crimson" rounded="xl">
            <v-icon size="64" color="white">mdi-youtube</v-icon>
          </v-sheet>
          <div class="card-title">YouTube → Telegram</div>
          <div class="card-subtitle">New video</div>
          <div class="card-description">Get a Telegram ping when a channel uploads</div>
        </div>
      </div>
    </v-container>

    <div v-if="showCreateModal" class="custom-modal-overlay" @click="showCreateModal = false">
      <div class="custom-modal-content" @click.stop>
        <CreateArea @close="showCreateModal = false" @save="showCreateModal = false" />
      </div>
    </div>

    <v-container class="mt-6">
      <div class="section-header">
        <div class="section-info">
          <h2 class="section-title">Recommended for you</h2>
          <p class="section-subtitle">Based on your usage patterns</p>
        </div>
        <button class="view-all-btn">
          <span>View All</span>
          <v-icon size="16">mdi-arrow-right</v-icon>
        </button>
      </div>
      <div class="cards-grid">
        <div class="card-col">
          <v-sheet class="area-card gradient-blue" rounded="xl">
            <v-icon size="64" color="white">mdi-weather-partly-cloudy</v-icon>
          </v-sheet>
          <div class="card-title">Weather → Telegram</div>
          <div class="card-subtitle">Daily reminder</div>
          <div class="card-description">Get weather forecast every morning on Telegram</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-pink" rounded="xl">
            <v-icon size="64" color="white">mdi-camera</v-icon>
          </v-sheet>
          <div class="card-title">Instagram → Dropbox</div>
          <div class="card-subtitle">Auto backup</div>
          <div class="card-description">Automatically backup your Instagram stories</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-teal" rounded="xl">
            <v-icon size="64" color="white">mdi-twitch</v-icon>
          </v-sheet>
          <div class="card-title">Twitch → Discord</div>
          <div class="card-subtitle">Go live</div>
          <div class="card-description">Alert your server when your stream starts</div>
        </div>
        <div class="card-col">
          <v-sheet class="area-card gradient-orange" rounded="xl">
            <v-icon size="64" color="white">mdi-newspaper-variant-outline</v-icon>
          </v-sheet>
          <div class="card-title">RSS → Notion</div>
          <div class="card-subtitle">Save articles</div>
          <div class="card-description">Append new posts to your Notion reading list</div>
        </div>
      </div>
    </v-container>

    <v-container>
      <div class="section-header">
        <div class="section-info">
          <h2 class="section-title">Create new AREA</h2>
          <p class="section-subtitle">Start building your automation</p>
        </div>
      </div>
      <div class="create-section">
        <div class="floating-icons">
          <div class="floating-card card-1">
            <v-icon size="24" color="white">mdi-email-outline</v-icon>
          </div>
          <div class="floating-card card-2">
            <v-icon size="24" color="white">mdi-music-note</v-icon>
          </div>
          <div class="floating-card card-3">
            <v-icon size="24" color="white">mdi-github</v-icon>
          </div>
          <div class="floating-card card-4">
            <v-icon size="24" color="white">mdi-chat</v-icon>
          </div>
          <div class="floating-card card-5">
            <v-icon size="24" color="white">mdi-calendar</v-icon>
          </div>
        </div>
        <div class="cards-grid">
          <CardButton @open="() => requireAuth(() => showCreateModal = true)" />
        </div>
      </div>
    </v-container>
    </div>
  </div>

  <!-- Dialog de confirmation de déconnexion -->
  <div v-if="showLogoutDialog" class="custom-modal-overlay" @click="showLogoutDialog = false">
    <div class="custom-modal-content logout-modal" @click.stop>
      <div class="logout-modal-header">
        <div class="logout-icon-container">
          <v-icon size="32" color="white">mdi-logout</v-icon>
        </div>
        <h3 class="logout-title">Sign Out</h3>
        <p class="logout-message">You will be redirected to the sign in page</p>
      </div>

      <div class="logout-modal-actions">
        <v-btn
          class="logout-cancel-btn"
          variant="outlined"
          @click="showLogoutDialog = false"
        >
          Cancel
        </v-btn>
        <v-btn
          class="logout-confirm-btn"
          variant="flat"
          @click="confirmLogout"
        >
          Sign Out
        </v-btn>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
import CreateArea from '../components/CreateArea/CreateArea.vue'
import SidebarButton from '../components/CreateArea/SidebarButton.vue'
import CardButton from '../components/CreateArea/CardButton.vue'
import { ref, watch, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useRouter } from 'vue-router'

const year = new Date().getFullYear()
const showCreateModal = ref(false)
const showLogoutDialog = ref(false)

const { isAuthenticated, currentUser, logout, refreshProfile, getProfileImageUrl } = useAuth()
const router = useRouter()

onMounted(async () => {
  await refreshProfile()
})


watch(isAuthenticated, (newValue) => {
  console.log('Authentication state changed:', newValue)
  console.log('Current user:', currentUser.value)
})

const goToLogin = () => {
  router.push('/login')
}

const requireAuth = (action: () => void) => {
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }
  action()
}

const confirmLogout = async () => {
  try {
    await logout()
    showLogoutDialog.value = false
    router.push('/login')
  } catch (error) {
    console.error('Error during sign out:', error)
  }
}

watch(showCreateModal, (isOpen) => {
  if (isOpen) {
    document.body.classList.add('modal-open')
  } else {
    document.body.classList.remove('modal-open')
  }
})
</script>

<style scoped>
.landing-dark {
  background: var(--gradient-bg-primary);
  color: var(--color-text-primary);
  min-height: 100vh;
  overflow-x: hidden;
  width: 100%;
  max-width: 100vw;
  box-sizing: border-box;
}

* {
  box-sizing: border-box;
}
.content {
  padding-left: 0;
  max-width: 100%;
  overflow-x: hidden;
  box-sizing: border-box;
}

.search-section {
  padding: 2rem 2rem 3rem 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 30vh;
  position: relative;
}

.search-container {
  width: 100%;
  max-width: 800px;
  text-align: center;
}

.search-header {
  margin-bottom: 2rem;
}

.search-title {
  font-size: 2.5rem;
  font-weight: 800;
  color: white;
  margin: 0 0 0.75rem 0;
  line-height: 1.1;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff 0%, #e5e7eb 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.search-subtitle {
  font-size: 1rem;
  color: #9ca3af;
  margin: 0;
  line-height: 1.6;
  font-weight: 400;
}

.search-bar {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.search-input-container {
  position: relative;
  display: flex;
  align-items: center;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-xl);
  padding: var(--spacing-md) var(--spacing-lg);
  backdrop-filter: blur(20px);
  transition: var(--transition-normal);
}

.search-input-container:focus-within {
  border-color: var(--color-border-focus);
  box-shadow: 0 0 0 3px var(--color-focus-ring);
}

.search-icon {
  margin-right: 1rem;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: var(--color-text-primary);
  font-size: 1rem;
  font-weight: 400;
}

.search-input::placeholder {
  color: var(--color-text-secondary);
}

.search-filter-btn {
  background: transparent;
  border: none;
  padding: var(--spacing-sm);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: var(--transition-normal);
  margin-left: var(--spacing-md);
}

.search-filter-btn:hover {
  background: var(--color-hover-bg);
}

.search-suggestions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  justify-content: center;
  flex-wrap: wrap;
}

.suggestion-label {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.suggestion-chip {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-2xl);
  padding: var(--spacing-sm) var(--spacing-md);
  color: var(--color-text-primary);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(20px);
}

.suggestion-chip:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-focus);
  transform: translateY(-1px);
}

.create-section {
  position: relative;
  min-height: 300px;
}

.floating-icons {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  z-index: 1;
}

.floating-card {
  position: absolute;
  width: 60px;
  height: 60px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.card-1 {
  background: var(--gradient-accent);
  top: 10%;
  left: 5%;
  animation-delay: 0s;
}

.card-2 {
  background: linear-gradient(135deg, var(--color-accent-secondary), var(--color-accent-tertiary));
  top: 20%;
  right: 10%;
  animation-delay: 1s;
}

.card-3 {
  background: linear-gradient(135deg, var(--color-accent-tertiary), var(--color-accent-primary));
  top: 60%;
  left: 15%;
  animation-delay: 2s;
}

.card-4 {
  background: var(--gradient-accent);
  top: 70%;
  right: 20%;
  animation-delay: 3s;
}

.card-5 {
  background: linear-gradient(135deg, var(--color-accent-secondary), var(--color-accent-tertiary));
  top: 40%;
  left: 50%;
  animation-delay: 4s;
}

.cards-grid {
  position: relative;
  z-index: 2;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
}

.user-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.user-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.profile-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.sidebar-user-section {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 0.75rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  margin-bottom: 0.5rem;
}

.sidebar-avatar {
  background: var(--gradient-accent);
  box-shadow: var(--shadow-glow);
}

.sidebar-profile-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.sidebar-user-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.sidebar-user-name {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-status {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  font-weight: 400;
}

/* Section Guest */
.guest-section {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.guest-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.guest-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-full);
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-glow);
}

.guest-text {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.guest-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.guest-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  font-weight: 400;
  margin: 0;
}

.guest-actions {
  display: flex;
  gap: 0.75rem;
}

.guest-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-lg);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-normal);
  border: 2px solid transparent;
}

.guest-btn.primary {
  background: var(--gradient-accent);
  color: var(--color-text-primary);
  border-color: transparent;
}

.guest-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow:
    var(--shadow-glow),
    0 8px 16px -5px rgba(6, 182, 212, 0.4);
}

.guest-btn.secondary {
  background: transparent;
  color: var(--color-text-primary);
  border-color: var(--color-border-primary);
}

.guest-btn.secondary:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: translateY(-1px);
}

.filter-tabs {
  display: flex;
  gap: var(--spacing-sm);
  background: var(--color-bg-card);
  border-radius: var(--radius-lg);
  padding: var(--spacing-xs);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
}

.filter-tab {
  padding: var(--spacing-sm) var(--spacing-lg);
  border: none;
  background: transparent;
  color: var(--color-text-secondary);
  font-weight: 500;
  font-size: 0.875rem;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: var(--transition-normal);
}

.filter-tab.active {
  background: var(--gradient-accent);
  color: var(--color-text-primary);
  box-shadow: var(--shadow-glow);
}

.filter-tab:hover:not(.active) {
  background: var(--color-hover-bg);
  color: var(--color-text-primary);
}

.action-buttons {
  display: flex;
  gap: var(--spacing-sm);
}

.action-btn-icon {
  width: 40px;
  height: 40px;
  border: none;
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
}

.action-btn-icon:hover {
  background: var(--color-hover-bg);
  transform: translateY(-1px);
}

.v-container {
  max-width: 100% !important;
  overflow-x: hidden;
}
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2rem;
}

.section-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.section-title {
  font-weight: 800;
  font-size: 2rem;
  color: var(--color-text-primary);
  margin: 0;
  letter-spacing: -0.02em;
}

.section-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 400;
}

.view-all-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-lg);
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-weight: 500;
  font-size: 0.875rem;
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(20px);
}

.view-all-btn:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: translateY(-1px);
}
.chips-row :deep(.v-chip) {
  font-weight: 600;
}
.gradient-avatar {
  background: var(--gradient-accent);
}
.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 24px;
  width: 100%;
  box-sizing: border-box;
}
.card-col { max-width: 320px; }
.area-card {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 24px;
  box-shadow: 0 6px 16px rgba(0,0,0,0.25);
  transform: translateY(0) scale(1);
  background-size: 130% 130%;
  background-position: 50% 50%;
  transition:
    transform .25s ease,
    box-shadow .25s ease,
    background-position .6s ease,
    filter .25s ease;
}
.area-card :deep(.v-icon) {
  transition: transform .25s ease, opacity .25s ease;
}
.area-card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 12px 28px rgba(0,0,0,0.35);
  background-position: 80% 20%;
}
.area-card:hover :deep(.v-icon) {
  transform: translateY(-2px) scale(1.06);
}
.area-card:active {
  transform: translateY(-2px) scale(0.99);
}

.cards-grid .card-col { animation: fadeUp .45s ease both; }
.cards-grid .card-col:nth-child(2) { animation-delay: .05s; }
.cards-grid .card-col:nth-child(3) { animation-delay: .1s; }
.cards-grid .card-col:nth-child(4) { animation-delay: .15s; }

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (prefers-reduced-motion: reduce) {
  .area-card,
  .area-card :deep(.v-icon),
  .cards-grid .card-col {
    transition: none !important;
    animation: none !important;
  }
}
.gradient-red { background: var(--color-area-red); }
.gradient-green { background: var(--color-area-green); }
.gradient-blue { background: var(--color-area-blue); }
.gradient-pink { background: var(--color-area-pink); }
.gradient-indigo { background: var(--color-area-blue); }
.gradient-crimson { background: var(--color-area-red); }
.gradient-teal { background: var(--color-area-green); }
.gradient-orange { background: var(--color-area-orange); }
.card-title { margin-top: 12px; font-weight: 800; font-size: 20px; }
.card-subtitle { color: rgba(255,255,255,0.85); font-weight: 700; }
.card-description { color: rgba(255,255,255,0.7); }
.bottom-nav {
  position: fixed;
  left: 50%;
  transform: translateX(-50%);
  bottom: 16px;
  width: min(920px, 92%);
  background: rgba(255,255,255,0.06);
  border: 1px solid rgba(255,255,255,0.1);
  backdrop-filter: blur(12px);
  border-radius: 28px;
  padding: 8px 6px;
}
.nav-inner { display: grid; grid-template-columns: repeat(5, 1fr); gap: 6px; }
.nav-btn { color: white !important; text-transform: none; }
.sidebar-desktop {
  display: none;
}
@media (min-width: 1280px) {
  .sidebar-desktop {
    display: block;
  }
  .content {
    padding-left: 80px;
  }
  .bottom-nav {
    display: none;
  }
}

@media (max-width: 1024px) {
  .search-section {
    padding: 2rem 1rem 2rem 1rem;
    min-height: 25vh;
  }

  .search-title {
    font-size: 2rem;
  }

  .search-suggestions {
    justify-content: flex-start;
  }

  .create-section {
    min-height: 250px;
  }

  .floating-card {
    width: 50px;
    height: 50px;
  }
}

@media (max-width: 768px) {
  .search-section {
    padding: 1.5rem 1rem 1.5rem 1rem;
    min-height: 20vh;
  }

  .search-title {
    font-size: 1.75rem;
  }

  .search-subtitle {
    font-size: 0.875rem;
  }

  .search-input-container {
    padding: 0.875rem 1rem;
  }

  .guest-section {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }

  .guest-content {
    flex-direction: column;
    gap: 0.75rem;
  }

  .guest-actions {
    flex-direction: column;
    width: 100%;
  }

  .guest-btn {
    width: 100%;
    justify-content: center;
  }

  .search-suggestions {
    gap: 0.5rem;
  }

  .suggestion-chip {
    padding: 0.375rem 0.75rem;
    font-size: 0.8125rem;
  }

  .filter-tabs {
    flex-wrap: wrap;
    gap: 0.25rem;
  }

  .filter-tab {
    padding: 0.5rem 1rem;
    font-size: 0.8125rem;
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .view-all-btn {
    align-self: flex-end;
  }
}

html, body {
  overflow-x: hidden !important;
  max-width: 100vw !important;
}


body.modal-open {
  overflow: hidden !important;
}

.custom-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--gradient-bg-modal);
  backdrop-filter: blur(12px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl);
}

.custom-modal-content {
  width: 100%;
  max-width: 960px;
  max-height: 90vh;
  overflow-y: auto;
  border-radius: 24px;
  background: transparent;
}

.custom-modal-content::-webkit-scrollbar {
  width: 8px;
}

.custom-modal-content::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.custom-modal-content::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, #3b82f6, #7c3aed, #ec4899);
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.custom-modal-content::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, #2563eb, #6d28d9, #db2777);
  box-shadow: 0 0 8px rgba(59, 130, 246, 0.4);
}

.custom-modal-content::-webkit-scrollbar-thumb:active {
  background: var(--color-scrollbar-thumb-active);
}

.custom-modal-content {
  scrollbar-width: thin;
  scrollbar-color: #3b82f6 rgba(255, 255, 255, 0.05);
}

/* Styles pour le modal de déconnexion - Cohérent avec CreateArea */
.logout-modal {
  max-width: 400px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-xl);
  padding: 0;
  overflow: hidden;
}

.logout-modal-header {
  padding: 2rem;
  text-align: center;
  background: var(--gradient-accent);
}

.logout-icon-container {
  width: 60px;
  height: 60px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem auto;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.logout-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logout-message {
  font-size: 1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-weight: 400;
  line-height: 1.5;
}

.logout-modal-actions {
  padding: 1.5rem 2rem 2rem 2rem;
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.logout-cancel-btn {
  background: var(--color-bg-card) !important;
  color: var(--color-text-primary) !important;
  border: 1px solid var(--color-border-primary) !important;
  border-radius: var(--radius-lg);
  font-weight: 500;
  text-transform: none;
  transition: var(--transition-normal);
}

.logout-cancel-btn:hover {
  background: var(--color-hover-bg) !important;
  border-color: var(--color-border-secondary) !important;
  transform: translateY(-1px);
}

.logout-confirm-btn {
  background: var(--gradient-accent) !important;
  color: var(--color-text-primary) !important;
  border: none !important;
  border-radius: var(--radius-lg);
  font-weight: 600;
  text-transform: none;
  transition: var(--transition-normal);
  box-shadow: var(--shadow-glow);
}

.logout-confirm-btn:hover {
  transform: translateY(-1px);
  box-shadow:
    var(--shadow-glow),
    0 8px 20px rgba(59, 130, 246, 0.3);
}

/* Responsive */
@media (max-width: 480px) {
  .logout-modal {
    margin: 1rem;
    max-width: calc(100vw - 2rem);
  }

  .logout-modal-header {
    padding: 1.5rem;
  }

  .logout-icon-container {
    width: 50px;
    height: 50px;
    margin-bottom: 0.75rem;
  }

  .logout-title {
    font-size: 1.25rem;
  }

  .logout-message {
    font-size: 0.875rem;
  }

  .logout-modal-actions {
    padding: 1rem 1.5rem 1.5rem 1.5rem;
    flex-direction: column;
  }

  .logout-cancel-btn,
  .logout-confirm-btn {
    width: 100%;
  }
}
</style>


