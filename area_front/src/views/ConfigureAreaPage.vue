<template>
  <div class="configure-area-page">
    <div class="page-header">
      <div class="header-content">
        <button class="back-button" @click="goBack">
          <v-icon size="20">mdi-arrow-left</v-icon>
          Back to Home
        </button>
        <div class="header-text">
          <h1 class="page-title">{{ isEditingExisting ? 'Edit Your Area' : 'Configure Your Area' }}</h1>
          <p class="page-subtitle">{{ isEditingExisting ? 'Modify your existing automation area' : 'Set up your automation with the selected template' }}</p>
        </div>
      </div>
    </div>

    <div class="page-content">
      <div class="template-section">
        <div class="section-header">
          <div class="section-icon">
            <v-icon :size="32" color="white">{{ getTriggerIcon(template?.triggerService || '') }}</v-icon>
          </div>
          <div class="section-info">
            <h2 class="section-title">{{ template?.title }}</h2>
            <p class="section-subtitle">{{ template?.subtitle }}</p>
            <p class="section-description">{{ template?.description }}</p>
          </div>
        </div>

        <div class="workflow-preview">
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

      <div class="configuration-section">
        <div class="section-header">
          <h3 class="section-title">Configuration</h3>
          <p class="section-subtitle">Configure your trigger and action settings</p>
          <p class="debug-info">Template: {{ template?.title }} | Trigger: {{ template?.triggerService }} | Action: {{ template?.actionService }}</p>
          <p class="debug-info">Form Data: {{ JSON.stringify(form) }}</p>
        </div>

        <div v-if="template && template.triggerService === 'Google Calendar'" class="config-card">
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
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">📅 Event Date</label>
                <div class="date-picker-container">
                  <input
                    v-model="form.triggerConfig.eventDate"
                    type="date"
                    class="form-input date-input"
                    :min="getTodayDate()"
                    required
                  />
                  <div class="date-presets">
                    <button
                      type="button"
                      class="preset-btn"
                      @click="setDatePreset('today')"
                    >
                      Today
                    </button>
                    <button
                      type="button"
                      class="preset-btn"
                      @click="setDatePreset('tomorrow')"
                    >
                      Tomorrow
                    </button>
                    <button
                      type="button"
                      class="preset-btn"
                      @click="setDatePreset('nextweek')"
                    >
                      Next Week
                    </button>
                  </div>
                </div>
                <small class="form-hint">Select the date for your event</small>
                <small class="debug-info">Current date: {{ form.triggerConfig.eventDate }}</small>
              </div>

              <div class="form-group">
                <label class="form-label">⏰ Event Time</label>
                <div class="time-picker-container">
                  <input
                    v-model="form.triggerConfig.eventTime"
                    type="time"
                    class="form-input time-input"
                    required
                  />
                  <div class="time-presets">
                    <button
                      type="button"
                      class="preset-btn"
                      @click="setTimePreset('09:00')"
                    >
                      9:00 AM
                    </button>
                    <button
                      type="button"
                      class="preset-btn"
                      @click="setTimePreset('12:00')"
                    >
                      12:00 PM
                    </button>
                    <button
                      type="button"
                      class="preset-btn"
                      @click="setTimePreset('15:00')"
                    >
                      3:00 PM
                    </button>
                    <button
                      type="button"
                      class="preset-btn"
                      @click="setTimePreset('18:00')"
                    >
                      6:00 PM
                    </button>
                  </div>
                </div>
                <small class="form-hint">Select the time when the reminder should be sent</small>
                <small class="debug-info">Current time: {{ form.triggerConfig.eventTime }}</small>
              </div>

              <div class="form-group">
                <label class="form-label">📝 Event Title Filter</label>
                <input
                  v-model="form.triggerConfig.eventTitle"
                  type="text"
                  class="form-input"
                  placeholder="e.g., Meeting, Appointment"
                  required
                />
                <small class="form-hint">Only trigger for events containing this text (leave empty for all events)</small>
              </div>

              <div class="form-group">
                <label class="form-label">📅 Calendar ID</label>
                <input
                  v-model="form.triggerConfig.calendarId"
                  type="text"
                  class="form-input"
                  placeholder="primary"
                  required
                />
                <small class="form-hint">Which calendar to monitor (usually 'primary')</small>
              </div>
            </div>
          </div>
        </div>

        <div v-if="template && template.triggerService === 'GitHub'" class="config-card">
          <div class="config-header">
            <div class="config-icon">
              <v-icon size="24" color="white">mdi-github</v-icon>
            </div>
            <div class="config-info">
              <h4 class="config-title">🐙 GitHub Repository Trigger</h4>
              <p class="config-subtitle">Configure which GitHub events should trigger this area</p>
            </div>
          </div>

          <div class="config-content">
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">📁 Repository</label>
                <input
                  v-model="form.triggerConfig.repository"
                  type="text"
                  class="form-input"
                  placeholder="owner/repository-name"
                  required
                />
                <small class="form-hint">Enter the repository in format: owner/repository-name (e.g., microsoft/vscode)</small>
              </div>

              <div class="form-group">
                <label class="form-label">🌿 Branch</label>
                <input
                  v-model="form.triggerConfig.branch"
                  type="text"
                  class="form-input"
                  placeholder="main"
                />
                <small class="form-hint">Specific branch to monitor (leave empty for all branches)</small>
              </div>

              <div class="form-group full-width">
                <label class="form-label">📋 Event Types</label>
                <div class="checkbox-group">
                  <label class="checkbox-item">
                    <input
                      v-model="form.triggerConfig.events"
                      type="checkbox"
                      value="push"
                    />
                    <span class="checkbox-label">Push Events</span>
                    <small class="checkbox-hint">Triggered when code is pushed to the repository</small>
                  </label>
                  <label class="checkbox-item">
                    <input
                      v-model="form.triggerConfig.events"
                      type="checkbox"
                      value="pull_request"
                    />
                    <span class="checkbox-label">Pull Request Events</span>
                    <small class="checkbox-hint">Triggered when PRs are opened, closed, or merged</small>
                  </label>
                  <label class="checkbox-item">
                    <input
                      v-model="form.triggerConfig.events"
                      type="checkbox"
                      value="issues"
                    />
                    <span class="checkbox-label">Issue Events</span>
                    <small class="checkbox-hint">Triggered when issues are opened, closed, or commented on</small>
                  </label>
                  <label class="checkbox-item">
                    <input
                      v-model="form.triggerConfig.events"
                      type="checkbox"
                      value="release"
                    />
                    <span class="checkbox-label">Release Events</span>
                    <small class="checkbox-hint">Triggered when new releases are published</small>
                  </label>
                </div>
                <small class="form-hint">Select which GitHub events should trigger this area</small>
              </div>

              <div class="form-group">
                <label class="form-label">🔐 Webhook Secret</label>
                <input
                  v-model="form.triggerConfig.webhookSecret"
                  type="password"
                  class="form-input"
                  placeholder="your-webhook-secret"
                />
                <small class="form-hint">Secret token for webhook verification (optional but recommended)</small>
              </div>
            </div>
          </div>
        </div>

        <div v-if="template && template.actionService === 'Gmail'" class="config-card">
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
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">📧 Send Email To</label>
                <input
                  v-model="form.actionConfig.toEmail"
                  type="email"
                  class="form-input"
                  placeholder="your-email@gmail.com"
                  required
                />
                <small class="form-hint">Enter the email address where you want to receive the reminder</small>
                <small class="debug-info">Current value: {{ form.actionConfig.toEmail }}</small>
              </div>

              <div class="form-group">
                <label class="form-label">📝 Email Subject</label>
                <input
                  v-model="form.actionConfig.subject"
                  class="form-input"
                  placeholder="Reminder: &#123;&#123;eventTitle&#125;&#125;"
                  required
                />
                <small class="form-hint">Use &#123;&#123;eventTitle&#125;&#125; to include the event name</small>
              </div>

              <div class="form-group full-width">
                <label class="form-label">📄 Email Body</label>
                <textarea
                  v-model="form.actionConfig.body"
                  class="form-textarea"
                  placeholder="Hello! This is a reminder about your upcoming event: &#123;&#123;eventTitle&#125;&#125; at &#123;&#123;eventTime&#125;&#125;."
                  rows="4"
                  required
                ></textarea>
                <small class="form-hint">Use &#123;&#123;eventTitle&#125;&#125;, &#123;&#123;eventTime&#125;&#125;, &#123;&#123;areaName&#125;&#125; as placeholders</small>
              </div>
            </div>
          </div>
        </div>

        <div v-if="template && template.actionService === 'Discord'" class="config-card">
          <div class="config-header">
            <div class="config-icon">
              <v-icon size="24" color="white">mdi-discord</v-icon>
            </div>
            <div class="config-info">
              <h4 class="config-title">💬 Discord Message</h4>
              <p class="config-subtitle">Configure the channel and content for your Discord notification</p>
            </div>
          </div>

          <div class="config-content">
            <div class="form-grid">
              <div class="form-group">
                <label class="form-label">🔗 Webhook URL</label>
                <input
                  v-model="form.actionConfig.webhookUrl"
                  type="url"
                  class="form-input"
                  placeholder="https://discord.com/api/webhooks/..."
                  required
                />
                <small class="form-hint">Collez ici l'URL du webhook Discord généré pour ce salon.</small>
              </div>

              <div class="form-group full-width">
                <label class="form-label">💬 Message Content</label>
                <textarea
                  v-model="form.actionConfig.message"
                  class="form-textarea"
                  placeholder="Reminder: &#123;&#123;eventTitle&#125;&#125; starts at &#123;&#123;eventTime&#125;&#125;. Area: &#123;&#123;areaName&#125;&#125;"
                  rows="4"
                  required
                ></textarea>
                <small class="form-hint">Use &#123;&#123;eventTitle&#125;&#125;, &#123;&#123;eventTime&#125;&#125;, &#123;&#123;areaName&#125;&#125; as placeholders</small>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="error" class="error-message">
        <v-icon size="16" color="#ef4444">mdi-alert-circle</v-icon>
        <span>{{ error }}</span>
      </div>

      <div class="action-buttons">
        <button class="btn btn-secondary" @click="goBack">
          <v-icon size="18">mdi-close</v-icon>
          Cancel
        </button>
        <button class="btn btn-primary" @click="createArea" :disabled="!isFormValid || isLoading">
          <v-icon size="18">{{ isEditingExisting ? 'mdi-content-save' : 'mdi-check' }}</v-icon>
          {{ isLoading ? 'Saving...' : (isEditingExisting ? 'Update Area' : 'Create Area') }}
        </button>

        <div v-if="isEditingExisting" class="edit-mode-banner">
          <v-icon size="16" color="#10b981">mdi-pencil</v-icon>
          <span>Editing existing area - Changes will be saved to your current configuration</span>
        </div>

        <div class="test-email-section" v-if="template?.actionService === 'Gmail'">
          <div class="test-email-info">
            <h4>📧 Test Email Configuration</h4>
            <p>Send a test email to verify your configuration works correctly.</p>
            <div class="email-preview">
              <strong>To:</strong> {{ form.actionConfig.toEmail || 'Enter email address' }}<br>
              <strong>Subject:</strong> {{ form.actionConfig.subject || 'Enter subject' }}<br>
              <strong>Body:</strong> {{ form.actionConfig.body || 'Enter message body' }}
            </div>
          </div>
          <button
            class="btn btn-test"
            @click="sendTestEmail"
            :disabled="!canSendTestEmail || isSendingTest"
          >
            <v-icon size="18">mdi-email-send</v-icon>
            {{ isSendingTest ? 'Sending...' : 'Send Test Email' }}
          </button>
          <div v-if="error" class="error-message">
            ❌ {{ error }}
          </div>
        </div>

        <div class="test-email-section" v-if="template?.actionService === 'Discord'">
          <div class="test-email-info">
            <h4>💬 Test Discord Message</h4>
            <p>Send a test message to verify your Discord configuration works correctly.</p>
            <div class="email-preview">
              <strong>Webhook URL:</strong> {{ form.actionConfig.webhookUrl || 'Enter webhook URL' }}<br>
              <strong>Message:</strong> {{ form.actionConfig.message || 'Enter message content' }}
            </div>
          </div>
          <button
            class="btn btn-test"
            @click="sendTestDiscord"
            :disabled="!canSendDiscordTest || isSendingDiscordTest"
          >
            <v-icon size="18">mdi-send</v-icon>
            {{ isSendingDiscordTest ? 'Sending...' : 'Send Test Message' }}
          </button>
          <div v-if="discordTestError" class="error-message">
            ❌ {{ discordTestError }}
          </div>
        </div>

        <div class="test-trigger-section" v-if="template?.triggerService === 'Google Calendar'">
          <div class="test-trigger-info">
            <h4>🕐 Test Calendar Trigger</h4>
            <p>Test if your calendar trigger is working correctly. This will simulate the trigger firing.</p>
            <div class="trigger-preview">
              <strong>Event Date:</strong> {{ form.triggerConfig.eventDate || 'Select date' }}<br>
              <strong>Event Time:</strong> {{ form.triggerConfig.eventTime || 'Select time' }}<br>
              <strong>Combined:</strong> {{ getCombinedDateTime() }}
            </div>
          </div>
          <button
            class="btn btn-test-trigger"
            @click="testTrigger"
            :disabled="!canTestTrigger || isTestingTrigger"
          >
            <v-icon size="18">mdi-calendar-clock</v-icon>
            {{ isTestingTrigger ? 'Testing...' : 'Test Trigger' }}
          </button>
          <div v-if="triggerError" class="error-message">
            ❌ {{ triggerError }}
          </div>
        </div>

        <div class="test-trigger-section" v-if="template?.triggerService === 'GitHub'">
          <div class="test-trigger-info">
            <h4>🐙 Test GitHub Trigger</h4>
            <p>Test if your GitHub trigger is working correctly. This will simulate a webhook event.</p>
            <div class="trigger-preview">
              <strong>Repository:</strong> {{ form.triggerConfig.repository || 'Enter repository' }}<br>
              <strong>Branch:</strong> {{ form.triggerConfig.branch || 'All branches' }}<br>
              <strong>Events:</strong> {{ form.triggerConfig.events?.join(', ') || 'Select events' }}
            </div>
          </div>
          <button
            class="btn btn-test-trigger"
            @click="testGitHubTrigger"
            :disabled="!canTestGitHubTrigger || isTestingTrigger"
          >
            <v-icon size="18">mdi-github</v-icon>
            {{ isTestingTrigger ? 'Testing...' : 'Test GitHub Trigger' }}
          </button>
          <div v-if="triggerError" class="error-message">
            ❌ {{ triggerError }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { areaService, type Area } from '@/services/area'
import { useAuth } from '@/composables/useAuth'

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

const router = useRouter()
const route = useRoute()
const { currentUser } = useAuth()

const template = ref<AreaTemplate | null>(null)
const existingArea = ref<Area | null>(null)
const isEditingExisting = ref(false)
const form = reactive({
  triggerConfig: {
    eventDate: '',
    eventTime: '',
    eventTitle: '',
    calendarId: 'primary',
    repository: '',
    branch: 'main'
  } as any,
  actionConfig: {
    toEmail: '',
    subject: '',
    body: '',
    webhookUrl: '',
    message: ''
  } as any,
})

const isLoading = ref(false)
const isSendingTest = ref(false)
const isSendingDiscordTest = ref(false)
const isTestingTrigger = ref(false)
const error = ref<string | null>(null)
const triggerError = ref<string | null>(null)
const discordTestError = ref<string | null>(null)

watch(() => template.value, (newTemplate) => {
  if (newTemplate && !isEditingExisting.value) {
    console.log('Initializing form for new template:', newTemplate)
    console.log('Trigger service:', newTemplate.triggerService)
    console.log('Action service:', newTemplate.actionService)

    if (newTemplate.triggerService === 'Google Calendar') {
      form.triggerConfig = {
        eventDate: '',
        eventTime: '',
        eventTitle: '',
        calendarId: 'primary'
      }
    } else if (newTemplate.triggerService === 'GitHub') {
      form.triggerConfig = {
        repository: '',
        branch: '',
        events: [],
        webhookSecret: ''
      }
    } else if (newTemplate.triggerService === 'Weather') {
      form.triggerConfig = {
        city: '',
        temperature: 30,
        condition: ''
      }
    } else {
      form.triggerConfig = {}
    }

    if (newTemplate.actionService === 'Gmail') {
      form.actionConfig = {
        toEmail: '',
        subject: 'Reminder: {{eventTitle}}',
        body: 'Hello! This is a reminder about your upcoming event: {{eventTitle}} at {{eventTime}}.\n\nArea: {{areaName}}'
      }
      discordTestError.value = null
    } else if (newTemplate.actionService === 'Discord') {
      form.actionConfig = {
        webhookUrl: '',
        message: 'Reminder: {{eventTitle}} starts at {{eventTime}}. Area: {{areaName}}'
      }
      discordTestError.value = null
    } else {
      form.actionConfig = {}
      discordTestError.value = null
    }

    console.log('Form initialized for new template:', form)
  }
}, { immediate: true })

watch(() => existingArea.value, (newArea) => {
  if (newArea && isEditingExisting.value) {
    console.log('Loading existing area configuration:', newArea)

    if (newArea.triggerConfig) {
      const triggerConfig = { ...form.triggerConfig, ...newArea.triggerConfig }

      if (triggerConfig.eventTime && typeof triggerConfig.eventTime === 'string' && triggerConfig.eventTime.includes('T')) {
        const dateTime = new Date(triggerConfig.eventTime)
        triggerConfig.eventDate = dateTime.toISOString().split('T')[0]
        triggerConfig.eventTime = dateTime.toTimeString().split(' ')[0].substring(0, 5)
        console.log('Parsed datetime:', { original: newArea.triggerConfig.eventTime, date: triggerConfig.eventDate, time: triggerConfig.eventTime })
      }

      form.triggerConfig = triggerConfig
      console.log('Trigger config loaded:', form.triggerConfig)
    }

    if (newArea.actionConfig) {
      form.actionConfig = { ...form.actionConfig, ...newArea.actionConfig }
      console.log('Action config loaded:', form.actionConfig)
    }

    console.log('Final form state:', { triggerConfig: form.triggerConfig, actionConfig: form.actionConfig })
  }
}, { immediate: true })

const isFormValid = computed(() => {
  if (!template.value) return false

  if (template.value.triggerService === 'Google Calendar' && template.value.actionService === 'Gmail') {
    return form.triggerConfig.eventDate &&
           form.triggerConfig.eventTime &&
           form.actionConfig.toEmail &&
           form.actionConfig.subject
  }

  if (template.value.triggerService === 'GitHub') {
    return form.triggerConfig.repository &&
           form.triggerConfig.events &&
           form.triggerConfig.events.length > 0
  }

  if (template.value.actionService === 'Discord') {
    const webhookUrl = (form.actionConfig.webhookUrl || form.actionConfig.webhookURL || '').trim()
    const message = (form.actionConfig.message || '').trim()
    return !!webhookUrl && !!message
  }

  return true
})

const canSendTestEmail = computed(() => {
  return form.actionConfig.toEmail &&
         form.actionConfig.subject &&
         form.actionConfig.body
})

const canSendDiscordTest = computed(() => {
  const webhookUrl = (form.actionConfig.webhookUrl || form.actionConfig.webhookURL || '').trim()
  const message = (form.actionConfig.message || '').trim()
  return !!webhookUrl && !!message
})

const canTestTrigger = computed(() => {
  return form.triggerConfig.eventDate &&
         form.triggerConfig.eventTime
})

const canTestGitHubTrigger = computed(() => {
  return form.triggerConfig.repository &&
         form.triggerConfig.events &&
         form.triggerConfig.events.length > 0
})

const getCombinedDateTime = () => {
  if (form.triggerConfig.eventDate && form.triggerConfig.eventTime) {
    const eventDateTime = new Date(`${form.triggerConfig.eventDate}T${form.triggerConfig.eventTime}:00`)
    return eventDateTime.toISOString()
  }
  return 'Not set'
}

onMounted(async () => {
  console.log('=== CONFIGURE AREA PAGE DEBUG ===')
  console.log('Route query:', route.query)
  console.log('Route params:', route.params)
  console.log('Current route:', route.path)
  console.log('Full URL:', window.location.href)

  if (route.query.areaId) {
    try {
      console.log('Loading existing area with ID:', route.query.areaId)
      console.log('Calling areaService.getAreaById...')

      const token = localStorage.getItem('authToken')
      console.log('Auth token exists:', !!token)
      console.log('Auth token preview:', token ? token.substring(0, 20) + '...' : 'null')

      existingArea.value = await areaService.getAreaById(route.query.areaId as string)
      console.log('Area loaded successfully:', existingArea.value)
      isEditingExisting.value = true

      template.value = {
        id: existingArea.value.id,
        title: existingArea.value.name,
        subtitle: `${existingArea.value.triggerService} → ${existingArea.value.actionService}`,
        description: existingArea.value.description,
        icon: '',
        gradientClass: '',
        triggerService: existingArea.value.triggerService,
        actionService: existingArea.value.actionService,
        isActive: existingArea.value.isActive
      }

      console.log('Template created:', template.value)
    } catch (error) {
      console.error('=== ERROR LOADING EXISTING AREA ===')
      console.error('Error details:', error)
      console.error('Error message:', error instanceof Error ? error.message : 'Unknown error')
      console.error('Error stack:', error instanceof Error ? error.stack : 'No stack trace')
      console.error('NOT redirecting - staying on page to debug')
    }
  } else if (route.query.template) {
    try {
      template.value = JSON.parse(route.query.template as string)
      console.log('Template loaded:', template.value)
      console.log('Template trigger service:', template.value?.triggerService)
      console.log('Template action service:', template.value?.actionService)
    } catch (error) {
      console.error('Error parsing template:', error)
      router.push('/')
    }
  } else {
    console.log('No template or areaId found in route, redirecting to home')
    router.push('/')
  }
})

const goBack = () => {
  router.push('/')
}

const setTimePreset = (time: string) => {
  form.triggerConfig.eventTime = time
}

const setDatePreset = (preset: string) => {
  const today = new Date()
  let targetDate = new Date()

  switch (preset) {
    case 'today':
      targetDate = new Date(today)
      break
    case 'tomorrow':
      targetDate = new Date(today)
      targetDate.setDate(today.getDate() + 1)
      break
    case 'nextweek':
      targetDate = new Date(today)
      targetDate.setDate(today.getDate() + 7)
      break
  }

  form.triggerConfig.eventDate = targetDate.toISOString().split('T')[0]
}

const getTodayDate = () => {
  return new Date().toISOString().split('T')[0]
}

const resolveActionType = (service: string) => {
  switch (service) {
    case 'Gmail':
      return 'SendEmail'
    case 'Discord':
      return 'SendDiscordMessage'
    default:
      return 'Action'
  }
}

const sendTestEmail = async () => {
  if (!canSendTestEmail.value) {
    console.log('Cannot send test email - form not valid')
    return
  }

  console.log('Sending test email with data:', {
    to: form.actionConfig.toEmail,
    subject: form.actionConfig.subject,
    body: form.actionConfig.body
  })

  isSendingTest.value = true
  error.value = null

  try {
    const testEmailData = {
      to: form.actionConfig.toEmail,
      subject: form.actionConfig.subject,
      body: form.actionConfig.body
    }

    console.log('Making request to backend...')
    const response = await fetch('http://localhost:8080/test/email', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(testEmailData)
    })

    console.log('Response status:', response.status)
    const result = await response.json()
    console.log('Response data:', result)

    if (response.ok) {
      alert('✅ Test email sent successfully! Check your inbox.')
      error.value = null
    } else {
      throw new Error(result.error || `Server error: ${response.status}`)
    }
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Failed to send test email'
    error.value = errorMessage
    console.error('Error sending test email:', err)
    alert('❌ Failed to send test email: ' + errorMessage)
  } finally {
    isSendingTest.value = false
  }
}

