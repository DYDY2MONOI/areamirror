<template>
  <div class="profile-image-upload">
    <div class="avatar-container">
      <img
        v-if="imageUrl"
        :src="imageUrl"
        alt="Photo de profil"
        class="profile-image"
      />
      <div v-else class="default-avatar">
        <v-icon size="48" color="white">mdi-account</v-icon>
      </div>

      <div v-if="isUploading" class="upload-overlay">
        <v-progress-circular indeterminate size="24" color="white"></v-progress-circular>
      </div>
    </div>

    <button
      class="upload-button"
      @click="handleUpload"
      :disabled="isUploading"
    >
      <v-icon size="16">mdi-camera</v-icon>
    </button>

    <input
      ref="fileInput"
      type="file"
      accept="image/*"
      @change="onFileSelected"
      style="display: none"
    />

    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '@/composables/useAuth'

const { uploadProfileImage, getProfileImageUrl } = useAuth()

const fileInput = ref<HTMLInputElement | null>(null)
const imageUrl = ref<string | null>(null)
const isUploading = ref(false)
const error = ref<string | null>(null)

const handleUpload = () => {
  fileInput.value?.click()
}

const onFileSelected = async (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (!file) return

  if (!file.type.startsWith('image/')) {
    error.value = 'Veuillez sélectionner un fichier image valide'
    return
  }

  if (file.size > 5 * 1024 * 1024) {
    error.value = 'L\'image ne doit pas dépasser 5MB'
    return
  }

  try {
    isUploading.value = true
    error.value = null

    await uploadProfileImage(file)

    imageUrl.value = getProfileImageUrl()

    if (target) {
      target.value = ''
    }
  } catch (err) {
    console.error('Erreur lors de l\'upload:', err)
    error.value = err instanceof Error ? err.message : 'Erreur lors de l\'upload de l\'image'
  } finally {
    isUploading.value = false
  }
}

imageUrl.value = getProfileImageUrl()
</script>

<style scoped>
.profile-image-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.avatar-container {
  position: relative;
  width: 80px;
  height: 80px;
  border-radius: 50%;
  overflow: hidden;
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-glow);
}

.profile-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.default-avatar {
  width: 100%;
  height: 100%;
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
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-button {
  position: absolute;
  bottom: -5px;
  right: -5px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--color-bg-card);
  border: 2px solid var(--color-border-primary);
  color: var(--color-text-primary);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: var(--transition-normal);
}

.upload-button:hover:not(:disabled) {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: scale(1.1);
}

.upload-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.error-message {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #ef4444;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  font-size: 0.875rem;
  text-align: center;
  max-width: 300px;
}
</style>
