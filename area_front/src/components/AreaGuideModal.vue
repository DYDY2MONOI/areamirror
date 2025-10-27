<template>
  <v-dialog v-model="dialogModel" max-width="900" persistent>
    <v-card class="guide-modal" flat>
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
          <transition :name="slideDirection" mode="out-in" appear>
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
                  <p>"When a new event is added to my Google Calendar (Trigger), send me an email reminder via Gmail (Action)"</p>
                </div>
              </div>

              <div v-else-if="currentSlide === 1" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#f59e0b">mdi-cog-outline</v-icon>
                </div>
                <h3 class="slide-title">Comment ça fonctionne ?</h3>
                
                <div class="explanation-section">
                  <div class="concept-card trigger-card">
                    <div class="concept-header">
                      <v-icon size="28" color="#3b82f6">mdi-play-circle</v-icon>
                      <h4>🎯 Trigger (Déclencheur)</h4>
                    </div>
                    <p class="concept-text">
                      Un <strong>trigger</strong> est l'<strong>événement qui démarre</strong> votre automatisation. 
                      C'est le "QUAND" de votre area.
                    </p>
                    <div class="concept-examples">
                      <div class="mini-example">📅 Quand un événement arrive</div>
                      <div class="mini-example">💬 Quand je reçois un message</div>
                      <div class="mini-example">⏰ Toutes les heures</div>
                    </div>
                  </div>
                  
                  <div class="arrow-divider">
                    <v-icon size="40" color="#3b82f6">mdi-arrow-right-bold</v-icon>
                  </div>
                  
                  <div class="concept-card action-card">
                    <div class="concept-header">
                      <v-icon size="28" color="#10b981">mdi-lightning-bolt</v-icon>
                      <h4>⚡ Action (Réaction)</h4>
                    </div>
                    <p class="concept-text">
                      Une <strong>action</strong> est ce qui se <strong>passe automatiquement</strong> quand le trigger se déclenche.
                      C'est le "ALORS" de votre area.
                    </p>
                    <div class="concept-examples">
                      <div class="mini-example">📧 Envoyer un email</div>
                      <div class="mini-example">💬 Poster sur Discord</div>
                      <div class="mini-example">📱 Notifier sur Telegram</div>
                    </div>
                  </div>
                </div>
                
                <div class="info-box">
                  <v-icon size="18" color="#3b82f6">mdi-information</v-icon>
                  <p>Le système vérifie automatiquement vos triggers toutes les 30 secondes et exécute les actions correspondantes !</p>
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
                <h3 class="slide-title">Étape 2 : Choisir un Trigger (Déclencheur)</h3>
                <p class="slide-text">
                  Le <strong>trigger</strong> est l'événement qui <strong>démarre</strong> votre automatisation. 
                  C'est la condition "QUAND" qui active votre area.
                </p>
                <div class="services-grid">
                  <div class="service-card">
                    <v-icon size="24" color="#3b82f6">mdi-calendar</v-icon>
                    <h4>Google Calendar</h4>
                    <p>Quand un événement approche</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#8b5cf6">mdi-github</v-icon>
                    <h4>GitHub</h4>
                    <p>Quand du code est poussé</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#f59e0b">mdi-weather-partly-cloudy</v-icon>
                    <h4>Weather</h4>
                    <p>Quand la météo change</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#10b981">mdi-table</v-icon>
                    <h4>Google Sheets</h4>
                    <p>Quand des données changent</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#06b6d4">mdi-clock-outline</v-icon>
                    <h4>Timer</h4>
                    <p>À intervalles réguliers</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#0088cc">mdi-send</v-icon>
                    <h4>Telegram</h4>
                    <p>Quand je reçois un message</p>
                  </div>
                </div>
                <div class="tip-box">
                  <v-icon size="18" color="#f59e0b">mdi-lightbulb</v-icon>
                  <p>💡 Chaque trigger a ses propres paramètres (ex: Chat ID pour Telegram, ville pour Weather)</p>
                </div>
              </div>

              <div v-else-if="currentSlide === 4" class="slide-content">
                <div class="slide-icon">
                  <v-icon size="64" color="#3b82f6">mdi-lightning-bolt-outline</v-icon>
                </div>
                <h3 class="slide-title">Étape 3 : Choisir une Action (Réaction)</h3>
                <p class="slide-text">
                  L'<strong>action</strong> est ce qui se <strong>passe automatiquement</strong> quand le trigger se déclenche.
                  C'est la partie "ALORS" qui exécute quelque chose.
                </p>
                <div class="services-grid">
                  <div class="service-card">
                    <v-icon size="24" color="#ea4335">mdi-gmail</v-icon>
                    <h4>Gmail</h4>
                    <p>Envoyer un email</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#5865f2">mdi-discord</v-icon>
                    <h4>Discord</h4>
                    <p>Poster un message</p>
                  </div>
                  <div class="service-card">
                    <v-icon size="24" color="#0088cc">mdi-send</v-icon>
                    <h4>Telegram</h4>
                    <p>Envoyer sur Telegram</p>
                  </div>
                </div>
                <div class="example-box">
                  <div class="example-title">
                    <v-icon size="18" color="#22c55e">mdi-lightbulb-on</v-icon>
                    Exemple Complet
                  </div>
                  <p>
                    <strong>QUAND</strong> je reçois un message Telegram (Trigger)<br/>
                    <strong>ALORS</strong> envoyer le message sur Discord (Action)
                  </p>
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
                <h3 class="slide-title">Les Variables Template avec <code v-pre>{{}}</code></h3>
                
                <div class="variables-explanation">
                  <p class="slide-text">
                    Les <strong>variables template</strong> sont des <strong>placeholders magiques</strong> qui sont 
                    <strong>automatiquement remplacés</strong> par des valeurs réelles quand votre area s'exécute.
                  </p>
                  
                  <div class="how-it-works-box">
                    <div class="example-row">
                      <div class="before">
                        <strong>❌ Vous écrivez :</strong>
                        <code v-pre>Message de {{firstName}}: {{messageText}}</code>
                      </div>
                      <v-icon size="24" color="#3b82f6">mdi-arrow-right</v-icon>
                      <div class="after">
                        <strong>✅ Le système envoie :</strong>
                        <code>Message de John: Bonjour !</code>
                      </div>
                    </div>
                  </div>
                  
                  <h4 class="variables-section-title">📋 Variables Disponibles :</h4>
                  
                  <div class="variables-categories">
                    <div class="variable-category">
                      <div class="category-title">
                        <v-icon size="20" color="#3b82f6">mdi-clock</v-icon>
                        Générales
                      </div>
                      <div class="variables-list">
                        <div class="var-item"><code v-pre>{{areaName}}</code> → Nom de votre area</div>
                        <div class="var-item"><code v-pre>{{eventTime}}</code> → Date et heure</div>
                        <div class="var-item"><code v-pre>{{triggerService}}</code> → Service trigger</div>
                      </div>
                    </div>
                    
                    <div class="variable-category">
                      <div class="category-title">
                        <v-icon size="20" color="#0088cc">mdi-send</v-icon>
                        Telegram
                      </div>
                      <div class="variables-list">
                        <div class="var-item"><code v-pre>{{messageText}}</code> → Texte du message</div>
                        <div class="var-item"><code v-pre>{{firstName}}</code> → Prénom de l'expéditeur</div>
                        <div class="var-item"><code v-pre>{{username}}</code> → @username</div>
                        <div class="var-item"><code v-pre>{{chatId}}</code> → ID du chat</div>
                      </div>
                    </div>
                    
                    <div class="variable-category">
                      <div class="category-title">
                        <v-icon size="20" color="#10b981">mdi-table</v-icon>
                        Google Sheets
                      </div>
                      <div class="variables-list">
                        <div class="var-item"><code v-pre>{{sheetName}}</code> → Nom de la feuille</div>
                        <div class="var-item"><code v-pre>{{rowNumber}}</code> → Numéro de ligne</div>
                        <div class="var-item"><code v-pre>{{rowData}}</code> → Données de la ligne</div>
                      </div>
                    </div>
                    
                    <div class="variable-category">
                      <div class="category-title">
                        <v-icon size="20" color="#06b6d4">mdi-timer</v-icon>
                        Timer
                      </div>
                      <div class="variables-list">
                        <div class="var-item"><code v-pre>{{triggerTime}}</code> → Heure du déclenchement</div>
                        <div class="var-item"><code v-pre>{{interval}}</code> → Intervalle configuré</div>
                      </div>
                    </div>
                  </div>
                </div>
                
                <div class="tip-box">
                  <v-icon size="18" color="#f59e0b">mdi-lightbulb</v-icon>
                  <p>💡 Utilisez autant de variables que vous voulez dans vos messages pour les rendre dynamiques et informatifs !</p>
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
                    <p>Google Calendar → Gmail</p>
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
        <button 
          class="nav-btn prev-btn" 
          @click="prevSlide" 
          :disabled="currentSlide === 0"
        >
          <v-icon size="20">mdi-chevron-left</v-icon>
          Précédent
        </button>
        
        <div class="footer-center">
          <div class="progress-dots">
            <span 
              v-for="(slide, index) in slides" 
              :key="index"
              class="progress-dot"
              :class="{ active: index === currentSlide }"
              @click="currentSlide = index"
            ></span>
          </div>
          <button v-if="currentSlide < slides.length - 1" class="secondary-btn" @click="$emit('close')">
            Passer le tutoriel
          </button>
        </div>
        
        <button 
          v-if="currentSlide < slides.length - 1"
          class="nav-btn next-btn" 
          @click="nextSlide"
        >
          Suivant
          <v-icon size="20">mdi-chevron-right</v-icon>
        </button>
        <button 
          v-else 
          class="nav-btn finish-btn" 
          @click="$emit('close')"
        >
          <v-icon size="18">mdi-check</v-icon>
          C'est parti !
        </button>
      </div>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const props = defineProps<{
  isOpen: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
}>()

