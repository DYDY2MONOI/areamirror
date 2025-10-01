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

              <div class="service-selection">
                <div v-if="!form.triggerService" class="service-grid">
                  <div
                    v-for="item in appItems.slice(0, 8)"
                    :key="item.value"
                    class="service-card"
                    @click="selectTrigger(item.value)"
                  >
                    <div class="service-card-icon">
                      <img :src="item.icon" :alt="item.title" class="service-icon" />
                    </div>
                    <span class="service-card-name">{{ item.title }}</span>
                  </div>
                  <div class="service-card more-services" @click="showAllTriggerServices = true">
                    <div class="service-card-icon">
                      <v-icon size="24" color="#3b82f6">mdi-plus</v-icon>
                    </div>
                    <span class="service-card-name">More...</span>
                  </div>
                </div>

                <div v-else class="selected-service-display">
                  <div class="selected-service-card">
                    <div class="service-avatar">
                      <img :src="getIconUrl(apps.find(a => a.name === form.triggerService)?.icon || '')" :alt="getServiceName(form.triggerService)" class="service-icon" />
                    </div>
                    <div class="service-info">
                      <span class="service-name">{{ getServiceName(form.triggerService) }}</span>
                      <span class="service-type">Trigger Service</span>
                    </div>
                    <button class="change-service-btn" @click="form.triggerService = ''">
                      <v-icon size="16">mdi-close</v-icon>
                    </button>
                  </div>
                </div>
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

              <div class="service-selection">
                <div v-if="!form.actionService" class="service-grid">
                  <div
                    v-for="item in appItems.slice(0, 8)"
                    :key="item.value"
                    class="service-card"
                    @click="selectAction(item.value)"
                  >
                    <div class="service-card-icon">
                      <img :src="item.icon" :alt="item.title" class="service-icon" />
                    </div>
                    <span class="service-card-name">{{ item.title }}</span>
                  </div>
                  <div class="service-card more-services" @click="showAllReactionServices = true">
                    <div class="service-card-icon">
                      <v-icon size="24" color="#3b82f6">mdi-plus</v-icon>
                    </div>
                    <span class="service-card-name">More...</span>
                  </div>
                </div>

                <div v-else class="selected-service-display">
                  <div class="selected-service-card">
                    <div class="service-avatar">
                      <img :src="getIconUrl(apps.find(a => a.name === form.actionService)?.icon || '')" :alt="getServiceName(form.actionService)" class="service-icon" />
                    </div>
                    <div class="service-info">
                      <span class="service-name">{{ getServiceName(form.actionService) }}</span>
                      <span class="service-type">Action Service</span>
                    </div>
                    <button class="change-service-btn" @click="form.actionService = ''">
                      <v-icon size="16">mdi-close</v-icon>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Configuration Section -->
        <div v-if="form.triggerService && form.actionService" class="form-section">
          <div class="section-label">
            <v-icon class="label-icon" size="20">mdi-cog-outline</v-icon>
            <span class="label-text">Configuration</span>
          </div>

          <!-- Calendar Trigger Configuration -->
          <div v-if="form.triggerService === 'Google Calendar'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('google-calendar.png')" alt="Google Calendar" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">📅 Calendar Event Trigger</h4>
                <p class="config-subtitle">Configure when this area should trigger</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">📅 Event Date & Time</label>
                  <input
                    v-model="form.triggerConfig.eventTime"
                    type="datetime-local"
                    class="modern-input"
                    :min="new Date().toISOString().slice(0, 16)"
                    required
                  />
                  <small class="input-hint">Select the date and time when you want to be reminded</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📝 Event Title (Optional)</label>
                  <input
                    v-model="form.triggerConfig.eventTitle"
                    class="modern-input"
                    placeholder="e.g., Meeting with John, Doctor Appointment"
                  />
                  <small class="input-hint">This will be used in the email subject and body</small>
                </div>

                <div class="input-container">
                  <label class="input-label">🗓️ Calendar ID</label>
                  <input
                    v-model="form.triggerConfig.calendarId"
                    class="modern-input"
                    placeholder="primary"
                  />
                  <small class="input-hint">Use 'primary' for your main calendar</small>
                </div>
              </div>
            </div>
          </div>

          <!-- Gmail Action Configuration -->
          <div v-if="form.actionService === 'Gmail'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <img :src="getIconUrl('gmail.png')" alt="Gmail" class="service-icon" />
              </div>
              <div class="config-info">
                <h4 class="config-title">📧 Gmail Action</h4>
                <p class="config-subtitle">Configure the email to be sent</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">📧 Send Email To</label>
                  <input
                    v-model="form.actionConfig.toEmail"
                    type="email"
                    class="modern-input"
                    placeholder="your-email@gmail.com"
                    required
                  />
                  <small class="input-hint">Enter the email address where you want to receive the reminder</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📝 Email Subject</label>
                  <input
                    v-model="form.actionConfig.subject"
                    class="modern-input"
                    placeholder="Reminder: {{eventTitle}}"
                    required
                  />
                  <small class="input-hint">Use &#123;&#123;eventTitle&#125;&#125; to include the event name</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📄 Email Body</label>
                  <textarea
                    v-model="form.actionConfig.body"
                    class="modern-textarea"
                    placeholder="Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}"
                    rows="4"
                    required
                  ></textarea>
                  <small class="input-hint">Use &#123;&#123;eventTitle&#125;&#125;, &#123;&#123;eventTime&#125;&#125;, and &#123;&#123;areaName&#125;&#125; as placeholders</small>
                </div>
              </div>
            </div>
          </div>

          <!-- Email Preview for Calendar → Gmail -->
          <div v-if="form.triggerService === 'Google Calendar' && form.actionService === 'Gmail'" class="preview-section">
            <div class="preview-header">
              <v-icon class="preview-icon" size="20">mdi-eye-outline</v-icon>
              <span class="preview-title">Email Preview</span>
            </div>
            <div class="preview-content">
              <div class="email-preview">
                <div class="email-header">
                  <strong>To:</strong> {{ form.actionConfig.toEmail || 'your-email@gmail.com' }}
                </div>
                <div class="email-header">
                  <strong>Subject:</strong> {{ form.actionConfig.subject || 'Reminder: Event Title' }}
                </div>
                <div class="email-body">
                  {{ form.actionConfig.body || 'Hello! This is a reminder about your upcoming event: Event Title at Event Time' }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="error" class="error-message">
        <v-icon size="20" color="#ef4444">mdi-alert-circle</v-icon>
        <span>{{ error }}</span>
      </div>

      <div class="card-actions">
        <button class="action-btn secondary" @click="$emit('close')">
          <v-icon size="18">mdi-close</v-icon>
          Cancel
        </button>
        <button class="action-btn primary" @click="createArea" :disabled="!isFormValid || isLoading">
          <v-icon size="18">mdi-check</v-icon>
          {{ isLoading ? 'Creating...' : 'Create Area' }}
        </button>
        
        <!-- Test Email Button for Calendar → Gmail -->
        <button 
          v-if="form.triggerService === 'Google Calendar' && form.actionService === 'Gmail'"
          class="action-btn test-email-btn" 
          @click="sendTestEmail" 
          :disabled="!canSendTestEmail || isSendingTest"
        >
          <v-icon size="18">mdi-email-send</v-icon>
          {{ isSendingTest ? 'Sending...' : 'Send Test Email' }}
        </button>
        
        <!-- Debug info for Calendar → Gmail -->
        <div v-if="form.triggerService === 'Google Calendar' && form.actionService === 'Gmail'" class="debug-info">
          <small style="color: #666; font-size: 0.75rem;">
            Debug: {{ isFormValid ? '✅ Ready to create' : '❌ Missing: ' + getMissingFields() }}
          </small>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, onMounted, onUnmounted, nextTick, watch } from 'vue'
