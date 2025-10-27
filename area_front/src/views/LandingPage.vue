<template>
  <div class="landing-dark">
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

    <div class="globe-background">
      <Globe class="globe-container" />
    </div>

    <v-navigation-drawer
      v-if="isDesktop"
      class="sidebar-desktop text-white"
      color="#0d0d0d"
      elevation="0"
      permanent
      rail
    >
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
        <SidebarButton
          v-if="isAuthenticated"
          tooltip="Create"
          @open="() => requireAuth(() => showCreateModal = true)"
        />
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
    <div class="search-section" id="search-section">
      <div class="search-container">
        <div class="search-header">
          <h1 class="search-title">Find Your Perfect Automation</h1>
          <p class="search-subtitle">Discover templates, browse services, or create something new</p>
        </div>
        <div class="search-bar">
          <div class="search-input-container">
            <v-icon size="20" color="#9ca3af" class="search-icon">mdi-magnify</v-icon>
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search automations, services, or templates..."
              class="search-input"
            >
            <button v-if="searchQuery" @click="searchQuery = ''" class="search-filter-btn">
              <v-icon size="16" color="#9ca3af">mdi-close</v-icon>
            </button>
          </div>
          <div class="search-suggestions">
            <span class="suggestion-label">Popular:</span>
            <button class="suggestion-chip" @click="searchQuery = 'Gmail'">Gmail</button>
            <button class="suggestion-chip" @click="searchQuery = 'Discord'">Discord</button>
            <button class="suggestion-chip" @click="searchQuery = 'Spotify'">Spotify</button>
            <button class="suggestion-chip" @click="searchQuery = 'GitHub'">GitHub</button>
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

    <v-container v-if="isAuthenticated && areas.length > 0" id="my-areas-section" class="mt-6">
      <div class="section-header">
        <div class="section-info">
          <h2 class="section-title">My Areas</h2>
          <p class="section-subtitle">Your created automation areas</p>
        </div>
        <button class="view-all-btn" @click="goToAllAreas">
          <span>View All</span>
          <v-icon size="16">mdi-arrow-right</v-icon>
        </button>
      </div>

      <div class="cards-grid">
        <CardSpotlight
          v-for="area in filteredAreas"
          :key="area.id"
          :area="area"
          :show-delete-button="true"
          @click="handleAreaClick"
          @delete="handleDeleteArea"
        />
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
          <CardButton
            @open="() => requireAuth(() => showCreateModal = true)"
          />
        </div>
      </div>
    </v-container>


    <div v-if="showCreateModal" class="custom-modal-overlay" @click="showCreateModal = false">
      <div class="custom-modal-content" @click.stop>
        <CreateArea
          :template="selectedArea"
          @close="showCreateModal = false"
          @save="handleAreaCreated"
        />
      </div>
    </div>

    <div class="bottom-nav">
      <div class="nav-inner">
        <v-btn class="nav-btn" variant="text" @click="() => requireAuth(() => showCreateModal = true)">
          <v-icon size="22">mdi-plus-circle</v-icon>
        </v-btn>
        <v-btn class="nav-btn" variant="text" @click="openProfileOrLogin">
          <v-icon size="22">mdi-account-circle</v-icon>
        </v-btn>
      </div>
    </div>

    </div>

    <footer class="site-footer">
      <div class="footer-content">
        <div class="footer-section">
          <div class="footer-logo">
            <h3 class="company-name">AREA</h3>
            <p class="company-tagline">Intelligent Automation Platform</p>
          </div>
          <p class="footer-description">
            Connect your favorite services with intelligent automation.
            Build powerful workflows that work for you.
          </p>
        </div>

        <div class="footer-section">
          <h4 class="footer-title">Product</h4>
          <ul class="footer-links">
            <li><a href="#" @click.prevent>Features</a></li>
            <li><a href="#" @click.prevent>Integrations</a></li>
            <li><a href="#" @click.prevent>API</a></li>
            <li><a href="#" @click.prevent>Documentation</a></li>
          </ul>
        </div>

        <div class="footer-section">
          <h4 class="footer-title">Company</h4>
          <ul class="footer-links">
            <li><a href="#" @click.prevent>About Us</a></li>
            <li><a href="#" @click.prevent>Careers</a></li>
            <li><a href="#" @click.prevent>Blog</a></li>
            <li><a href="#" @click.prevent>Press</a></li>
          </ul>
        </div>

        <div class="footer-section">
          <h4 class="footer-title">Support</h4>
          <ul class="footer-links">
            <li><a href="#" @click.prevent>Help Center</a></li>
            <li><a href="#" @click.prevent>Community</a></li>
            <li><a href="#" @click.prevent>Contact</a></li>
            <li><a href="#" @click.prevent>Status</a></li>
          </ul>
        </div>

        <div class="footer-section">
          <h4 class="footer-title">Contact</h4>
          <div class="contact-info">
            <div class="contact-item">
              <v-icon size="16" color="var(--color-accent-primary)">mdi-email</v-icon>
              <span>contact@area.com</span>
            </div>
            <div class="contact-item">
              <v-icon size="16" color="var(--color-accent-primary)">mdi-phone</v-icon>
              <span>+33 7 41 61 72 18</span>
            </div>
            <div class="contact-item">
              <v-icon size="16" color="var(--color-accent-primary)">mdi-map-marker</v-icon>
              <span>Paname, France</span>
            </div>
          </div>
          <div class="social-links">
            <a href="#" class="social-link" @click.prevent>
              <v-icon size="20">mdi-twitter</v-icon>
            </a>
            <a href="#" class="social-link" @click.prevent>
              <v-icon size="20">mdi-github</v-icon>
            </a>
            <a href="#" class="social-link" @click.prevent>
              <v-icon size="20">mdi-linkedin</v-icon>
            </a>
            <a href="#" class="social-link" @click.prevent>
              <v-icon size="20">mdi-discord</v-icon>
            </a>
          </div>
        </div>
      </div>

      <div class="footer-bottom">
        <div class="footer-bottom-content">
          <p class="copyright">
            © {{ year }} AREA. All rights reserved.
          </p>
          <div class="footer-bottom-links">
            <a href="#" @click.prevent>Privacy Policy</a>
            <a href="#" @click.prevent>Terms of Service</a>
            <a href="#" @click.prevent>Cookie Policy</a>
          </div>
        </div>
      </div>
    </footer>
  </div>

  <div v-if="showLogoutDialog" class="logout-modal-overlay" @click="showLogoutDialog = false">
    <div class="logout-modal-container" @click.stop>
      <div class="logout-modal-content">
        <div class="logout-modal-header">
          <div class="logout-icon-wrapper">
            <div class="logout-icon-bg">
              <v-icon size="24" color="#ff3b30">mdi-logout-variant</v-icon>
            </div>
          </div>
          <h2 class="logout-modal-title">Sign Out</h2>
          <p class="logout-modal-subtitle">Are you sure you want to sign out of your account?</p>
        </div>

        <div class="logout-modal-actions">
          <button
            class="logout-action-btn logout-cancel-btn"
            @click="showLogoutDialog = false"
          >
            Cancel
          </button>
          <button
            class="logout-action-btn logout-confirm-btn"
            @click="confirmLogout"
          >
            Sign Out
          </button>
        </div>
      </div>
    </div>
  </div>

  <div v-if="showAreaModal && selectedArea" class="custom-modal-overlay" @click="showAreaModal = false">
    <div class="custom-modal-content area-modal" @click.stop>
      <div class="area-modal-header">
        <div class="area-icon-container">
          <v-icon :size="48" color="white">{{ getTriggerIcon(selectedArea?.triggerService) }}</v-icon>
        </div>
        <h3 class="area-modal-title">{{ selectedArea?.title }}</h3>
        <p class="area-modal-subtitle">{{ selectedArea?.subtitle }}</p>
      </div>

      <div class="area-modal-content">
        <div class="area-description">
          <h4 class="description-title">Description</h4>
          <p class="description-text">{{ selectedArea?.description }}</p>
        </div>

        <div class="area-workflow">
          <h4 class="workflow-title">How it works</h4>
          <div class="workflow-steps">
            <div class="workflow-step">
              <div class="step-icon trigger-icon">
                <v-icon size="20" color="white">{{ getTriggerIcon(selectedArea?.triggerService) }}</v-icon>
              </div>
              <div class="step-content">
                <div class="step-label">Trigger</div>
                <div class="step-service">{{ selectedArea?.triggerService }}</div>
              </div>
            </div>
            <div class="workflow-arrow">
              <v-icon size="24" color="#9ca3af">mdi-arrow-right</v-icon>
            </div>
            <div class="workflow-step">
              <div class="step-icon action-icon">
                <v-icon size="20" color="white">{{ getActionIcon(selectedArea?.actionService) }}</v-icon>
              </div>
              <div class="step-content">
                <div class="step-label">Action</div>
                <div class="step-service">{{ selectedArea?.actionService }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="area-modal-actions">
        <v-btn
          class="area-modal-close-btn"
          variant="outlined"
          @click="showAreaModal = false"
        >
          Close
        </v-btn>
        <v-btn
          class="area-modal-create-btn"
          variant="flat"
          @click="createAreaFromTemplate"
        >
          Use This Area
        </v-btn>
      </div>
    </div>
  </div>



</template>

<script setup lang="ts">
import CreateArea from '../components/CreateArea/CreateArea.vue'
import SidebarButton from '../components/CreateArea/SidebarButton.vue'
import CardButton from '../components/CreateArea/CardButton.vue'
import CardSpotlight from '../components/CardSpotlight.vue'
import Globe from '../components/Globe.vue'
import { ref, watch, onMounted, computed } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useAreas } from '@/composables/useAreas'
import { useRouter } from 'vue-router'
import { githubService, type GitHubRepository } from '@/services/github'
import { type Area, areaService } from '@/services/area'

