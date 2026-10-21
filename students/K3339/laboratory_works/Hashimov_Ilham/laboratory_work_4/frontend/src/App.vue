<template>
  <v-app>
    <v-app-bar color="primary" title="Учебная часть">
      <template #append>
        <v-btn v-if="token" to="/students">Студенты</v-btn>
        <v-btn v-if="token" to="/groups">Группы</v-btn>
        <v-btn v-if="token" to="/grades">Оценки</v-btn>
        <v-btn to="/health">API</v-btn>
        <v-btn v-if="token" to="/account">{{ user?.username || 'Профиль' }}</v-btn>
        <v-btn v-if="token" :loading="leaving" @click="signOut">Выйти</v-btn>
        <v-btn v-else to="/login">Войти</v-btn>
      </template>
    </v-app-bar>

    <v-main>
      <v-container class="py-8" style="max-width: 1100px">
        <RouterView />
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { RouterView, useRouter } from 'vue-router'
import { token, user, loadUser, logout } from '@/services/auth'

const router = useRouter()
const leaving = ref(false)

onMounted(async () => {
  if (token.value) {
    try { await loadUser() } catch { if (!token.value) router.push('/login') }
  }
})

async function signOut() {
  leaving.value = true
  try { await logout() } catch { /* выход из приложения всё равно выполнится */ }
  leaving.value = false
  router.push('/login')
}
</script>
