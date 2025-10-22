<template>
  <div class="all-areas-page">
    <div class="space-background">
      <div class="stars"></div>
      <div class="stars2"></div>
      <div class="stars3"></div>
      <div class="nebula"></div>
    </div>

    <div class="page-header">
      <button class="back-btn" @click="goBack">
        <v-icon size="20">mdi-arrow-left</v-icon>
        <span>Back</span>
      </button>
      <div class="header-content">
        <h1 class="page-title">My Areas</h1>
        <p class="page-subtitle">Manage all your automation areas</p>
      </div>
    </div>

    <div class="content-container">
      <div class="search-filter-section">
        <div class="search-wrapper">
          <v-icon class="search-icon" size="20">mdi-magnify</v-icon>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search areas by name, description, trigger or action service..."
            class="search-input"
          />
          <button v-if="searchQuery" @click="searchQuery = ''" class="clear-btn">
            <v-icon size="18">mdi-close</v-icon>
          </button>
        </div>

        <div class="filter-options">
          <button
            :class="['filter-chip', { active: statusFilter === 'all' }]"
            @click="statusFilter = 'all'"
          >
            All ({{ areas.length }})
          </button>
          <button
            :class="['filter-chip', { active: statusFilter === 'active' }]"
            @click="statusFilter = 'active'"
          >
            Active ({{ activeAreasCount }})
          </button>
          <button
            :class="['filter-chip', { active: statusFilter === 'inactive' }]"
            @click="statusFilter = 'inactive'"
          >
            Inactive ({{ inactiveAreasCount }})
          </button>
        </div>
      </div>

      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon">
            <v-icon size="24">mdi-vector-square</v-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ areas.length }}</div>
            <div class="stat-label">Total Areas</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon active">
            <v-icon size="24">mdi-check-circle</v-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ activeAreasCount }}</div>
            <div class="stat-label">Active</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon inactive">
            <v-icon size="24">mdi-pause-circle</v-icon>
          </div>
          <div class="stat-content">
            <div class="stat-number">{{ inactiveAreasCount }}</div>
            <div class="stat-label">Inactive</div>
          </div>
        </div>
      </div>

      <div v-if="isLoading" class="loading-state">
        <v-progress-circular indeterminate size="48" color="primary"></v-progress-circular>
        <p>Loading your areas...</p>
      </div>

      <div v-else-if="filteredAreas.length === 0" class="empty-state">
        <v-icon size="64" color="rgba(255, 255, 255, 0.3)">mdi-magnify-close</v-icon>
        <h3>No areas found</h3>
        <p v-if="searchQuery">Try adjusting your search terms</p>
        <p v-else>You haven't created any areas yet</p>
        <button class="create-btn" @click="goToCreate">
          <v-icon size="20">mdi-plus</v-icon>
          <span>Create Your First Area</span>
        </button>
      </div>

      <div v-else class="areas-grid">
        <div v-for="area in filteredAreas" :key="area.id" class="area-card" @click="editArea(area)">
          <div class="area-card-header">
            <div class="area-status">
              <div :class="['status-indicator', { active: (area as any).is_active || (area as any).isActive }]"></div>
              <span class="status-text">{{ ((area as any).is_active || (area as any).isActive) ? 'Active' : 'Inactive' }}</span>
            </div>
            <div class="area-actions">
              <button class="action-btn" @click.stop="toggleArea(area)" :title="((area as any).is_active || (area as any).isActive) ? 'Deactivate' : 'Activate'">
                <v-icon size="18">{{ ((area as any).is_active || (area as any).isActive) ? 'mdi-pause' : 'mdi-play' }}</v-icon>
              </button>
              <button class="action-btn" @click.stop="editArea(area)" title="Edit">
                <v-icon size="18">mdi-pencil</v-icon>
              </button>
              <button class="action-btn danger" @click.stop="confirmDelete(area)" title="Delete">
                <v-icon size="18">mdi-delete</v-icon>
              </button>
            </div>
          </div>

          <div class="area-card-content">
            <h3 class="area-name">{{ area.name }}</h3>
            <p class="area-description">{{ area.description || 'No description' }}</p>

            <div class="area-services">
              <div class="service-flow">
                <div class="service-badge trigger">
                  <v-icon size="16">mdi-lightning-bolt</v-icon>
                  <span>{{ (area as any).trigger_service || (area as any).triggerService }}</span>
                </div>
                <v-icon size="20" class="arrow-icon">mdi-arrow-right</v-icon>
                <div class="service-badge action">
                  <v-icon size="16">mdi-cog</v-icon>
                  <span>{{ (area as any).action_service || (area as any).actionService }}</span>
                </div>
              </div>
            </div>

            <div class="area-meta">
              <div class="meta-item">
                <v-icon size="14">mdi-calendar</v-icon>
                <span>{{ formatDate((area as any).created_at || (area as any).createdAt) }}</span>
              </div>
              <div v-if="(area as any).run_count || (area as any).runCount" class="meta-item">
                <v-icon size="14">mdi-counter</v-icon>
                <span>{{ (area as any).run_count || (area as any).runCount }} runs</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <v-dialog v-model="showDeleteDialog" max-width="400">
      <div class="delete-dialog">
        <div class="dialog-header">
          <v-icon size="48" color="#ef4444">mdi-alert-circle</v-icon>
          <h2>Delete Area?</h2>
        </div>
        <p>Are you sure you want to delete "{{ areaToDelete?.name }}"? This action cannot be undone.</p>
        <div class="dialog-actions">
          <button class="dialog-btn cancel" @click="showDeleteDialog = false">Cancel</button>
          <button class="dialog-btn delete" @click="deleteArea">Delete</button>
        </div>
      </div>
    </v-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useAreas } from '@/composables/useAreas'
