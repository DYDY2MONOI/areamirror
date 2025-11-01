export interface SpotifyOAuthConfig {
  clientId: string
  redirectUri: string
  scope: string
  showDialog?: boolean
}

export class SpotifyOAuth {
  private config: SpotifyOAuthConfig

  constructor(config: SpotifyOAuthConfig) {
    this.config = config
  }

  private generateState(): string {
    return Math.random().toString(36).slice(2) + Math.random().toString(36).slice(2)
  }

  generateAuthUrl(): string {
    const params = new URLSearchParams({
      client_id: this.config.clientId,
      response_type: 'code',
      redirect_uri: this.config.redirectUri,
      scope: this.config.scope,
      state: this.generateState(),
    })

    if (this.config.showDialog) {
      params.set('show_dialog', 'true')
    }

    return `https://accounts.spotify.com/authorize?${params.toString()}`
  }

  redirectToSpotify(): void {
    const authUrl = this.generateAuthUrl()
    window.location.href = authUrl
  }
}

interface SpotifyOAuthOptions {
  redirectUri?: string
  scope?: string
  showDialog?: boolean
}

export const createSpotifyOAuth = (options: SpotifyOAuthOptions = {}): SpotifyOAuth => {
  const clientId = import.meta.env.VITE_SPOTIFY_CLIENT_ID
  const fallbackRedirect =
    typeof window !== 'undefined'
      ? `${window.location.origin}/oauth2/spotify/callback`
      : 'http://localhost:3000/oauth2/spotify/callback'
  const defaultRedirect = (import.meta.env.VITE_SPOTIFY_REDIRECT_URI || '').trim() || fallbackRedirect

  if (!clientId || clientId === 'your_spotify_client_id_here') {
    throw new Error('Spotify OAuth client ID is not configured')
  }

  const scope =
    options.scope ||
    'user-read-email user-read-private user-read-currently-playing user-read-playback-state'
  const redirectUri = options.redirectUri || defaultRedirect

  return new SpotifyOAuth({
    clientId,
    redirectUri,
    scope,
    showDialog: options.showDialog ?? true,
  })
}

export const initiateSpotifyLogin = (): void => {
  const spotifyOAuth = createSpotifyOAuth()
  spotifyOAuth.redirectToSpotify()
}

export const initiateSpotifyLink = (redirectUri: string, scope?: string): void => {
  const spotifyOAuth = createSpotifyOAuth({ redirectUri, scope })
  spotifyOAuth.redirectToSpotify()
}
