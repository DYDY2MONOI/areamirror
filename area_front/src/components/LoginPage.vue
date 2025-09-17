<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'

const formData = reactive({
  email: '',
  password: '',
  rememberMe: false
})

const showPassword = ref(false)
const isLoading = ref(false)
const isPageLoaded = ref(false)

onMounted(() => {
  setTimeout(() => {
    isPageLoaded.value = true
  }, 100)
})

const handleLogin = async () => {
  isLoading.value = true
  console.log('Login attempt:', formData)

  setTimeout(() => {
    isLoading.value = false
    alert('Login functionality would be implemented here!')
  }, 2000)
}
</script>

<template>
  <div class="min-h-screen bg-slate-900 relative overflow-hidden">
    <div class="geometric-shape top-0 left-0 w-96 h-96 bg-gradient-to-br from-blue-500 to-purple-600 rounded-full transform -translate-x-48 -translate-y-48"></div>
    <div class="geometric-shape bottom-0 right-0 w-80 h-80 bg-gradient-to-tl from-cyan-500 to-purple-600 rounded-full transform translate-x-40 translate-y-40"></div>
    <div class="geometric-shape top-1/2 left-1/4 w-64 h-64 bg-gradient-to-r from-purple-600 to-blue-500 rounded-full transform -translate-y-32 opacity-10"></div>
    <div class="relative z-10 flex items-center justify-center min-h-screen px-4">
      <div class="w-full max-w-md">
        <div class="text-center mb-8" :class="{ 'animate-fade-in-down': isPageLoaded }">
          <h1 class="text-4xl font-bold text-white mb-2" :class="{ 'logo-animate': isPageLoaded }">AREA.</h1>
          <h2 class="text-2xl font-semibold text-white mb-2" :class="{ 'animate-fade-in-up delay-200': isPageLoaded }">LOGIN TO YOUR ACCOUNT</h2>
          <p class="text-gray-400" :class="{ 'animate-fade-in-up delay-300': isPageLoaded }">Enter your login information.</p>
        </div>
        <div class="bg-slate-800/50 backdrop-blur-sm rounded-2xl p-8 border border-slate-600/50" :class="{ 'animate-fade-in-up delay-400': isPageLoaded }">
          <form @submit.prevent="handleLogin" class="space-y-6">
            <div :class="{ 'animate-fade-in-left delay-500': isPageLoaded }">
              <label for="email" class="block text-sm font-medium text-gray-300 mb-2">
                Email Address
              </label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                </div>
                <input
                  type="email"
                  v-model="formData.email"
                  class="input-field pl-10"
                  placeholder="oliverross@gmail.com"
                  required
                />
              </div>
            </div>

            <div :class="{ 'animate-fade-in-right delay-600': isPageLoaded }">
              <label for="password" class="block text-sm font-medium text-gray-300 mb-1">
                Password
              </label>
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                </div>
                <input
                  :type="showPassword ? 'text' : 'password'"
                  v-model="formData.password"
                  class="input-field pl-10 pr-10"
                  placeholder="YourCoolPassword123"
                  required
                />
                <button
                  type="button"
                  @click="showPassword = !showPassword"
                  class="absolute inset-y-0 right-0 pr-3 flex items-center"
                >
                  <svg v-if="showPassword" class="h-5 w-5 text-gray-400 hover:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                  <svg v-else class="h-5 w-5 text-gray-400 hover:text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"></path>
                  </svg>
                </button>
              </div>
            </div>

            <div class="flex items-center justify-between" :class="{ 'animate-fade-in-up delay-700': isPageLoaded }">
              <label class="flex items-center">
                <input
                  type="checkbox"
                  v-model="formData.rememberMe"
                  class="h-4 w-4 text-blue-500 bg-slate-800 border-slate-600 rounded focus:ring-blue-500 focus:ring-2"
                />
                <span class="ml-2 text-sm text-gray-300">Remember me</span>
              </label>
              <a href="#" class="text-sm text-blue-500 hover:text-blue-400 transition-colors link-animate">
                Forgot password?
              </a>
            </div>

            <button type="submit" class="btn-primary" :class="{ 'animate-fade-in-up delay-800': isPageLoaded }">
              <span v-if="isLoading" class="loading-dots">LOGIN</span>
              <span v-else>LOGIN</span>
            </button>

            <div class="relative">
              <div class="absolute inset-0 flex items-center">
                <div class="w-full border-t border-gray-600"></div>
              </div>
            <div class="relative flex justify-center text-sm">
              <span class="px-2 bg-slate-800 text-gray-400">Or</span>
            </div>
            </div>

            <div class="space-y-3" :class="{ 'animate-fade-in-up delay-900': isPageLoaded }">
              <button type="button" class="btn-secondary flex items-center justify-center space-x-3">
                <svg class="h-5 w-5" viewBox="0 0 24 24">
                  <path fill="currentColor" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                  <path fill="currentColor" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                  <path fill="currentColor" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                  <path fill="currentColor" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
                </svg>
                <span>GOOGLE</span>
              </button>

              <button type="button" class="btn-secondary flex items-center justify-center space-x-3">
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M13 3.5c.73-.83 1.94-1.46 2.94-1.5.13 1.17-.34 2.35-1.04 3.19-.69.85-1.83 1.51-2.95 1.42-.15-1.15.41-2.35 1.05-3.11z"/>
                </svg>
                <span>APPLE</span>
              </button>
            </div>
          </form>

          <div class="mt-6 text-center" :class="{ 'animate-fade-in-up delay-1000': isPageLoaded }">
            <p class="text-gray-400">
              Don't have an account?
              <a href="#" class="text-blue-500 hover:text-blue-400 font-medium transition-colors link-animate">
                Sign Up
              </a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
