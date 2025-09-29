import { createRouter, createWebHistory } from 'vue-router'
import LandingPage from '../views/LandingPage.vue'
import LoginPage from '../views/LoginPage.vue'
import RegisterPage from '../views/RegisterPage.vue'
import ProfilePage from '../views/ProfilePage.vue'
import ServiceLinkPage from '../views/ServiceLinkPage.vue'
import GitHubCallbackPage from '../views/GitHubCallbackPage.vue'
import EditProfilePage from '../views/EditProfilePage.vue'
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
      path: '/services',
      name: 'services',
      component: ServiceLinkPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/github-link',
      name: 'github-link',
      component: ServiceLinkPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/auth/github/callback',
      name: 'github-callback',
      component: GitHubCallbackPage,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile/edit',
      name: 'edit-profile',
      component: EditProfilePage,
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


