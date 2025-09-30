const getApiUrl = () => {
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
