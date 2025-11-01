export interface GitHubOAuthConfig {
  clientId: string
  redirectUri: string
  scope: string
}

export class GitHubOAuth {
  private config: GitHubOAuthConfig

  constructor(config: GitHubOAuthConfig) {
    this.config = config
  }

  generateAuthUrl(): string {
    const params = new URLSearchParams({
      client_id: this.config.clientId,
      redirect_uri: this.config.redirectUri,
      scope: this.config.scope,
      state: this.generateState()
    })

    return `https://github.com/login/oauth/authorize?${params.toString()}`
  }

  private generateState(): string {
    return Math.random().toString(36).substring(2, 15) +
           Math.random().toString(36).substring(2, 15)
  }

  redirectToGitHub(): void {
    const authUrl = this.generateAuthUrl()
    window.location.href = authUrl
  }
}

export const createGitHubOAuth = (): GitHubOAuth => {
  const config: GitHubOAuthConfig = {
    clientId: import.meta.env.VITE_GITHUB_CLIENT_ID || 'Ov23liQ7GPEEWs0hVzyM',
    redirectUri: `${window.location.origin}/oauth2/github/callback`,
    scope: 'user:email'
  }

  return new GitHubOAuth(config)
}

export const initiateGitHubLogin = (): void => {
  const gitHubOAuth = createGitHubOAuth()
  gitHubOAuth.redirectToGitHub()
}
