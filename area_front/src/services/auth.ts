export interface User {
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
  onedrive_id?: string | null
  profile_image?: string | null
}

export interface AuthResponse {
  message: string
  token: string
  user: User
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  first_name?: string
  last_name?: string
}

export interface ProfileUpdateRequest {
  first_name?: string
  last_name?: string
  phone?: string
  country?: string
  current_password?: string
  new_password?: string
}

export interface ProfileResponse {
  user: User
}

import { API_BASE_URL, API_ENDPOINTS, AUTH_ENDPOINTS } from '@/config/api'

const BASE_URL = API_BASE_URL

class AuthService {
  private token: string | null = null
  private user: User | null = null

  constructor() {
    this.token = localStorage.getItem('authToken')
    this.user = this.getStoredUser()
  }

  get isAuthenticated(): boolean {
    return !!this.token && !!this.user
  }

  get currentUser(): User | null {
    return this.user
  }

  get authToken(): string | null {
    return this.token
  }

  get isAdmin(): boolean {
    return this.user?.role === 'admin'
  }

  get isMember(): boolean {
    return this.user?.role === 'member'
  }

  get canCreateAreas(): boolean {
    return this.isAuthenticated
  }

  async login(credentials: LoginRequest): Promise<AuthResponse> {
    try {
      const response = await fetch(`${BASE_URL}${AUTH_ENDPOINTS.LOGIN}`, {
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

      const data = await response.json()
      this.handleSuccessfulAuth(data)
      return data
    } catch (error) {
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new Error('Unable to connect to server. Please check that the backend is running.')
      }
      throw error
    }
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    console.log('🔐 Service: Début de l\'enregistrement', userData)
    try {
      console.log('🌐 Service: Envoi de la requête vers', `${BASE_URL}${API_ENDPOINTS.REGISTER}`)

      const controller = new AbortController()
      const timeoutId = setTimeout(() => controller.abort(), 10000)

      const response = await fetch(`${BASE_URL}${AUTH_ENDPOINTS.REGISTER}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
        signal: controller.signal
      })

      clearTimeout(timeoutId)
      console.log('📡 Service: Réponse reçue', response.status, response.statusText)

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Registration error' }))
        throw new Error(errorData.error || `HTTP error ${response.status}`)
      }

      const data = await response.json()
      console.log('✅ Service: Données reçues', data)
      this.handleSuccessfulAuth(data)
      return data
    } catch (error) {
      console.error('💥 Service: Erreur capturée', error)
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new Error('Unable to connect to server. Please check that the backend is running.')
      }
      if (error instanceof Error && error.name === 'AbortError') {
        throw new Error('La requête a expiré. Le serveur met trop de temps à répondre.')
      }
      throw error
    }
  }

  async logout(): Promise<void> {
    this.token = null
    this.user = null
    localStorage.removeItem('authToken')
    localStorage.removeItem('user')
  }

  async fetchProfile(): Promise<User> {
    if (!this.token) {
      throw new Error('Authentication token missing')
    }

    try {
      const response = await fetch(`${BASE_URL}${AUTH_ENDPOINTS.PROFILE}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`,
        },
      })

      const data: ProfileResponse = await response.json()

      if (!response.ok) {
        throw new Error(data.user ? 'Error retrieving profile' : 'Profile error')
      }

      this.user = data.user
      this.storeUser(data.user)
      return data.user
    } catch (error) {
      if (error instanceof Error && error.message.includes('401')) {
        await this.logout()
      }
      throw error
    }
  }

  async updateProfile(updateData: ProfileUpdateRequest): Promise<User> {
    if (!this.token) {
      throw new Error('Authentication token missing')
    }

    try {
      const response = await fetch(`${BASE_URL}${API_ENDPOINTS.PROFILE}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`,
        },
        body: JSON.stringify(updateData),
      })

      const data: ProfileResponse = await response.json()

      if (!response.ok) {
        throw new Error(data.user ? 'Error updating profile' : 'Update error')
      }

      this.user = data.user
      this.storeUser(data.user)
      return data.user
    } catch (error) {
      throw error
    }
  }


  private handleSuccessfulAuth(authResponse: AuthResponse): void {
    this.token = authResponse.token
    this.user = authResponse.user
    localStorage.setItem('authToken', authResponse.token)
    this.storeUser(authResponse.user)
  }

  private storeUser(user: User): void {
    localStorage.setItem('user', JSON.stringify(user))
  }

  private getStoredUser(): User | null {
    try {
      const storedUser = localStorage.getItem('user')
      return storedUser ? JSON.parse(storedUser) : null
    } catch {
      return null
    }
  }

  async checkAuthStatus(): Promise<boolean> {
    if (!this.token) {
      return false
    }

    try {
      await this.fetchProfile()
      return true
    } catch {
      await this.logout()
      return false
    }
  }

  async linkGitHubAccount(code: string): Promise<{ github_username: string }> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/github/link`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ code })
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to link GitHub account')
    }

    const data = await response.json()
    await this.fetchProfile()
    return data
  }

