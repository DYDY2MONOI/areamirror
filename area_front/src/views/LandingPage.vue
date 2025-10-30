<template>
  <div class="landing-page">
    <AnimatedBackground />

    <div v-if="isDesktop" class="globe-background">
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

        <v-divider class="my-2"></v-divider>

        <div class="sidebar-theme-toggle">
          <ThemeToggle />
          <HighContrastToggle @open="openDaltonismModal" />
        </div>
      </v-list>
    </v-navigation-drawer>

    <div class="content">
      <SearchSection
        v-model="searchQuery"
        :suggestions="searchSuggestions"
        @suggestion-click="handleSuggestionClick"
        @clear="searchQuery = ''"
      />

      <v-container class="pt-6 pb-4">
        <UserSection
          :is-authenticated="isAuthenticated"
          :profile-image-url="getProfileImageUrl() || undefined"
          :user-display-name="userDisplayName"
          :user-status="userStatus"
          :action-buttons="actionButtons"
          @login="goToLogin"
          @register="() => router.push('/register')"
          @action-click="handleActionClick"
        />
      </v-container>

      <v-container v-if="isAuthenticated && areas.length > 0" id="my-areas-section" class="mt-6">
        <SectionHeader
          title="My Areas"
          subtitle="Your created automation areas"
          :action-button="{ text: 'View All', icon: 'mdi-arrow-right' }"
          @action-click="goToAllAreas"
        />

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
        <SectionHeader
          title="Create new AREA"
          subtitle="Start building your automation"
        />

        <div class="create-section">
          <div class="floating-icons">
            <div
              v-for="(icon, index) in floatingIcons"
              :key="index"
              class="floating-card"
              :class="`card-${index + 1}`"
            >
              <v-icon size="24" color="white">{{ icon }}</v-icon>
            </div>
          </div>
          <div class="cards-grid">
            <CardButton
              @open="() => requireAuth(() => showCreateModal = true)"
            />
          </div>
        </div>
      </v-container>

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

    <AppFooter @link-click="handleFooterLinkClick" />

    <AppModal
      v-model="showLogoutDialog"
      title="Sign Out"
      subtitle="Are you sure you want to sign out of your account?"
      icon="mdi-logout-variant"
      :cancel-button="{ text: 'Cancel' }"
      :confirm-button="{ text: 'Sign Out' }"
      @confirm="confirmLogout"
    />

    <div v-if="showCreateModal" class="custom-modal-overlay" @click="showCreateModal = false">
      <div class="custom-modal-content" @click.stop>
        <CreateArea
          :template="selectedArea"
          @close="showCreateModal = false"
          @save="handleAreaCreated"
        />
      </div>
    </div>

    <div v-if="showDeleteDialog" class="logout-modal-overlay" @click="showDeleteDialog = false">
      <div class="logout-modal-container" @click.stop>
        <div class="logout-modal-content">
          <div class="logout-modal-header">
            <div class="logout-icon-wrapper">
              <div class="logout-icon-bg delete-icon-bg">
                <v-icon size="24" color="#ef4444">mdi-delete-outline</v-icon>
              </div>
            </div>
            <h2 class="logout-modal-title">Delete Area</h2>
            <p class="logout-modal-subtitle">Are you sure you want to delete "{{ areaToDelete?.name || 'this area' }}"? This action cannot be undone.</p>
          </div>

          <div class="logout-modal-actions">
            <button
              class="logout-action-btn logout-cancel-btn"
              @click="showDeleteDialog = false"
            >
              Cancel
            </button>
            <button
              class="logout-action-btn delete-confirm-btn"
              @click="confirmDelete"
            >
              Delete
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

    <OnboardingTutorial :is-open="showOnboarding" @close="closeOnboarding" />

    <!-- Daltonism Modal -->
    <div v-if="showDaltonismModal" class="daltonism-modal-overlay" @click="onDaltonismCancel">
      <div class="daltonism-modal-container" @click.stop>
        <div class="daltonism-modal-content">
          <div class="daltonism-modal-header">
            <div class="daltonism-icon-wrapper">
              <div class="daltonism-icon-bg">
                <v-icon size="24" color="white">mdi-eye</v-icon>
              </div>
            </div>
            <h2 class="daltonism-modal-title">Daltonism modes</h2>
            <p class="daltonism-modal-subtitle">Choose a color-vision-friendly filter</p>
          </div>

          <div class="daltonism-options">
            <button
              v-for="opt in daltonismOptions"
              :key="opt.value"
              class="daltonism-option"
              :class="{ active: currentDaltonismMode === opt.value }"
              @click="onDaltonismSelect(opt.value)"
            >
              <div class="option-left">
                <div class="option-radio" :class="{ checked: currentDaltonismMode === opt.value }">
                  <div class="dot" />
                </div>
                <span class="option-title">{{ opt.title }}</span>
              </div>
              <v-icon size="18" v-if="currentDaltonismMode === opt.value">mdi-check</v-icon>
            </button>
          </div>

          <div class="daltonism-modal-actions">
            <button class="daltonism-action-btn cancel" @click="onDaltonismCancel">Cancel</button>
          </div>
        </div>
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
import SearchSection from '../components/SearchSection.vue'
import UserSection from '../components/UserSection.vue'
import SectionHeader from '../components/SectionHeader.vue'
import AppFooter from '../components/AppFooter.vue'
import AppModal from '../components/AppModal.vue'
import AnimatedBackground from '../components/AnimatedBackground.vue'
import OnboardingTutorial from '../components/OnboardingTutorial.vue'
import ThemeToggle from '../components/ThemeToggle.vue'
import HighContrastToggle from '../components/HighContrastToggle.vue'
import { ref, watch, onMounted, computed } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useAreas } from '@/composables/useAreas'
import { useRouter } from 'vue-router'
import { type Area } from '@/services/area'
import { useTheme, type DaltonismMode } from '@/composables/useTheme'

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