import appsJson from '../../assets/apps.json'
import { areaService } from '../../services/area'

type AppDef = { name: string; icon: string }
const apps = (Array.isArray(appsJson) ? appsJson : (appsJson as any).apps ?? []) as AppDef[]

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

const props = defineProps<{
  template?: AreaTemplate | null
}>()

const ICONS_DIR = 'app-icons'

const getIconUrl = (file: string) =>
  new URL(`../../assets/${ICONS_DIR}/${file}`, import.meta.url).href

const appItems = computed(() =>
  apps.map(a => ({ title: a.name, value: a.name, icon: getIconUrl(a.icon) }))
)

const form = reactive({
  areaName: '',
  description: '',
  triggerService: '' as string | null,
  actionService: '' as string | null,
  triggerConfig: {} as any,
  actionConfig: {} as any,
})

watch(() => props.template, (newTemplate) => {
  if (newTemplate) {
    form.areaName = newTemplate.title
    form.description = newTemplate.description
    form.triggerService = newTemplate.triggerService
    form.actionService = newTemplate.actionService
    
    if (newTemplate.triggerService === 'Google Calendar' && newTemplate.actionService === 'Gmail') {
      form.triggerConfig = {
        eventTime: '',
        eventTitle: '',
        calendarId: 'primary'
      }
      form.actionConfig = {
        toEmail: '',
        subject: 'Reminder: {{eventTitle}}',
        body: 'Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}.\n\nArea: {{areaName}}'
      }
    }
  }
}, { immediate: true })

