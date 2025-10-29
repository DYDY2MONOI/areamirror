import { ref, watch } from 'vue'

export type Theme = 'dark' | 'light' | 'high-contrast'

const getStoredTheme = (): Theme => {
  const raw = localStorage.getItem('area_theme')
  if (raw === 'dark' || raw === 'light' || raw === 'high-contrast') return raw
  return 'dark'
}

const getStoredPrevNonHC = (): Exclude<Theme, 'high-contrast'> => {
  const raw = localStorage.getItem('area_theme_prev')
  return raw === 'light' ? 'light' : 'dark'
}

const currentTheme = ref<Theme>(getStoredTheme())
const prevNonHCTheme = ref<Exclude<Theme, 'high-contrast'>>(
  currentTheme.value === 'high-contrast' ? getStoredPrevNonHC() : (currentTheme.value as Exclude<Theme, 'high-contrast'>)
)

export function useTheme() {
  const setTheme = (theme: Theme) => {
    if (theme === 'dark' || theme === 'light') {
      prevNonHCTheme.value = theme
      localStorage.setItem('area_theme_prev', theme)
    }
    currentTheme.value = theme
    localStorage.setItem('area_theme', theme)
    document.documentElement.setAttribute('data-theme', theme)
  }

  // Toggle only between light and dark.
  // If currently in high-contrast, flip the stored base theme without exiting high-contrast.
  const toggleTheme = () => {
    if (currentTheme.value === 'high-contrast') {
      const nextBase = prevNonHCTheme.value === 'dark' ? 'light' : 'dark'
      prevNonHCTheme.value = nextBase
      localStorage.setItem('area_theme_prev', nextBase)
      // Do not change currentTheme; stay in high-contrast
      return
    }
    const next = currentTheme.value === 'dark' ? 'light' : 'dark'
    setTheme(next)
  }

  const toggleHighContrast = () => {
    if (currentTheme.value === 'high-contrast') {
      setTheme(prevNonHCTheme.value)
    } else {
      setTheme('high-contrast')
    }
  }

  const enableHighContrast = () => setTheme('high-contrast')
  const disableHighContrast = () => setTheme(prevNonHCTheme.value)

  const isDark = () => currentTheme.value === 'dark'
  const isLight = () => currentTheme.value === 'light'
  const isHighContrast = () => currentTheme.value === 'high-contrast'

  if (!document.documentElement.hasAttribute('data-theme')) {
    document.documentElement.setAttribute('data-theme', currentTheme.value)
  }

  watch(currentTheme, (newTheme) => {
    document.documentElement.setAttribute('data-theme', newTheme)
  })

  return {
    currentTheme,
    prevNonHCTheme,
    setTheme,
    toggleTheme,
    toggleHighContrast,
    enableHighContrast,
    disableHighContrast,
    isDark,
    isLight,
    isHighContrast,
  }
}






