<template>
  <div class="edit-profile-page">
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
        <h1 class="profile-title">Edit Profile</h1>
        <div class="header-spacer"></div>
      </div>

      <div class="profile-card">
          <div class="profile-section">
            <h3 class="section-title">Profile Photo</h3>
            <div class="avatar-section">
              <div class="profile-avatar" @click="handleImageUpload" :class="{ 'uploading': isUploading }">
                <img
                  v-if="profileImageUrl"
                  :src="profileImageUrl"
                  alt="Profile picture"
                  class="profile-image"
                />
                <v-icon v-else size="48" color="white">mdi-account</v-icon>
                <div v-if="isUploading" class="upload-overlay">
                  <v-progress-circular indeterminate size="24" color="white"></v-progress-circular>
                </div>
                <div v-else class="change-overlay">
                  <v-icon size="20" color="white">mdi-camera</v-icon>
                  <span class="change-text">Change</span>
                </div>
              </div>
              <p class="upload-hint">Click on the photo to change it (PNG, JPG, max 5MB)</p>
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
          </div>

        <form @submit.prevent="handleSave" class="edit-form">
          <div class="profile-section">
            <h3 class="section-title">Personal Information</h3>

            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">First Name</label>
                <div class="input-container">
                  <v-icon class="input-icon" size="20">mdi-account-outline</v-icon>
                  <input
                    v-model="form.first_name"
                    type="text"
                    class="form-input"
                    placeholder="Your first name"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">Last Name</label>
                <div class="input-container">
                  <v-icon class="input-icon" size="20">mdi-account-outline</v-icon>
                  <input
                    v-model="form.last_name"
                    type="text"
                    class="form-input"
                    placeholder="Your last name"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">Phone</label>
                <div class="input-container">
                  <v-icon class="input-icon" size="20">mdi-phone-outline</v-icon>
                  <input
                    v-model="form.phone"
                    type="tel"
                    class="form-input"
                    placeholder="Your phone number"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">Country</label>
                <div class="input-container">
                  <v-icon class="input-icon" size="20">mdi-earth-outline</v-icon>
                  <input
                    v-model="form.country"
                    type="text"
                    class="form-input"
                    placeholder="Your country"
                  />
                </div>
              </div>
            </div>
          </div>

          <div class="profile-section">
            <h3 class="section-title">Change Password</h3>

            <div class="password-section">
              <div class="form-group">
                <label class="form-label">Current Password</label>
                <div class="input-container">
                  <v-icon class="input-icon" size="20">mdi-lock-outline</v-icon>
                  <input
                    v-model="passwordForm.current_password"
                    type="password"
                    class="form-input"
                    placeholder="Your current password"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">New Password</label>
                <div class="input-container">
                  <v-icon class="input-icon" size="20">mdi-lock-outline</v-icon>
                  <input
                    v-model="passwordForm.new_password"
                    type="password"
                    class="form-input"
                    placeholder="New password"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">Confirm New Password</label>
                <div class="input-container">
                  <v-icon class="input-icon" size="20">mdi-lock-check-outline</v-icon>
                  <input
                    v-model="passwordForm.confirm_password"
                    type="password"
                    class="form-input"
                    placeholder="Confirm new password"
                  />
                </div>
              </div>
            </div>
          </div>

          <div v-if="error" class="error-message">
            <v-icon size="16" class="error-icon">mdi-alert-circle</v-icon>
            {{ error }}
          </div>

          <div v-if="successMessage" class="success-message">
            <v-icon size="16" class="success-icon">mdi-check-circle</v-icon>
            {{ successMessage }}
          </div>

          <div class="action-buttons">
            <button type="button" class="cancel-button" @click="goBack">
              <v-icon size="20">mdi-close</v-icon>
              <span>Cancel</span>
            </button>
            <button type="submit" class="save-button" :disabled="isLoading">
              <div v-if="isLoading" class="loading-spinner"></div>
              <v-icon v-else size="20">mdi-check</v-icon>
              <span>{{ isLoading ? 'Saving...' : 'Save' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { currentUser, updateProfile, uploadProfileImage, getProfileImageUrl, isLoading } = useAuth()

const form = ref({
  first_name: '',
  last_name: '',
  phone: '',
  country: ''
})

const passwordForm = ref({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const fileInput = ref<HTMLInputElement | null>(null)
const profileImageUrl = ref<string | null>(null)
const isUploading = ref(false)
const uploadError = ref<string | null>(null)
const error = ref('')
const successMessage = ref('')

const goBack = () => {
  router.push('/profile')
}

const handleImageUpload = () => {
  fileInput.value?.click()
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

    successMessage.value = 'Profile photo updated successfully!'
    setTimeout(() => successMessage.value = '', 3000)
  } catch (error) {
    uploadError.value = error instanceof Error ? error.message : 'Error uploading image'
  } finally {
    isUploading.value = false
  }
}

const handleSave = async () => {
  error.value = ''
  successMessage.value = ''

  if (passwordForm.value.current_password || passwordForm.value.new_password || passwordForm.value.confirm_password) {
    if (!passwordForm.value.current_password || !passwordForm.value.new_password || !passwordForm.value.confirm_password) {
      error.value = 'Please fill in all password fields'
      return
    }

    if (passwordForm.value.new_password !== passwordForm.value.confirm_password) {
      error.value = 'New passwords do not match'
      return
    }

    if (passwordForm.value.new_password.length < 6) {
      error.value = 'New password must contain at least 6 characters'
      return
    }
  }

  try {
    const updateData = {
      first_name: form.value.first_name || undefined,
      last_name: form.value.last_name || undefined,
      phone: form.value.phone || undefined,
      country: form.value.country || undefined,
      current_password: passwordForm.value.current_password || undefined,
      new_password: passwordForm.value.new_password || undefined
    }

    await updateProfile(updateData)

    successMessage.value = 'Profile updated successfully!'
    setTimeout(() => {
      successMessage.value = ''
      router.push('/profile')
    }, 2000)

    passwordForm.value = {
      current_password: '',
      new_password: '',
      confirm_password: ''
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Error updating profile'
  }
}

onMounted(async () => {
  if (currentUser.value) {
    form.value = {
      first_name: currentUser.value.first_name || '',
      last_name: currentUser.value.last_name || '',
      phone: currentUser.value.phone || '',
      country: currentUser.value.country || ''
    }

    profileImageUrl.value = getProfileImageUrl()
  }
})
</script>

<style scoped>
.edit-profile-page {
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
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.profile-avatar {
  width: 100px;
  height: 100px;
  border-radius: var(--radius-full);
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  box-shadow: var(--shadow-glow);
  overflow: hidden;
  cursor: pointer;
  transition: var(--transition-normal);
}

.profile-avatar:hover:not(.uploading) {
  transform: scale(1.05);
  box-shadow:
    var(--shadow-glow),
    0 10px 20px -5px rgba(6, 182, 212, 0.5);
}

.profile-avatar.uploading {
  cursor: not-allowed;
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

.change-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  opacity: 0;
  transition: var(--transition-normal);
}

.profile-avatar:hover .change-overlay {
  opacity: 1;
}

.change-text {
  font-size: 0.75rem;
  font-weight: 600;
  color: white;
  margin-top: 0.25rem;
}


.upload-hint {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin: 0;
  line-height: 1.4;
  text-align: center;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-secondary);
  margin-bottom: 0.5rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 1rem;
  color: var(--color-text-secondary);
  z-index: 1;
}

.form-input {
  width: 100%;
  padding: 1rem 1rem 1rem 3rem;
  background: rgba(15, 23, 42, 0.4);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-size: 1rem;
  transition: var(--transition-normal);
}

.form-input:focus {
  outline: none;
  border-color: var(--color-border-secondary);
  background: rgba(15, 23, 42, 0.6);
}

.form-input::placeholder {
  color: var(--color-text-secondary);
}

.password-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.error-message {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.success-message {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: #22c55e;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  margin-bottom: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.action-buttons {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 2rem;
}

.cancel-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  background: transparent;
  border: 2px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-normal);
}

.cancel-button:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: translateY(-1px);
}

.save-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 1.5rem;
  background: var(--gradient-accent);
  border: none;
  border-radius: var(--radius-lg);
  color: var(--color-text-primary);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-normal);
}

.save-button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow:
    var(--shadow-glow),
    0 10px 20px -5px rgba(6, 182, 212, 0.5);
}

.save-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
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
    align-items: center;
    text-align: center;
  }

  .form-grid,
  .password-section {
    grid-template-columns: 1fr;
  }

  .action-buttons {
    flex-direction: column;
  }
}
</style>
