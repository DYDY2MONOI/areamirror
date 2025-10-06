<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'

const currentView = ref<'login' | 'register'>('login')
const isAnimating = ref(false)

const loginData = reactive({
  email: '',
  password: '',
  rememberMe: false
})

const registerData = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  confirmPassword: '',
  acceptTerms: false
})

const showPassword = ref(false)
const showConfirmPassword = ref(false)
const isLoading = ref(false)
const isPageLoaded = ref(false)

onMounted(() => {
  setTimeout(() => {
    isPageLoaded.value = true
  }, 100)
})

const switchView = (view: 'login' | 'register') => {
  if (isAnimating.value || currentView.value === view) return

  isAnimating.value = true
  setTimeout(() => {
    currentView.value = view
    setTimeout(() => {
      isAnimating.value = false
    }, 50)
  }, 300)
}

const handleLogin = async () => {
  isLoading.value = true
  console.log('Login attempt:', loginData)

  setTimeout(() => {
    isLoading.value = false
    alert('Login functionality would be implemented here!')
  }, 2000)
}

const handleRegister = async () => {
  if (registerData.password !== registerData.confirmPassword) {
    alert('Passwords do not match!')
    return
  }

  isLoading.value = true
  console.log('Register attempt:', registerData)

  setTimeout(() => {
    isLoading.value = false
    alert('Registration functionality would be implemented here!')
  }, 2000)
}
</script>

