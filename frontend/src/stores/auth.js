import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '../api/auth'
import router from '../router'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)
  const loading = ref(false)
  const error = ref(null)

  const isLoggedIn = computed(() => !!token.value)

  async function register(name, email, password) {
    loading.value = true
    error.value = null
    try {
      const res = await authApi.register({ name, email, password })
      setAuth(res.data)
      router.push('/dashboard')
    } catch (e) {
      error.value = e.response?.data?.error || 'Registration failed'
    } finally {
      loading.value = false
    }
  }

  async function login(email, password) {
    loading.value = true
    error.value = null
    try {
      const res = await authApi.login({ email, password })
      setAuth(res.data)
      router.push('/dashboard')
    } catch (e) {
      error.value = e.response?.data?.error || 'Login failed'
    } finally {
      loading.value = false
    }
  }

  async function fetchMe() {
    try {
      const res = await authApi.me()
      user.value = res.data
      // อัปเดต role ใน localStorage ด้วย
      localStorage.setItem('role', res.data.role)
    } catch {
      logout()
    }
  }

  function setAuth(data) {
    token.value = data.token
    user.value = data.user
    localStorage.setItem('token', data.token)
    localStorage.setItem('role', data.user.role)
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('role')
    router.push('/login')
  }

  return { user, token, loading, error, isLoggedIn, register, login, fetchMe, logout }
})
