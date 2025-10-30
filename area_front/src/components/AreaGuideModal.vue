<template>
  <v-dialog :model-value="isOpen" @update:model-value="$emit('close')" max-width="900" persistent>
    <div class="guide-modal">
      <div class="guide-header">
        <div class="header-content">
          <div class="header-icon">
            <v-icon size="32" color="#3b82f6">mdi-information</v-icon>
          </div>
          <div class="header-text">
            <h2 class="guide-title">📚 Guide: Creating an Area</h2>
            <p class="guide-subtitle">Step {{ currentSlide + 1 }} of {{ slides.length }}</p>
          </div>
        </div>
        <button class="close-btn" @click="$emit('close')">
          <v-icon size="24">mdi-close</v-icon>
        </button>
      </div>

      <div class="guide-content">
        <div class="slides-container">
          <transition :name="slideDirection" mode="out-in">
            <div :key="currentSlide" class="slide">
              <div v-if="currentSlide === 0" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#8b5cf6">mdi-help-circle-outline</v-icon>
                </div>
                <h3 class="slide-title">What is an Area?</h3>
                <p class="slide-text">
                  An <strong>Area</strong> is an automation that connects two services together. 
                  When something happens in one service (the <span class="highlight trigger">Trigger</span>), 
                  it automatically performs an action in another service (the <span class="highlight action">Action</span>).
                </p>
                <div class="example-box">
                  <div class="example-title">
                    <v-icon size="18" color="#22c55e">mdi-lightbulb-on</v-icon>
                    Example
                  </div>
                  <p>"When a new event is added to my Date Timer (Trigger), send me an email reminder via Gmail (Action)"</p>
                </div>
              </div>

              <div v-else-if="currentSlide === 1" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#f59e0b">mdi-cog-outline</v-icon>
                </div>
                <h3 class="slide-title">How Does it Work?</h3>
                <div class="workflow-diagram">
                  <div class="workflow-step">
                    <div class="step-icon trigger">
                      <v-icon size="24" color="white">mdi-play-circle</v-icon>
                    </div>
                    <div class="step-content">
                      <h4>1. Trigger Service</h4>
                      <p>The event that starts everything</p>
                    </div>
                  </div>
                  <div class="workflow-arrow">
                    <v-icon size="32" color="#3b82f6">mdi-arrow-right</v-icon>
                  </div>
                  <div class="workflow-step">
                    <div class="step-icon action">
                      <v-icon size="24" color="white">mdi-lightning-bolt</v-icon>
                    </div>
                    <div class="step-content">
                      <h4>2. Action Service</h4>
                      <p>What happens automatically</p>
                    </div>
                  </div>
                </div>
                <div class="info-box">
                  <p>The scheduler checks every 30 seconds for triggers and executes the corresponding actions automatically!</p>
                </div>
              </div>

              <div v-else-if="currentSlide === 2" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#ec4899">mdi-tag-outline</v-icon>
                </div>
                <h3 class="slide-title">Step 1: Name Your Area</h3>
                <p class="slide-text">Choose a clear, descriptive name that explains what your automation does.</p>
                <div class="examples-list">
                  <div class="example-item">
                    <v-icon size="18" color="#22c55e">mdi-check</v-icon>
                    <span>"Calendar to Email Reminder"</span>
                  </div>
                  <div class="example-item">
                    <v-icon size="18" color="#22c55e">mdi-check</v-icon>
                    <span>"GitHub Push Notifications"</span>
                  </div>
                  <div class="example-item">
                    <v-icon size="18" color="#22c55e">mdi-check</v-icon>
                    <span>"Weather Alert to Discord"</span>
                  </div>
                  <div class="example-item">
                    <v-icon size="18" color="#22c55e">mdi-check</v-icon>
                    <span>"Daily Report via Telegram"</span>
                  </div>
                </div>
                <div class="tip-box">
                  <v-icon size="18" color="#f59e0b">mdi-lightbulb</v-icon>
                  <p>Use names that make it easy to find your areas later!</p>
                </div>
              </div>

              <div v-else-if="currentSlide === 3" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#ef4444">mdi-play-circle-outline</v-icon>
                </div>
                <h3 class="slide-title">Step 2: Select a Trigger</h3>
                <p class="slide-text">Choose the service that will start your automation:</p>
                <div class="services-grid">
                  <div class="service-card">
                    <v-icon size="24" color="#3b82f6">mdi-calendar</v-icon>
                    <h4>Date Timer</h4>
                    <p>When an event is coming up</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#8b5cf6">mdi-github</v-icon>
                    <h4>GitHub</h4>
                    <p>When code is pushed</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#f59e0b">mdi-weather-partly-cloudy</v-icon>
                    <h4>Weather</h4>
                    <p>When conditions change</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#10b981">mdi-table</v-icon>
                    <h4>Google Sheets</h4>
                    <p>When data is modified</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#06b6d4">mdi-clock-outline</v-icon>
                    <h4>Timer</h4>
                    <p>At regular intervals</p>
                  </div>
                </div>
              </div>

              <div v-else-if="currentSlide === 4" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#3b82f6">mdi-lightning-bolt-outline</v-icon>
                </div>
                <h3 class="slide-title">Step 3: Select an Action</h3>
                <p class="slide-text">Choose what should happen when the trigger fires:</p>
                <div class="services-grid">
                  <div class="service-card">
                    <v-icon size="24" color="#ea4335">mdi-gmail</v-icon>
                    <h4>Gmail</h4>
                    <p>Send an email notification</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#5865f2">mdi-discord</v-icon>
                    <h4>Discord</h4>
                    <p>Post to a Discord channel</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#0088cc">mdi-send</v-icon>
                    <h4>Telegram</h4>
                    <p>Send a Telegram message</p>
                  </div>
                </div>
              </div>

              <div v-else-if="currentSlide === 5" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#22c55e">mdi-cog-outline</v-icon>
                </div>
                <h3 class="slide-title">Step 4: Configure Details</h3>
                <p class="slide-text">Set up the specific settings for your trigger and action:</p>
                <div class="config-examples">
                  <div class="config-item">
                    <div class="config-header">
                      <v-icon size="20" color="#ef4444">mdi-play-circle</v-icon>
                      <strong>Trigger Configuration</strong>
                    </div>
                    <ul>
                      <li>Event times & dates</li>
                      <li>Repository selections</li>
                      <li>City & weather conditions</li>
                      <li>Time intervals (5m, 1h, 24h)</li>
                    </ul>
                  </div>
                  <div class="config-item">
                    <div class="config-header">
                      <v-icon size="20" color="#3b82f6">mdi-lightning-bolt</v-icon>
                      <strong>Action Configuration</strong>
                    </div>
                    <ul>
                      <li>Email addresses</li>
                      <li>Discord webhook URLs</li>
                      <li>Message content & formatting</li>
                      <li>Telegram chat IDs</li>
                    </ul>
                  </div>
                </div>
              </div>

              <div v-else-if="currentSlide === 6" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#8b5cf6">mdi-code-braces</v-icon>
                </div>
                <h3 class="slide-title">Step 5: Use Template Variables</h3>
                <p class="slide-text">Make your messages dynamic with these variables:</p>
                <div class="variables-showcase">
                  <div class="variable-item">
                    <code v-pre>{{areaName}}</code>
                    <span>Your area name</span>
                  </div>
                  <div class="variable-item">
                    <code v-pre>{{eventTime}}</code>
                    <span>Event date & time</span>
                  </div>
                  <div class="variable-item">
                    <code v-pre>{{triggerTime}}</code>
                    <span>When triggered</span>
                  </div>
                  <div class="variable-item">
                    <code v-pre>{{interval}}</code>
                    <span>Timer interval</span>
                  </div>
                </div>
                <div class="example-box">
                  <div class="example-title">
                    <v-icon size="18" color="#22c55e">mdi-lightbulb-on</v-icon>
                    Example Message
                  </div>
                  <p v-pre>⏰ Reminder for {{areaName}} at {{eventTime}}</p>
                </div>
              </div>

              <div v-else-if="currentSlide === 7" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#10b981">mdi-star-outline</v-icon>
                </div>
                <h3 class="slide-title">Popular Examples</h3>
                <p class="slide-text">Get inspired by these common automations:</p>
                <div class="popular-examples">
                  <div class="popular-item">
                    <div class="popular-header">
                      <v-icon size="20" color="#3b82f6">mdi-calendar</v-icon>
                      <span>Calendar Reminder</span>
                    </div>
                    <p>Date Timer → Gmail</p>
                  </div>
                  <div class="popular-item">
                    <div class="popular-header">
                      <v-icon size="20" color="#f59e0b">mdi-weather-partly-cloudy</v-icon>
                      <span>Weather Alert</span>
                    </div>
                    <p>Weather (>30°C) → Discord</p>
                  </div>
                  <div class="popular-item">
                    <div class="popular-header">
                      <v-icon size="20" color="#8b5cf6">mdi-github</v-icon>
                      <span>Code Push Notification</span>
                    </div>
                    <p>GitHub Push → Gmail</p>
                  </div>
                  <div class="popular-item">
                    <div class="popular-header">
                      <v-icon size="20" color="#06b6d4">mdi-clock-outline</v-icon>
                      <span>Daily Report</span>
                    </div>
                    <p>Timer (24h) → Telegram</p>
                  </div>
                </div>
              </div>

              <div v-else-if="currentSlide === 8" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#06b6d4">mdi-lightbulb-on-outline</v-icon>
                </div>
                <h3 class="slide-title">Tips & Best Practices</h3>
                <div class="tips-list-slides">
                  <div class="tip-item-slide">
                    <v-icon size="24" color="#22c55e">mdi-check-circle</v-icon>
                    <div>
                      <strong>Test First</strong>
                      <p>Always use test buttons before activating</p>
                    </div>
                  </div>
                  <div class="tip-item-slide">
                    <v-icon size="24" color="#22c55e">mdi-check-circle</v-icon>
                    <div>
                      <strong>Clear Names</strong>
                      <p>Use descriptive names to find areas easily</p>
                    </div>
                  </div>
                  <div class="tip-item-slide">
                    <v-icon size="24" color="#22c55e">mdi-check-circle</v-icon>
                    <div>
                      <strong>Start Simple</strong>
                      <p>Begin with basic areas, add complexity later</p>
                    </div>
                  </div>
                  <div class="tip-item-slide">
                    <v-icon size="24" color="#f59e0b">mdi-alert</v-icon>
                    <div>
                      <strong>Service Limits</strong>
                      <p>Some services have rate limits</p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </transition>
        </div>

        <div class="navigation-arrows">
          <button 
            class="nav-arrow left" 
            @click="prevSlide" 
            :disabled="currentSlide === 0"
            :class="{ disabled: currentSlide === 0 }"
          >
            <v-icon size="32">mdi-chevron-left</v-icon>
          </button>
          
          <div class="progress-dots">
            <div 
              v-for="(slide, index) in slides" 
              :key="index" 
              class="dot"
              :class="{ active: index === currentSlide }"
              @click="goToSlide(index)"
            ></div>
          </div>

          <button 
            class="nav-arrow right" 
            @click="nextSlide" 
            :disabled="currentSlide === slides.length - 1"
            :class="{ disabled: currentSlide === slides.length - 1 }"
          >
            <v-icon size="32">mdi-chevron-right</v-icon>
          </button>
        </div>
      </div>

      <div class="guide-footer">
        <button v-if="currentSlide < slides.length - 1" class="secondary-btn" @click="$emit('close')">
          Skip Tutorial
        </button>
        <button v-else class="primary-btn" @click="$emit('close')">
          <v-icon size="18">mdi-check</v-icon>
          Got it, let's create!
        </button>
      </div>
    </div>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  isOpen: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()

