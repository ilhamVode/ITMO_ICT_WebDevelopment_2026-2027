<template>
  <v-card class="mx-auto" max-width="460">
    <v-card-title>{{ isRegister ? 'Регистрация' : 'Вход' }}</v-card-title>
    <v-card-text>
      <v-alert v-if="route.query.changed" type="success" class="mb-4">Данные изменены. Войдите заново.</v-alert>
      <v-alert v-if="error" type="error" class="mb-4">{{ error }}</v-alert>
      <v-form @submit.prevent="submit">
        <v-text-field v-model.trim="username" label="Имя пользователя" autocomplete="username" required />
        <v-text-field v-model="password" label="Пароль" type="password" :autocomplete="isRegister ? 'new-password' : 'current-password'" required />
        <v-btn type="submit" color="primary" block :loading="saving">{{ isRegister ? 'Создать аккаунт' : 'Войти' }}</v-btn>
      </v-form>
    </v-card-text>
    <v-card-actions>
      <v-btn :to="isRegister ? '/login' : '/register'">{{ isRegister ? 'Уже есть аккаунт?' : 'Зарегистрироваться' }}</v-btn>
    </v-card-actions>
  </v-card>
</template>

<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api, { apiError } from '@/services/api'
import { login } from '@/services/auth'

const props = defineProps({ mode: { type: String, required: true } })
const route = useRoute()
const router = useRouter()
const isRegister = computed(() => props.mode === 'register')
const username = ref('')
const password = ref('')
const saving = ref(false)
const error = ref('')
watch(isRegister, () => { error.value = ''; password.value = '' })

async function submit() {
  error.value = ''
  if (!/^[\p{L}\p{N}_-]{3,50}$/u.test(username.value)) {
    error.value = 'Имя: 3–50 букв, цифр, _ или -'
    return
  }
  const passwordBytes = new TextEncoder().encode(password.value).length
  if (isRegister.value && (passwordBytes < 8 || passwordBytes > 72)) {
    error.value = 'Пароль: от 8 до 72 байт'
    return
  }
  saving.value = true
  try {
    if (isRegister.value) await api.post('/api/auth/users', { username: username.value, password: password.value })
    await login(username.value, password.value)
    const redirect = typeof route.query.redirect === 'string' && route.query.redirect.startsWith('/') && !route.query.redirect.startsWith('//') ? route.query.redirect : '/students'
    router.replace(redirect)
  } catch (err) { error.value = apiError(err) } finally { saving.value = false }
}
</script>