  async unlinkGitHubAccount(): Promise<void> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/github/unlink`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to unlink GitHub account')
    }

    await this.fetchProfile()
  }

  async linkGoogleAccount(code: string, redirectUri?: string): Promise<{ google_email?: string; message?: string; user?: User }> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const payload: Record<string, unknown> = { code }
    if (redirectUri) {
      payload.redirect_uri = redirectUri
    }

    const response = await fetch(`${API_BASE_URL}/profile/google/link`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to link Google account')
    }

    const data = await response.json()

    if (data.user) {
      this.user = data.user
      this.storeUser(data.user)
    } else {
      await this.fetchProfile()
    }
    return data
  }

  async unlinkGoogleAccount(): Promise<void> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/google/unlink`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to unlink Google account')
    }

    await this.fetchProfile()
  }

  async uploadProfileImage(imageFile: File): Promise<User> {
    if (!this.token) {
      throw new Error('Authentication token missing')
    }

    try {
      console.log('📸 Service: Starting image upload')
      const formData = new FormData()
      formData.append('image', imageFile)

      console.log('📸 Service: Sending to', `${BASE_URL}/profile/image`)
      const response = await fetch(`${BASE_URL}/profile/image`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${this.token}`,
        },
        body: formData,
      })

      console.log('📸 Service: Response received', response.status)
      const data: ProfileResponse = await response.json()
      console.log('📸 Service: Data received', data)

      if (!response.ok) {
        throw new Error(data.user ? 'Error uploading image' : 'Upload error')
      }

      this.user = data.user
      this.storeUser(data.user)
      console.log('📸 Service: Upload completed successfully')
      return data.user
    } catch (error) {
      console.error('📸 Service: Error during upload:', error)
      throw error
    }
  }

  getProfileImageUrl(): string | null {
    if (!this.user?.profile_image) {
      return null
    }

    if (this.user.profile_image.startsWith('uploads/')) {
      return `${BASE_URL}/${this.user.profile_image}`
    }

    return this.user.profile_image
  }

  async linkFacebookAccount(code: string): Promise<{ facebook_email: string }> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/facebook/link`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ code })
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to link Facebook account')
    }

    const data = await response.json()
    await this.fetchProfile()
    return data
  }

  async unlinkFacebookAccount(): Promise<void> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/facebook/unlink`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to unlink Facebook account')
    }

    await this.fetchProfile()
  }

  async linkOneDriveAccount(code: string): Promise<{ onedrive_email: string }> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/onedrive/link`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ code })
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to link OneDrive account')
    }

    const data = await response.json()
    await this.fetchProfile()
    return data
  }

  async unlinkOneDriveAccount(): Promise<void> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/onedrive/unlink`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to unlink OneDrive account')
    }

    await this.fetchProfile()
  }

  async linkSpotifyAccount(code: string): Promise<{ spotify_email: string }> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/spotify/link`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ code })
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to link Spotify account')
    }

    const data = await response.json()
    await this.fetchProfile()
    return data
  }

  async unlinkSpotifyAccount(): Promise<void> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/spotify/unlink`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to unlink Spotify account')
    }

    await this.fetchProfile()
  }

  async linkTwitterAccount(code: string, codeVerifier: string): Promise<{ twitter_username: string }> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/twitter/link`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      },
      body: JSON.stringify({ code, code_verifier: codeVerifier })
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to link Twitter/X account')
    }

    const data = await response.json()
    await this.fetchProfile()
    return data
  }

  async unlinkTwitterAccount(): Promise<void> {
    const token = localStorage.getItem('authToken')
    if (!token) {
      throw new Error('No authentication token found')
    }

    const response = await fetch(`${API_BASE_URL}/profile/twitter/unlink`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${token}`
      }
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to unlink Twitter/X account')
    }

    await this.fetchProfile()
  }
}

export const authService = new AuthService()
