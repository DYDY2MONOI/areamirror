export interface GitHubRepository {
  id: number
  name: string
  full_name: string
  description: string | null
  private: boolean
  html_url: string
  clone_url: string
  default_branch: string
  created_at: string
  updated_at: string
  pushed_at: string
}

export interface GitHubRepositoriesResponse {
  repositories: GitHubRepository[]
}

import { API_BASE_URL } from '@/config/api'

class GitHubService {
  private token: string | null = null

  constructor() {
    this.updateToken()
  }

  private updateToken() {
    this.token = localStorage.getItem('authToken')
  }

  private getAuthHeaders() {
    if (!this.token) {
      throw new Error('Authentication token missing')
    }
    return {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${this.token}`,
    }
  }

  async getRepositories(): Promise<GitHubRepository[]> {
    try {
      this.updateToken()
      console.log('🔍 GitHub Service: Fetching repositories from', `/api/github/repositories`)

      const response = await fetch(`/api/github/repositories`, {
        method: 'GET',
        headers: this.getAuthHeaders(),
      })

      console.log('📡 GitHub Service: Response status', response.status, response.statusText)

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Failed to fetch repositories' }))
        console.error('❌ GitHub Service: API Error', errorData)
        throw new Error(errorData.error || `HTTP error ${response.status}`)
      }

      const data: GitHubRepositoriesResponse = await response.json()
      console.log('✅ GitHub Service: Repositories loaded', data.repositories?.length || 0, 'repositories')
      return data.repositories || []
    } catch (error) {
      console.error('💥 GitHub Service: Error fetching repositories:', error)
      if (error instanceof TypeError && error.message.includes('fetch')) {
        throw new Error('Impossible de se connecter au serveur. Vérifiez que le backend est démarré.')
      }
      throw error
    }
  }

  async createArea(repositoryId: number, destinationEmail: string, notificationTypes: string[]): Promise<void> {
    try {
      this.updateToken()
      const response = await fetch(`/api/areas/github-gmail`, {
        method: 'POST',
        headers: this.getAuthHeaders(),
        body: JSON.stringify({
          repository_id: repositoryId,
          destination_email: destinationEmail,
          notification_types: notificationTypes
        }),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ error: 'Failed to create area' }))
        throw new Error(errorData.error || `HTTP error ${response.status}`)
      }
    } catch (error) {
      console.error('Error creating GitHub-Gmail area:', error)
      throw error
    }
  }
}

export const githubService = new GitHubService()
