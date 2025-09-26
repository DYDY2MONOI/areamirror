<template>
  <div class="profile-page">
    <div class="profile-background">
      <div class="geometric-shape shape-1"></div>
      <div class="geometric-shape shape-2"></div>
      <div class="geometric-shape shape-3"></div>
    </div>

    <div class="profile-container">
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
              <img
                v-if="profileImageUrl"
                :src="profileImageUrl"
                alt="Photo de profil"
                class="profile-image"
              />
              <v-icon v-else size="48" color="white">mdi-account</v-icon>
              <div v-if="isUploading" class="upload-overlay">
                <v-progress-circular indeterminate size="24" color="white"></v-progress-circular>
              </div>
            </div>
            <button class="edit-avatar-btn" @click="handleImageUpload" :disabled="isUploading">
              <v-icon size="16">mdi-camera</v-icon>
            </button>
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              @change="onFileSelected"
              style="display: none"
            />
          </div>

          <div v-if="uploadError" class="error-message">
            {{ uploadError }}
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
          <h3 class="section-title">Account Actions</h3>
          <div class="actions-grid">
            <button class="action-button primary" @click="editProfile">
              <v-icon size="20">mdi-pencil</v-icon>
              <span>Edit Profile</span>
            </button>
            <button class="action-button secondary" @click="changePassword">
              <v-icon size="20">mdi-key</v-icon>
              <span>Change Password</span>
            </button>
            <button class="action-button secondary" @click="manageNotifications">
              <v-icon size="20">mdi-bell</v-icon>
              <span>Notifications</span>
            </button>
            <button class="action-button danger" @click="deleteAccount">
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

  // Vérifier le type de fichier
  if (!file.type.startsWith('image/')) {
    uploadError.value = 'Veuillez sélectionner un fichier image valide'
    return
  }

  // Vérifier la taille (5MB max)
  if (file.size > 5 * 1024 * 1024) {
    uploadError.value = 'L\'image ne doit pas dépasser 5MB'
    return
  }

  try {
    isUploading.value = true
    uploadError.value = null

    await uploadProfileImage(file)

    // Mettre à jour l'URL de l'image
    profileImageUrl.value = getProfileImageUrl()

    // Réinitialiser l'input
    if (target) {
      target.value = ''
    }
  } catch (error) {
    console.error('Erreur lors de l\'upload:', error)
    uploadError.value = error instanceof Error ? error.message : 'Erreur lors de l\'upload de l\'image'
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
    // TODO: Implémenter le changement de mot de passe
    console.log('Changement de mot de passe')
  })
}

const manageNotifications = () => {
  requireAuth(() => {
    // TODO: Implémenter la gestion des notifications
    console.log('Gestion des notifications')
  })
}

const deleteAccount = () => {
  requireAuth(() => {
    // TODO: Implémenter la suppression du compte
    console.log('Suppression du compte')
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
  overflow: hidden;
}

.profile-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: var(--radius-full);
}

.upload-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
}

.error-message {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  margin-top: 1rem;
  text-align: center;
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

.edit-avatar-btn:hover:not(:disabled) {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: scale(1.1);
}

.edit-avatar-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
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

@media (prefers-reduced-motion: reduce) {
  .geometric-shape,
  .profile-card,
  .stat-card,
  .action-button {
    animation: none !important;
  }

  .action-button:hover,
  .stat-card:hover,
  .back-button:hover {
    transform: none;
  }
}
</style>