const isFormValid = computed(() => {
  const hasBasicInfo = form.areaName.trim() !== '' &&
                      form.triggerService !== '' &&
                      form.actionService !== ''
  
  if (form.triggerService === 'Google Calendar' && form.actionService === 'Gmail') {
    return hasBasicInfo &&
           form.triggerConfig.eventTime &&
           form.actionConfig.toEmail &&
           form.actionConfig.subject
  }
  
  return hasBasicInfo
})

const showAllTriggerServices = ref(false)
const showAllReactionServices = ref(false)


const selectTrigger = (serviceId: string) => {
  form.triggerService = serviceId
  showAllTriggerServices.value = false
  
  if (serviceId === 'Google Calendar') {
    form.triggerConfig = {
      eventTime: '',
      eventTitle: '',
      calendarId: 'primary'
    }
  } else {
    form.triggerConfig = {}
  }
}

const selectAction = (serviceId: string) => {
  form.actionService = serviceId
  showAllReactionServices.value = false
  
  if (serviceId === 'Gmail') {
    form.actionConfig = {
      toEmail: '',
      subject: 'Reminder: {{eventTitle}}',
      body: 'Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}.\n\nArea: {{areaName}}'
    }
  }
}

const getServiceName = (serviceId: string) => {
  const service = appItems.value.find(item => item.value === serviceId)
  return service?.title || ''
}

const getMissingFields = () => {
  const missing = []
  if (!form.areaName.trim()) missing.push('Area Name')
  if (!form.triggerConfig.eventTime) missing.push('Event Time')
  if (!form.actionConfig.toEmail) missing.push('Email Address')
  if (!form.actionConfig.subject) missing.push('Email Subject')
  return missing.join(', ')
}

const isLoading = ref(false)
const error = ref<string | null>(null)
const isSendingTest = ref(false)

const canSendTestEmail = computed(() => {
  return form.actionConfig.toEmail && 
         form.actionConfig.subject && 
         form.actionConfig.body
})

