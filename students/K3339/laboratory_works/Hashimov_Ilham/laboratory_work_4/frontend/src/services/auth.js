import { ref } from 'vue'
import api from './api'

export const token = ref(localStorage.getItem('auth_token'))
export const user = ref(null)

export async function loadUser() {
  if (!token.value) return null
  try {
    const { data } = await api.get('/api/auth/users/me')
    user.value = data
    return data
  } catch (error) {
    if (error.response?.status === 401) clearSession()
    throw error
  }
}

export function clearSession() {
  localStorage.removeItem('auth_token')
  token.value = null
  user.value = null
}

export async function login(username, password) {
  const { data } = await api.post('/api/auth/token/login', { username, password })
  localStorage.setItem('auth_token', data.auth_token)
  token.value = data.auth_token
  return loadUser()
}

export async function logout() {
  try {
    await api.post('/api/auth/token/logout')
  } finally {
    clearSession()
  }
}
