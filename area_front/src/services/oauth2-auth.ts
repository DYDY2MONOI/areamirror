import { API_BASE_URL } from '@/config/api'

export interface OAuth2User {
  id: number
  email: string
  first_name?: string | null
  last_name?: string | null
  created_at?: string | null
  updated_at?: string | null
  phone?: string | null
  birthday?: string | null
  gender?: string | null
  country?: string | null
  lang?: string | null
  login_provider?: string | null
  role?: string | null
  is_active?: boolean
  github_id?: string | null
  github_username?: string | null
  google_id?: string | null
  google_email?: string | null
  facebook_id?: string | null
  facebook_email?: string | null
  discord_id?: string | null
  discord_username?: string | null
  spotify_id?: string | null
  spotify_email?: string | null
  twitter_id?: string | null
  twitter_username?: string | null
  profile_image?: string | null
}

export interface OAuth2LoginRequest {
  email: string
  password: string
}

export interface OAuth2TokenResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
  user: OAuth2User
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface RefreshTokenResponse {
  access_token: string
  token_type: string
  expires_in: number
}

export interface MeResponse {
  user: OAuth2User
}

class OAuth2AuthService {
  private accessToken: string | null = null
  private refreshToken: string | null = null
  private user: OAuth2User | null = null
  private tokenExpiry: number | null = null

  constructor() {
    this.loadTokensFromStorage()
  }

  get isAuthenticated(): boolean {
    return !!this.accessToken && !!this.user && this.isTokenValid()
  }

  get currentUser(): OAuth2User | null {
    return this.user
  }

  get authToken(): string | null {
    return this.accessToken
  }

  get isAdmin(): boolean {
    return this.user?.role === 'admin'
  }

  get isMember(): boolean {
    return this.user?.role === 'member'
  }

  get canCreateAreas(): boolean {
    return this.isAdmin
  }

  private isTokenValid(): boolean {
    if (!this.tokenExpiry) return false
    return Date.now() < this.tokenExpiry
  }

  private loadTokensFromStorage(): void {
    this.accessToken = localStorage.getItem('oauth2_access_token')
    this.refreshToken = localStorage.getItem('oauth2_refresh_token')
    this.user = this.getStoredUser()
    this.tokenExpiry = this.getStoredTokenExpiry()
  }

  storeTokens(accessToken: string, refreshToken: string, expiresIn: number): void {
    this.accessToken = accessToken
    this.refreshToken = refreshToken
    this.tokenExpiry = Date.now() + (expiresIn * 1000)

    localStorage.setItem('oauth2_access_token', accessToken)
    localStorage.setItem('oauth2_refresh_token', refreshToken)
    localStorage.setItem('oauth2_token_expiry', this.tokenExpiry.toString())
  }

  private getStoredTokenExpiry(): number | null {
    const stored = localStorage.getItem('oauth2_token_expiry')
    return stored ? parseInt(stored, 10) : null
  }

  storeUser(user: OAuth2User): void {
    console.log('📦 storeUser called with:', user)
    console.log('📦 User github_username:', user?.github_username)

    if (!user) {
      console.error('❌ storeUser called with null/undefined user!')
      return
    }

    this.user = user
    try {
      localStorage.setItem('oauth2_user', JSON.stringify(user))
      console.log('✅ User stored successfully in localStorage')
    } catch (error) {
      console.error('❌ Error storing user in localStorage:', error)
    }
  }

  private getStoredUser(): OAuth2User | null {
    try {
      const storedUser = localStorage.getItem('oauth2_user')
      return storedUser ? JSON.parse(storedUser) : null
    } catch {
      return null
    }
  }

  async login(credentials: OAuth2LoginRequest): Promise<OAuth2TokenResponse> {
    try {
      const response = await fetch(`${API_BASE_URL}/oauth2/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Connection error' }))
        throw new Error(errorData.error || `HTTP error ${response.status}`)
      }

      const data: OAuth2TokenResponse = await response.json()
      this.handleSuccessfulAuth(data)
      return data
    } catch (error) {
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new Error('Unable to connect to server. Please check that the backend is running.')
      }
      throw error
    }
  }

  async refreshAccessToken(): Promise<RefreshTokenResponse> {
    if (!this.refreshToken) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await fetch(`${API_BASE_URL}/oauth2/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: this.refreshToken }),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Token refresh failed' }))
        throw new Error(errorData.error || `HTTP error ${response.status}`)
      }