const currentSlide = ref(0)
const slideDirection = ref('slide-left')

const slides = [
  'What is an Area?',
  'How it Works',
  'Name Your Area',
  'Select Trigger',
  'Select Action',
  'Configure Details',
  'Template Variables',
  'Popular Examples',
  'Tips & Best Practices'
]

const nextSlide = () => {
  if (currentSlide.value < slides.length - 1) {
    slideDirection.value = 'slide-left'
    currentSlide.value++
  }
}

const prevSlide = () => {
  if (currentSlide.value > 0) {
    slideDirection.value = 'slide-right'
    currentSlide.value--
  }
}

const goToSlide = (index: number) => {
  if (index > currentSlide.value) {
    slideDirection.value = 'slide-left'
  } else {
    slideDirection.value = 'slide-right'
  }
  currentSlide.value = index
}
</script>

<style scoped>
.guide-modal {
  background: linear-gradient(135deg, #0f0e1e 0%, #1a1632 100%);
  border-radius: 24px;
  overflow: hidden;
  max-height: 85vh;
  display: flex;
  flex-direction: column;
}

.guide-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.02);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
}

.header-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: rgba(59, 130, 246, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-text {
  flex: 1;
}

.guide-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: white;
  margin: 0 0 0.25rem 0;
}

.guide-subtitle {
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.close-btn {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: white;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.2);
}

