import { API_BASE_URL } from '@/config/api'

const transformAreaData = (backendArea: any): Area => {
  return {
    id: backendArea.id,
    name: backendArea.name,
    description: backendArea.description,
    triggerService: backendArea.trigger_service || backendArea.triggerService,
    triggerType: backendArea.trigger_type || backendArea.triggerType,
    actionService: backendArea.action_service || backendArea.actionService,
    actionType: backendArea.action_type || backendArea.actionType,
    isActive: backendArea.is_active !== undefined ? backendArea.is_active : backendArea.isActive,
    isPublic: backendArea.is_public !== undefined ? backendArea.is_public : backendArea.isPublic,
    createdAt: backendArea.created_at || backendArea.createdAt,
    updatedAt: backendArea.updated_at || backendArea.updatedAt,
    triggerIconUrl: backendArea.trigger_icon_url || backendArea.triggerIconUrl,
    actionIconUrl: backendArea.action_icon_url || backendArea.actionIconUrl,
    triggerConfig: backendArea.trigger_config || backendArea.triggerConfig,
    actionConfig: backendArea.action_config || backendArea.actionConfig
  }
}

export interface Area {
  id: string
  name: string
  description: string
  triggerService: string
  triggerType: string
  actionService: string
  actionType: string
  isActive: boolean
  isPublic: boolean
  createdAt: string
  updatedAt: string
  triggerIconUrl?: string
  actionIconUrl?: string
  triggerConfig?: any
  actionConfig?: any
}

export interface DiscordLog {
  id: number
  areaId: string
  filePath?: string | null
  sheetName?: string | null
  changeType?: string | null
  rowNumber?: number | null
  message: string
  rowData?: Record<string, any> | null
  createdAt: string
}

export interface GoogleSheetsTestResponse {
  message: string
  rowCount: number
  previewRows: string[][]
}

export interface AreaTemplate {
  id: string
  title: string
  subtitle: string
  description: string
  icon: string
  gradientClass: string
  triggerService: string
  actionService: string
  triggerIconUrl?: string
  actionIconUrl?: string
  isActive: boolean
}

export interface CreateAreaRequest {
  name: string
  description: string
  triggerService: string
  triggerType: string
  actionService: string
  actionType: string
  triggerConfig?: any
  actionConfig?: any
}

class AreaService {
  private baseURL = `${API_BASE_URL}/areas`