const sendTestDiscord = async () => {
  if (!canSendDiscordTest.value) {
    console.log('Cannot send test Discord message - form not valid')
    return
  }

  const webhookUrl = (form.actionConfig.webhookUrl || form.actionConfig.webhookURL || '').trim()
  const message = (form.actionConfig.message || '').trim() || `Test message from ${template.value?.title || 'AREAmirror'}`

  console.log('Sending test Discord message with data:', {
    webhookUrl,
    message
  })

  isSendingDiscordTest.value = true
  discordTestError.value = null

  try {
    const response = await fetch('http://localhost:8080/test/discord', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        webhookUrl,
        message
      })
    })

    console.log('Discord test response status:', response.status)
    const result = await response.json()
    console.log('Discord test response data:', result)

    if (response.ok) {
      alert('✅ Test Discord message sent successfully! Check your channel.')
      discordTestError.value = null
    } else {
      throw new Error(result.error || `Server error: ${response.status}`)
    }
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Failed to send test Discord message'
    discordTestError.value = errorMessage
    console.error('Error sending test Discord message:', err)
    alert('❌ Failed to send test Discord message: ' + errorMessage)
  } finally {
    isSendingDiscordTest.value = false
  }
}

const testTrigger = async () => {
  if (!canTestTrigger.value) {
    console.log('Cannot test trigger - date and time not set')
    return
  }

  console.log('Testing trigger with data:', {
    eventDate: form.triggerConfig.eventDate,
    eventTime: form.triggerConfig.eventTime,
    combined: getCombinedDateTime()
  })

  isTestingTrigger.value = true
  triggerError.value = null

  try {
    const areaData = {
      name: `Test Area - ${template.value?.title || 'Unknown'}`,
      description: 'Temporary test area',
      triggerService: template.value?.triggerService || 'Google Calendar',
      triggerType: 'Event',
      actionService: template.value?.actionService || 'Gmail',
      actionType: resolveActionType(template.value?.actionService || 'Gmail'),
      triggerConfig: {
        ...form.triggerConfig,
        eventTime: getCombinedDateTime()
      },
      actionConfig: form.actionConfig
    }

    console.log('Creating test area...')
    const createdArea = await areaService.createArea(areaData)
    console.log('Test area created:', createdArea)

    console.log('Testing scheduler for area ID:', createdArea.id)
    const response = await fetch(`http://localhost:8080/test/scheduler/${createdArea.id}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      }
    })

    const result = await response.json()
    console.log('Scheduler test response:', result)

    if (response.ok) {
      alert('✅ Trigger test successful! Check your email inbox.')
      triggerError.value = null

      try {
        await fetch(`http://localhost:8080/areas/${createdArea.id}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          }
        })
        console.log('Test area cleaned up successfully')
      } catch (cleanupErr) {
        console.warn('Failed to clean up test area:', cleanupErr)
      }
    } else {
      throw new Error(result.error || `Server error: ${response.status}`)
    }
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Failed to test trigger'
    triggerError.value = errorMessage
    console.error('Error testing trigger:', err)
    alert('❌ Trigger test failed: ' + errorMessage)
  } finally {
    isTestingTrigger.value = false
  }
}