.guide-content {
  flex: 1;
  padding: 2rem;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.slides-container {
  flex: 1;
  position: relative;
  overflow: hidden;
}

.slide {
  position: absolute;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.slide-content {
  width: 100%;
  text-align: center;
  padding: 2rem;
}

.slide-icon {
  margin-bottom: 1.5rem;
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-10px); }
}

.slide-title {
  font-size: 2rem;
  font-weight: 700;
  color: white;
  margin: 0 0 1rem 0;
}

.slide-text {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin: 0 auto 2rem;
  max-width: 600px;
}

.slide-text strong {
  color: white;
  font-weight: 600;
}

.highlight {
  padding: 0.125rem 0.5rem;
  border-radius: 6px;
  font-weight: 600;
}

.highlight.trigger {
  background: rgba(239, 68, 68, 0.2);
  color: #fca5a5;
}

.highlight.action {
  background: rgba(59, 130, 246, 0.2);
  color: #93c5fd;
}

.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.4s ease;
}

.slide-left-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-left-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

.slide-right-enter-from {
  transform: translateX(-100%);
  opacity: 0;
}

.slide-right-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

.navigation-arrows {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 2rem;
  margin-top: 2rem;
  padding: 1rem 0;
}

.nav-arrow {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  color: white;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.nav-arrow:hover:not(.disabled) {
  transform: scale(1.1);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.nav-arrow.disabled {
  opacity: 0.3;
  cursor: not-allowed;
  background: rgba(255, 255, 255, 0.1);
  box-shadow: none;
}

.progress-dots {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  cursor: pointer;
  transition: all 0.3s ease;
}

.dot.active {
  width: 12px;
  height: 12px;
  background: linear-gradient(135deg, #667eea, #764ba2);
  box-shadow: 0 0 10px rgba(102, 126, 234, 0.5);
}

.dot:hover {
  background: rgba(255, 255, 255, 0.4);
}

.example-box,
.info-box,
.tip-box {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.2);
  border-radius: 12px;
  padding: 1rem;
  margin: 1.5rem auto 0;
  max-width: 500px;
}

.info-box {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.2);
}

.tip-box {
  background: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.2);
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.example-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
  color: #22c55e;
  margin-bottom: 0.5rem;
}