interface AreaTemplate {
  id: string
  title: string
  subtitle: string
  description: string
  icon: string
  gradientClass: string
  triggerService: string
  actionService: string
  isActive: boolean
}

const year = new Date().getFullYear()
const showCreateModal = ref(false)
const showLogoutDialog = ref(false)
const showAreaModal = ref(false)
const selectedArea = ref<AreaTemplate | null>(null)
const searchQuery = ref('')
const isDesktop = ref(typeof window !== 'undefined' ? window.innerWidth >= 1280 : true)

const { isAuthenticated, currentUser, logout, refreshProfile, getProfileImageUrl } = useAuth()
const { areas, popularAreas, recommendedAreas, fetchPopularAreas, fetchRecommendedAreas, fetchUserAreas, deleteArea } = useAreas()
const router = useRouter()

const filteredAreas = computed(() => {
  if (!searchQuery.value.trim()) {
    return areas.value
  }

  const query = searchQuery.value.toLowerCase().trim()

  return areas.value.filter((area: any) => {
    const name = area.name?.toLowerCase() || ''
    const description = area.description?.toLowerCase() || ''
    const triggerService = area.trigger_service?.toLowerCase() || ''
    const actionService = area.action_service?.toLowerCase() || ''

    return (
      name.includes(query) ||
      description.includes(query) ||
      triggerService.includes(query) ||
      actionService.includes(query)
    )
  })
})

