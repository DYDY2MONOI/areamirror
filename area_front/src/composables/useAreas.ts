import { ref, computed } from 'vue'
import { areaService, type Area, type AreaTemplate } from '@/services/area'

const areas = ref<Area[]>([])
const popularAreas = ref<AreaTemplate[]>([])
const recommendedAreas = ref<AreaTemplate[]>([])
const isLoading = ref(false)
const error = ref<string | null>(null)
const togglingAreaIds = ref<string[]>([])

export function useAreas() {
  const fetchAreas = async () => {
    isLoading.value = true
    error.value = null

    try {
      const data = await areaService.getAreas()
      areas.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch areas'
      console.error('Error fetching areas:', err)
    } finally {
      isLoading.value = false
    }
  }

  const fetchPopularAreas = async () => {
    isLoading.value = true
    error.value = null

    try {
      const data = await areaService.getPopularAreas()
      popularAreas.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch popular areas'
      console.error('Error fetching popular areas:', err)
    } finally {
      isLoading.value = false
    }
  }

  const fetchRecommendedAreas = async () => {
    isLoading.value = true
    error.value = null

    try {
      const data = await areaService.getRecommendedAreas()
      recommendedAreas.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch recommended areas'
      console.error('Error fetching recommended areas:', err)
    } finally {
      isLoading.value = false
    }
  }

  const fetchUserAreas = async (userId: number) => {
    isLoading.value = true
    error.value = null

    try {
      const data = await areaService.getUserAreas(userId)
      areas.value = data
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch user areas'
      console.error('Error fetching user areas:', err)
    } finally {
      isLoading.value = false
    }
  }

  const createArea = async (areaData: {
    name: string
    description: string
    triggerService: string
    triggerType: string
    actionService: string
    actionType: string
  }) => {
    isLoading.value = true
    error.value = null

    try {
      const newArea = await areaService.createArea(areaData)
      areas.value.push(newArea)
      return newArea
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create area'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const updateArea = async (id: string, areaData: Partial<{
    name: string
    description: string
    triggerService: string
    triggerType: string
    actionService: string
    actionType: string
  }>) => {
    isLoading.value = true
    error.value = null

    try {
      const updatedArea = await areaService.updateArea(id, areaData)
      const index = areas.value.findIndex(area => area.id === id)
      if (index !== -1) {
        areas.value[index] = updatedArea
      }
      return updatedArea
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update area'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const deleteArea = async (id: string) => {
    isLoading.value = true
    error.value = null

    try {
      await areaService.deleteArea(id)
      areas.value = areas.value.filter(area => area.id !== id)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete area'
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const toggleArea = async (id: string) => {
    if (togglingAreaIds.value.includes(id)) {
      return areas.value.find(area => area.id === id)
    }

    togglingAreaIds.value = [...togglingAreaIds.value, id]
    error.value = null

    try {
      const updatedArea = await areaService.toggleArea(id)
      const index = areas.value.findIndex(area => area.id === id)
      if (index !== -1) {
        areas.value[index] = updatedArea
      }
      return updatedArea
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to toggle area'
      throw err
    } finally {
      togglingAreaIds.value = togglingAreaIds.value.filter(areaId => areaId !== id)
    }
  }

  const getAreaById = (id: string) => {
    return areas.value.find(area => area.id === id)
  }

  const getActiveAreas = computed(() => {
    return areas.value.filter(area => area.isActive)
  })

  const getPublicAreas = computed(() => {
    return areas.value.filter(area => area.isPublic)
  })

  return {
    areas: computed(() => areas.value),
    popularAreas: computed(() => popularAreas.value),
    recommendedAreas: computed(() => recommendedAreas.value),
    isLoading: computed(() => isLoading.value),
    error: computed(() => error.value),
    togglingAreaIds: computed(() => togglingAreaIds.value),
    activeAreas: getActiveAreas,
    publicAreas: getPublicAreas,

    fetchAreas,
    fetchPopularAreas,
    fetchRecommendedAreas,
    fetchUserAreas,
    createArea,
    updateArea,
    deleteArea,
    toggleArea,
    getAreaById
  }
}
