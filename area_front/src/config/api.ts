const getApiUrl = () => {
  const envUrl = import.meta.env.VITE_API_BASE_URL?.trim().replace(/\/$/, '')

  const isLocalEnv =
    envUrl &&
    (envUrl.includes('localhost') || envUrl.includes('127.0.0.1'))

  if (envUrl && !isLocalEnv) {
    return envUrl
  }

  if (typeof window !== 'undefined') {
    const isLocalHost =
      window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'

    if (!isLocalHost) {
      return `${window.location.origin}/api`
    }
  }

  if (envUrl) {
    return envUrl
  }

  return 'http://localhost:8080'
}

export const API_BASE_URL = getApiUrl()

export const API_ENDPOINTS = {
  LOGIN: '/login',
  REGISTER: '/register',
  PROFILE: '/profile',
  SERVICES: '/services',
  ACTIONS: '/actions',
  REACTIONS: '/reactions',
  AREAS: '/areas',
  APPLETS: '/applets'
} as const

export const AUTH_ENDPOINTS = {
  LOGIN: '/login',
  REGISTER: '/register',
  PROFILE: '/profile'
} as const
