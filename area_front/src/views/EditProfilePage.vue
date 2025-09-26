<template>
  <div class="edit-profile-page">
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
        <span>Back to Profile</span>
      </button>
      <div class="header-content">
        <h1 class="page-title">Edit Your Profile</h1>
        <p class="page-subtitle">Customize your personal information and preferences</p>
      </div>
    </div>

    <div class="content-container">
      <div class="photo-section">
        <div class="photo-card">
          <div class="photo-header">
            <h3 class="section-title">
              <v-icon size="24" class="title-icon">mdi-camera</v-icon>
              Profile Photo
            </h3>
          </div>
          <div class="photo-content">
            <div class="avatar-container">
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
                    <v-icon size="24" color="white">mdi-camera-plus</v-icon>
                    <span class="change-text">Change Photo</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="photo-info">
              <p class="upload-hint">Click on your photo to upload a new one</p>
              <p class="file-info">Supports PNG, JPG up to 5MB</p>
            </div>
            <input
              ref="fileInput"
              type="file"
              accept="image/*"
              @change="onFileSelected"
              style="display: none"
            />
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
              <v-icon size="24" class="title-icon">mdi-account-edit</v-icon>
              Personal Information
            </h3>
            <p class="section-description">Update your basic profile details</p>
          </div>
          <div class="card-content">
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">First Name</label>
                <div class="input-wrapper">
                  <div class="input-icon">
                    <v-icon size="20">mdi-account-outline</v-icon>
                  </div>
                  <input
                    v-model="form.first_name"
                    type="text"
                    class="form-input"
                    placeholder="Enter your first name"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">Last Name</label>
                <div class="input-wrapper">
                  <div class="input-icon">
                    <v-icon size="20">mdi-account-outline</v-icon>
                  </div>
                  <input
                    v-model="form.last_name"
                    type="text"
                    class="form-input"
                    placeholder="Enter your last name"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">Phone Number</label>
                <div class="input-wrapper">
                  <div class="input-icon">
                    <v-icon size="20">mdi-phone-outline</v-icon>
                  </div>
                  <input
                    v-model="form.phone"
                    type="tel"
                    class="form-input"
                    placeholder="Enter your phone number"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">Country</label>
                <div class="input-wrapper">
                  <div class="input-icon">
                    <v-icon size="20">mdi-earth-outline</v-icon>
                  </div>
                  <input
                    v-model="form.country"
                    type="text"
                    class="form-input"
                    placeholder="Enter your country"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="password-section">
        <div class="password-card">
          <div class="card-header">
            <h3 class="section-title">
              <v-icon size="24" class="title-icon">mdi-lock-reset</v-icon>
              Change Password
            </h3>
            <p class="section-description">Update your account password</p>
          </div>
          <div class="card-content">
            <div class="password-grid">
              <div class="form-group">
                <label class="form-label">Current Password</label>
                <div class="input-wrapper">
                  <div class="input-icon">
                    <v-icon size="20">mdi-lock-outline</v-icon>
                  </div>
                  <input
                    v-model="passwordForm.current_password"
                    type="password"
                    class="form-input"
                    placeholder="Enter your current password"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">New Password</label>
                <div class="input-wrapper">
                  <div class="input-icon">
                    <v-icon size="20">mdi-lock-outline</v-icon>
                  </div>
                  <input
                    v-model="passwordForm.new_password"
                    type="password"
                    class="form-input"
                    placeholder="Enter your new password"
                  />
                </div>
              </div>

              <div class="form-group">
                <label class="form-label">Confirm New Password</label>
                <div class="input-wrapper">
                  <div class="input-icon">
                    <v-icon size="20">mdi-lock-check-outline</v-icon>
                  </div>
                  <input
                    v-model="passwordForm.confirm_password"
                    type="password"
                    class="form-input"
                    placeholder="Confirm your new password"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="error" class="message error-message">
        <v-icon size="20">mdi-alert-circle</v-icon>
        <span>{{ error }}</span>
      </div>

      <div v-if="successMessage" class="message success-message">
        <v-icon size="20">mdi-check-circle</v-icon>
        <span>{{ successMessage }}</span>
      </div>

      <!-- Action Buttons -->
      <div class="actions-section">
        <div class="action-buttons">
          <button type="button" class="btn btn-secondary" @click="goBack">
            <v-icon size="20">mdi-close</v-icon>
            <span>Cancel</span>
          </button>
          <button type="submit" class="btn btn-primary" :disabled="isLoading" @click="handleSave">
            <div v-if="isLoading" class="btn-spinner"></div>
            <v-icon v-else size="20">mdi-content-save</v-icon>
            <span>{{ isLoading ? 'Saving Changes...' : 'Save Changes' }}</span>
          </button>
        </div>
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
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #334155 100%);
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
  background: linear-gradient(45deg, rgba(6, 182, 212, 0.1), rgba(59, 130, 246, 0.1));
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
  background: radial-gradient(circle at 30% 20%, rgba(6, 182, 212, 0.1) 0%, transparent 50%),
              radial-gradient(circle at 70% 80%, rgba(59, 130, 246, 0.1) 0%, transparent 50%);
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
  background: linear-gradient(135deg, #06b6d4, #3b82f6, #8b5cf6);
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

/* Cards */
.photo-card,
.info-card,
.password-card {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 24px;
  padding: 2rem;
  backdrop-filter: blur(20px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.photo-card:hover,
.info-card:hover,
.password-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 30px 60px rgba(0, 0, 0, 0.15);
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
  color: #06b6d4;
}

.section-description {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
  margin: 0;
}

.photo-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
}

.avatar-container {
  position: relative;
}

.profile-avatar {
  width: 120px;
  height: 120px;
  border-radius: 50%;
  background: linear-gradient(135deg, #06b6d4, #3b82f6);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 10px 30px rgba(6, 182, 212, 0.3);
}

.profile-avatar:hover:not(.uploading) {
  transform: scale(1.1);
  box-shadow: 0 20px 40px rgba(6, 182, 212, 0.4);
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
  font-size: 0.875rem;
  font-weight: 600;
  color: white;
}

.photo-info {
  text-align: center;
}

.upload-hint {
  font-size: 1rem;
  color: white;
  margin: 0 0 0.5rem 0;
  font-weight: 500;
}

.file-info {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.form-grid,
.password-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.form-group {
  margin-bottom: 0;
}

.form-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: white;
  margin-bottom: 0.75rem;
  letter-spacing: 0.025em;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 1rem;
  color: rgba(255, 255, 255, 0.5);
  z-index: 2;
}

.form-input {
  width: 100%;
  padding: 1rem 1rem 1rem 3rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 16px;
  color: white;
  font-size: 1rem;
  transition: all 0.3s ease;
  backdrop-filter: blur(10px);
}

.form-input:focus {
  outline: none;
  border-color: #06b6d4;
  background: rgba(255, 255, 255, 0.15);
  box-shadow: 0 0 0 3px rgba(6, 182, 212, 0.1);
}

.form-input::placeholder {
  color: rgba(255, 255, 255, 0.5);
}

.message {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 1.5rem;
  border-radius: 16px;
  font-weight: 500;
  margin: 1rem 0;
}

.error-message {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
}

.success-message {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.3);
  color: #22c55e;
}

.actions-section {
  margin-top: 2rem;
}

.action-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem 2rem;
  border-radius: 16px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
  text-decoration: none;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: white;
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.2);
  border-color: rgba(255, 255, 255, 0.3);
  transform: translateY(-2px);
}

.btn-primary {
  background: linear-gradient(135deg, #06b6d4, #3b82f6);
  color: white;
  box-shadow: 0 10px 30px rgba(6, 182, 212, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-3px);
  box-shadow: 0 20px 40px rgba(6, 182, 212, 0.4);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

.btn-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* Responsive Design */
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

  .photo-card,
  .info-card,
  .password-card {
    padding: 1.5rem;
  }

  .form-grid,
  .password-grid {
    grid-template-columns: 1fr;
  }

  .action-buttons {
    flex-direction: column;
  }

  .btn {
    justify-content: center;
  }
}

@media (max-width: 480px) {
  .page-title {
    font-size: 1.75rem;
  }

  .profile-avatar {
    width: 100px;
    height: 100px;
  }

  .section-title {
    font-size: 1.25rem;
  }
}
</style>