      const data: RefreshTokenResponse = await response.json()
      this.accessToken = data.access_token
      this.tokenExpiry = Date.now() + (data.expires_in * 1000)

      localStorage.setItem('oauth2_access_token', data.access_token)
      localStorage.setItem('oauth2_token_expiry', this.tokenExpiry.toString())

      return data
    } catch (error) {
      await this.logout()
      throw error
    }
  }

  async fetchMe(): Promise<OAuth2User> {
    if (!this.accessToken) {
      throw new Error('Authentication token missing')
    }

    try {
      const response = await fetch(`${API_BASE_URL}/oauth2/me`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.accessToken}`,
        },
      })

      if (response.status === 401) {
        try {
          await this.refreshAccessToken()
          const retryResponse = await fetch(`${API_BASE_URL}/oauth2/me`, {
            method: 'GET',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${this.accessToken}`,
            },
          })

          if (!retryResponse.ok) {
            throw new Error('Authentication failed after token refresh')
          }

          const data: MeResponse = await retryResponse.json()
          this.user = data.user
          this.storeUser(data.user)
          return data.user
        } catch (refreshError) {
          await this.logout()
          throw new Error('Session expired. Please log in again.')
        }
      }

      if (!response.ok) {
        throw new Error('Failed to fetch user profile')
      }

      const data: MeResponse = await response.json()
      this.user = data.user
      this.storeUser(data.user)
      return data.user
    } catch (error) {
      throw error
    }
  }

  async logout(): Promise<void> {
    this.accessToken = null
    this.refreshToken = null
    this.user = null
    this.tokenExpiry = null

    localStorage.removeItem('oauth2_access_token')
    localStorage.removeItem('oauth2_refresh_token')
    localStorage.removeItem('oauth2_user')
    localStorage.removeItem('oauth2_token_expiry')
  }

  handleSuccessfulAuth(authResponse: OAuth2TokenResponse): void {
    console.log('🔐 handleSuccessfulAuth called with:', authResponse)
    console.log('🔐 Auth response user:', authResponse.user)

    if (!authResponse.user) {
      console.error('❌ handleSuccessfulAuth: user is missing from response!')
      console.error('❌ Full response:', authResponse)
      return
    }

    this.storeTokens(authResponse.access_token, authResponse.refresh_token, authResponse.expires_in)
    this.storeUser(authResponse.user)
  }

  async checkAuthStatus(): Promise<boolean> {
    if (!this.accessToken || !this.refreshToken) {
      return false
    }

    try {
      if (!this.isTokenValid()) {
        await this.refreshAccessToken()
      }

      await this.fetchMe()
      return true
    } catch {
      await this.logout()
      return false
    }
  }

  async makeAuthenticatedRequest(url: string, options: RequestInit = {}): Promise<Response> {
    if (!this.accessToken) {
      throw new Error('No access token available')
    }

    if (this.tokenExpiry && Date.now() > (this.tokenExpiry - 120000)) {
      try {
        await this.refreshAccessToken()
      } catch (error) {
        throw new Error('Failed to refresh token')
      }
    }

    const headers = {
      ...options.headers,
      'Authorization': `Bearer ${this.accessToken}`,
    }

    const response = await fetch(url, {
      ...options,
      headers,
    })

    if (response.status === 401) {
      try {
        await this.refreshAccessToken()
        const retryResponse = await fetch(url, {
          ...options,
          headers: {
            ...options.headers,
            'Authorization': `Bearer ${this.accessToken}`,
          },
        })
        return retryResponse
      } catch (error) {
        await this.logout()
        throw new Error('Authentication failed')
      }
    }

    return response
  }

  getProfileImageUrl(): string | null {
    if (!this.user?.profile_image) {
      return null
    }

    if (this.user.profile_image.startsWith('uploads/')) {
      return `${API_BASE_URL}/${this.user.profile_image}`
    }

    return this.user.profile_image
  }
}

export const oauth2AuthService = new OAuth2AuthService()
