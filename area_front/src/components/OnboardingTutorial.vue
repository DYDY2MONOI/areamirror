<template>
  <v-dialog :model-value="isOpen" @update:model-value="handleClose" max-width="900" persistent>
    <v-card class="onboarding-card">
      <div v-if="!startedTutorial" class="welcome-modal">
      <div class="welcome-content">
        <div class="welcome-icon"></div>
        <h1 class="welcome-title">Welcome to AREA!</h1>
        <p class="welcome-subtitle">Connect your favorite services with intelligent automation</p>
        
        <div class="welcome-question">
          <p>Is this your first time on AREA?</p>
          <div class="welcome-buttons">
            <button class="btn-primary" @click="startTutorial">
              <v-icon>mdi-school</v-icon>
              Yes, show me how it works
            </button>
            <button class="btn-secondary" @click="skipTutorial">
              <v-icon>mdi-close</v-icon>
              No, I already know
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="tutorial-modal">
      <div class="tutorial-header">
        <div class="tutorial-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: progressWidth }"></div>
          </div>
          <span class="progress-text">{{ currentSlide + 1 }} / {{ slides.length }}</span>
        </div>
        <button class="close-btn" @click="skipTutorial">
          <v-icon>mdi-close</v-icon>
        </button>
      </div>

      <div class="tutorial-content">
        <transition :name="slideDirection" mode="out-in">
          <div :key="currentSlide" class="slide">
            <div class="slide-image">
              <img :src="getImageUrl(slides[currentSlide].image)" :alt="slides[currentSlide].title" />
            </div>
            <div class="slide-info">
              <div class="slide-number">Step {{ currentSlide + 1 }}</div>
              <h2 class="slide-title">{{ slides[currentSlide].title }}</h2>
              <p class="slide-description">{{ slides[currentSlide].description }}</p>
              <ul v-if="slides[currentSlide].points" class="slide-points">
                <li v-for="(point, index) in slides[currentSlide].points" :key="index">
                  {{ point }}
                </li>
              </ul>
            </div>
          </div>
        </transition>
      </div>

      <div class="tutorial-footer">
        <button 
          class="nav-btn prev" 
          @click="prevSlide" 
          :disabled="currentSlide === 0"
        >
          <v-icon>mdi-chevron-left</v-icon>
          Previous
        </button>
        
        <div class="dots-indicator">
          <span 
            v-for="(slide, index) in slides" 
            :key="index"
            class="dot"
            :class="{ active: index === currentSlide }"
            @click="goToSlide(index)"
          ></span>
        </div>
        
        <button 
          class="nav-btn next" 
          @click="nextSlide"
        >
          {{ currentSlide === slides.length - 1 ? 'Finish' : 'Next' }}
          <v-icon v-if="currentSlide !== slides.length - 1">mdi-chevron-right</v-icon>
          <v-icon v-else>mdi-check</v-icon>
        </button>
      </div>
    </div>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

const startedTutorial = ref(false)
const currentSlide = ref(0)
const slideDirection = ref('slide-left')

watch(() => props.isOpen, (newValue) => {
  console.log(' OnboardingTutorial isOpen changed:', newValue)
  if (newValue) {
    startedTutorial.value = false
    currentSlide.value = 0
  }
})

const slides = [
  {
    image: '1.png',
    title: 'Create Your First Area',
    description: 'Start by clicking the + button to create a new automation',
    points: [
      'The + button is located at the top right of the homepage',
      'You can create as many areas as you want',
      'Each area connects a trigger to an action'
    ]
  },
  {
    image: '2.png',
    title: 'Discover Areas',
    description: 'Areas are automations that link your favorite services',
    points: [
      'An area = a trigger + an action (reaction)',
      'Example: "When I receive an email → send a Discord message"',
      'Check out examples for inspiration'
    ]
  },
  {
    image: '3.png',
    title: 'Choose Your Services',
    description: 'Select a trigger service and an action service',
    points: [
      'Left: the service that triggers the action (WHEN)',
      'Right: the service that executes the action (THEN)',
      'Dozens of services available: Gmail, Discord, Telegram, GitHub...'
    ]
  },
  {
    image: '4.png',
    title: 'Configure the Settings',
    description: 'Once services are selected, configure the details of your automation',
    points: [
      'Each service has its own parameters',
      'Fields are displayed automatically based on your choices',
      'All fields are validated in real time'
    ]
  },
  {
    image: '5.png',
    title: 'Configuration Example',
    description: 'Here\'s a complete example: Telegram to Discord',
    points: [
      'Chat ID: identifier of your Telegram conversation',
      'Webhook URL: Discord link to receive messages',
      'Use variables like {{messageText}} to customize'
    ]
  },
  {
    image: '6.png',
    title: 'Create Your Area',
    description: 'When everything is configured, click "Create Area"',
    points: [
      'Make sure all fields are filled in',
      'Give a descriptive name to your area',
      'The button becomes active when everything is ready'
    ]
  },
  {
    image: '7.png',
    title: 'Area Created Successfully!',
    description: 'Your area appears on the homepage and runs automatically',
    points: [
      'Activate/deactivate your areas at any time',
      'Modify or delete them if needed',
      'Check the execution history to track your automations'
    ]
  }
]

