export interface ServiceConfig {
  id: string
  name: string
  icon: string
  color: string
  description: string
  authUrl?: string
  callbackPath: string
  scopes: string[]
  isEnabled: boolean
}

export const SERVICES_CONFIG: ServiceConfig[] = [
  {
    id: 'github',
    name: 'GitHub',
    icon: 'github',
    color: '#24292e',
    description: 'Connect your GitHub account to access repositories, issues, and pull requests',
    authUrl: 'https://github.com/login/oauth/authorize',
    callbackPath: '/auth/github/callback',
    scopes: ['user:email'],
    isEnabled: true
  },
  {
    id: 'google',
    name: 'Google',
    icon: 'google',
    color: '#4285f4',
    description: 'Connect your Google account to access Gmail, Calendar, and Drive',
    authUrl: 'https://accounts.google.com/o/oauth2/v2/auth',
    callbackPath: '/callback',
    scopes: ['openid', 'email', 'profile'],
    isEnabled: true
  },
  {
    id: 'discord',
    name: 'Discord',
    icon: 'discord',
    color: '#5865f2',
    description: 'Connect your Discord account to send messages and manage servers',
    callbackPath: '/auth/discord/callback',
    scopes: ['identify', 'guilds'],
    isEnabled: false // Not implemented yet
  },
  {
    id: 'spotify',
    name: 'Spotify',
    icon: 'spotify',
    color: '#1db954',
    description: 'Connect your Spotify account to control music and playlists',
    callbackPath: '/auth/spotify/callback',
    scopes: ['user-read-email', 'user-read-private'],
    isEnabled: false // Not implemented yet
  }
]

export const getServiceConfig = (serviceId: string): ServiceConfig | undefined => {
  return SERVICES_CONFIG.find(service => service.id === serviceId)
}

export const getEnabledServices = (): ServiceConfig[] => {
  return SERVICES_CONFIG.filter(service => service.isEnabled)
}