<template>
  <div class="min-h-screen bg-slate-900 relative overflow-hidden">
    <div class="geometric-shape top-0 left-0 w-96 h-96 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full transform -translate-x-48 -translate-y-48"></div>
    <div class="geometric-shape bottom-0 right-0 w-80 h-80 bg-gradient-to-tl from-cyan-500 to-purple-600 rounded-full transform translate-x-40 translate-y-40"></div>
    <div class="geometric-shape top-1/2 left-1/4 w-64 h-64 bg-gradient-to-r from-purple-600 to-blue-500 rounded-full transform -translate-y-32 opacity-10"></div>

    <div class="relative z-10 flex items-center justify-center min-h-screen px-4">
      <div class="w-full max-w-md perspective-container">
        <div class="auth-container" :class="{ 'slide-out': isAnimating }">

          <div v-if="currentView === 'login'" class="auth-panel" :class="{ 'slide-in': !isAnimating }">
            <div class="text-center mb-8">
              <h1 class="text-4xl font-bold text-white mb-2">AREA.</h1>
              <h2 class="text-2xl font-semibold text-white mb-2">LOGIN TO YOUR ACCOUNT</h2>
              <p class="text-gray-400">Enter your login information.</p>
            </div>

            <div class="form-card">
              <form @submit.prevent="handleLogin" class="space-y-6">
                <div class="form-group">
                  <label class="form-label">
                    <svg class="w-4 h-4 inline-block mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"/>
                    </svg>
                    Email Address
                  </label>
                  <div class="input-wrapper">
                    <input
                      type="email"
                      v-model="loginData.email"
                      class="input-field-modern"
                      placeholder="oliverross@gmail.com"
                      required
                    />
                    <div class="input-focus-effect"></div>
                  </div>
                </div>

                <div class="form-group">
                  <label class="form-label">
                    <svg class="w-4 h-4 inline-block mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                    </svg>
                    Password
                  </label>
                  <div class="input-wrapper">
                    <input
                      :type="showPassword ? 'text' : 'password'"
                      v-model="loginData.password"
                      class="input-field-modern pr-12"
                      placeholder="Enter your password"
                      required
                    />
                    <button
                      type="button"
                      @click="showPassword = !showPassword"
                      class="password-toggle"
                    >
                      <svg v-if="showPassword" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                      </svg>
                      <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                      </svg>
                    </button>
                    <div class="input-focus-effect"></div>
                  </div>
                </div>

                <div class="flex items-center justify-between px-1">
                  <label class="custom-checkbox">
                    <input
                      type="checkbox"
                      v-model="loginData.rememberMe"
                      class="checkbox-input"
                    />
                    <span class="checkbox-mark"></span>
                    <span class="checkbox-label">Remember me</span>
                  </label>
                  <a href="#" class="forgot-link">
                    Forgot password?
                  </a>
                </div>

                <button type="submit" class="btn-primary-modern">
                  <span v-if="isLoading" class="flex items-center justify-center">
                    <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Logging in...
                  </span>
                  <span v-else class="btn-text">
                    LOGIN
                    <svg class="w-5 h-5 ml-2 inline-block group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"/>
                    </svg>
                  </span>
                </button>

                <div class="divider-container">
                  <div class="divider-line"></div>
                  <span class="divider-text">OR</span>
                  <div class="divider-line"></div>
                </div>

                <div class="social-buttons">
                  <button type="button" class="social-btn google-btn">
                    <svg class="h-5 w-5" viewBox="0 0 24 24">
                      <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                      <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                      <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                      <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                    </svg>
                    <span>Google</span>
                  </button>
                  <button type="button" class="social-btn apple-btn">
                    <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
                    </svg>
                    <span>Apple</span>
                  </button>
                </div>
              </form>

              <div class="mt-6 text-center">
                <p class="text-gray-400">
                  Don't have an account?
                  <button
                    @click="switchView('register')"
                    class="text-blue-500 hover:text-blue-400 font-medium transition-colors"
                  >
                    Sign Up
                  </button>
                </p>
              </div>
            </div>
          </div>

          <div v-if="currentView === 'register'" class="auth-panel" :class="{ 'slide-in': !isAnimating }">
            <div class="text-center mb-8">
              <h1 class="text-4xl font-bold text-white mb-2">AREA.</h1>
              <h2 class="text-2xl font-semibold text-white mb-2">CREATE ACCOUNT</h2>
              <p class="text-gray-400">Join us and start automating!</p>
            </div>

            <div class="form-card">
              <form @submit.prevent="handleRegister" class="space-y-6">
                <div class="grid grid-cols-2 gap-4">
                  <div class="form-group">
                    <label class="form-label">
                      <svg class="w-4 h-4 inline-block mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                      </svg>
                      First Name
                    </label>
                    <div class="input-wrapper">
                      <input
                        type="text"
                        v-model="registerData.firstName"
                        class="input-field-modern"
                        placeholder="John"
                        required
                      />
                      <div class="input-focus-effect"></div>
                    </div>
                  </div>
                  <div class="form-group">
                    <label class="form-label">
                      <svg class="w-4 h-4 inline-block mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                      </svg>
                      Last Name
                    </label>
                    <div class="input-wrapper">
                      <input
                        type="text"
                        v-model="registerData.lastName"
                        class="input-field-modern"
                        placeholder="Doe"
                        required
                      />
                      <div class="input-focus-effect"></div>
                    </div>
                  </div>
                </div>

                <div class="form-group">
                  <label class="form-label">
                    <svg class="w-4 h-4 inline-block mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"/>
                    </svg>
                    Email Address
                  </label>
                  <div class="input-wrapper">
                    <input
                      type="email"
                      v-model="registerData.email"
                      class="input-field-modern"
                      placeholder="john.doe@gmail.com"
                      required
                    />
                    <div class="input-focus-effect"></div>
                  </div>
                </div>

                <div class="form-group">
                  <label class="form-label">
                    <svg class="w-4 h-4 inline-block mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                    </svg>
                    Password
                  </label>
                  <div class="input-wrapper">
                    <input
                      :type="showPassword ? 'text' : 'password'"
                      v-model="registerData.password"
                      class="input-field-modern pr-12"
                      placeholder="Create a strong password"
                      required
                    />
                    <button
                      type="button"
                      @click="showPassword = !showPassword"
                      class="password-toggle"
                    >
                      <svg v-if="showPassword" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                      </svg>
                      <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                      </svg>
                    </button>
                    <div class="input-focus-effect"></div>
                  </div>
                </div>

                <div class="form-group">
                  <label class="form-label">
                    <svg class="w-4 h-4 inline-block mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    Confirm Password
                  </label>
                  <div class="input-wrapper">
                    <input
                      :type="showConfirmPassword ? 'text' : 'password'"
                      v-model="registerData.confirmPassword"
                      class="input-field-modern pr-12"
                      placeholder="Re-enter your password"
                      required
                    />
                    <button
                      type="button"
                      @click="showConfirmPassword = !showConfirmPassword"
                      class="password-toggle"
                    >
                      <svg v-if="showConfirmPassword" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                      </svg>
                      <svg v-else class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                      </svg>
                    </button>
                    <div class="input-focus-effect"></div>
                  </div>
                </div>

                <div class="terms-container">
                  <label class="custom-checkbox">
                    <input
                      type="checkbox"
                      v-model="registerData.acceptTerms"
                      class="checkbox-input"
                      required
                    />
                    <span class="checkbox-mark"></span>
                    <span class="checkbox-label">
                      I agree to the
                      <a href="#" class="terms-link">Terms of Service</a>
                      and
                      <a href="#" class="terms-link">Privacy Policy</a>
                    </span>
                  </label>
                </div>

                <button type="submit" class="btn-primary-modern">
                  <span v-if="isLoading" class="flex items-center justify-center">
                    <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Creating Account...
                  </span>
                  <span v-else class="btn-text">
                    CREATE ACCOUNT
                    <svg class="w-5 h-5 ml-2 inline-block group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"/>
                    </svg>
                  </span>
                </button>

                <div class="divider-container">
                  <div class="divider-line"></div>
                  <span class="divider-text">OR SIGN UP WITH</span>
                  <div class="divider-line"></div>
                </div>

                <div class="social-buttons">
                  <button type="button" class="social-btn google-btn">
                    <svg class="h-5 w-5" viewBox="0 0 24 24">
                      <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                      <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                      <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                      <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                    </svg>
                    <span>Google</span>
                  </button>
                  <button type="button" class="social-btn apple-btn">
                    <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
                    </svg>
                    <span>Apple</span>
                  </button>
                </div>
              </form>

              <div class="mt-6 text-center">
                <p class="text-gray-400">
                  Already have an account?
                  <button
                    @click="switchView('login')"
                    class="text-blue-500 hover:text-blue-400 font-medium transition-colors"
                  >
                    Sign In
                  </button>
                </p>
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.perspective-container {
  perspective: 1000px;
}

.auth-container {
  position: relative;
  transition: opacity 0.3s ease;
}

.auth-panel {
  animation: slideIn 0.5s ease-out;
}

.slide-out {
  opacity: 0.3;
  animation: slideOut 0.3s ease-out;
}

.slide-in {
  animation: slideIn 0.5s ease-out;
}

@keyframes slideOut {
  0% {
    transform: translateX(0) scale(1);
    opacity: 1;
  }
  100% {
    transform: translateX(-50px) scale(0.95);
    opacity: 0;
  }
}

@keyframes slideIn {
  0% {
    transform: translateX(50px) scale(0.95);
    opacity: 0;
  }
  100% {
    transform: translateX(0) scale(1);
    opacity: 1;
  }
}
</style>
