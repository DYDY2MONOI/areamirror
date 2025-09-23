<template>
  <div class="create-area-container">
    <div class="main-card">
      <div class="card-header">
        <div class="header-content">
          <div class="header-text">
            <h1 class="card-title">Create New Area</h1>
            <p class="card-subtitle">Connect your favorite services with intelligent automation</p>
          </div>
          <button class="close-button" @click="$emit('close')">
            <v-icon size="24" color="white">mdi-close</v-icon>
          </button>
        </div>
      </div>
      <div class="card-content">
        <div class="form-section">
          <div class="section-label">
            <v-icon class="label-icon" size="20">mdi-tag-outline</v-icon>
            <span class="label-text">Area Details</span>
          </div>

          <div class="input-group">
            <div class="input-container">
              <label class="input-label">Area Name</label>
              <input
                v-model="form.areaName"
                class="modern-input"
                placeholder="Enter a name for your area"
                required
              />
            </div>

            <div class="input-container">
              <label class="input-label">Description</label>
              <textarea
                v-model="form.description"
                class="modern-textarea"
                placeholder="Describe what this area does..."
                rows="3"
              ></textarea>
            </div>
          </div>
        </div>

        <div class="form-section">
          <div class="section-label">
            <v-icon class="label-icon" size="20">mdi-link-variant</v-icon>
            <span class="label-text">Service Connection</span>
          </div>

          <div class="connection-flow">
            <div class="service-selector">
              <div class="selector-header">
                <div class="selector-icon">
                  <v-icon size="24" color="#3b82f6">mdi-play-circle</v-icon>
                </div>
                <div class="selector-info">
                  <h3 class="selector-title">When this happens</h3>
                  <p class="selector-subtitle">Choose the trigger service</p>
                </div>
              </div>

              <div class="service-dropdown">
                <v-select
                  v-model="form.actionService"
                  :items="appItems"
                  item-title="title"
                  item-value="value"
                  placeholder="Select action service"
                  class="modern-select"
                  variant="outlined"
                  hide-details
                >
                  <template #selection="{ item }">
                    <div class="selected-service" v-if="item.raw">
                      <div class="service-avatar">
                        <img :src="item.raw.icon" :alt="item.raw.title" class="service-icon" />
                      </div>
                      <span class="service-name">{{ item.raw.title }}</span>
                    </div>
                    <span v-else class="placeholder-text">Select service</span>
                  </template>
                  <template #item="{ props, item }">
                    <v-list-item v-bind="props" class="service-item">
                      <template #prepend>
                        <div class="service-avatar">
                          <img :src="item.raw.icon" :alt="item.raw.title" class="service-icon" />
                        </div>
                      </template>
                      <v-list-item-title class="service-item-title">{{ item.raw.title }}</v-list-item-title>
                    </v-list-item>
                  </template>
                </v-select>
              </div>
            </div>

            <div class="connection-arrow">
              <div class="arrow-line"></div>
              <div class="arrow-icon">
                <v-icon size="20" color="#3b82f6">mdi-arrow-down</v-icon>
              </div>
              <div class="arrow-line"></div>
            </div>

            <div class="service-selector">
              <div class="selector-header">
                <div class="selector-icon">
                  <v-icon size="24" color="#3b82f6">mdi-lightning-bolt</v-icon>
                </div>
                <div class="selector-info">
                  <h3 class="selector-title">Then do this</h3>
                  <p class="selector-subtitle">Choose the action service</p>
                </div>
              </div>

              <div class="service-dropdown">
                <v-select
                  v-model="form.reactionService"
                  :items="appItems"
                  item-title="title"
                  item-value="value"
                  placeholder="Select reaction service"
                  class="modern-select"
                  variant="outlined"
                  hide-details
                >
                  <template #selection="{ item }">
                    <div class="selected-service" v-if="item.raw">
                      <div class="service-avatar">
                        <img :src="item.raw.icon" :alt="item.raw.title" class="service-icon" />
                      </div>
                      <span class="service-name">{{ item.raw.title }}</span>
                    </div>
                    <span v-else class="placeholder-text">Select service</span>
                  </template>
                  <template #item="{ props, item }">
                    <v-list-item v-bind="props" class="service-item">
                      <template #prepend>
                        <div class="service-avatar">
                          <img :src="item.raw.icon" :alt="item.raw.title" class="service-icon" />
                        </div>
                      </template>
                      <v-list-item-title class="service-item-title">{{ item.raw.title }}</v-list-item-title>
                    </v-list-item>
                  </template>
                </v-select>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="card-actions">
        <button class="action-btn secondary" @click="$emit('close')">
          <v-icon size="18">mdi-close</v-icon>
          Cancel
        </button>
        <button class="action-btn primary" @click="$emit('save')" :disabled="!isFormValid">
          <v-icon size="18">mdi-check</v-icon>
          Create Area
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive } from 'vue'
import appsJson from '../../assets/apps.json'

