import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '../views/LandingPage.vue'
import LoginPage from '../views/LoginPage.vue'
import RegisterPage from '../views/RegisterPage.vue'
import ProfilePage from '../views/ProfilePage.vue'
import GitHubCallbackPage from '../views/GitHubCallbackPage.vue'
import GoogleCallbackPage from '../views/GoogleCallbackPage.vue'
import GoogleOAuthCallback from '../views/GoogleOAuthCallback.vue'
import FacebookCallbackPage from '../views/FacebookCallbackPage.vue'
import OneDriveCallbackPage from '../views/OneDriveCallbackPage.vue'
import HomeCallback from '../views/HomeCallback.vue'
import EditProfilePage from '../views/EditProfilePage.vue'
import ConfigureAreaPage from '../views/ConfigureAreaPage.vue'
import AllAreasPage from '../views/AllAreasPage.vue'
import { authService } from '../services/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: LandingPage,
      meta: { requiresAuth: false }
    },
    {
      path: '/callback',
      name: 'home-callback',
      component: HomeCallback,
      meta: { requiresAuth: false }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginPage,
      meta: { requiresAuth: false }
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterPage,
      meta: { requiresAuth: false }
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfilePage,
      meta: { requiresAuth: true }
    },
    {
      path: '/auth/github/callback',
      name: 'github-callback',
      component: GitHubCallbackPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/auth/google/callback',
      name: 'google-callback',
      component: GoogleCallbackPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/auth/facebook/callback',
      name: 'facebook-callback',
      component: FacebookCallbackPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/auth/onedrive/callback',
      name: 'onedrive-callback',
      component: OneDriveCallbackPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/google-oauth-callback',
      name: 'google-oauth-callback',
      component: GoogleOAuthCallback,
      meta: { requiresAuth: false }
    },
    {
      path: '/oauth2/github/callback',
      name: 'oauth2-github-callback',
      component: () => import('@/views/OAuth2GitHubCallback.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/oauth2/google/callback',
      name: 'oauth2-google-callback',
      component: () => import('@/views/OAuth2GoogleCallback.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/oauth2/facebook/callback',
      name: 'oauth2-facebook-callback',
      component: () => import('@/views/OAuth2FacebookCallback.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/profile/edit',
      name: 'edit-profile',
      component: EditProfilePage,
      meta: { requiresAuth: true }
    },
    {
      path: '/configure-area',
      name: 'configure-area',
      component: ConfigureAreaPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/areas',
      name: 'all-areas',
      component: AllAreasPage,
      meta: { requiresAuth: true }
    },
  ],
})

router.beforeEach(async (to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)

  if (requiresAuth) {
    const isAuthenticated = await authService.checkAuthStatus()

    if (!isAuthenticated) {
      next({ name: 'login' })
    } else {
      next()
    }
  } else {
    if ((to.name === 'login' || to.name === 'register')) {
      const isAuthenticated = await authService.checkAuthStatus()

      if (isAuthenticated) {
        next({ name: 'landing' })
      } else {
        next()
      }
    } else {
      next()
    }
  }
})

export default router


