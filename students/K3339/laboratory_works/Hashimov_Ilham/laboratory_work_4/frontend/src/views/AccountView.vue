<template>
  <v-card max-width="600">
    <v-card-title>Учётные данные</v-card-title>
    <v-card-text>
      <v-alert v-if="error" type="error" class="mb-4">{{ error }}</v-alert>
      <v-form @submit.prevent="save">
        <v-text-field v-model.trim="username" label="Новое имя пользователя" autocomplete="username" required />
        <v-text-field v-model="currentPassword" label="Текущий пароль" type="password" autocomplete="current-password" required />
        <v-text-field v-model="newPassword" label="Новый пароль" type="password" autocomplete="new-password" required />
        <v-btn type="submit" color="primary" :loading="saving">Сохранить</v-btn>
      </v-form>
      <p class="mt-4">После сохранения потребуется войти снова.</p>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import api, { apiError } from '@/services/api'
import { user, clearSession } from '@/services/auth'

const router = useRouter()
const username = ref('')
const currentPassword = ref('')
const newPassword = ref('')
const saving = ref(false)
const error = ref('')
watch(user, (value) => { if (value) username.value = value.username }, { immediate: true })

async function save() {
  error.value = ''
  const bytes = new TextEncoder().encode(newPassword.value).length
  if (!/^[\p{L}\p{N}_-]{3,50}$/u.test(username.value) || bytes < 8 || bytes > 72) {
    error.value = 'Проверьте имя (3–50 символов) и новый пароль (8–72 байта)'
    return
  }
  saving.value = true
  try {
    await api.patch('/api/auth/users/me', { username: username.value, current_password: currentPassword.value, new_password: newPassword.value })
    clearSession()
    router.replace({ path: '/login', query: { changed: '1' } })
  } catch (err) { error.value = apiError(err) } finally { saving.value = false }
}
</script>