type AppDef = { name: string; icon: string }
const apps = (Array.isArray(appsJson) ? appsJson : (appsJson as any).apps ?? []) as AppDef[]

const ICONS_DIR = 'app-icons'

const getIconUrl = (file: string) =>
  new URL(`../../assets/${ICONS_DIR}/${file}`, import.meta.url).href

const appItems = computed(() =>
  apps.map(a => ({ title: a.name, value: a.name, icon: getIconUrl(a.icon) }))
)

const form = reactive({
  areaName: '',
  description: '',
  actionService: '' as string | null,
  reactionService: '' as string | null,
})

const isFormValid = computed(() => {
  return form.areaName.trim() !== '' &&
         form.actionService !== '' &&
         form.reactionService !== ''
})

defineEmits<{ (e: 'close'): void; (e: 'save'): void }>()
</script>

<style scoped>
:root {
  --primary-blue: #3b82f6;
  --secondary-purple: #7c3aed;
  --accent-pink: #ec4899;
  --dark-bg: #0f1419;
  --darker-bg: #0a0e13;
  --light-gray: #9ca3af;
  --white: #ffffff;
  --card-bg: rgba(255, 255, 255, 0.05);
  --border: rgba(255, 255, 255, 0.1);
}

.create-area-container {
  position: relative;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  width: 100%;
  background: linear-gradient(135deg,
    var(--darker-bg) 0%,
    var(--dark-bg) 25%,
    #1a1f2e 50%,
    var(--dark-bg) 75%,
    var(--darker-bg) 100%);
  border-radius: 24px;
  overflow: hidden;
}


.main-card {
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(20px);
  border: 1px solid var(--border);
  border-radius: 24px;
  padding: 0;
  width: 100%;
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.3),
    0 0 0 1px rgba(255, 255, 255, 0.05),
    inset 0 1px 0 rgba(255, 255, 255, 0.1);
  overflow: hidden;
  position: relative;
}

.card-header {
  padding: 2rem 2rem 1rem 2rem;
  border-bottom: 1px solid var(--spotify-border);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.close-button {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--border);
  border-radius: 8px;
  padding: 0.5rem;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.2);
}

.icon-container {
  width: 48px;
  height: 48px;
  background: linear-gradient(135deg, var(--spotify-green), var(--spotify-dark-green));
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 12px rgba(29, 185, 84, 0.3);
}

.header-icon {
  color: white !important;
}

.header-text {
  flex: 1;
}

.card-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--white);
  margin: 0 0 0.25rem 0;
  letter-spacing: -0.02em;
}

.card-subtitle {
  font-size: 1rem;
  color: var(--light-gray);
  margin: 0;
  font-weight: 400;
}

.card-content {
  padding: 2rem;
}

.form-section {
  margin-bottom: 2.5rem;
}

.form-section:last-child {
  margin-bottom: 0;
}

.section-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.label-icon {
  color: #3b82f6 !important;
}

.label-text {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--white);
  letter-spacing: -0.01em;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.input-container {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.input-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--light-gray);
  letter-spacing: 0.01em;
}

.modern-input,
.modern-textarea {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--border);
  border-radius: 12px;
  padding: 1rem;
  color: var(--white);
  font-size: 1rem;
  font-weight: 400;
  transition: all 0.2s ease;
  outline: none;
}

.modern-input:focus,
.modern-textarea:focus {
  border-color: #3b82f6;
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.modern-input::placeholder,
.modern-textarea::placeholder {
  color: var(--light-gray);
  opacity: 0.7;
}

.modern-textarea {
  resize: vertical;
  min-height: 80px;
  font-family: inherit;
}

.connection-flow {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.service-selector {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--border);
  border-radius: 16px;
  padding: 1.5rem;
  transition: all 0.2s ease;
}

.service-selector:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(59, 130, 246, 0.3);
}

.selector-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.selector-icon {
  width: 40px;
  height: 40px;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.selector-info {
  flex: 1;
}

.selector-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--white);
  margin: 0 0 0.25rem 0;
  letter-spacing: -0.01em;
}

.selector-subtitle {
  font-size: 0.875rem;
  color: var(--light-gray);
  margin: 0;
}