const showCreateModal = ref(false)
const showLogoutDialog = ref(false)
const showDeleteDialog = ref(false)
const areaToDelete = ref<Area | null>(null)
const showAreaModal = ref(false)
const selectedArea = ref<AreaTemplate | null>(null)
const showOnboarding = ref(false)
const searchQuery = ref('')
const isDesktop = ref(typeof window !== 'undefined' ? window.innerWidth >= 1280 : true)

const { isAuthenticated, currentUser, logout, refreshProfile, getProfileImageUrl } = useAuth()
const { areas, fetchUserAreas, fetchPopularAreas, fetchRecommendedAreas, deleteArea } = useAreas()
const router = useRouter()

const searchSuggestions = ['Gmail', 'Discord', 'Spotify', 'GitHub']

const floatingIcons = [
  'mdi-email-outline',
  'mdi-music-note',
  'mdi-github',
  'mdi-chat',
  'mdi-calendar'
]

const userDisplayName = computed(() => {
  if (!currentUser.value) return 'User Name'
  return `${currentUser.value.first_name || 'User'} ${currentUser.value.last_name || 'Name'}`
})

const userStatus = computed(() => 'Premium Member')

const actionButtons = computed(() => [
  { id: 'search', icon: 'mdi-magnify', tooltip: 'Search' },
  { id: 'notifications', icon: 'mdi-bell-outline', tooltip: 'Notifications' }
])

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

const showDaltonismModal = ref(false)
const previousDaltonismMode = ref<DaltonismMode>('none')
const { daltonismMode: currentDaltonismMode } = useTheme()
const { setDaltonismMode } = useTheme()

const daltonismOptions = [
  { title: 'No filter', value: 'none' as DaltonismMode },
  { title: 'Protanopia', value: 'protanopia' as DaltonismMode },
  { title: 'Deuteranopia', value: 'deuteranopia' as DaltonismMode },
  { title: 'Tritanopia', value: 'tritanopia' as DaltonismMode },
  { title: 'Monochrome', value: 'monochrome' as DaltonismMode },
]

const openDaltonismModal = () => {
  previousDaltonismMode.value = currentDaltonismMode.value
  showDaltonismModal.value = true
  document.body.classList.add('modal-open')
}

const onDaltonismCancel = () => {
  setDaltonismMode(previousDaltonismMode.value)
  showDaltonismModal.value = false
  document.body.classList.remove('modal-open')
}

const onDaltonismSelect = (mode: DaltonismMode) => {
  setDaltonismMode(mode)
  showDaltonismModal.value = false
  document.body.classList.remove('modal-open')
}

