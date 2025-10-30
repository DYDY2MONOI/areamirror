<template>
  <div class="user-section-container">
    <div class="user-section" v-if="isAuthenticated">
      <v-avatar size="48" class="gradient-avatar">
        <img
          v-if="profileImageUrl"
          :src="profileImageUrl"
          alt="Profile picture"
          class="profile-image"
        />
        <v-icon v-else color="white">mdi-account</v-icon>
      </v-avatar>
      <div class="user-info">
        <span class="user-name">{{ userDisplayName }}</span>
        <span class="user-status">{{ userStatus }}</span>
      </div>
    </div>

    <div class="guest-section" v-else>
      <div class="guest-content">
        <div class="guest-icon">
          <v-icon size="32" color="white">mdi-account-plus</v-icon>
        </div>
        <div class="guest-text">
          <h3 class="guest-title">{{ guestTitle }}</h3>
          <p class="guest-subtitle">{{ guestSubtitle }}</p>
        </div>
      </div>
      <div class="guest-actions">
        <button class="guest-btn primary" @click="$emit('login')">
          <v-icon size="16">mdi-login</v-icon>
          <span>{{ loginButtonText }}</span>
        </button>
        <button class="guest-btn secondary" @click="$emit('register')">
          <v-icon size="16">mdi-account-plus</v-icon>
          <span>{{ registerButtonText }}</span>
        </button>
      </div>
    </div>

    <div class="action-buttons">
      <button
        v-for="action in actionButtons"
        :key="action.id"
        class="action-btn-icon"
        @click="$emit('action-click', action.id)"
        :title="action.tooltip"
      >
        <v-icon size="20">{{ action.icon }}</v-icon>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
interface ActionButton {
  id: string
  icon: string
  tooltip?: string
}

interface Props {
  isAuthenticated: boolean
  profileImageUrl?: string
  userDisplayName?: string
  userStatus?: string
  guestTitle?: string
  guestSubtitle?: string
  loginButtonText?: string
  registerButtonText?: string
  actionButtons?: ActionButton[]
}

interface Emits {
  (e: 'login'): void
  (e: 'register'): void
  (e: 'action-click', actionId: string): void
}

withDefaults(defineProps<Props>(), {
  userDisplayName: 'User Name',
  userStatus: 'Premium Member',
  guestTitle: 'Join AREA Today',
  guestSubtitle: 'Start automating your workflow',
  loginButtonText: 'Sign In',
  registerButtonText: 'Join Us',
  actionButtons: () => [
    { id: 'search', icon: 'mdi-magnify', tooltip: 'Search' },
    { id: 'notifications', icon: 'mdi-bell-outline', tooltip: 'Notifications' }
  ]
})

defineEmits<Emits>()
</script>

<style scoped>
.user-section-container {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.user-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.user-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
}

.profile-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 50%;
}

.user-status {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  font-weight: 400;
}

.guest-section {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.guest-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.guest-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-full);
  background: var(--gradient-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--shadow-glow);
}

.guest-text {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.guest-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.guest-subtitle {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  font-weight: 400;
  margin: 0;
}

.guest-actions {
  display: flex;
  gap: 0.75rem;
}

.guest-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1rem;
  border-radius: var(--radius-lg);
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: var(--transition-normal);
  border: 2px solid transparent;
}

.guest-btn.primary {
  background: var(--gradient-accent);
  color: var(--color-text-primary);
  border-color: transparent;
}

.guest-btn.primary:hover {
  transform: translateY(-2px);
  box-shadow:
    var(--shadow-glow),
    0 8px 16px -5px rgba(6, 182, 212, 0.4);
}

.guest-btn.secondary {
  background: transparent;
  color: var(--color-text-primary);
  border-color: var(--color-border-primary);
}

.guest-btn.secondary:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-secondary);
  transform: translateY(-1px);
}

.action-buttons {
  display: flex;
  gap: var(--spacing-sm);
}

.action-btn-icon {
  width: 40px;
  height: 40px;
  border: none;
  background: var(--color-bg-card);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(20px);
  border: 1px solid var(--color-border-primary);
}

.action-btn-icon:hover {
  background: var(--color-hover-bg);
  transform: translateY(-1px);
}

.gradient-avatar {
  background: var(--gradient-accent);
}

@media (max-width: 1024px) {
  .user-section-container {
    flex-wrap: wrap;
    gap: 1rem;
  }
}

@media (max-width: 768px) {
  .guest-section {
    flex-direction: column;
    gap: 1rem;
    text-align: center;
  }

  .guest-content {
    flex-direction: column;
    gap: 0.75rem;
  }

  .guest-actions {
    flex-direction: column;
    width: 100%;
  }

  .guest-btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
