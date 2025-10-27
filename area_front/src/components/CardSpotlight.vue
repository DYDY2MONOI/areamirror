<template>
  <div
    ref="cardRef"
    class="card-spotlight-container"
    @click="handleClick"
    @mousemove="handleMouseMove"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
    :class="{ 'has-delete': showDeleteButton }"
  >
    <div class="card-spotlight" :style="spotlightStyle">
      <div class="spotlight"></div>
      <div class="card-content">
        <div class="service-icons">
          <div class="service-icon trigger-icon">
            <img
              v-if="area.triggerIconUrl"
              :src="getIconUrl(area.triggerIconUrl)"
              :alt="area.triggerService"
              class="service-logo"
            />
            <v-icon v-else :size="36" color="white">{{ getTriggerIcon(area.triggerService) }}</v-icon>
            <div v-if="!area.triggerIconUrl" class="icon-fallback">{{ getTriggerEmoji(area.triggerService) }}</div>
          </div>
          <div class="service-arrow">
            <v-icon size="20" color="white">mdi-arrow-right</v-icon>
          </div>
          <div class="service-icon action-icon">
            <img
              v-if="area.actionIconUrl"
              :src="getIconUrl(area.actionIconUrl)"
              :alt="area.actionService"
              class="service-logo"
            />
            <v-icon v-else :size="36" color="white">{{ getActionIcon(area.actionService) }}</v-icon>
            <div v-if="!area.actionIconUrl" class="icon-fallback">{{ getActionEmoji(area.actionService) }}</div>
          </div>
        </div>

        <div class="card-info">
          <h3 class="card-title">{{ 'title' in area ? area.title : area.name }}</h3>
          <p v-if="'subtitle' in area" class="card-subtitle">{{ area.subtitle }}</p>
          <p class="card-description">{{ area.description }}</p>
        </div>
      </div>

      <button
        v-if="showDeleteButton"
        class="delete-button"
        @click.stop="handleDelete"
        type="button"
      >
        <v-icon size="16" color="white">mdi-delete</v-icon>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref, computed } from 'vue'
import { type Area } from '@/services/area'

import githubIcon from '@/assets/app-icons/github.png'
import discordIcon from '@/assets/app-icons/discord.png'
import gmailIcon from '@/assets/app-icons/gmail.png'
import slackIcon from '@/assets/app-icons/slack.png'
import twitterIcon from '@/assets/app-icons/twitter.png'
import notionIcon from '@/assets/app-icons/notion.png'
import instagramIcon from '@/assets/app-icons/instagram.png'
import youtubeIcon from '@/assets/app-icons/youtube.png'
import spotifyIcon from '@/assets/app-icons/spotify.png'
import telegramIcon from '@/assets/app-icons/telegram.png'
import twitchIcon from '@/assets/app-icons/twitch.png'
import dropboxIcon from '@/assets/app-icons/dropbox.png'
import weatherIcon from '@/assets/app-icons/weather.png'
import googleCalendarIcon from '@/assets/app-icons/google-calendar.png'

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
  area: AreaTemplate | Area
  showDeleteButton?: boolean
}>()

const emit = defineEmits<{
  click: [area: AreaTemplate | Area]
  delete: [area: AreaTemplate | Area]
}>()

const cardRef = ref<HTMLElement | null>(null)
const mousePosition = ref({ x: 50, y: 50 })

const spotlightStyle = computed(() => ({
  '--mouse-x': `${mousePosition.value.x}%`,
  '--mouse-y': `${mousePosition.value.y}%`
}))

const handleMouseMove = (event: MouseEvent) => {
  if (!cardRef.value) return

  const rect = cardRef.value.getBoundingClientRect()
  const x = ((event.clientX - rect.left) / rect.width) * 100
  const y = ((event.clientY - rect.top) / rect.height) * 100

  mousePosition.value = { x, y }
}

const handleMouseEnter = () => {
}

const handleMouseLeave = () => {
  mousePosition.value = { x: 50, y: 50 }
}

const handleClick = () => {
  emit('click', props.area)
}

