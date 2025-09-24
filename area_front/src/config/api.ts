const getApiUrl = () => {
  if (import.meta.env.DEV) {
    return 'http://localhost:8080'
  }

  return 'http://10.68.246.178:8080'
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
