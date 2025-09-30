<template>
  <div class="card-col" @click="handleClick">
    <v-sheet :class="`area-card ${area.gradientClass}`" rounded="xl">
      <div class="area-icons-container">
        <div class="service-icon trigger-icon">
          <img 
            v-if="area.triggerIconUrl" 
            :src="area.triggerIconUrl" 
            :alt="area.triggerService"
            class="service-logo"
          />
          <v-icon v-else :size="32" color="white">{{ getTriggerIcon(area.triggerService) }}</v-icon>
          <div v-if="!area.triggerIconUrl" class="icon-fallback">{{ getTriggerEmoji(area.triggerService) }}</div>
        </div>
        <div class="service-arrow">
          <v-icon size="20" color="white">mdi-arrow-right</v-icon>
          <div class="arrow-fallback">→</div>
        </div>
        <div class="service-icon action-icon">
          <img 
            v-if="area.actionIconUrl" 
            :src="area.actionIconUrl" 
            :alt="area.actionService"
            class="service-logo"
          />
          <v-icon v-else :size="32" color="white">{{ getActionIcon(area.actionService) }}</v-icon>
          <div v-if="!area.actionIconUrl" class="icon-fallback">{{ getActionEmoji(area.actionService) }}</div>
        </div>
      </div>
    </v-sheet>
    <div class="card-title">{{ area.title }}</div>
    <div class="card-subtitle">{{ area.subtitle }}</div>
    <div class="card-description">{{ area.description }}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

interface AreaTemplate {
  id: string
  title: string
  subtitle: string
  description: string
  icon: string
  gradientClass: string
  triggerService: string
  actionService: string
  triggerIconUrl?: string
  actionIconUrl?: string
  isActive: boolean
}

const props = defineProps<{
  area: AreaTemplate
}>()

const emit = defineEmits<{
  click: [area: AreaTemplate]
}>()

const handleClick = () => {
  emit('click', props.area)
}

const getTriggerIcon = (service: string) => {
  switch (service) {
    case "Google Calendar":
      return "mdi-calendar"
    case "GitHub":
      return "mdi-github"
    case "Gmail":
      return "mdi-email-outline"
    case "Discord":
      return "mdi-discord"
    case "Slack":
      return "mdi-slack"
    case "Weather":
      return "mdi-weather-partly-cloudy"
    case "Instagram":
      return "mdi-instagram"
    case "Twitter":
      return "mdi-twitter"
    case "YouTube":
      return "mdi-youtube"
    case "Spotify":
      return "mdi-music"
    case "Telegram":
      return "mdi-telegram"
    case "Twitch":
      return "mdi-twitch"
    case "Dropbox":
      return "mdi-dropbox"
    case "Notion":
      return "mdi-notebook"
    default:
      return "mdi-cog"
  }
}

const getActionIcon = (service: string) => {
  switch (service) {
    case "Gmail":
      return "mdi-email-outline"
    case "Discord":
      return "mdi-discord"
    case "Slack":
      return "mdi-slack"
    case "GitHub":
      return "mdi-github"
    case "Telegram":
      return "mdi-telegram"
    case "Twitter":
      return "mdi-twitter"
    case "Instagram":
      return "mdi-instagram"
    case "Dropbox":
      return "mdi-dropbox"
    case "Notion":
      return "mdi-notebook"
    default:
      return "mdi-cog"
  }
}

const getTriggerEmoji = (service: string) => {
  switch (service) {
    case "Google Calendar":
      return "📅"
    case "GitHub":
      return "🐙"
    case "Gmail":
      return "📧"
    case "Discord":
      return "💬"
    case "Slack":
      return "💼"
    case "Weather":
      return "🌤️"
    case "Instagram":
      return "📸"
    case "Twitter":
      return "🐦"
    case "YouTube":
      return "📺"
    case "Spotify":
      return "🎵"
    case "Telegram":
      return "✈️"
    case "Twitch":
      return "🎮"
    case "Dropbox":
      return "📁"
    case "Notion":
      return "📝"
    default:
      return "⚙️"
  }
}

const getActionEmoji = (service: string) => {
  switch (service) {
    case "Gmail":
      return "📧"
    case "Discord":
      return "💬"
    case "Slack":
      return "💼"
    case "GitHub":
      return "🐙"
    case "Telegram":
      return "✈️"
    case "Twitter":
      return "🐦"
    case "Instagram":
      return "📸"
    case "Dropbox":
      return "📁"
    case "Notion":
      return "📝"
    default:
      return "⚙️"
  }
}
</script>

<style scoped>
.card-col {
  max-width: 320px;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.card-col:hover {
  transform: translateY(-2px);
}

.area-card {
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 24px;
  box-shadow: 0 6px 16px rgba(0,0,0,0.25);
  transform: translateY(0) scale(1);
  background-size: 130% 130%;
  background-position: 50% 50%;
  transition:
    transform .25s ease,
    box-shadow .25s ease,
    background-position .6s ease,
    filter .25s ease;
}

.area-icons-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}

.service-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

.trigger-icon {
  background: rgba(255, 255, 255, 0.2);
}

.action-icon {
  background: rgba(255, 255, 255, 0.2);
}

.service-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-fallback {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 24px;
  color: white;
  z-index: 10;
}

.arrow-fallback {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  font-size: 20px;
  color: white;
  z-index: 10;
}

.service-logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
  filter: brightness(0) invert(1);
}

.area-card :deep(.v-icon) {
  transition: transform .25s ease, opacity .25s ease;
}

.area-card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 12px 28px rgba(0,0,0,0.35);
  background-position: 80% 20%;
}

.area-card:hover :deep(.v-icon) {
  transform: translateY(-2px) scale(1.06);
}

.area-card:active {
  transform: translateY(-2px) scale(0.99);
}

.card-title {
  margin-top: 12px;
  font-weight: 800;
  font-size: 20px;
  color: white;
}

.card-subtitle {
  color: rgba(255,255,255,0.85);
  font-weight: 700;
  font-size: 14px;
  margin-top: 4px;
}

.card-description {
  color: rgba(255,255,255,0.7);
  font-size: 14px;
  margin-top: 8px;
  line-height: 1.4;
}
</style>