const handleDelete = () => {
  emit('delete', props.area)
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

const getIconUrl = (iconName: string) => {
  const iconMap: { [key: string]: string } = {
    'github.png': githubIcon,
    'discord.png': discordIcon,
    'gmail.png': gmailIcon,
    'slack.png': slackIcon,
    'twitter.png': twitterIcon,
    'notion.png': notionIcon,
    'instagram.png': instagramIcon,
    'youtube.png': youtubeIcon,
    'spotify.png': spotifyIcon,
    'telegram.png': telegramIcon,
    'twitch.png': twitchIcon,
    'dropbox.png': dropboxIcon,
    'weather.png': weatherIcon,
    'google-calendar.png': googleCalendarIcon,
  }

  return iconMap[iconName] || ''
}
</script>

<style scoped>
.card-spotlight-container {
  position: relative;
  cursor: pointer;
  max-width: 100%;
  width: 100%;
  transition: transform 0.3s ease;
}

.card-spotlight-container:hover {
  transform: translateY(-4px);
}

.card-spotlight {
  position: relative;
  background: linear-gradient(135deg, #1a1f2e 0%, #0f1419 100%);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 32px;
  height: 280px;
  overflow: hidden;
  transition: all 0.3s ease;
}

[data-theme="light"] .card-spotlight {
  background: linear-gradient(135deg, #ffffff 0%, #f8f9fa 100%);
  border: 2px solid rgba(0, 0, 0, 0.15);
}

.card-spotlight:hover {
  border-color: rgba(255, 255, 255, 0.2);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

[data-theme="light"] .card-spotlight:hover {
  border-color: rgba(59, 130, 246, 0.4);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12);
}

.spotlight {
  position: absolute;
  top: -50%;
  left: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(
    circle at var(--mouse-x, 50%) var(--mouse-y, 50%),
    rgba(59, 130, 246, 0.15) 0%,
    rgba(139, 92, 246, 0.1) 30%,
    transparent 70%
  );
  opacity: 0;
  transition: opacity 0.3s ease;
  pointer-events: none;
}

.card-spotlight-container:hover .spotlight {
  opacity: 1;
}

.card-content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.service-icons {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 20px;
  margin-bottom: 24px;
}

.service-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  transition: all 0.3s ease;
}

[data-theme="light"] .service-icon {
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid rgba(0, 0, 0, 0.15);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.service-icon:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.05);
}

[data-theme="light"] .service-icon:hover {
  background: rgba(255, 255, 255, 1);
  border-color: rgba(59, 130, 246, 0.4);
}

.service-arrow {
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0.7;
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


.service-logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  padding: 2px;
}

.card-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.card-title {
  font-weight: 700;
  font-size: 22px;
  color: white;
  margin: 0 0 6px 0;
  line-height: 1.2;
}

[data-theme="light"] .card-title {
  color: #1a1a1a;
}

.card-subtitle {
  color: rgba(255, 255, 255, 0.8);
  font-weight: 500;
  font-size: 16px;
  margin: 0 0 10px 0;
}

[data-theme="light"] .card-subtitle {
  color: rgba(0, 0, 0, 0.7);
}

.card-description {
  color: rgba(255, 255, 255, 0.6);
  font-size: 15px;
  margin: 0;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  line-clamp: 3;
  overflow: hidden;
}

[data-theme="light"] .card-description {
  color: rgba(0, 0, 0, 0.6);
}

.delete-button {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: rgba(239, 68, 68, 0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 10;
}

[data-theme="light"] .delete-button {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.delete-button:hover {
  background: rgba(239, 68, 68, 0.4);
  transform: scale(1.1);
}

[data-theme="light"] .delete-button:hover {
  background: rgba(239, 68, 68, 0.25);
}

[data-theme="light"] .card-spotlight :deep(.v-icon) {
  color: rgba(0, 0, 0, 0.8) !important;
}

[data-theme="light"] .icon-fallback {
  filter: brightness(0);
}

.delete-button:active {
  transform: scale(0.95);
}

.card-spotlight-container.has-delete .card-spotlight {
  padding-top: 48px;
}

@media (max-width: 768px) {
  .card-spotlight-container {
    max-width: 100%;
  }

  .card-spotlight {
    height: 240px;
    padding: 24px;
  }

  .card-title {
    font-size: 20px;
  }

  .card-description {
    font-size: 14px;
  }

  .service-icon {
    width: 48px;
    height: 48px;
  }

  .service-logo {
    width: 28px;
    height: 28px;
  }
}
</style>