import type { Area } from '@/services/area'

const router = useRouter()
const { currentUser, isAuthenticated } = useAuth()
const { areas, isLoading, fetchUserAreas, toggleArea: toggleAreaService, deleteArea: deleteAreaService } = useAreas()

const searchQuery = ref('')
const statusFilter = ref<'all' | 'active' | 'inactive'>('all')
const showDeleteDialog = ref(false)
const areaToDelete = ref<Area | null>(null)

const filteredAreas = computed(() => {
  let filtered = areas.value

  if (statusFilter.value === 'active') {
    filtered = filtered.filter((area: any) => area.is_active || area.isActive)
  } else if (statusFilter.value === 'inactive') {
    filtered = filtered.filter((area: any) => !(area.is_active || area.isActive))
  }

  if (searchQuery.value.trim()) {
    const query = searchQuery.value.toLowerCase().trim()
    filtered = filtered.filter((area: any) => {
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
  }

  return filtered
})

const activeAreasCount = computed(() => {
  return areas.value.filter((area: any) => area.is_active || area.isActive).length
})

const inactiveAreasCount = computed(() => {
  return areas.value.filter((area: any) => !(area.is_active || area.isActive)).length
})

const goBack = () => {
  router.push('/')
}

const goToCreate = () => {
  router.push('/')
}

const formatDate = (dateString?: string) => {
  if (!dateString) return 'Unknown'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const toggleArea = function(area: Area) {
  return toggleAreaService(area.id).catch((error: unknown) => {
    console.error('Error toggling area:', error)
  })
}


const editArea = (area: Area) => {
  router.push({
    name: 'configure-area',
    query: { areaId: area.id }
  })
}

const confirmDelete = (area: Area) => {
  areaToDelete.value = area
  showDeleteDialog.value = true
}

const deleteArea = async () => {
  if (!areaToDelete.value) return

  try {
    await deleteAreaService(areaToDelete.value.id)
    showDeleteDialog.value = false
    areaToDelete.value = null
  } catch (error) {
    console.error('Error deleting area:', error)
  }
}

onMounted(() => {
  if (!isAuthenticated.value) {
    router.push('/login')
    return
  }

  if (currentUser.value?.id) {
    fetchUserAreas(currentUser.value.id)
  }
})
</script>

<style scoped>
.all-areas-page {
  min-height: 100vh;
  background: var(--gradient-bg-primary);
  position: relative;
  overflow-x: hidden;
}

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

.stars, .stars2, .stars3 {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: radial-gradient(2px 2px at 20px 30px, #fff, transparent),
    radial-gradient(2px 2px at 40px 70px, rgba(255,255,255,0.8), transparent),
    radial-gradient(1px 1px at 90px 40px, #fff, transparent);
  background-repeat: repeat;
  background-size: 200px 100px;
  animation: twinkle 4s ease-in-out infinite alternate;
}

.stars2 {
  background-size: 180px 120px;
  animation-duration: 6s;
}

.stars3 {
  background-size: 160px 140px;
  animation-duration: 8s;
}

@keyframes twinkle {
  0% { opacity: 0.8; }
  50% { opacity: 1; }
  100% { opacity: 0.6; }
}

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
  max-width: 1400px;
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
}

.content-container {
  position: relative;
  z-index: 10;
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}

.search-filter-section {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 1.5rem;
  margin-bottom: 2rem;
  backdrop-filter: blur(20px);
}

.search-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  margin-bottom: 1rem;
  transition: all 0.3s ease;
}

.search-wrapper:focus-within {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(135, 81, 209, 0.5);
  box-shadow: 0 0 0 3px rgba(135, 81, 209, 0.1);
}

.search-icon {
  position: absolute;
  left: 1rem;
  color: rgba(255, 255, 255, 0.5);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 1rem 3rem;
  background: transparent;
  border: none;
  color: white;
  font-size: 0.95rem;
  outline: none;
}

.search-input::placeholder {
  color: rgba(255, 255, 255, 0.4);
}

.clear-btn {
  position: absolute;
  right: 0.5rem;
  background: rgba(255, 255, 255, 0.1);
  border: none;
  border-radius: 12px;
  padding: 0.5rem;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.6);
  transition: all 0.3s ease;
}

.clear-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.filter-options {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.filter-chip {
  padding: 0.5rem 1rem;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.3s ease;
}

.filter-chip:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.3);
}

.filter-chip.active {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-color: transparent;
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

/* Stats Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  backdrop-filter: blur(20px);
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.stat-icon.active {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
}

.stat-icon.inactive {
  background: linear-gradient(135deg, #ee0979 0%, #ff6a00 100%);
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 1.75rem;
  font-weight: 700;
  color: white;
}

.stat-label {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
}

.empty-state h3 {
  font-size: 1.5rem;
  color: white;
  margin: 1rem 0 0.5rem;
}

.empty-state p {
  margin: 0 0 1.5rem;
}

.create-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  color: white;
  padding: 0.75rem 1.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.create-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.areas-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, 400px);
  gap: 1.5rem;
  justify-content: center;
}

.area-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 20px;
  padding: 1.5rem;
  backdrop-filter: blur(20px);
  transition: all 0.3s ease;
  cursor: pointer;
}

.area-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  border-color: rgba(255, 255, 255, 0.2);
}

.area-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.area-status {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-indicator {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #ef4444;
}

.status-indicator.active {
  background: #10b981;
  box-shadow: 0 0 10px rgba(16, 185, 129, 0.5);
}

.status-text {
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.6);
  text-transform: uppercase;
  font-weight: 600;
}

.area-actions {
  display: flex;
  gap: 0.5rem;
}

.action-btn {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  background: rgba(255, 255, 255, 0.05);
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: scale(1.1);
}

.action-btn.danger:hover {
  background: rgba(239, 68, 68, 0.2);
  border-color: #ef4444;
}

.area-card-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.area-name {
  font-size: 1.25rem;
  font-weight: 700;
  color: white;
  margin: 0;
}

.area-description {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
  line-height: 1.5;
}

.area-services {
  margin: 0.5rem 0;
}

.service-flow {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.service-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.75rem;
  border-radius: 10px;
  font-size: 0.75rem;
  font-weight: 600;
  flex: 1;
}

.service-badge.trigger {
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #fca5a5;
}

.service-badge.action {
  background: rgba(59, 130, 246, 0.2);
  border: 1px solid rgba(59, 130, 246, 0.3);
  color: #93c5fd;
}

.arrow-icon {
  color: rgba(255, 255, 255, 0.4);
  flex-shrink: 0;
}

.area-meta {
  display: flex;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.75rem;
  color: rgba(255, 255, 255, 0.5);
}

.delete-dialog {
  background: linear-gradient(135deg, #0f0e1e 0%, #1a1632 100%);
  border-radius: 20px;
  padding: 2rem;
  color: white;
}

.dialog-header {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.dialog-header h2 {
  margin: 0;
  font-size: 1.5rem;
}

.delete-dialog p {
  text-align: center;
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 2rem;
}

.dialog-actions {
  display: flex;
  gap: 1rem;
}

.dialog-btn {
  flex: 1;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.dialog-btn.cancel {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.dialog-btn.cancel:hover {
  background: rgba(255, 255, 255, 0.2);
}

.dialog-btn.delete {
  background: linear-gradient(135deg, #ee0979 0%, #ff6a00 100%);
  color: white;
}

.dialog-btn.delete:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(238, 9, 121, 0.4);
}

@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }

  .content-container {
    padding: 1rem;
  }

  .stats-grid {
    grid-template-columns: 1fr;
  }

  .areas-grid {
    grid-template-columns: 1fr;
  }
}
</style>