  async getAreas(): Promise<Area[]> {
    try {
      const response = await fetch(this.baseURL, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch areas: ${response.statusText}`)
      }

      const data = await response.json()
      return data.data || []
    } catch (error) {
      console.error('Error fetching areas:', error)
      throw error
    }
  }

  async getAreaById(id: string): Promise<Area> {
    try {
      const token = localStorage.getItem('authToken')

      const response = await fetch(`${API_BASE_URL}/user/me/areas`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch areas: ${response.statusText}`)
      }

      const data = await response.json()
      const area = data.data.find((area: any) => area.id === id)

      if (!area) {
        throw new Error('Area not found')
      }

      return transformAreaData(area)
    } catch (error) {
      console.error('Error fetching area:', error)
      throw error
    }
  }

  async getPopularAreas(): Promise<AreaTemplate[]> {
    try {
      const response = await fetch(`${this.baseURL}/popular`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch popular areas: ${response.statusText}`)
      }

      const data = await response.json()
      return data.data || []
    } catch (error) {
      console.error('Error fetching popular areas:', error)
      throw error
    }
  }

  async getRecommendedAreas(): Promise<AreaTemplate[]> {
    try {
      const response = await fetch(`${this.baseURL}/recommended`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch recommended areas: ${response.statusText}`)
      }

      const data = await response.json()
      return data.data || []
    } catch (error) {
      console.error('Error fetching recommended areas:', error)
      throw error
    }
  }

  async getUserAreas(userId: number): Promise<Area[]> {
    try {
      const token = localStorage.getItem('authToken')

      const response = await fetch(`${API_BASE_URL}/user/me/areas`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch user areas: ${response.statusText}`)
      }

      const data = await response.json()
      return (data.data || []).map((area: any) => transformAreaData(area))
    } catch (error) {
      console.error('Error fetching user areas:', error)
      throw error
    }
  }

  async createArea(areaData: CreateAreaRequest): Promise<Area> {
    try {
      const token = localStorage.getItem('authToken')

      const response = await fetch(this.baseURL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(areaData),
      })

      if (!response.ok) {
        throw new Error(`Failed to create area: ${response.statusText}`)
      }

      const data = await response.json()
      return data.data
    } catch (error) {
      console.error('Error creating area:', error)
      throw error
    }
  }

  async updateArea(id: string, areaData: Partial<CreateAreaRequest>): Promise<Area> {
    try {
      const token = localStorage.getItem('authToken')

      const response = await fetch(`${this.baseURL}/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify(areaData),
      })

      if (!response.ok) {
        throw new Error(`Failed to update area: ${response.statusText}`)
      }

      const data = await response.json()
      return data.data
    } catch (error) {
      console.error('Error updating area:', error)
      throw error
    }
  }

  async deleteArea(id: string): Promise<void> {
    try {
      const token = localStorage.getItem('authToken')

      if (!token) {
        throw new Error('No authentication token found')
      }

      const response = await fetch(`${this.baseURL}/${id}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        const errorText = await response.text()
        throw new Error(`Failed to delete area: ${response.status} ${response.statusText}`)
      }
    } catch (error) {
      console.error('Error deleting area:', error)
      throw error
    }
  }

  async toggleArea(id: string): Promise<Area> {
    try {
      const token = localStorage.getItem('authToken')

      const response = await fetch(`${this.baseURL}/${id}/toggle`, {
        method: 'PATCH',
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to toggle area: ${response.statusText}`)
      }

      const data = await response.json()
      return data.data
    } catch (error) {
      console.error('Error toggling area:', error)
      throw error
    }
  }

  async getDiscordLogs(id: string, limit = 50): Promise<DiscordLog[]> {
    try {
      const token = localStorage.getItem('authToken')

      const response = await fetch(`${this.baseURL}/${id}/discord-logs?limit=${limit}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch Discord logs: ${response.statusText}`)
      }

      const data = await response.json()
      const logs = data.data || []
      return logs.map((log: any) => transformDiscordLog(log))
    } catch (error) {
      console.error('Error fetching Discord logs:', error)
      throw error
    }
  }

  async testGoogleSheets(config: { spreadsheetId: string; range: string }): Promise<GoogleSheetsTestResponse> {
    try {
      const token = localStorage.getItem('authToken')
      const headers: Record<string, string> = {
        'Content-Type': 'application/json',
      }

      if (token) {
        headers['Authorization'] = `Bearer ${token}`
      }

      const response = await fetch(`${API_BASE_URL}/test/google-sheets`, {
        method: 'POST',
        headers,
        body: JSON.stringify(config),
      })

      if (!response.ok) {
        const bodyText = await response.text()
        let message = `Failed to test Google Sheets: ${response.statusText}`

        if (bodyText) {
          try {
            const parsed = JSON.parse(bodyText)
            message = parsed.error || parsed.message || message
          } catch {
            message = bodyText
          }
        }

        throw new Error(message)
      }

      return await response.json()
    } catch (error) {
      console.error('Error testing Google Sheets trigger:', error)
      throw error
    }
  }
}

export const areaService = new AreaService()

const transformDiscordLog = (log: any): DiscordLog => {
  const rawRowData = log.row_data || log.rowData || null
  let rowData: Record<string, any> | null = rawRowData

  if (rawRowData && typeof rawRowData === 'string') {
    try {
      rowData = JSON.parse(rawRowData)
    } catch (error) {
      console.warn('Failed to parse row data JSON', error)
      rowData = null
    }
  }

  return {
    id: log.id,
    areaId: log.area_id || log.areaId,
    filePath: log.file_path || log.filePath || null,
    sheetName: log.sheet_name || log.sheetName || null,
    changeType: log.change_type || log.changeType || null,
    rowNumber: log.row_number !== undefined ? log.row_number : log.rowNumber,
    message: log.message,
    rowData: rowData,
    createdAt: log.created_at || log.createdAt,
  }
}