const testGitHubTrigger = async () => {
  if (!canTestGitHubTrigger.value) {
    console.log('Cannot test GitHub trigger - repository and events not set')
    return
  }

  console.log('Testing GitHub trigger with data:', {
    repository: form.triggerConfig.repository,
    branch: form.triggerConfig.branch,
    events: form.triggerConfig.events
  })

  isTestingTrigger.value = true
  triggerError.value = null

  try {
    const areaData = {
      name: `Test GitHub Area - ${template.value?.title || 'Unknown'}`,
      description: 'Temporary test area for GitHub trigger',
      triggerService: template.value?.triggerService || 'GitHub',
      triggerType: 'Webhook',
      actionService: template.value?.actionService || 'Gmail',
      actionType: resolveActionType(template.value?.actionService || 'Gmail'),
      triggerConfig: {
        ...form.triggerConfig,
        events: Array.isArray(form.triggerConfig.events) ? form.triggerConfig.events : [form.triggerConfig.events]
      },
      actionConfig: form.actionConfig
    }

    console.log('Creating test GitHub area...')
    const createdArea = await areaService.createArea(areaData)
    console.log('Test GitHub area created:', createdArea)

    console.log('Testing GitHub webhook for area ID:', createdArea.id)
    const response = await fetch(`http://localhost:8080/test/github/${createdArea.id}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        repository: form.triggerConfig.repository,
        branch: form.triggerConfig.branch,
        event: form.triggerConfig.events[0] || 'push'
      })
    })

    const result = await response.json()
    console.log('GitHub webhook test response:', result)

    if (response.ok) {
      alert('✅ GitHub trigger test successful! Check your configured action.')
      triggerError.value = null

      try {
        await fetch(`http://localhost:8080/areas/${createdArea.id}`, {
          method: 'DELETE',
          headers: {
            'Content-Type': 'application/json',
          }
        })
        console.log('Test GitHub area cleaned up successfully')
      } catch (cleanupErr) {
        console.warn('Failed to clean up test GitHub area:', cleanupErr)
      }
    } else {
      throw new Error(result.error || `Server error: ${response.status}`)
    }
  } catch (err) {
    const errorMessage = err instanceof Error ? err.message : 'Failed to test GitHub trigger'
    triggerError.value = errorMessage
    console.error('Error testing GitHub trigger:', err)
    alert('❌ GitHub trigger test failed: ' + errorMessage)
  } finally {
    isTestingTrigger.value = false
  }
}

