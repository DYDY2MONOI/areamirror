<template>
  <div class="profile-page">
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

const router = useRouter()
const { currentUser, isAuthenticated, uploadProfileImage, getProfileImageUrl, refreshProfile, isLoading } = useAuth()

const fileInput = ref<HTMLInputElement | null>(null)
const profileImageUrl = ref<string | null>(null)
const isUploading = ref(false)
const uploadError = ref<string | null>(null)

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

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-30px) rotate(10deg); }
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

.page-subtitle {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0;
  font-weight: 400;
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

.profile-card:hover,
.info-card:hover,
.stats-card:hover,
.actions-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 30px 60px rgba(0, 0, 0, 0.15);
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

.profile-email {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 1rem 0;
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

.title-icon {
  color: var(--color-accent-primary);
}

.section-description {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
  margin: 0;
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

.info-item:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-2px);
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

.info-value {
  font-size: 1rem;
  font-weight: 600;
  color: white;
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

.stat-card:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
  transform: translateY(-3px);
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

.stat-label {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
  font-weight: 500;
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

.action-btn.secondary:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
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
</style>