const progressWidth = computed(() => {
  return `${((currentSlide.value + 1) / slides.length) * 100}%`
})

const getImageUrl = (imageName: string) => {
  return new URL(`../assets/${imageName}`, import.meta.url).href
}

const startTutorial = () => {
  startedTutorial.value = true
}

const skipTutorial = () => {
  emit('close')
}

const nextSlide = () => {
  if (currentSlide.value < slides.length - 1) {
    slideDirection.value = 'slide-left'
    currentSlide.value++
  } else {
    skipTutorial()
  }
}

const prevSlide = () => {
  if (currentSlide.value > 0) {
    slideDirection.value = 'slide-right'
    currentSlide.value--
  }
}

const goToSlide = (index: number) => {
  slideDirection.value = index > currentSlide.value ? 'slide-left' : 'slide-right'
  currentSlide.value = index
}

const handleClose = (value: boolean) => {
  if (!value) {
    skipTutorial()
  }
}
</script>

<style scoped>
.onboarding-card {
  background: transparent !important;
  box-shadow: none !important;
  border-radius: 24px !important;
  overflow: hidden;
}

.welcome-modal {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  border-radius: 24px;
  padding: 3rem;
  text-align: center;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.welcome-content {
  max-width: 600px;
  margin: 0 auto;
}

.welcome-icon {
  font-size: 5rem;
  margin-bottom: 1.5rem;
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

.welcome-title {
  font-size: 2.5rem;
  font-weight: 800;
  color: white;
  margin-bottom: 1rem;
}

.welcome-subtitle {
  font-size: 1.25rem;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 3rem;
}

.welcome-question {
  background: rgba(0, 0, 0, 0.3);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 2rem;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.welcome-question p {
  font-size: 1.25rem;
  color: white;
  margin-bottom: 1.5rem;
}

.welcome-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.btn-primary, .btn-secondary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 1rem 2rem;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.btn-primary {
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  color: white;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(59, 130, 246, 0.4);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.05);
  color: white;
  border: 2px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.3);
}

.tutorial-modal {
  background: #0d1117;
  border-radius: 24px;
  overflow: hidden;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.tutorial-header {
  padding: 1.5rem 2rem;
  background: rgba(255, 255, 255, 0.05);
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.tutorial-progress {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #3b82f6 0%, #8b5cf6 100%);
  transition: width 0.3s ease;
  border-radius: 4px;
}

.progress-text {
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.875rem;
  font-weight: 600;
  min-width: 60px;
  text-align: right;
}

.close-btn {
  background: none;
  border: none;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 8px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.tutorial-content {
  flex: 1;
  overflow-y: auto;
  padding: 2rem;
}

.slide {
  display: flex;
  gap: 2rem;
  align-items: flex-start;
}

.slide-image {
  flex: 1.5;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.slide-image img {
  width: 100%;
  height: auto;
  display: block;
}

.slide-info {
  flex: 1;
  color: white;
}

.slide-number {
  color: #3b82f6;
  font-size: 0.875rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  margin-bottom: 0.5rem;
}

.slide-title {
  font-size: 1.75rem;
  font-weight: 700;
  margin-bottom: 1rem;
  color: white;
}

.slide-description {
  font-size: 1.125rem;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin-bottom: 1.5rem;
}

.slide-points {
  list-style: none;
  padding: 0;
}

.slide-points li {
  padding: 0.75rem 0;
  padding-left: 2rem;
  position: relative;
  color: rgba(255, 255, 255, 0.7);
  font-size: 0.95rem;
  line-height: 1.5;
}

.slide-points li::before {
  content: '';
  position: absolute;
  left: 0;
  color: #3b82f6;
  font-weight: bold;
  font-size: 1.2rem;
}

.tutorial-footer {
  padding: 1.5rem 2rem;
  background: rgba(255, 255, 255, 0.05);
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.nav-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  color: white;
}

.nav-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(59, 130, 246, 0.4);
}

.nav-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.nav-btn.prev {
  background: rgba(255, 255, 255, 0.1);
}

.dots-indicator {
  display: flex;
  gap: 0.5rem;
}

.dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  cursor: pointer;
  transition: all 0.3s ease;
}

.dot.active {
  background: #3b82f6;
  transform: scale(1.3);
}

.dot:hover {
  background: rgba(59, 130, 246, 0.7);
}

.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.3s ease;
}

.slide-left-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.slide-left-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-right-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-right-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

@media (max-width: 768px) {
  .slide {
    flex-direction: column;
  }
  
  .welcome-modal {
    padding: 2rem 1rem;
  }
  
  .welcome-title {
    font-size: 2rem;
  }
  
  .welcome-buttons {
    flex-direction: column;
  }
  
  .btn-primary, .btn-secondary {
    width: 100%;
  }
}
</style>