const createArea = async () => {
  if (!isFormValid.value || !template.value) return

  isLoading.value = true
  error.value = null

  try {
    let triggerConfig = { ...form.triggerConfig }

    if (template.value.triggerService === 'Google Calendar' && form.triggerConfig.eventDate && form.triggerConfig.eventTime) {
      const eventDateTime = new Date(`${form.triggerConfig.eventDate}T${form.triggerConfig.eventTime}:00`)
      triggerConfig.eventTime = eventDateTime.toISOString()
      console.log('Combined event time:', triggerConfig.eventTime)
    }

    if (template.value.triggerService === 'GitHub') {
      if (typeof triggerConfig.events === 'string') {
        triggerConfig.events = [triggerConfig.events]
      }
      console.log('GitHub trigger config:', triggerConfig)
    }

    if (template.value.triggerService === 'Weather') {
      console.log('Weather trigger config:', triggerConfig)
    }

    const areaData = {
      name: template.value.title || 'Untitled Area',
      description: template.value.description || '',
      triggerService: template.value.triggerService || 'Unknown',
      triggerType: template.value.triggerService === 'Google Calendar' ? 'Event' : 'Webhook',
      actionService: template.value.actionService || 'Unknown',
      actionType: resolveActionType(template.value.actionService || 'Unknown'),
      triggerConfig: triggerConfig,
      actionConfig: form.actionConfig
    }

    if (isEditingExisting.value && existingArea.value) {
      console.log('Updating existing area with data:', areaData)
      try {
        const updatedArea = await areaService.updateArea(existingArea.value.id, areaData)
        console.log('Area updated successfully:', updatedArea)

        alert(`✅ Area "${updatedArea.name}" updated successfully!`)

        router.push('/')
        return
      } catch (updateError) {
        console.error('Update failed:', updateError)
        error.value = `Failed to update area: ${updateError instanceof Error ? updateError.message : 'Unknown error'}`

        alert(`❌ Failed to update area: ${updateError instanceof Error ? updateError.message : 'Unknown error'}`)
        return
      }
    } else {
      console.log('Creating new area with data:', areaData)
      const createdArea = await areaService.createArea(areaData)
      console.log('New area created successfully:', createdArea)

      alert(`✅ Area "${createdArea.name}" created successfully!`)

      router.push('/')
    }
  } catch (err) {
    error.value = err instanceof Error ? err.message : `Failed to ${isEditingExisting.value ? 'update' : 'create'} area`
    console.error(`Error ${isEditingExisting.value ? 'updating' : 'creating'} area:`, err)
  } finally {
    isLoading.value = false
  }
}