onMounted(async () => {
  const onResize = () => {
    isDesktop.value = window.innerWidth >= 1280
  }
  window.addEventListener('resize', onResize)

  await refreshProfile()
  await fetchPopularAreas()
  await fetchRecommendedAreas()
  if (currentUser.value?.id) {
    await fetchUserAreas(currentUser.value.id)
  }

  const isNewUser = localStorage.getItem('area_new_user') === 'true'

  console.log('🔍 Onboarding Check:', {
    isNewUser,
    currentUser: currentUser.value,
    userId: currentUser.value?.id
  })

  if (isNewUser && currentUser.value) {
    const tutorialKey = `area_tutorial_completed_${currentUser.value.id}`
    const tutorialCompleted = localStorage.getItem(tutorialKey) === 'true'

    console.log('📚 Tutorial Status:', {
      tutorialKey,
      tutorialCompleted,
      willShowTutorial: !tutorialCompleted
    })

    if (!tutorialCompleted) {
      setTimeout(() => {
        console.log('🎉 Affichage du tutoriel!')
        showOnboarding.value = true
        localStorage.removeItem('area_new_user')
      }, 1000)
    } else {
      console.log('✅ Tutoriel déjà complété pour cet utilisateur')
      localStorage.removeItem('area_new_user')
    }
  } else {
    console.log('❌ Conditions non remplies pour afficher le tutoriel')
  }

  window.addEventListener('beforeunload', () => {
    window.removeEventListener('resize', onResize)
  })
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

const closeOnboarding = () => {
  showOnboarding.value = false
  if (currentUser.value) {
    const tutorialKey = `area_tutorial_completed_${currentUser.value.id}`
    localStorage.setItem(tutorialKey, 'true')
  }
  localStorage.removeItem('area_new_user')
}

const handleDeleteArea = (area: AreaTemplate | Area) => {
  requireAuth(() => {
    areaToDelete.value = area as Area
    showDeleteDialog.value = true
  })
}

const confirmDelete = async () => {
  if (!areaToDelete.value) return

  try {
    await deleteArea(areaToDelete.value.id)
    showDeleteDialog.value = false
    areaToDelete.value = null
  } catch (error) {
    showDeleteDialog.value = false
    areaToDelete.value = null
    console.error('Failed to delete area:', error)
  }
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
  if (currentUser.value?.id) {
    await fetchUserAreas(currentUser.value.id)
  }
}

const handleSuggestionClick = (suggestion: string) => {
  searchQuery.value = suggestion
}

const handleActionClick = (actionId: string) => {
  requireAuth(() => {
    console.log('Action clicked:', actionId)
  })
}

const handleFooterLinkClick = (linkId: string) => {
  console.log('Footer link clicked:', linkId)
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
  const iconMap: Record<string, string> = {
    "Gmail": "mdi-email-outline",
    "Discord": "mdi-discord",
    "Slack": "mdi-slack",
    "GitHub": "mdi-github",
    "Telegram": "mdi-telegram",
    "Twitter": "mdi-twitter",
    "Instagram": "mdi-instagram",
    "Dropbox": "mdi-dropbox",
    "Notion": "mdi-notebook"
  }
  return iconMap[service] || "mdi-cog"
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
.landing-page {
  background: var(--gradient-bg-primary);
  color: var(--color-text-primary);
  min-height: 100vh;
  overflow-x: hidden;
  width: 100%;
  max-width: 100vw;
  box-sizing: border-box;
  position: relative;
}

.content {
  padding-left: 0;
  max-width: 100%;
  overflow-x: hidden;
  box-sizing: border-box;
  position: relative;
  z-index: 2;
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

.sidebar-desktop {
  display: none;
  position: relative;
  z-index: 3;
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

.sidebar-theme-toggle {
  display: flex;
  flex-direction: column;
  align-items: stretch;
  gap: 8px;
  padding: 8px 12px;
}

.create-section {
  position: relative;
  min-height: 300px;
  display: grid;
  grid-template-columns: auto 1fr;
  gap: 32px;
  align-items: center;
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
  display: grid;
  grid-template-columns: repeat(3, 1fr);
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

.create-section .cards-grid {
  grid-template-columns: 1fr;
  justify-items: start;
}

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

.nav-inner {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 6px;
}

.nav-btn {
  color: white !important;
  text-transform: none;
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

.daltonism-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--gradient-bg-modal, rgba(0,0,0,0.6));
  backdrop-filter: blur(12px);
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  margin: 0;
  width: 100vw;
  height: 100vh;
}

.daltonism-modal-container { max-width: 520px; width: 100%; }

.daltonism-modal-content {
  background: var(--color-bg-card, var(--bg-card));
  border: 1px solid var(--color-border-primary, var(--border-primary));
  border-radius: 20px;
  overflow: hidden;
}

.daltonism-modal-header {
  padding: 24px;
  text-align: center;
  background: var(--gradient-accent);
}

.daltonism-icon-wrapper { display: flex; justify-content: center; margin-bottom: 12px; }
.daltonism-icon-bg {
  width: 56px; height: 56px; border-radius: 9999px;
  background: rgba(255, 255, 255, 0.2);
  display: flex; align-items: center; justify-content: center;
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.daltonism-modal-title { color: var(--color-text-primary); margin: 8px 0 6px; font-size: 1.25rem; font-weight: 700; }
.daltonism-modal-subtitle { color: rgba(255,255,255,0.9); margin: 0; font-size: 0.95rem; }

.daltonism-options { padding: 16px; display: grid; gap: 8px; background: var(--color-bg-card, var(--bg-card)); }
.daltonism-option {
  width: 100%;
  display: flex; align-items: center; justify-content: space-between;
  padding: 12px 14px;
  border-radius: 12px;
  border: 1px solid var(--color-border-primary, var(--border-primary));
  background: var(--color-bg-card, var(--bg-card));
  color: var(--color-text-primary, var(--text-primary));
  cursor: pointer;
  transition: var(--transition-colors);
}
.daltonism-option:hover { background: var(--color-hover-bg, var(--overlay-hover)); }
.daltonism-option.active { box-shadow: 0 0 0 2px var(--color-border-focus, var(--border-accent)); }

.option-left { display: inline-flex; align-items: center; gap: 10px; }
.option-radio { width: 18px; height: 18px; border-radius: 50%; border: 2px solid var(--color-border-primary, var(--border-primary)); display: flex; align-items: center; justify-content: center; }
.option-radio .dot { width: 8px; height: 8px; border-radius: 50%; background: transparent; }
.option-radio.checked .dot { background: var(--color-accent-primary, var(--accent-primary)); }
.option-title { font-weight: 600; }

.daltonism-modal-actions { padding: 16px; display: flex; justify-content: center; gap: 12px; }
.daltonism-action-btn {
  padding: 10px 16px; border-radius: 10px; font-weight: 600; cursor: pointer; transition: var(--transition-colors);
  border: 1px solid var(--color-border-primary, var(--border-primary)); background: var(--color-bg-card, var(--bg-card)); color: var(--color-text-primary, var(--text-primary));
}
.daltonism-action-btn.cancel:hover { background: var(--color-hover-bg, var(--overlay-hover)); }

.logout-modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--gradient-bg-modal, rgba(0, 0, 0, 0.6));
  backdrop-filter: blur(12px);
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px;
  animation: modalOverlayFadeIn 0.3s ease-out;
}

.logout-modal-container {
  max-width: 520px;
  width: 100%;
  max-height: 90vh;
  overflow-y: auto;
}

.logout-modal-content {
  background: var(--color-bg-card, var(--bg-card));
  border: 1px solid var(--color-border-primary, var(--border-primary));
  border-radius: 20px;
  overflow: hidden;
  animation: modalContentSlideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.logout-modal-header {
  padding: 24px;
  text-align: center;
  background: var(--gradient-accent);
}

.logout-icon-wrapper {
  display: flex;
  justify-content: center;
  margin-bottom: 12px;
}

.logout-icon-bg {
  width: 56px;
  height: 56px;
  border-radius: 9999px;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.delete-icon-bg {
  background: rgba(239, 68, 68, 0.2);
}

.logout-modal-title {
  color: var(--color-text-primary);
  margin: 8px 0 6px;
  font-size: 1.25rem;
  font-weight: 700;
}

.logout-modal-subtitle {
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-size: 0.95rem;
  line-height: 1.5;
}

.logout-modal-actions {
  padding: 16px 24px 24px 24px;
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}

.logout-action-btn {
  padding: 10px 16px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-colors);
  border: 1px solid var(--color-border-primary, var(--border-primary));
  background: var(--color-bg-card, var(--bg-card));
  color: var(--color-text-primary, var(--text-primary));
  font-size: 0.95rem;
}

.logout-cancel-btn:hover {
  background: var(--color-hover-bg, var(--overlay-hover));
}

.delete-confirm-btn {
  background: #ef4444 !important;
  color: white !important;
  border-color: #dc2626 !important;
}

.delete-confirm-btn:hover {
  background: #dc2626 !important;
  box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(5deg); }
}

@keyframes fadeUp {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes modalOverlayFadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes modalContentSlideUp {
  from { opacity: 0; transform: translateY(30px) scale(0.9); }
  to { opacity: 1; transform: translateY(0) scale(1); }
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

  .create-section {
    min-height: 250px;
  }

  .floating-card {
    width: 50px;
    height: 50px;
  }

  .cards-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .create-section {
    grid-template-columns: 1fr;
    gap: 16px;
  }

  .create-section .cards-grid {
    justify-items: center;
  }
}

@media (max-width: 768px) {
  .globe-background {
    display: none;
  }

  .cards-grid {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .workflow-steps {
    flex-direction: column;
    gap: 1rem;
  }

  .workflow-arrow {
    transform: rotate(90deg);
  }

  .area-modal-actions {
    flex-direction: column;
  }

  .area-modal-close-btn,
  .area-modal-create-btn {
    width: 100%;
  }
}

@media (prefers-reduced-motion: reduce) {
  .cards-grid .card-spotlight-container {
    transition: none !important;
    animation: none !important;
  }
}

html, body {
  overflow-x: hidden !important;
  max-width: 100vw !important;
}

body.modal-open {
  overflow: hidden !important;
}
</style>








