<template>
  <div class="configure-area-container">
    <div class="main-card">
      <div class="card-header">
        <div class="header-content">
          <div class="header-text">
            <h1 class="card-title">Configure Your Area</h1>
            <p class="card-subtitle">Set up your automation with the selected template</p>
          </div>
          <button class="close-button" @click="$emit('close')">
            <v-icon size="24" color="white">mdi-close</v-icon>
          </button>
        </div>
      </div>
      <div class="card-content">
        <div class="form-section">
          <div class="section-label">
            <v-icon class="label-icon" size="20">mdi-information-outline</v-icon>
            <span class="label-text">Template Information</span>
          </div>

          <div class="template-info-card">
            <div class="template-header">
              <div class="template-icon">
                <v-icon :size="32" color="white">{{ getTriggerIcon(template?.triggerService || '') }}</v-icon>
              </div>
              <div class="template-details">
                <h3 class="template-title">{{ template?.title }}</h3>
                <p class="template-subtitle">{{ template?.subtitle }}</p>
                <p class="template-description">{{ template?.description }}</p>
              </div>
            </div>

            <div class="template-workflow">
              <div class="workflow-step">
              <div class="step-icon trigger-icon">
                <v-icon size="20" color="white">{{ getTriggerIcon(template?.triggerService || '') }}</v-icon>
              </div>
                <div class="step-content">
                  <div class="step-label">Trigger</div>
                  <div class="step-service">{{ template?.triggerService }}</div>
                </div>
              </div>
              <div class="workflow-arrow">
                <v-icon size="24" color="#9ca3af">mdi-arrow-right</v-icon>
              </div>
              <div class="workflow-step">
              <div class="step-icon action-icon">
                <v-icon size="20" color="white">{{ getActionIcon(template?.actionService || '') }}</v-icon>
              </div>
                <div class="step-content">
                  <div class="step-label">Action</div>
                  <div class="step-service">{{ template?.actionService }}</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="form-section">
          <div class="section-label">
            <v-icon class="label-icon" size="20">mdi-cog-outline</v-icon>
            <span class="label-text">Configuration</span>
          </div>

          <div v-if="template?.triggerService === 'Date Timer'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <v-icon size="24" color="white">mdi-calendar</v-icon>
              </div>
              <div class="config-info">
                <h4 class="config-title">📅 Calendar Event Trigger</h4>
                <p class="config-subtitle">Configure when this area should trigger</p>
              </div>
            </div>

            <div class="config-content">
              <div class="input-group">
                <div class="input-container">
                  <label class="input-label">⏰ Event Time</label>
                  <input
                    v-model="form.triggerConfig.eventTime"
                    type="text"
                    class="modern-input"
                    placeholder="e.g., 09:00, 2:30 PM, tomorrow"
                    required
                  />
                  <small class="input-hint">When should the reminder be sent?</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📝 Event Title Filter</label>
                  <input
                    v-model="form.triggerConfig.eventTitle"
                    type="text"
                    class="modern-input"
                    placeholder="e.g., Meeting, Appointment"
                    required
                  />
                  <small class="input-hint">Only trigger for events containing this text (leave empty for all events)</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📅 Calendar ID</label>
                  <input
                    v-model="form.triggerConfig.calendarId"
                    type="text"
                    class="modern-input"
                    placeholder="primary"
                    required
                  />
                  <small class="input-hint">Which calendar to monitor (usually 'primary')</small>
                </div>
              </div>
            </div>
          </div>

          <div v-if="template?.actionService === 'Gmail'" class="config-section">
            <div class="config-header">
              <div class="config-icon">
                <v-icon size="24" color="white">mdi-email</v-icon>
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
                    placeholder="Reminder: &#123;&#123;eventTitle&#125;&#125;"
                    required
                  />
                  <small class="input-hint">Use &#123;&#123;eventTitle&#125;&#125; to include the event name</small>
                </div>

                <div class="input-container">
                  <label class="input-label">📄 Email Body</label>
                  <textarea
                    v-model="form.actionConfig.body"
                    class="modern-textarea"
                    placeholder="Hello! This is a reminder about your upcoming event: &#123;&#123;eventTitle&#125;&#125; at &#123;&#123;eventTime&#125;&#125;."
                    rows="4"
                    required
                  ></textarea>
                  <small class="input-hint">Use &#123;&#123;eventTitle&#125;&#125;, &#123;&#123;eventTime&#125;&#125;, &#123;&#123;areaName&#125;&#125; as placeholders</small>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div v-if="error" class="error-message">
          <v-icon size="16" color="#ef4444">mdi-alert-circle</v-icon>
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

          <button
            v-if="template?.triggerService === 'Date Timer' && template?.actionService === 'Gmail'"
            class="action-btn test-email-btn"
            @click="sendTestEmail"
            :disabled="!canSendTestEmail || isSendingTest"
          >
            <v-icon size="18">mdi-email-send</v-icon>
            {{ isSendingTest ? 'Sending...' : 'Send Test Email' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { areaService } from '../../services/area'
import { useAuth } from '@/composables/useAuth'
import { API_BASE_URL } from '@/config/api'

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
  template: AreaTemplate | null
}>()

const { currentUser } = useAuth()

const form = reactive({
  triggerConfig: {} as any,
  actionConfig: {} as any,
})

watch(() => props.template, (newTemplate) => {
  if (newTemplate) {
    if (newTemplate.triggerService === 'Date Timer' && newTemplate.actionService === 'Gmail') {
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
    } else if (newTemplate.triggerService === 'Weather') {
      form.triggerConfig = {
        city: '',
        temperature: 30,
        condition: ''
      }
    }
  }
}, { immediate: true })

const isFormValid = computed(() => {
  if (!props.template) return false

  if (props.template.triggerService === 'Date Timer' && props.template.actionService === 'Gmail') {
    return form.triggerConfig.eventTime &&
           form.actionConfig.toEmail &&
           form.actionConfig.subject
  }

  if (props.template.triggerService === 'Weather') {
    return form.triggerConfig.city &&
           form.triggerConfig.temperature !== undefined
  }

  return true
})

const canSendTestEmail = computed(() => {
  return form.actionConfig.toEmail &&
         form.actionConfig.subject &&
         form.actionConfig.body
})

const isLoading = ref(false)
const isSendingTest = ref(false)
const error = ref<string | null>(null)

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

    const response = await fetch(`${API_BASE_URL}/test/email`, {
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
  if (!isFormValid.value || !props.template) return

  isLoading.value = true
  error.value = null

  try {
    const areaData = {
      name: props.template.title,
      description: props.template.description,
      triggerService: props.template.triggerService,
      triggerType: props.template.triggerService === 'Date Timer' ? 'Event' : 'Webhook',
      actionService: props.template.actionService,
      actionType: props.template.actionService === 'Gmail' ? 'SendEmail' : 'Action',
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

const getTriggerIcon = (service: string) => {
  switch (service) {
    case "Date Timer": return "mdi-calendar"
    case "GitHub": return "mdi-github"
    case "Gmail": return "mdi-email-outline"
    case "Discord": return "mdi-discord"
    case "Slack": return "mdi-slack"
    case "Weather": return "mdi-weather-partly-cloudy"
    default: return "mdi-cog"
  }
}

const getActionIcon = (service: string) => {
  switch (service) {
    case "Gmail": return "mdi-email"
    case "Slack": return "mdi-slack"
    case "Discord": return "mdi-discord"
    case "GitHub": return "mdi-github"
    default: return "mdi-cog"
  }
}

const emit = defineEmits<{ (e: 'close'): void; (e: 'save'): void }>()
</script>

<style scoped>
.configure-area-container {
  position: relative;
  width: 100%;
  max-width: 800px;
  margin: 0 auto;
}

.main-card {
  background: var(--color-bg-card);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
  border-radius: 24px;
  overflow: hidden;
  box-shadow: var(--shadow-glow);
}

.card-header {
  background: var(--gradient-accent);
  padding: 2rem;
  position: relative;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.header-text {
  flex: 1;
}

.card-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: white;
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.card-subtitle {
  font-size: 1rem;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  font-weight: 400;
}

.close-button {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  padding: 0.75rem;
  color: white;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-button:hover {
  background: rgba(255, 255, 255, 0.3);
  border-color: rgba(255, 255, 255, 0.5);
}

.card-content {
  padding: 2rem;
}

.form-section {
  margin-bottom: 2rem;
}

.section-label {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  font-weight: 600;
  color: var(--color-text-primary);
  font-size: 1.125rem;
}

.label-icon {
  color: var(--color-accent-primary);
}

.template-info-card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--color-border-primary);
  border-radius: 16px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.template-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.template-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.template-details {
  flex: 1;
}

.template-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
}

.template-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0 0 0.5rem 0;
}

.template-description {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
  line-height: 1.5;
}

.template-workflow {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
}

.workflow-step {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
}

.step-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.trigger-icon {
  background: var(--gradient-blue);
}

.action-icon {
  background: var(--gradient-green);
}

.step-content {
  flex: 1;
}

.step-label {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.step-service {
  font-size: 0.875rem;
  color: var(--color-text-primary);
  font-weight: 500;
}

.workflow-arrow {
  color: var(--color-text-secondary);
}

.config-section {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid var(--color-border-primary);
  border-radius: 16px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
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
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
}

.config-info {
  flex: 1;
}

.config-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.25rem 0;
}

.config-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.config-content {
  padding-left: 3.5rem;
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
  color: var(--color-text-primary);
}

.modern-input,
.modern-textarea {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  padding: 0.875rem 1rem;
  color: var(--color-text-primary);
  font-size: 0.875rem;
  transition: all 0.2s ease;
  width: 100%;
}

.modern-input:focus,
.modern-textarea:focus {
  outline: none;
  border-color: var(--color-border-focus);
  box-shadow: var(--shadow-glow);
}

.modern-textarea {
  resize: vertical;
  min-height: 80px;
}

.input-hint {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-top: 0.25rem;
  display: block;
  opacity: 0.8;
  font-style: italic;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 12px;
  padding: 1rem;
  color: #ef4444;
  font-size: 0.875rem;
  margin: 1rem 0;
}

.card-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
  margin-top: 2rem;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.875rem 1.5rem;
  border-radius: 12px;
  font-weight: 500;
  font-size: 0.875rem;
  transition: all 0.2s ease;
  cursor: pointer;
  border: none;
  text-decoration: none;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.secondary {
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border-primary);
}

.action-btn.secondary:hover:not(:disabled) {
  background: var(--color-bg-card-hover);
  border-color: var(--color-border-focus);
}

.action-btn.primary {
  background: var(--gradient-accent);
  color: white;
  border: 1px solid transparent;
}

.action-btn.primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-glow);
}

.action-btn.test-email-btn {
  background: var(--gradient-green);
  color: white;
  border: 1px solid transparent;
}

.action-btn.test-email-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-glow);
}
</style>