const getTriggerIcon = (service: string) => {
  switch (service) {
    case "Google Calendar": return "mdi-calendar"
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
</script>

<style scoped>
.configure-area-page {
  min-height: 100vh;
  background: var(--color-bg-primary);
  padding: 2rem;
}

.page-header {
  margin-bottom: 3rem;
}

.header-content {
  display: flex;
  align-items: center;
  gap: 2rem;
}

.back-button {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  color: var(--color-text-primary);
  text-decoration: none;
  transition: all 0.2s ease;
  cursor: pointer;
}

.back-button:hover {
  background: var(--color-bg-card-hover);
  border-color: var(--color-border-focus);
}

.header-text {
  flex: 1;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
  letter-spacing: -0.02em;
}

.page-subtitle {
  font-size: 1.125rem;
  color: var(--color-text-secondary);
  margin: 0;
}

.page-content {
  max-width: 1000px;
  margin: 0 auto;
}

.template-section {
  background: var(--color-bg-card);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
  border-radius: 20px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: var(--shadow-glow);
}

.section-header {
  display: flex;
  align-items: flex-start;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.section-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.section-info {
  flex: 1;
}

.section-title {
  font-size: 1.75rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0 0 0.5rem 0;
}

.section-subtitle {
  font-size: 1rem;
  color: var(--color-text-secondary);
  margin: 0 0 0.75rem 0;
}

.section-description {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  margin: 0;
  line-height: 1.6;
}

.workflow-preview {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1.5rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 16px;
  border: 1px solid var(--color-border-primary);
}

.workflow-step {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
}

.step-icon {
  width: 40px;
  height: 40px;
  border-radius: 10px;
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
  font-size: 1rem;
  color: var(--color-text-primary);
  font-weight: 500;
}

.workflow-arrow {
  color: var(--color-text-secondary);
}

.configuration-section {
  background: var(--color-bg-card);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
  border-radius: 20px;
  padding: 2rem;
  margin-bottom: 2rem;
  box-shadow: var(--shadow-glow);
}

.config-card {
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

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-group.full-width {
  grid-column: 1 / -1;
}

.form-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--color-text-primary);
}

.form-input,
.form-textarea {
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  padding: 0.875rem 1rem;
  color: var(--color-text-primary);
  font-size: 0.875rem;
  transition: all 0.2s ease;
  width: 100%;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--color-border-focus);
  box-shadow: var(--shadow-glow);
}

.form-textarea {
  resize: vertical;
  min-height: 100px;
}

.form-hint {
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

.action-buttons {
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
  justify-content: center;
}

.edit-mode-banner {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.3);
  border-radius: 8px;
  color: #10b981;
  font-size: 0.875rem;
  font-weight: 500;
  margin-top: 1rem;
  width: 100%;
  justify-content: center;
}

.btn {
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

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--color-bg-secondary);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border-primary);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--color-bg-card-hover);
  border-color: var(--color-border-focus);
}

