const getApiUrl = () => {
  return '/api'
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
