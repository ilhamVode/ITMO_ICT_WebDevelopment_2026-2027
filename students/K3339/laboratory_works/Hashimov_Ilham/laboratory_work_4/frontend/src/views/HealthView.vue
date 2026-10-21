<template>
  <v-card max-width="600">
    <v-card-title>Проверка подключения к Go API</v-card-title>
    <v-card-text>
      <v-alert :type="status === 'ok' ? 'success' : status === 'Проверяем...' ? 'info' : 'error'" class="mb-4">{{ status }}</v-alert>
      <v-btn color="primary" :loading="checking" @click="checkServer">Проверить сервер</v-btn>
    </v-card-text>
  </v-card>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/services/api'

const status = ref('Проверяем...')
const checking = ref(false)

async function checkServer() {
  checking.value = true
  try {
    const response = await api.get('/health')
    status.value = response.data.status
  } catch (error) {
    status.value = 'Ошибка подключения'
    console.error(error)
  } finally {
    checking.value = false
  }
}

onMounted(checkServer)
</script>
