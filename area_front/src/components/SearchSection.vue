<template>
  <div class="search-section" :id="sectionId">
    <div class="search-container">
      <div class="search-header">
        <h1 class="search-title">{{ title }}</h1>
        <p class="search-subtitle">{{ subtitle }}</p>
      </div>
      <div class="search-bar">
        <div class="search-input-container">
          <v-icon size="20" color="#9ca3af" class="search-icon">mdi-magnify</v-icon>
          <input
            :value="modelValue"
            @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
            type="text"
            :placeholder="placeholder"
            class="search-input"
          >
          <button v-if="modelValue" @click="$emit('clear')" class="search-filter-btn">
            <v-icon size="16" color="#9ca3af">mdi-close</v-icon>
          </button>
        </div>
        <div v-if="suggestions.length > 0" class="search-suggestions">
          <span class="suggestion-label">{{ suggestionLabel }}</span>
          <button
            v-for="suggestion in suggestions"
            :key="suggestion"
            class="suggestion-chip"
            @click="$emit('suggestion-click', suggestion)"
          >
            {{ suggestion }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  modelValue: string
  title?: string
  subtitle?: string
  placeholder?: string
  suggestions?: string[]
  suggestionLabel?: string
  sectionId?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'clear'): void
  (e: 'suggestion-click', suggestion: string): void
}

withDefaults(defineProps<Props>(), {
  title: 'Find Your Perfect Automation',
  subtitle: 'Discover templates, browse services, or create something new',
  placeholder: 'Search automations, services, or templates...',
  suggestions: () => ['Gmail', 'Discord', 'Spotify', 'GitHub'],
  suggestionLabel: 'Popular:',
  sectionId: 'search-section'
})

defineEmits<Emits>()
</script>

<style scoped>
.search-section {
  padding: 2rem 2rem 3rem 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 30vh;
  position: relative;
}

.search-container {
  width: 100%;
  max-width: 800px;
  text-align: center;
}

.search-header {
  margin-bottom: 2rem;
}

.search-title {
  font-size: 2.5rem;
  font-weight: 800;
  color: white;
  margin: 0 0 0.75rem 0;
  line-height: 1.1;
  letter-spacing: -0.02em;
  background: linear-gradient(135deg, #ffffff 0%, #f0f9ff 50%, #e0f2fe 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
  filter: drop-shadow(0 4px 12px rgba(255, 255, 255, 0.1));
}

.search-subtitle {
  font-size: 1rem;
  color: #9ca3af;
  margin: 0;
  line-height: 1.6;
  font-weight: 400;
}

.search-bar {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.search-input-container {
  position: relative;
  display: flex;
  align-items: center;
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-xl);
  padding: var(--spacing-md) var(--spacing-lg);
  backdrop-filter: blur(20px);
  transition: var(--transition-normal);
}

.search-input-container:focus-within {
  border-color: var(--color-border-focus);
  box-shadow: 0 0 0 3px var(--color-focus-ring);
}

.search-icon {
  margin-right: 1rem;
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  color: var(--color-text-primary);
  font-size: 1rem;
  font-weight: 400;
}

.search-input::placeholder {
  color: var(--color-text-secondary);
}

.search-filter-btn {
  background: transparent;
  border: none;
  padding: var(--spacing-sm);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: var(--transition-normal);
  margin-left: var(--spacing-md);
}

.search-filter-btn:hover {
  background: var(--color-hover-bg);
}

.search-suggestions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  justify-content: center;
  flex-wrap: wrap;
}

.suggestion-label {
  font-size: 0.875rem;
  color: var(--color-text-secondary);
  font-weight: 500;
}

.suggestion-chip {
  background: var(--color-bg-card);
  border: 1px solid var(--color-border-primary);
  border-radius: var(--radius-2xl);
  padding: var(--spacing-sm) var(--spacing-md);
  color: var(--color-text-primary);
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: var(--transition-normal);
  backdrop-filter: blur(20px);
}

.suggestion-chip:hover {
  background: var(--color-hover-bg);
  border-color: var(--color-border-focus);
  transform: translateY(-1px);
}

@media (max-width: 1024px) {
  .search-section {
    padding: 2rem 1rem 2rem 1rem;
    min-height: 25vh;
  }

  .search-title {
    font-size: 2rem;
  }

  .search-suggestions {
    justify-content: flex-start;
  }
}

@media (max-width: 768px) {
  .search-section {
    padding: 1.5rem 1rem 1.5rem 1rem;
    min-height: 20vh;
  }

  .search-title {
    font-size: 1.75rem;
  }

  .search-subtitle {
    font-size: 0.875rem;
  }

  .search-input-container {
    padding: 0.875rem 1rem;
  }

  .search-suggestions {
    gap: 0.5rem;
  }

  .suggestion-chip {
    padding: 0.375rem 0.75rem;
    font-size: 0.8125rem;
  }
}
</style>