import appsJson from '../assets/apps.json'

type AppDef = { name: string; icon: string }
const apps = (Array.isArray(appsJson) ? appsJson : (appsJson as any).apps ?? []) as AppDef[]

const getIconUrl = (file: string) =>
  new URL(`../assets/app-icons/${file}`, import.meta.url).href

const getServiceIcon = (serviceName: string | undefined) => {
  if (!serviceName || serviceName === 'undefined' || serviceName === 'null') {
    const githubApp = apps.find(a => a.name === 'GitHub')
    if (githubApp) {
      return getIconUrl(githubApp.icon)
    }
    return ''
  }

  const app = apps.find(a => a.name === serviceName)

  if (app) {
    return getIconUrl(app.icon)
  } else {
    const caseInsensitiveApp = apps.find(a => a.name.toLowerCase() === serviceName.toLowerCase())
    if (caseInsensitiveApp) {
      return getIconUrl(caseInsensitiveApp.icon)
    }

    const githubApp = apps.find(a => a.name === 'GitHub')
    if (githubApp) {
      return getIconUrl(githubApp.icon)
    }
  }
  return ''
}



onMounted(() => {
  const onResize = () => {
    isDesktop.value = window.innerWidth >= 1280
  }
  window.addEventListener('resize', onResize)
  refreshProfile()
    .then(() => fetchPopularAreas())
    .then(() => fetchRecommendedAreas())
    .then(() => {
      if (currentUser.value && currentUser.value.id) {
        return fetchUserAreas(currentUser.value.id)
      }
    })
  window.addEventListener('beforeunload', () => {
    window.removeEventListener('resize', onResize)
  })
})

onMounted(async () => {
  await refreshProfile()
  await fetchPopularAreas()
  await fetchRecommendedAreas()
  if (currentUser.value?.id) {
    await fetchUserAreas(currentUser.value.id)
  }
})