const sendTestEmail = async () => {
  if (!canSendTestEmail.value) return
  
  isSendingTest.value = true
  error.value = null
  
  try {
    const testEmailData = {
      to: form.actionConfig.toEmail,
      subject: form.actionConfig.subject,
      body: form.actionConfig.body
    }
    
    const response = await fetch('http://localhost:8080/test/email', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testEmailData)
    })
    
    const result = await response.json()
    
    if (response.ok) {
      alert('✅ Test email sent successfully!')
    } else {
      throw new Error(result.error || 'Failed to send test email')
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to send test email'
    console.error('Error sending test email:', err)
    alert('❌ Failed to send test email: ' + (err instanceof Error ? err.message : 'Unknown error'))
  } finally {
    isSendingTest.value = false
  }
}

const createArea = async () => {
  if (!isFormValid.value) return
  
  isLoading.value = true
  error.value = null
  
  try {
    const areaData = {
      name: form.areaName,
      description: form.description,
      triggerService: form.triggerService!,
      triggerType: form.triggerService === 'Google Calendar' ? 'Event' : 'Webhook',
      actionService: form.actionService!,
      actionType: form.actionService === 'Gmail' ? 'SendEmail' : 'Action',
      triggerConfig: form.triggerConfig,
      actionConfig: form.actionConfig
    }
    
    await areaService.createArea(areaData)
    emit('save')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to create area'
    console.error('Error creating area:', err)
  } finally {
    isLoading.value = false
  }
}

const emit = defineEmits<{ (e: 'close'): void; (e: 'save'): void }>()
</script>

<style scoped>

.create-area-container {
  position: relative;
  padding: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  width: 100%;
  background: var(--gradient-bg-primary);
  border-radius: 24px;
  overflow: hidden;
}


.main-card {
  background: rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
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
  border-bottom: 1px solid var(--color-border-primary);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.close-button {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--color-border-primary);
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
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
  letter-spacing: -0.02em;
}

.card-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
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
  color: var(--color-accent-primary) !important;
}

.label-text {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
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
  color: var(--color-text-secondary);
  letter-spacing: 0.01em;
}

.modern-input,
.modern-textarea {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  padding: 1rem;
  color: var(--color-text-primary);
  font-size: 1rem;
  font-weight: 400;
  transition: all 0.2s ease;
  outline: none;
}

.modern-input:focus,
.modern-textarea:focus {
  border-color: var(--color-accent-primary);
  background: rgba(255, 255, 255, 0.08);
  box-shadow: 0 0 0 3px var(--color-focus-ring);
}

.modern-input::placeholder,
.modern-textarea::placeholder {
  color: var(--color-text-secondary);
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
  border: 1px solid var(--color-border-primary);
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
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
  letter-spacing: -0.01em;
}

.selector-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
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
  color: var(--color-text-primary);
}

.placeholder-text {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
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
  color: var(--color-text-primary) !important;
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
  background: linear-gradient(90deg, transparent, var(--color-accent-primary), transparent);
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
  border-top: 1px solid var(--color-border-primary);
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
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border-primary);
}

.action-btn.secondary:hover {
  background: rgba(255, 255, 255, 0.08);
  color: var(--color-text-primary);
}

.action-btn.primary {
  background: var(--gradient-accent);
  color: white;
  box-shadow: var(--shadow-glow);
}

.action-btn.primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-glow);
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

/* Service Selection Styles */
.service-selection {
  width: 100%;
}

.service-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-top: 16px;
}

.service-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px 12px;
  background: rgba(26, 31, 46, 0.6);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(10px);
  text-align: center;
}

.service-card:hover {
  background: rgba(26, 31, 46, 0.8);
  border-color: var(--color-border-focus);
  transform: translateY(-2px);
  box-shadow: var(--shadow-lg);
}

.service-card.more-services {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
}

.service-card.more-services:hover {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.5);
}

.service-card-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 8px;
  border-radius: var(--radius-md);
  background: rgba(255, 255, 255, 0.05);
}

.service-card-icon .service-icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.service-card-name {
  color: var(--color-text-primary);
  font-size: 12px;
  font-weight: 500;
  line-height: 1.2;
}

.selected-service-display {
  margin-top: 16px;
}