.example-box p,
.info-box p,
.tip-box p {
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
  line-height: 1.5;
}

.workflow-diagram {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 2rem;
  background: rgba(255, 255, 255, 0.03);
  border-radius: 16px;
  margin: 1.5rem auto;
  max-width: 600px;
}

.workflow-step {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.step-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.step-icon.trigger {
  background: linear-gradient(135deg, #ef4444, #dc2626);
}

.step-icon.action {
  background: linear-gradient(135deg, #3b82f6, #2563eb);
}

.step-content h4 {
  color: white;
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0 0 0.25rem 0;
}

.step-content p {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.875rem;
  margin: 0;
}

.workflow-arrow {
  color: rgba(59, 130, 246, 0.5);
  animation: pulse 2s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; }
  50% { opacity: 1; }
}

.examples-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  max-width: 400px;
  margin: 2rem auto;
}

.example-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 1rem;
}

.services-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 1rem;
  margin: 2rem auto 0;
  max-width: 600px;
}

.service-card {
  padding: 1.5rem 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  text-align: center;
  transition: all 0.3s ease;
  cursor: pointer;
}

.service-card:hover {
  transform: translateY(-4px);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}

.service-card h4 {
  color: white;
  font-size: 0.875rem;
  font-weight: 600;
  margin: 0.75rem 0 0.25rem;
}

