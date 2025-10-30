import { ref, watch } from 'vue'

export type Theme = 'dark' | 'light' | 'high-contrast'
export type DaltonismMode = 'none' | 'protanopia' | 'deuteranopia' | 'tritanopia' | 'monochrome'

const getStoredTheme = (): Theme => {
  const raw = localStorage.getItem('area_theme')
  if (raw === 'dark' || raw === 'light' || raw === 'high-contrast') return raw
  return 'dark'
}

const getStoredPrevNonHC = (): Exclude<Theme, 'high-contrast'> => {
  const raw = localStorage.getItem('area_theme_prev')
  return raw === 'light' ? 'light' : 'dark'
}

const getStoredDaltonism = (): DaltonismMode => {
  const raw = localStorage.getItem('area_daltonism_mode') as DaltonismMode | null
  return raw === 'protanopia' || raw === 'deuteranopia' || raw === 'tritanopia' || raw === 'monochrome' ? raw : 'none'
}

const currentTheme = ref<Theme>(getStoredTheme())
const prevNonHCTheme = ref<Exclude<Theme, 'high-contrast'>>(
  currentTheme.value === 'high-contrast' ? getStoredPrevNonHC() : (currentTheme.value as Exclude<Theme, 'high-contrast'>)
)
const daltonismMode = ref<DaltonismMode>(getStoredDaltonism())

function getDaltonismRoot(): HTMLElement {
  return (document.getElementById('app') as HTMLElement) || document.documentElement
}

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

  const toggleTheme = () => {
    if (currentTheme.value === 'high-contrast') {
      const nextBase = prevNonHCTheme.value === 'dark' ? 'light' : 'dark'
      prevNonHCTheme.value = nextBase
      localStorage.setItem('area_theme_prev', nextBase)
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

  const setDaltonismMode = (mode: DaltonismMode) => {
    daltonismMode.value = mode
    localStorage.setItem('area_daltonism_mode', mode)
    const root = getDaltonismRoot()
    if (mode === 'none') {
      root.removeAttribute('data-daltonism')
    } else {
      root.setAttribute('data-daltonism', mode)
    }
  }

  const isDark = () => currentTheme.value === 'dark'
  const isLight = () => currentTheme.value === 'light'
  const isHighContrast = () => currentTheme.value === 'high-contrast'

  if (!document.documentElement.hasAttribute('data-theme')) {
    document.documentElement.setAttribute('data-theme', currentTheme.value)
  }
  if (daltonismMode.value !== 'none') {
    const root = getDaltonismRoot()
    root.setAttribute('data-daltonism', daltonismMode.value)
  }

  watch(currentTheme, (newTheme) => {
    document.documentElement.setAttribute('data-theme', newTheme)
  })

  watch(daltonismMode, (newMode) => {
    const root = getDaltonismRoot()
    if (newMode === 'none') {
      root.removeAttribute('data-daltonism')
    } else {
      root.setAttribute('data-daltonism', newMode)
    }
  })

  return {
    currentTheme,
    prevNonHCTheme,
    daltonismMode,
    setTheme,
    toggleTheme,
    toggleHighContrast,
    enableHighContrast,
    disableHighContrast,
    setDaltonismMode,
    isDark,
    isLight,
    isHighContrast,
  }
}