.service-dropdown {
  margin-top: 1rem;
}

.modern-select {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.selected-service {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.service-avatar {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.service-icon {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.service-name {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--white);
}

.placeholder-text {
  font-size: 0.875rem;
  color: var(--light-gray);
  opacity: 0.7;
  font-style: italic;
}

.service-item {
  padding: 0.75rem 1rem !important;
  border-radius: 8px;
  margin: 0.25rem 0;
  transition: all 0.2s ease;
}

.service-item:hover {
  background: rgba(59, 130, 246, 0.1) !important;
}

.service-item-title {
  font-size: 0.875rem !important;
  font-weight: 500 !important;
  color: var(--spotify-white) !important;
}

.connection-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0.5rem 0;
}

.arrow-line {
  flex: 1;
  height: 2px;
  background: linear-gradient(90deg, transparent, #3b82f6, transparent);
  opacity: 0.5;
}

.arrow-icon {
  width: 32px;
  height: 32px;
  background: rgba(59, 130, 246, 0.1);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 1rem;
}

.card-actions {
  padding: 1.5rem 2rem;
  border-top: 1px solid var(--spotify-border);
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 12px;
  font-size: 0.875rem;
  font-weight: 600;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  letter-spacing: 0.01em;
}

.action-btn.secondary {
  background: rgba(255, 255, 255, 0.05);
  color: var(--light-gray);
  border: 1px solid var(--border);
}

.action-btn.secondary:hover {
  background: rgba(255, 255, 255, 0.08);
  color: var(--white);
}

.action-btn.primary {
  background: linear-gradient(135deg, #3b82f6, #7c3aed, #ec4899);
  color: white;
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.action-btn.primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 6px 16px rgba(59, 130, 246, 0.4);
}

.action-btn.primary:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

@media (max-width: 768px) {
  .create-area-container {
    padding: 1rem;
  }

  .main-card {
    border-radius: 16px;
  }

  .card-header,
  .card-content,
  .card-actions {
    padding: 1.5rem;
  }

  .header-content {
    flex-direction: column;
    text-align: center;
    gap: 1rem;
  }

  .connection-flow {
    gap: 1rem;
  }

  .service-selector {
    padding: 1rem;
  }

  .card-actions {
    flex-direction: column;
  }

  .action-btn {
    width: 100%;
    justify-content: center;
  }
}

@media (max-height: 800px) {
  .main-card {
    max-height: 80vh;
    overflow-y: auto;
  }

  .card-header {
    padding: 1.5rem 2rem 1rem 2rem;
  }

  .card-content {
    padding: 1.5rem 2rem;
  }

  .card-actions {
    padding: 1rem 2rem;
  }
}

.main-card {
  animation: slideInUp 0.4s ease-out;
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

:deep(.v-field__outline) {
  --v-field-border-opacity: 0.1;
}

:deep(.v-field--variant-outlined .v-field__outline) {
  color: var(--spotify-border);
}

:deep(.v-field--focused .v-field__outline) {
  color: #3b82f6;
}

:deep(.v-list) {
  background: rgba(0, 0, 0, 0.8) !important;
  backdrop-filter: blur(20px);
  border: 1px solid var(--border);
  border-radius: 12px;
}

:deep(.v-list-item) {
  color: var(--spotify-white) !important;
}

:deep(.v-list-item:hover) {
  background: rgba(29, 185, 84, 0.1) !important;
}

.create-area-container::-webkit-scrollbar,
.main-card::-webkit-scrollbar {
  width: 8px;
}

.create-area-container::-webkit-scrollbar-track,
.main-card::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
  border-radius: 4px;
}

.create-area-container::-webkit-scrollbar-thumb,
.main-card::-webkit-scrollbar-thumb {
  background: linear-gradient(180deg, #3b82f6, #7c3aed, #ec4899);
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.create-area-container::-webkit-scrollbar-thumb:hover,
.main-card::-webkit-scrollbar-thumb:hover {
  background: linear-gradient(180deg, #2563eb, #6d28d9, #db2777);
  box-shadow: 0 0 8px rgba(59, 130, 246, 0.4);
}

.create-area-container::-webkit-scrollbar-thumb:active,
.main-card::-webkit-scrollbar-thumb:active {
  background: linear-gradient(180deg, #1d4ed8, #5b21b6, #be185d);
}

.create-area-container,
.main-card {
  scrollbar-width: thin;
  scrollbar-color: #3b82f6 rgba(255, 255, 255, 0.05);
}
</style>