.service-card p {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.75rem;
  margin: 0;
}

.config-examples {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin: 2rem auto 0;
  max-width: 600px;
}

.config-item {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1.5rem;
  text-align: left;
}

.config-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  color: white;
}

.config-item ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.config-item li {
  color: rgba(255, 255, 255, 0.7);
  padding: 0.5rem 0;
  padding-left: 1.5rem;
  position: relative;
}

.config-item li::before {
  content: "•";
  position: absolute;
  left: 0;
  color: #3b82f6;
  font-weight: bold;
}

.variables-showcase {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1rem;
  margin: 2rem auto;
  max-width: 600px;
}

.variable-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem;
  background: rgba(59, 130, 246, 0.1);
  border: 1px solid rgba(59, 130, 246, 0.2);
  border-radius: 12px;
}

.variable-item code {
  background: rgba(59, 130, 246, 0.2);
  color: #93c5fd;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  font-family: 'Courier New', monospace;
  font-size: 0.875rem;
}

.variable-item span {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.75rem;
}

.popular-examples {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin: 2rem auto 0;
  max-width: 600px;
}

.popular-item {
  padding: 1.5rem 1rem;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  transition: all 0.3s ease;
}

.popular-item:hover {
  transform: translateY(-2px);
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.2);
}

.popular-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
  color: white;
  font-weight: 600;
  font-size: 0.875rem;
}

.popular-item p {
  color: rgba(255, 255, 255, 0.6);
  font-size: 0.75rem;
  margin: 0;
}

.tips-list-slides {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin: 2rem auto 0;
  max-width: 500px;
}

.tip-item-slide {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.25rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  text-align: left;
}

.tip-item-slide strong {
  color: white;
  font-size: 1rem;
  display: block;
  margin-bottom: 0.25rem;
}

.tip-item-slide p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.875rem;
  margin: 0;
}

.guide-footer {
  padding: 1.5rem 2rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  background: rgba(255, 255, 255, 0.02);
  display: flex;
  justify-content: center;
  gap: 1rem;
}

.primary-btn,
.secondary-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.875rem 2rem;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  border: none;
}

.primary-btn {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
}

.primary-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
}

.secondary-btn {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.secondary-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  color: white;
}

@media (max-width: 768px) {
  .slide-content {
    padding: 1rem;
  }

  .slide-title {
    font-size: 1.5rem;
  }

  .slide-text {
    font-size: 1rem;
  }

  .services-grid,
  .config-examples,
  .variables-showcase {
    grid-template-columns: 1fr;
  }

  .workflow-diagram {
    flex-direction: column;
  }

  .navigation-arrows {
    gap: 1rem;
  }

  .nav-arrow {
    width: 40px;
    height: 40px;
  }
}
</style>