watch(isAuthenticated, (newValue) => {
  if (newValue && currentUser.value?.id) {
    fetchUserAreas(currentUser.value.id)
  }
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



const openProfileOrLogin = () => {
  if (isAuthenticated.value) {
    router.push('/profile')
  } else {
    router.push('/login')
  }
}


const confirmLogout = async () => {
  try {
    await logout()
    showLogoutDialog.value = false
    router.push('/login')
  } catch (error) {
  }
}

const handleAreaClick = (area: AreaTemplate | Area) => {
  requireAuth(() => {
    router.push({
      name: 'configure-area',
      query: {
        areaId: area.id
      }
    })
  })
}

const handleDeleteArea = async (area: AreaTemplate | Area) => {
  requireAuth(async () => {
    const areaName = 'name' in area ? area.name : area.title
    if (confirm(`Are you sure you want to delete "${areaName}"? This action cannot be undone.`)) {
      try {
        await deleteArea(area.id)
      } catch (error) {
        alert(`Failed to delete area: ${error instanceof Error ? error.message : 'Unknown error'}`)
      }
    }
  })
}

const createAreaFromTemplate = () => {
  showAreaModal.value = false

  router.push({
    name: 'configure-area',
    query: {
      template: JSON.stringify(selectedArea.value)
    }
  })
}

const handleAreaCreated = async () => {
  showCreateModal.value = false
  await fetchPopularAreas()
  await fetchRecommendedAreas()
  if (currentUser.value?.id) {
    await fetchUserAreas(currentUser.value.id)
  }
}

const scrollToMyAreas = () => {
  requireAuth(() => {
    const element = document.getElementById('my-areas-section')
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' })
    } else {
      showCreateModal.value = true
    }
  })
}

const goToAllAreas = () => {
  requireAuth(() => {
    router.push('/areas')
  })
}

const getTriggerIcon = (service: string) => {
  switch (service) {
    case "Date Timer":
      return "mdi-calendar"
    case "GitHub":
      return "mdi-github"
    case "Gmail":
      return "mdi-email-outline"
    case "Discord":
      return "mdi-discord"
    case "Slack":
      return "mdi-slack"
    case "Weather":
      return "mdi-weather-partly-cloudy"
    case "Instagram":
      return "mdi-instagram"
    case "Twitter":
      return "mdi-twitter"
    case "YouTube":
      return "mdi-youtube"
    case "Spotify":
      return "mdi-music"
    case "Telegram":
      return "mdi-telegram"
    case "Twitch":
      return "mdi-twitch"
    case "Dropbox":
      return "mdi-dropbox"
    case "Notion":
      return "mdi-notebook"
    default:
      return "mdi-cog"
  }
}

const getActionIcon = (service: string) => {
  switch (service) {
    case "Gmail":
      return "mdi-email-outline"
    case "Discord":
      return "mdi-discord"
    case "Slack":
      return "mdi-slack"
    case "GitHub":
      return "mdi-github"
    case "Telegram":
      return "mdi-telegram"
    case "Twitter":
      return "mdi-twitter"
    case "Instagram":
      return "mdi-instagram"
    case "Dropbox":
      return "mdi-dropbox"
    case "Notion":
      return "mdi-notebook"
    default:
      return "mdi-cog"
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
  background: radial-gradient(circle at 30% 20%, rgba(87, 128, 232, 0.1) 0%, transparent 50%),
              radial-gradient(circle at 70% 80%, rgba(135, 81, 209, 0.1) 0%, transparent 50%);
}

.globe-background {
  position: absolute;
  top: 0;
  right: -200px;
  width: 600px;
  height: 100vh;
  z-index: 999;
  pointer-events: none;
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.globe-container {
  opacity: 0.4;
  filter: blur(0.5px);
}

@media (max-width: 1280px) {
  .globe-background {
    right: -150px;
    width: 480px;
  }
}

@media (max-width: 1024px) {
  .globe-background {
    right: -120px;
    width: 360px;
  }
}

@media (max-width: 768px) {
  .globe-background {
    display: none;
  }
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-30px) rotate(10deg); }
}



.stars3 {
  background-image:
    radial-gradient(1px 1px at 31px 59px, rgba(255,255,255,0.5), transparent),
    radial-gradient(1px 1px at 83px 27px, rgba(255,255,255,0.7), transparent),
    radial-gradient(1px 1px at 119px 87px, rgba(255,255,255,0.6), transparent),
    radial-gradient(1px 1px at 157px 13px, rgba(255,255,255,0.4), transparent),
    radial-gradient(1px 1px at 191px 71px, rgba(255,255,255,0.8), transparent);
  background-repeat: repeat;
  background-size: 400px 200px;
  animation: twinkle 15s ease-in-out infinite;
}

@keyframes twinkle {
  0%, 100% {
    opacity: 0.3;
    transform: translateY(0px) scale(1);
  }
  25% {
    opacity: 0.8;
    transform: translateY(-15px) scale(1.1);
  }
  50% {
    opacity: 0.4;
    transform: translateY(-8px) scale(0.9);
  }
  75% {
    opacity: 0.9;
    transform: translateY(-22px) scale(1.05);
  }
}

.planet-container {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  pointer-events: none;
  z-index: 1;
}

.planet {
  position: absolute;
  border-radius: 50%;
  animation: float 20s ease-in-out infinite;
  border: 2px solid rgba(255, 255, 255, 0.3);
}

.planet-surface {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  position: relative;
  overflow: hidden;
}

.planet-atmosphere {
  position: absolute;
  top: -10%;
  left: -10%;
  width: 120%;
  height: 120%;
  border-radius: 50%;
  opacity: 0.3;
  animation: atmosphere 8s ease-in-out infinite;
}

.planet-glow {
  position: absolute;
  top: -20%;
  left: -20%;
  width: 140%;
  height: 140%;
  border-radius: 50%;
  opacity: 0.2;
  animation: glow 6s ease-in-out infinite;
}

.planet-left {
  position: fixed;
  width: 400px;
  height: 400px;
  top: 15%;
  left: -200px;
  animation-delay: 0s;
  z-index: 1;
}

.planet-right {
  position: fixed;
  width: 500px;
  height: 500px;
  top: 10%;
  right: -250px;
  animation-delay: 2s;
  z-index: 1;
}

.planet-core {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background: radial-gradient(circle at 40% 40%,
    #0f0f23 0%,
    #1a1a2e 30%,
    #16213e 60%,
    #0a0a0f 100%);
  box-shadow:
    inset -30px -30px 80px rgba(0, 0, 0, 0.6),
    inset 20px 20px 60px rgba(255, 255, 255, 0.05);
}

.planet-surface {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background:
    radial-gradient(ellipse 120px 60px at 30% 40%, rgba(26, 26, 46, 0.8), transparent),
    radial-gradient(ellipse 80px 40px at 60% 20%, rgba(22, 33, 62, 0.6), transparent),
    radial-gradient(ellipse 100px 50px at 70% 70%, rgba(15, 15, 35, 0.7), transparent),
    radial-gradient(ellipse 60px 30px at 40% 80%, rgba(26, 26, 46, 0.5), transparent);
  animation: planet-drift 25s ease-in-out infinite;
}

.planet-glow-edge {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  background:
    linear-gradient(90deg,
      transparent 0%,
      transparent 40%,
      rgba(255, 100, 50, 0.8) 50%,
      rgba(255, 150, 80, 0.6) 60%,
      transparent 70%,
      transparent 100%);
  box-shadow:
    0 0 60px rgba(255, 100, 50, 0.4),
    0 0 120px rgba(255, 150, 80, 0.2);
  animation: glow-pulse 4s ease-in-out infinite;
}

.planet-left .planet-glow-edge {
  background:
    linear-gradient(90deg,
      transparent 0%,
      transparent 30%,
      rgba(255, 120, 60, 0.9) 40%,
      rgba(255, 180, 100, 0.7) 50%,
      transparent 60%,
      transparent 100%);
  box-shadow:
    0 0 80px rgba(255, 120, 60, 0.5),
    0 0 160px rgba(255, 180, 100, 0.3);
}

.planet-right .planet-glow-edge {
  background:
    linear-gradient(270deg,
      transparent 0%,
      transparent 30%,
      rgba(255, 80, 40, 0.9) 40%,
      rgba(255, 140, 70, 0.7) 50%,
      transparent 60%,
      transparent 100%);
  box-shadow:
    0 0 100px rgba(255, 80, 40, 0.6),
    0 0 200px rgba(255, 140, 70, 0.4);
}

.planet-atmosphere {
  position: absolute;
  top: -20%;
  left: -20%;
  width: 140%;
  height: 140%;
  border-radius: 50%;
  background:
    radial-gradient(circle at 30% 30%,
      rgba(26, 26, 46, 0.3) 0%,
      rgba(22, 33, 62, 0.2) 50%,
      transparent 100%);
  animation: atmosphere-drift 15s ease-in-out infinite;
}

.nebula {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background:
    radial-gradient(ellipse at 20% 20%,
      rgba(59, 130, 246, 0.1) 0%,
      transparent 50%),
    radial-gradient(ellipse at 80% 80%,
      rgba(139, 92, 246, 0.08) 0%,
      transparent 50%),
    radial-gradient(ellipse at 60% 30%,
      rgba(236, 72, 153, 0.05) 0%,
      transparent 50%);
  animation: nebula-drift 30s ease-in-out infinite;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  25% {
    transform: translateY(-20px) rotate(5deg);
  }
  50% {
    transform: translateY(-10px) rotate(-3deg);
  }
  75% {
    transform: translateY(-30px) rotate(8deg);
  }
}

@keyframes atmosphere {
  0%, 100% {
    transform: rotate(0deg) scale(1);
    opacity: 0.3;
  }
  50% {
    transform: rotate(180deg) scale(1.1);
    opacity: 0.5;
  }
}

@keyframes glow {
  0%, 100% {
    opacity: 0.2;
    transform: scale(1);
  }
  50% {
    opacity: 0.4;
    transform: scale(1.1);
  }
}

@keyframes nebula-drift {
  0%, 100% {
    transform: translateX(0px) translateY(0px);
  }
  25% {
    transform: translateX(20px) translateY(-10px);
  }
  50% {
    transform: translateX(-10px) translateY(20px);
  }
  75% {
    transform: translateX(-20px) translateY(-5px);
  }
}

@keyframes planet-drift {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  25% {
    transform: translateY(-15px) rotate(2deg);
  }
  50% {
    transform: translateY(-8px) rotate(-1deg);
  }
  75% {
    transform: translateY(-20px) rotate(3deg);
  }
}

@keyframes atmosphere-drift {
  0%, 100% {
    transform: scale(1) rotate(0deg);
    opacity: 0.2;
  }
  50% {
    transform: scale(1.08) rotate(180deg);
    opacity: 0.4;
  }
}

@keyframes glow-pulse {
  0%, 100% {
    opacity: 0.6;
    transform: scale(1);
  }
  50% {
    opacity: 0.9;
    transform: scale(1.05);
  }
}

@media (max-width: 1024px) {
  .planet-left {
    width: 300px;
    height: 300px;
    top: 20%;
    left: -150px;
  }

  .planet-right {
    width: 400px;
    height: 400px;
    top: 15%;
    right: -200px;
  }
}

@media (max-width: 768px) {
  .planet-left {
    width: 250px;
    height: 250px;
    top: 25%;
    left: -125px;
  }

  .planet-right {
    width: 300px;
    height: 300px;
    top: 20%;
    right: -150px;
  }

  .stars, .stars2, .stars3 {
    background-size: 150px 75px;
  }
}

@media (max-width: 480px) {
  .planet-left {
    width: 200px;
    height: 200px;
    top: 30%;
    left: -100px;
  }

  .planet-right {
    width: 250px;
    height: 250px;
    top: 25%;
    right: -125px;
  }

  .stars, .stars2, .stars3 {
    background-size: 100px 50px;
  }
}

.landing-dark {
  background: var(--gradient-bg-primary);
  color: var(--color-text-primary);
  min-height: 100vh;
  overflow-x: hidden;
  width: 100%;
  max-width: 100vw;
  box-sizing: border-box;
  position: relative;
}

* {
  box-sizing: border-box;
}
.content {
  padding-left: 0;
  max-width: 100%;
  overflow-x: hidden;
  box-sizing: border-box;
  position: relative;
  z-index: 2;
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
  background: linear-gradient(135deg, #ffffff 0%, #f0f9ff 50%, #e0f2fe 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  filter: drop-shadow(0 4px 12px rgba(255, 255, 255, 0.1));
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
  color: #ffffff;
  margin: 0;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff 0%, #e2e8f0 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  filter: drop-shadow(0 2px 8px rgba(255, 255, 255, 0.1));
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
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 32px;
  width: 100%;
  box-sizing: border-box;
  justify-items: center;
}
.cards-grid .card-spotlight-container {
  animation: fadeUp .45s ease both;
}
.cards-grid .card-spotlight-container:nth-child(2) { animation-delay: .05s; }
.cards-grid .card-spotlight-container:nth-child(3) { animation-delay: .1s; }
.cards-grid .card-spotlight-container:nth-child(4) { animation-delay: .15s; }

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (max-width: 768px) {
  .cards-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .cards-grid .card-spotlight-container {
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
.gradient-purple { background: linear-gradient(135deg, #8b5cf6, #a855f7); }
.gradient-gray { background: linear-gradient(135deg, #6b7280, #9ca3af); }
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
  z-index: 1001;
}
.nav-inner { display: grid; grid-template-columns: repeat(5, 1fr); gap: 6px; }
.nav-btn { color: white !important; text-transform: none; }
.sidebar-desktop {
  display: none;
  position: relative;
  z-index: 3;
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
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--spacing-xl);
  animation: modalOverlayFadeIn 0.3s ease-out;
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

.logout-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(20px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  animation: fadeIn 0.3s ease-out;
}

.logout-modal-container {
  max-width: 400px;
  width: 90%;
  margin: 0 auto;
  animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.logout-modal-content {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(40px);
  border-radius: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow:
    0 20px 40px rgba(0, 0, 0, 0.15),
    0 0 0 1px rgba(255, 255, 255, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
}

.logout-modal-header {
  padding: 32px 24px 24px 24px;
  text-align: center;
}

.logout-icon-wrapper {
  margin-bottom: 20px;
}

.logout-icon-bg {
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: linear-gradient(135deg, #ff3b30, #ff6b6b);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  box-shadow: 0 8px 20px rgba(255, 59, 48, 0.3);
}

.logout-modal-title {
  font-size: 24px;
  font-weight: 600;
  color: #1d1d1f;
  margin: 0 0 8px 0;
  letter-spacing: -0.02em;
}

.logout-modal-subtitle {
  font-size: 16px;
  color: #86868b;
  margin: 0;
  line-height: 1.4;
  font-weight: 400;
}

.logout-modal-actions {
  padding: 0 24px 24px 24px;
  display: flex;
  gap: 12px;
}

.logout-action-btn {
  flex: 1;
  padding: 14px 20px;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
  letter-spacing: -0.01em;
}

.logout-cancel-btn {
  background: rgba(142, 142, 147, 0.12);
  color: #1d1d1f;
}

.logout-cancel-btn:hover {
  background: rgba(142, 142, 147, 0.18);
  transform: translateY(-1px);
}

.logout-confirm-btn {
  background: linear-gradient(135deg, #ff3b30, #ff6b6b);
  color: white;
  box-shadow: 0 4px 12px rgba(255, 59, 48, 0.3);
}

.logout-confirm-btn:hover {
  background: linear-gradient(135deg, #ff2d55, #ff5252);
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(255, 59, 48, 0.4);
}

@media (prefers-color-scheme: dark) {
  .logout-modal-content {
    background: rgba(28, 28, 30, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.1);
  }

  .logout-modal-title {
    color: #f2f2f7;
  }

  .logout-modal-subtitle {
    color: #8e8e93;
  }

  .logout-cancel-btn {
    background: rgba(142, 142, 147, 0.2);
    color: #f2f2f7;
  }

  .logout-cancel-btn:hover {
    background: rgba(142, 142, 147, 0.3);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes modalOverlayFadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes modalContentSlideUp {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.9);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@media (max-width: 480px) {
  .logout-modal-container {
    width: 95%;
    margin: 0 16px;
  }

  .logout-modal-header {
    padding: 24px 20px 20px 20px;
  }

  .logout-modal-actions {
    padding: 0 20px 20px 20px;
    flex-direction: column;
  }

  .logout-action-btn {
    width: 100%;
  }

  .logout-icon-bg {
    width: 56px;
    height: 56px;
  }

  .logout-modal-title {
    font-size: 20px;
  }

  .logout-modal-subtitle {
    font-size: 14px;
  }
}


.github-gmail-modal {
  max-width: 600px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-xl);
  padding: 0;
  overflow: hidden;
}

.github-gmail-modal-header {
  padding: 2rem;
  text-align: center;
  background: var(--gradient-accent);
}

.github-gmail-icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.arrow-icon {
  opacity: 0.8;
}

.github-gmail-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.github-gmail-message {
  font-size: 1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-weight: 400;
  line-height: 1.5;
}

.github-gmail-modal-content {
  padding: 2rem;
}

.form-section {
  margin-bottom: 1.5rem;
}

.form-section:last-child {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 0.5rem;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: 0.5rem;
  padding: 0.75rem;
  background: rgba(244, 67, 54, 0.1);
  border: 1px solid rgba(244, 67, 54, 0.3);
  border-radius: var(--radius-md);
  color: #f44336;
  font-size: 0.875rem;
  font-weight: 500;
}

.repository-select,
.email-input {
  width: 100%;
}

.repository-select :deep(.v-field),
.email-input :deep(.v-field) {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
}

.repository-select :deep(.v-menu) {
  z-index: 1001 !important;
}

.repository-select :deep(.v-list) {
  z-index: 1001 !important;
}

.repository-select :deep(.v-overlay) {
  z-index: 1001 !important;
}

.repository-select :deep(.v-overlay__content) {
  z-index: 1001 !important;
}

.repository-select :deep(.v-field--focused),
.email-input :deep(.v-field--focused) {
  border-color: var(--color-border-focus);
  box-shadow: 0 0 0 3px var(--color-focus-ring);
}

.repository-select :deep(.v-field__input),
.email-input :deep(.v-field__input) {
  color: var(--color-text-primary);
}

.repository-select :deep(.v-field__input::placeholder),
.email-input :deep(.v-field__input::placeholder) {
  color: var(--color-text-secondary);
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.checkbox-group :deep(.v-checkbox) {
  margin: 0;
}

.checkbox-group :deep(.v-checkbox .v-label) {
  color: var(--color-text-primary);
  font-weight: 500;
}

.github-gmail-modal-actions {
  padding: 1.5rem 2rem 2rem 2rem;
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.github-gmail-cancel-btn {
  background: var(--color-bg-card) !important;
  color: var(--color-text-primary) !important;
  border: 1px solid var(--color-border-primary) !important;
  border-radius: var(--radius-lg);
  font-weight: 500;
  text-transform: none;
  transition: var(--transition-normal);
}

.github-gmail-cancel-btn:hover {
  background: var(--color-hover-bg) !important;
  border-color: var(--color-border-secondary) !important;
  transform: translateY(-1px);
}

.github-gmail-confirm-btn {
  background: var(--gradient-accent) !important;
  color: var(--color-text-primary) !important;
  border: none !important;
  border-radius: var(--radius-lg);
  font-weight: 600;
  text-transform: none;
  transition: var(--transition-normal);
  box-shadow: var(--shadow-glow);
}

.github-gmail-confirm-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow:
    var(--shadow-glow),
    0 8px 20px rgba(59, 130, 246, 0.3);
}

.github-gmail-confirm-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

@media (max-width: 480px) {
  .logout-modal,
  .github-gmail-modal {
    margin: 1rem;
    max-width: calc(100vw - 2rem);
  }

  .logout-modal-header,
  .github-gmail-modal-header {
    padding: 1.5rem;
  }

  .logout-icon-container {
    width: 50px;
    height: 50px;
    margin-bottom: 0.75rem;
  }

  .logout-title,
  .github-gmail-title {
    font-size: 1.25rem;
  }

  .logout-message,
  .github-gmail-message {
    font-size: 0.875rem;
  }

  .logout-modal-actions,
  .github-gmail-modal-actions {
    padding: 1rem 1.5rem 1.5rem 1.5rem;
    flex-direction: column;
  }

  .logout-cancel-btn,
  .logout-confirm-btn,
  .github-gmail-cancel-btn,
  .github-gmail-confirm-btn {
    width: 100%;
  }

  .github-gmail-modal-content {
    padding: 1.5rem;
  }

  .github-gmail-icon-container {
    gap: 0.75rem;
  }
}

.area-modal {
  max-width: 600px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-xl);
  padding: 0;
  overflow: hidden;
  animation: modalContentSlideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.area-modal-header {
  padding: 2rem;
  text-align: center;
  background: var(--gradient-accent);
}

.area-icon-container {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 1rem auto;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.area-modal-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.area-modal-subtitle {
  font-size: 1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-weight: 500;
}

.area-modal-content {
  padding: 2rem;
}

.area-description {
  margin-bottom: 2rem;
}

.description-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.75rem 0;
}

.description-text {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0;
  line-height: 1.6;
}

.area-workflow {
  margin-bottom: 1rem;
}

.workflow-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 1rem 0;
}

.workflow-steps {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
}

.workflow-step {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
}

.step-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
}

.trigger-icon {
  background: var(--gradient-accent);
}

.action-icon {
  background: linear-gradient(135deg, #ef4444, #dc2626);
}

.step-content {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.step-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
}

.step-service {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.workflow-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.area-modal-actions {
  padding: 1.5rem 2rem 2rem 2rem;
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.area-modal-close-btn {
  background: var(--color-bg-card) !important;
  color: var(--color-text-primary) !important;
  border: 1px solid var(--color-border-primary) !important;
  border-radius: var(--radius-lg);
  font-weight: 500;
  text-transform: none;
  transition: var(--transition-normal);
}

.area-modal-close-btn:hover {
  background: var(--color-hover-bg) !important;
  border-color: var(--color-border-secondary) !important;
  transform: translateY(-1px);
}

.area-modal-create-btn {
  background: var(--gradient-accent) !important;
  color: var(--color-text-primary) !important;
  border: none !important;
  border-radius: var(--radius-lg);
  font-weight: 600;
  text-transform: none;
  transition: var(--transition-normal);
  box-shadow: var(--shadow-glow);
}

.area-modal-create-btn:hover {
  transform: translateY(-1px);
  box-shadow:
    var(--shadow-glow),
    0 8px 20px rgba(59, 130, 246, 0.3);
}

.area-modal-create-btn:disabled {
  background: var(--color-bg-card) !important;
  color: var(--color-text-secondary) !important;
  border: 1px solid var(--color-border-primary) !important;
  cursor: not-allowed;
  transform: none !important;
  box-shadow: none !important;
}

@media (max-width: 480px) {
  .area-modal {
    margin: 1rem;
    max-width: calc(100vw - 2rem);
  }

  .area-modal-header {
    padding: 1.5rem;
  }

  .area-icon-container {
    width: 60px;
    height: 60px;
    margin-bottom: 0.75rem;
  }

  .area-modal-title {
    font-size: 1.5rem;
  }

  .area-modal-content {
    padding: 1.5rem;
  }

  .workflow-steps {
    flex-direction: column;
    gap: 1rem;
  }

  .workflow-arrow {
    transform: rotate(90deg);
  }

  .area-modal-actions {
    padding: 1rem 1.5rem 1.5rem 1.5rem;
    flex-direction: column;
  }

  .area-modal-close-btn,
  .area-modal-create-btn {
    width: 100%;
  }
}

.site-footer {
  background: var(--color-bg-secondary);
  border-top: 1px solid var(--color-border-primary);
  margin-top: 4rem;
  padding: 3rem 0 0 0;
  backdrop-filter: blur(20px);
}

.footer-content {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr 1.5fr;
  gap: 3rem;
  margin-bottom: 3rem;
  padding: 0 2rem;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
}

.footer-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.footer-logo {
  margin-bottom: 1rem;
}

.company-name {
  font-size: 1.5rem;
  font-weight: 700;
  margin: 0 0 0.5rem 0;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 50%, #ec4899 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
  filter: drop-shadow(0 2px 8px rgba(59, 130, 246, 0.3));
}

.company-tagline {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
  font-weight: 500;
}

.footer-description {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 0;
}

.footer-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 1rem 0;
}

.footer-links {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.footer-links li a {
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s ease;
}

.footer-links li a:hover {
  color: var(--color-accent-primary);
}

.contact-info {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: var(--color-text-secondary);
  font-size: 0.875rem;
}

.social-links {
  display: flex;
  gap: 1rem;
}

.social-link {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  text-decoration: none;
  transition: var(--transition-normal);
}

.social-link:hover {
  background: var(--color-bg-card-hover);
  border-color: var(--color-border-focus);
  color: var(--color-accent-primary);
  transform: translateY(-2px);
  box-shadow: var(--shadow-glow);
}

.footer-bottom {
  border-top: 1px solid var(--color-border-primary);
  padding: 2rem 0;
}

.footer-bottom-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 1rem;
  padding: 0 2rem;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
}

.copyright {
  color: var(--color-text-secondary);
  font-size: 0.875rem;
  margin: 0;
}

.footer-bottom-links {
  display: flex;
  gap: 2rem;
}

.footer-bottom-links a {
  color: var(--color-text-secondary);
  text-decoration: none;
  font-size: 0.875rem;
  transition: color 0.2s ease;
}

.footer-bottom-links a:hover {
  color: var(--color-accent-primary);
}

@media (max-width: 1024px) {
  .footer-content {
    grid-template-columns: 1fr 1fr 1fr;
    gap: 2rem;
  }

  .footer-section:first-child {
    grid-column: 1 / -1;
  }
}

@media (max-width: 768px) {
  .footer-content {
    grid-template-columns: 1fr;
    gap: 2rem;
    padding: 0 1rem;
  }

  .footer-bottom-content {
    flex-direction: column;
    text-align: center;
    padding: 0 1rem;
  }

  .footer-bottom-links {
    justify-content: center;
  }

  .social-links {
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .site-footer {
    padding: 2rem 0 0 0;
  }

  .footer-content {
    gap: 1.5rem;
    padding: 0 1rem;
  }

  .footer-bottom {
    padding: 1.5rem 0;
  }

  .footer-bottom-content {
    padding: 0 1rem;
  }

  .footer-bottom-links {
    flex-direction: column;
    gap: 1rem;
  }
}

.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, 400px);
  gap: 2rem;
  margin-top: 2rem;
  justify-content: center;
}

@media (max-width: 1024px) {
  .v-container > .d-flex {
    flex-wrap: wrap;
    gap: 1rem;
  }
}

.area-glare-card {
  width: 400px;
  height: 420px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.area-glare-card:hover {
  transform: translateY(-4px);
}

.card-content {
  padding: 1.5rem;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
}

.delete-button {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 0, 0, 0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 10;
}

.delete-button:hover {
  background: rgba(255, 0, 0, 0.4);
  transform: scale(1.1);
}

.delete-button:active {
  transform: scale(0.95);
}

.card-header {
  margin-bottom: 1rem;
}

.card-title {
  color: white;
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 0.5rem;
}

.card-description {
  color: rgba(255, 255, 255, 0.8);
  font-size: 0.9rem;
  line-height: 1.4;
}

.card-services {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.service-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.service-icon {
  width: 24px;
  height: 24px;
  flex-shrink: 0;
}

.service-icon img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: 4px;
}

.service-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}

.service-label {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.service-name {
  color: white;
  font-size: 0.85rem;
  font-weight: 600;
}

</style>


