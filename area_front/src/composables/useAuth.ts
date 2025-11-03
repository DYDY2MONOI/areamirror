import { ref, computed } from 'vue'
import { authService, type User } from '@/services/auth'
import { oauth2AuthService, type OAuth2User } from '@/services/oauth2-auth'

const isAuthenticated = ref(false)
const currentUser = ref<User | null>(null)
const isLoading = ref(false)

const initAuth = async () => {
  isLoading.value = true
  try {
    if (oauth2AuthService.isAuthenticated) {
      const isAuth = await oauth2AuthService.checkAuthStatus()
      isAuthenticated.value = isAuth
      if (oauth2AuthService.currentUser) {
        const oauthUser = oauth2AuthService.currentUser
        currentUser.value = {
          id: oauthUser.id,
          email: oauthUser.email,
          first_name: oauthUser.first_name,
          last_name: oauthUser.last_name,
          created_at: oauthUser.created_at,
          updated_at: oauthUser.updated_at,
          phone: oauthUser.phone,
          birthday: oauthUser.birthday,
          gender: oauthUser.gender,
          country: oauthUser.country,
          lang: oauthUser.lang,
          login_provider: oauthUser.login_provider,
          role: oauthUser.role || null,
          is_active: oauthUser.is_active || false,
          github_id: oauthUser.github_id,
          github_username: oauthUser.github_username,
          google_id: oauthUser.google_id,
          google_email: oauthUser.google_email,
          facebook_id: oauthUser.facebook_id,
          facebook_email: oauthUser.facebook_email,
          discord_id: oauthUser.discord_id,
          discord_username: oauthUser.discord_username,
          spotify_id: oauthUser.spotify_id,
          spotify_email: oauthUser.spotify_email,
          twitter_id: oauthUser.twitter_id,
          twitter_username: oauthUser.twitter_username,
          slack_id: oauthUser.slack_id,
          slack_team_id: oauthUser.slack_team_id,
          onedrive_id: oauthUser.onedrive_id,
          profile_image: oauthUser.profile_image,
        }
      }
    } else {
      const isAuth = await authService.checkAuthStatus()
      isAuthenticated.value = isAuth
      currentUser.value = authService.currentUser
    }
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
    console.log(' Composable: Début de l\'enregistrement')
    isLoading.value = true
    try {
      console.log(' Composable: Appel du service authService.register')
      const response = await authService.register(userData)
      console.log(' Composable: Service terminé, mise à jour de l\'état')
      isAuthenticated.value = authService.isAuthenticated
      currentUser.value = authService.currentUser
      console.log(' Composable: État mis à jour', { isAuthenticated: isAuthenticated.value, currentUser: currentUser.value })
      return response
    } catch (error) {
      console.error(' Composable: Erreur capturée', error)
      throw error
    } finally {
      console.log(' Composable: Fin du processus')
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
      if (oauth2AuthService.isAuthenticated) {
        console.log(' Using OAuth2 auth service for profile refresh')
        const isAuth = await oauth2AuthService.checkAuthStatus()
        isAuthenticated.value = isAuth
        if (oauth2AuthService.currentUser) {
          const oauthUser = oauth2AuthService.currentUser
          currentUser.value = {
            id: oauthUser.id,
            email: oauthUser.email,
            first_name: oauthUser.first_name,
            last_name: oauthUser.last_name,
            created_at: oauthUser.created_at,
            updated_at: oauthUser.updated_at,
            phone: oauthUser.phone,
            birthday: oauthUser.birthday,
            gender: oauthUser.gender,
            country: oauthUser.country,
            lang: oauthUser.lang,
            login_provider: oauthUser.login_provider,
            role: oauthUser.role || null,
            is_active: oauthUser.is_active || false,
            github_id: oauthUser.github_id,
            github_username: oauthUser.github_username,
            google_id: oauthUser.google_id,
            google_email: oauthUser.google_email,
            facebook_id: oauthUser.facebook_id,
            facebook_email: oauthUser.facebook_email,
            discord_id: oauthUser.discord_id,
            discord_username: oauthUser.discord_username,
            spotify_id: oauthUser.spotify_id,
            spotify_email: oauthUser.spotify_email,
            twitter_id: oauthUser.twitter_id,
            twitter_username: oauthUser.twitter_username,
            slack_id: oauthUser.slack_id,
            slack_team_id: oauthUser.slack_team_id,
            onedrive_id: oauthUser.onedrive_id,
            profile_image: oauthUser.profile_image,
          }
        }
        return currentUser.value
      }

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

  const linkGitHubAccount = async (code: string, redirectUri?: string) => {
    try {
      const result = await authService.linkGitHubAccount(code, redirectUri)
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

  const linkGoogleAccount = async (code: string, redirectUri?: string) => {
    try {
      const result = await authService.linkGoogleAccount(code, redirectUri)
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

  const linkFacebookAccount = async (code: string) => {
    try {
      const result = await authService.linkFacebookAccount(code)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('Facebook link error:', error)
      throw error
    }
  }

  const unlinkFacebookAccount = async () => {
    try {
      await authService.unlinkFacebookAccount()
      await refreshProfile()
    } catch (error) {
      console.error('Facebook unlink error:', error)
      throw error
    }
  }

  const linkOneDriveAccount = async (code: string) => {
    try {
      const result = await authService.linkOneDriveAccount(code)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('OneDrive link error:', error)
      throw error
    }
  }

  const linkSpotifyAccount = async (code: string) => {
    try {
      const result = await authService.linkSpotifyAccount(code)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('Spotify link error:', error)
      throw error
    }
  }

  const unlinkOneDriveAccount = async () => {
    try {
      await authService.unlinkOneDriveAccount()
      await refreshProfile()
    } catch (error) {
      console.error('OneDrive unlink error:', error)
      throw error
    }
  }

  const unlinkSpotifyAccount = async () => {
    try {
      await authService.unlinkSpotifyAccount()
      await refreshProfile()
    } catch (error) {
      console.error('Spotify unlink error:', error)
      throw error
    }
  }

  const linkTwitterAccount = async (code: string, codeVerifier: string) => {
    try {
      const result = await authService.linkTwitterAccount(code, codeVerifier)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('Twitter link error:', error)
      throw error
    }
  }

  const unlinkTwitterAccount = async () => {
    try {
      await authService.unlinkTwitterAccount()
      await refreshProfile()
    } catch (error) {
      console.error('Twitter unlink error:', error)
      throw error
    }
  }

  const linkSlackAccount = async (code: string) => {
    try {
      const result = await authService.linkSlackAccount(code)
      await refreshProfile()
      return result
    } catch (error) {
      console.error('Slack link error:', error)
      throw error
    }
  }

  const unlinkSlackAccount = async () => {
    try {
      await authService.unlinkSlackAccount()
      await refreshProfile()
    } catch (error) {
      console.error('Slack unlink error:', error)
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
    isAuthenticated: computed(() => isAuthenticated.value || oauth2AuthService.isAuthenticated),
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
    linkFacebookAccount,
    unlinkFacebookAccount,
    linkOneDriveAccount,
    unlinkOneDriveAccount,
    linkSpotifyAccount,
    unlinkSpotifyAccount,
    linkTwitterAccount,
    unlinkTwitterAccount,
    linkSlackAccount,
    unlinkSlackAccount,
    uploadProfileImage,
    getProfileImageUrl,
    updateProfile,
    initAuth
  }
}

initAuth()
