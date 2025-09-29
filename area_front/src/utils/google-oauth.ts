export interface GoogleOAuthConfig {
  clientId: string
  redirectUri: string
  scope: string
}

export class GoogleOAuth {
  private config: GoogleOAuthConfig

  constructor(config: GoogleOAuthConfig) {
    this.config = config
  }

  generateAuthUrl(): string {
    const params = new URLSearchParams({
      client_id: this.config.clientId,
      redirect_uri: this.config.redirectUri,
      scope: this.config.scope,
      response_type: 'code',
      access_type: 'offline',
      prompt: 'consent',
      state: this.generateState()
    })

    return `https://accounts.google.com/o/oauth2/v2/auth?${params.toString()}`
  }

  private generateState(): string {
    return Math.random().toString(36).substring(2, 15) + 
           Math.random().toString(36).substring(2, 15)
  }

  redirectToGoogle(): void {
    const authUrl = this.generateAuthUrl()
    window.location.href = authUrl
  }
}

export const createGoogleOAuth = (): GoogleOAuth => {
  const config: GoogleOAuthConfig = {
    clientId: import.meta.env.VITE_GOOGLE_CLIENT_ID || 'your_google_client_id_here',
    redirectUri: 'http://localhost:3000/callback',
    scope: 'openid email profile'
  }

  return new GoogleOAuth(config)
}

export const initiateGoogleLogin = (): void => {
  const googleOAuth = createGoogleOAuth()
  googleOAuth.redirectToGoogle()
}
