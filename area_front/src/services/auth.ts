export interface User {
  id: number
  email: string
  first_name?: string
  last_name?: string
  created_at?: string
  updated_at?: string
  phone?: string
  birthday?: string
  gender?: string
  country?: string
  lang?: string
  login_provider?: string
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

import { API_BASE_URL, API_ENDPOINTS } from '@/config/api'

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

  async login(credentials: LoginRequest): Promise<AuthResponse> {
    try {
      const response = await fetch(`${BASE_URL}${API_ENDPOINTS.LOGIN}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(credentials),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Erreur de connexion' }))
        throw new Error(errorData.error || `Erreur HTTP ${response.status}`)
      }

      const data = await response.json()
      this.handleSuccessfulAuth(data)
      return data
    } catch (error) {
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new Error('Impossible de se connecter au serveur. Vérifiez que le backend est démarré.')
      }
      throw error
    }
  }

  async register(userData: RegisterRequest): Promise<AuthResponse> {
    try {
      const response = await fetch(`${BASE_URL}${API_ENDPOINTS.REGISTER}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Erreur lors de l\'inscription' }))
        throw new Error(errorData.error || `Erreur HTTP ${response.status}`)
      }

      const data = await response.json()
      this.handleSuccessfulAuth(data)
      return data
    } catch (error) {
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new Error('Impossible de se connecter au serveur. Vérifiez que le backend est démarré.')
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
      throw new Error('Token d\'authentification manquant')
    }

    try {
      const response = await fetch(`${BASE_URL}${API_ENDPOINTS.PROFILE}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${this.token}`,
        },
      })

      const data: ProfileResponse = await response.json()

      if (!response.ok) {
        throw new Error(data.user ? 'Erreur lors de la récupération du profil' : 'Erreur de profil')
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
      throw new Error('Token d\'authentification manquant')
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
        throw new Error(data.user ? 'Erreur lors de la mise à jour du profil' : 'Erreur de mise à jour')
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
}

export const authService = new AuthService()
