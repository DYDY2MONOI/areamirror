import { ref, watch } from 'vue'

type Theme = 'dark' | 'light'

const currentTheme = ref<Theme>((localStorage.getItem('area_theme') as Theme) || 'dark')

export function useTheme() {
  const setTheme = (theme: Theme) => {
    currentTheme.value = theme
    localStorage.setItem('area_theme', theme)
    document.documentElement.setAttribute('data-theme', theme)
  }

  const toggleTheme = () => {
    const newTheme = currentTheme.value === 'dark' ? 'light' : 'dark'
    setTheme(newTheme)
  }

  const isDark = () => currentTheme.value === 'dark'
  const isLight = () => currentTheme.value === 'light'

  if (!document.documentElement.hasAttribute('data-theme')) {
    document.documentElement.setAttribute('data-theme', currentTheme.value)
  }

  watch(currentTheme, (newTheme) => {
    document.documentElement.setAttribute('data-theme', newTheme)
  })

  return {
    currentTheme,
    setTheme,
    toggleTheme,
    isDark,
    isLight
  }
}