// Pont v-model : quand le dialog veut se fermer, on émet 'close'
const dialogModel = computed({
  get: () => props.isOpen,
  set: (val: boolean) => { if (!val) emit('close') }
})

const currentSlide = ref(0)
const slideDirection = ref<'slide-left' | 'slide-right'>('slide-left')

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
  min-height: 500px;
}

.slide {
  position: absolute;
  inset: 0;
  width: 100%;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  overflow-y: auto;
}

.slide-content {
  width: 100%;
  text-align: center;
  padding: 2rem;
  max-width: 900px;
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

/* Slide LEFT */
.slide-left-enter-active,
.slide-left-leave-active {
  transition: transform 0.4s ease, opacity 0.4s ease;
}

.slide-left-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.slide-left-enter-to {
  transform: translateX(0);
  opacity: 1;
}

.slide-left-leave-from {
  transform: translateX(0);
  opacity: 1;
}

.slide-left-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

/* Slide RIGHT */
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 0.4s ease, opacity 0.4s ease;
}

.slide-right-enter-from {
  transform: translateX(-100%);
  opacity: 0;
}

.slide-right-enter-to {
  transform: translateX(0);
  opacity: 1;
}

.slide-right-leave-from {
  transform: translateX(0);
  opacity: 1;
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
  justify-content: space-between;
  align-items: center;
  gap: 1rem;
}

