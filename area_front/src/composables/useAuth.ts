import { ref, computed } from 'vue'
import { authService, type User } from '@/services/auth'

const isAuthenticated = ref(false)
const currentUser = ref<User | null>(null)
const isLoading = ref(false)

const initAuth = async () => {
  isLoading.value = true
  try {
    const isAuth = await authService.checkAuthStatus()
    isAuthenticated.value = isAuth
    currentUser.value = authService.currentUser
  } catch (error) {
    console.error('Erreur lors de l\'initialisation de l\'authentification:', error)
    isAuthenticated.value = false
    currentUser.value = null
  } finally {
    isLoading.value = false
  }
}

export function useAuth() {
  const login = async (email: string, password: string) => {
    isLoading.value = true
    try {
      const response = await authService.login({ email, password })
      isAuthenticated.value = authService.isAuthenticated
      currentUser.value = authService.currentUser
      return response
    } catch (error) {
      throw error
    } finally {
      isLoading.value = false
    }
  }

  const register = async (userData: { email: string; password: string; first_name?: string; last_name?: string }) => {
    console.log('🔄 Composable: Début de l\'enregistrement')
    isLoading.value = true
    try {
      console.log('🔄 Composable: Appel du service authService.register')
      const response = await authService.register(userData)
      console.log('✅ Composable: Service terminé, mise à jour de l\'état')
      isAuthenticated.value = authService.isAuthenticated
      currentUser.value = authService.currentUser
      console.log('✅ Composable: État mis à jour', { isAuthenticated: isAuthenticated.value, currentUser: currentUser.value })
      return response
    } catch (error) {
      console.error('❌ Composable: Erreur capturée', error)
      throw error
    } finally {
      console.log('🏁 Composable: Fin du processus')
      isLoading.value = false
    }
  }

  const logout = async () => {
    isLoading.value = true
    try {
      await authService.logout()
      isAuthenticated.value = authService.isAuthenticated
      currentUser.value = authService.currentUser
    } catch (error) {
      console.error('Erreur lors de la déconnexion:', error)
    } finally {
      isLoading.value = false
    }
  }

  const refreshProfile = async () => {
    try {
      const isAuth = await authService.checkAuthStatus()
      isAuthenticated.value = isAuth
      currentUser.value = authService.currentUser
      return authService.currentUser
    } catch (error) {
      console.error('Erreur lors du rafraîchissement du profil:', error)
      isAuthenticated.value = false
      currentUser.value = null
      return null
    }
  }

  const linkGitHubAccount = async (code: string) => {
    try {
      const result = await authService.linkGitHubAccount(code)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('GitHub link error:', error)
      throw error
    }
  }

  const unlinkGitHubAccount = async () => {
    try {
      await authService.unlinkGitHubAccount()
      await refreshProfile()
    } catch (error) {
      console.error('GitHub unlink error:', error)
      throw error
    }
  }

  const linkGoogleAccount = async (code: string) => {
    try {
      const result = await authService.linkGoogleAccount(code)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('Google link error:', error)
      throw error
    }
  }

  const unlinkGoogleAccount = async () => {
    try {
      await authService.unlinkGoogleAccount()
      await refreshProfile()
    } catch (error) {
      console.error('Google unlink error:', error)
      throw error
    }
  }

  const uploadProfileImage = async (imageFile: File) => {
    try {
      const result = await authService.uploadProfileImage(imageFile)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('Profile image upload error:', error)
      throw error
    }
  }

  const getProfileImageUrl = () => {
    return authService.getProfileImageUrl()
  }

  const updateProfile = async (updateData: any) => {
    try {
      const result = await authService.updateProfile(updateData)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('Profile update error:', error)
      throw error
    }
  }

  return {
    isAuthenticated: computed(() => isAuthenticated.value),
    currentUser: computed(() => currentUser.value),
    isLoading: computed(() => isLoading.value),

    login,
    register,
    logout,
    refreshProfile,
    linkGitHubAccount,
    unlinkGitHubAccount,
    linkGoogleAccount,
    unlinkGoogleAccount,
    uploadProfileImage,
    getProfileImageUrl,
    updateProfile,
    initAuth
  }
}

initAuth()
