<template>
  <div class="card-col" @click="handleClick">
    <v-sheet :class="`area-card ${getGradientClass(area)}`" rounded="xl">
      <div class="area-icons-container">
        <div class="service-icon trigger-icon">
          <img
            v-if="area.triggerIconUrl"
            :src="getIconUrl(area.triggerIconUrl)"
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
            :src="getIconUrl(area.actionIconUrl)"
            :alt="area.actionService"
            class="service-logo"
          />
          <v-icon v-else :size="32" color="white">{{ getActionIcon(area.actionService) }}</v-icon>
          <div v-if="!area.actionIconUrl" class="icon-fallback">{{ getActionEmoji(area.actionService) }}</div>
        </div>
      </div>
      <button
        v-if="showDeleteButton"
        class="delete-button"
        @click.stop="handleDelete"
        type="button"
      >
        <v-icon size="18" color="white">mdi-delete</v-icon>
      </button>
    </v-sheet>
    <div class="card-title">{{ 'title' in area ? area.title : area.name }}</div>
    <div v-if="'subtitle' in area" class="card-subtitle">{{ area.subtitle }}</div>
    <div class="card-description">{{ area.description }}</div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'
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

const handleClick = () => {
  emit('click', props.area)
}

const handleDelete = () => {
  emit('delete', props.area)
}

const getTriggerIcon = (service: string) => {
  switch (service) {
    case "Date Timer":
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
    case "Date Timer":
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

const getGradientClass = (area: AreaTemplate | Area) => {
  const trigger = area.triggerService
  const action = area.actionService

  if (trigger === "GitHub" && action === "Discord") {
    return "gradient-github-discord"
  }

  if (trigger === "Discord" && action === "Notion") {
    return "gradient-discord-notion"
  }

  if (trigger === "GitHub" && action === "Gmail") {
    return "gradient-github-gmail"
  }

  return "gradient-default"
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
.card-col {
  max-width: 400px;
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
  background-size: 100% 100%;
  transition:
    transform .25s ease,
    box-shadow .25s ease,
    filter .25s ease;
  position: relative;
}

/* GitHub to Discord gradient */
.gradient-github-discord {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

/* Discord to Notion gradient */
.gradient-discord-notion {
  background: linear-gradient(135deg, #5865f2 0%, #00d4aa 100%);
}

/* GitHub to Gmail gradient */
.gradient-github-gmail {
  background: linear-gradient(135deg, #333333 0%, #ea4335 100%);
}

/* Default gradient for other combinations */
.gradient-default {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  background: rgba(255, 255, 255, 0.1);
  border-radius: 4px;
  padding: 2px;
}

.area-card :deep(.v-icon) {
  transition: transform .25s ease, opacity .25s ease;
}

.area-card:hover {
  transform: translateY(-2px) scale(1.01);
  box-shadow: 0 8px 20px rgba(0,0,0,0.2);
}

.area-card:hover :deep(.v-icon) {
  transform: translateY(-1px) scale(1.03);
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

.delete-button {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: rgba(255, 0, 0, 0.2);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  z-index: 10;
}

.delete-button:hover {
  background: rgba(255, 0, 0, 0.4);
  transform: scale(1.1);
}

.delete-button:active {
  transform: scale(0.95);
}
</style>