.btn-primary {
  background: var(--gradient-accent);
  color: white;
  border: 1px solid transparent;
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-glow);
}

.btn-test {
  background: var(--gradient-green);
  color: white;
  border: 1px solid transparent;
}

.btn-test:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: var(--shadow-glow);
}

.debug-info {
  font-size: 0.75rem;
  color: #666;
  background: rgba(255, 255, 255, 0.1);
  padding: 0.5rem;
  border-radius: 8px;
  margin: 0.5rem 0;
  font-family: monospace;
}

.date-picker-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.date-input {
  width: 200px;
}

.time-picker-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.time-input {
  width: 200px;
}

.date-presets {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.time-presets {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.preset-btn {
  padding: 0.5rem 1rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 8px;
  color: var(--color-text-primary);
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s ease;
}

.preset-btn:hover {
  background: var(--color-bg-card-hover);
  border-color: var(--color-border-focus);
  transform: translateY(-1px);
}

.preset-btn:active {
  transform: translateY(0);
}

.test-email-section {
  margin-top: 2rem;
  padding: 1.5rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
}

.test-email-info h4 {
  margin: 0 0 0.5rem 0;
  color: var(--color-text-primary);
  font-size: 1.1rem;
}

.test-email-info p {
  margin: 0 0 1rem 0;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.email-preview {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: 8px;
  padding: 1rem;
  margin: 1rem 0;
  font-size: 0.9rem;
  color: var(--color-text-secondary);
}

.error-message {
  margin-top: 1rem;
  padding: 0.75rem;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.3);
  border-radius: 8px;
  color: #ef4444;
  font-size: 0.9rem;
}

.test-trigger-section {
  margin-top: 2rem;
  padding: 1.5rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
}

.test-trigger-info h4 {
  margin: 0 0 0.5rem 0;
  color: var(--color-text-primary);
  font-size: 1.1rem;
}

.test-trigger-info p {
  margin: 0 0 1rem 0;
  color: var(--color-text-secondary);
  font-size: 0.9rem;
}

.trigger-preview {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: 8px;
  padding: 1rem;
  margin: 1rem 0;
  font-size: 0.9rem;
  color: var(--color-text-secondary);
}

.btn-test-trigger {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: white;
  border: none;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-test-trigger:hover:not(:disabled) {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-1px);
  box-shadow: var(--shadow-glow);
}

.btn-test-trigger:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.checkbox-group {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 0.5rem;
}

.checkbox-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  padding: 1rem;
  background: var(--color-bg-secondary);
  border: 1px solid var(--color-border-primary);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.checkbox-item:hover {
  background: var(--color-bg-card-hover);
  border-color: var(--color-border-focus);
}

.checkbox-item input[type="checkbox"] {
  margin-right: 0.75rem;
  transform: scale(1.2);
}

.checkbox-label {
  font-weight: 500;
  color: var(--color-text-primary);
  font-size: 0.875rem;
}

.checkbox-hint {
  font-size: 0.75rem;
  color: var(--color-text-secondary);
  margin-left: 1.5rem;
  opacity: 0.8;
  font-style: italic;
}
</style>