.footer-center {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.progress-dots {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.progress-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  cursor: pointer;
  transition: all 0.3s ease;
}

.progress-dot.active {
  background: #667eea;
  transform: scale(1.3);
}

.progress-dot:hover {
  background: rgba(102, 126, 234, 0.7);
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.nav-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(102, 126, 234, 0.4);
}

.nav-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
  transform: none;
}

.prev-btn {
  background: rgba(255, 255, 255, 0.1);
}

.finish-btn {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
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

/* Explanation Section pour Triggers et Actions */
.explanation-section {
  display: flex;
  gap: 1rem;
  align-items: center;
  margin: 1.5rem 0;
}

.concept-card {
  flex: 1;
  padding: 1.5rem;
  border-radius: 12px;
  border: 2px solid;
}

.trigger-card {
  background: rgba(59, 130, 246, 0.1);
  border-color: rgba(59, 130, 246, 0.3);
}

.action-card {
  background: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.3);
}

.concept-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.concept-header h4 {
  font-size: 1.1rem;
  font-weight: 700;
  color: white;
  margin: 0;
}

.concept-text {
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.6;
  margin-bottom: 1rem;
  font-size: 0.95rem;
}

.concept-examples {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.mini-example {
  padding: 0.5rem 0.75rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 6px;
  font-size: 0.875rem;
  color: rgba(255, 255, 255, 0.7);
}

.arrow-divider {
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Variables Template Section */
.variables-explanation {
  width: 100%;
}

.how-it-works-box {
  margin: 1.5rem 0;
  padding: 1.5rem;
  background: rgba(59, 130, 246, 0.1);
  border: 2px solid rgba(59, 130, 246, 0.3);
  border-radius: 12px;
}

.example-row {
  display: flex;
  align-items: center;
  gap: 1rem;
  justify-content: space-between;
}

.before, .after {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.before code, .after code {
  display: block;
  padding: 0.75rem;
  background: rgba(0, 0, 0, 0.3);
  border-radius: 6px;
  font-size: 0.875rem;
  color: #a3e635;
  font-family: 'Courier New', monospace;
}

.variables-section-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: white;
  margin: 1.5rem 0 1rem 0;
}

.variables-categories {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin: 1rem 0;
}

.variable-category {
  padding: 1rem;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 8px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.category-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 700;
  color: white;
  margin-bottom: 0.75rem;
  font-size: 0.95rem;
}

.variables-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.var-item {
  font-size: 0.85rem;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.5;
}

.var-item code {
  color: #a3e635;
  background: rgba(0, 0, 0, 0.3);
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
  font-size: 0.8rem;
}

@media (max-width: 768px) {
  .explanation-section {
    flex-direction: column;
  }
  
  .arrow-divider {
    transform: rotate(90deg);
  }
  
  .example-row {
    flex-direction: column;
  }
  
  .variables-categories {
    grid-template-columns: 1fr;
  }
}
</style>