.selected-service-card {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  background: rgba(26, 31, 46, 0.8);
  border: 2px solid var(--color-accent-primary);
  border-radius: var(--radius-xl);
  backdrop-filter: blur(20px);
  box-shadow: var(--shadow-glow);
}

.service-info {
  flex: 1;
  margin-left: 12px;
}

.service-name {
  color: var(--color-text-primary);
  font-size: 16px;
  font-weight: 600;
  display: block;
}

.service-type {
  color: var(--color-text-secondary);
  font-size: 12px;
  font-weight: 400;
  display: block;
  margin-top: 2px;
}

.change-service-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: var(--transition-normal);
}

.change-service-btn:hover {
  background: rgba(255, 255, 255, 0.2);
  color: var(--color-text-primary);
}

@media (max-width: 768px) {
  .service-grid {
    grid-template-columns: repeat(3, 1fr);
    gap: 8px;
  }

  .service-card {
    padding: 12px 8px;
  }

  .service-card-icon {
    width: 32px;
    height: 32px;
    margin-bottom: 6px;
  }

  .service-card-icon .service-icon {
    width: 20px;
    height: 20px;
  }

  .service-card-name {
    font-size: 11px;
  }
}

@media (max-width: 480px) {
  .service-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

:deep(.v-field__outline) {
  --v-field-border-opacity: 0.1;
}

:deep(.v-field--variant-outlined .v-field__outline) {
  color: var(--color-border-primary);
}

:deep(.v-field--focused .v-field__outline) {
  color: var(--color-accent-primary);
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 12px;
  color: #ef4444;
  font-size: 0.875rem;
  margin: 1rem 0;
}

:deep(.v-list) {
  background: var(--color-bg-card) !important;
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
}

:deep(.v-list-item) {
  color: var(--color-text-primary) !important;
}

:deep(.v-list-item:hover) {
  background: var(--color-hover-bg) !important;
}

:deep(.v-field__input) {
  color: var(--color-text-primary) !important;
}

/* Configuration Section Styles */
.config-section {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--color-border-primary);
  border-radius: 16px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
  transition: all 0.2s ease;
}

.config-section:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(59, 130, 246, 0.3);
}

.config-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.config-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  overflow: hidden;
  background: rgba(255, 255, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.config-icon .service-icon {
  width: 24px;
  height: 24px;
  object-fit: contain;
}

.config-info {
  flex: 1;
}

.config-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
  letter-spacing: -0.01em;
}

.config-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.config-content {
  padding-left: 3.5rem;
}

:deep(.v-field__outline) {
  color: var(--color-border-primary) !important;
}

.input-hint {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-top: 0.25rem;
  display: block;
  opacity: 0.8;
  font-style: italic;
}

.preview-section {
  margin-top: 2rem;
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  border: 1px solid var(--color-border-primary);
}

.preview-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.preview-icon {
  color: var(--color-accent-primary);
}

.preview-title {
  font-weight: 600;
  color: var(--color-text-primary);
}

.email-preview {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  padding: 1rem;
  border: 1px solid var(--color-border-primary);
}

.email-header {
  margin-bottom: 0.5rem;
  color: var(--color-text-primary);
  font-size: 0.875rem;
}

.email-body {
  color: var(--color-text-secondary);
  line-height: 1.5;
  white-space: pre-wrap;
}

:deep(.v-field--focused .v-field__outline) {
  color: var(--color-accent-primary) !important;
}

.action-btn.test-email-btn {
  background: linear-gradient(135deg, #10b981, #059669);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-left: 0.5rem;
}

.action-btn.test-email-btn:hover:not(:disabled) {
  background: linear-gradient(135deg, #059669, #047857);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.4);
}

.action-btn.test-email-btn:disabled {
  background: #6b7280;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

:deep(.v-select .v-field__input) {
  color: var(--color-text-primary) !important;
}

:deep(.v-select .v-field__outline) {
  color: var(--color-border-primary) !important;
}

</style>
