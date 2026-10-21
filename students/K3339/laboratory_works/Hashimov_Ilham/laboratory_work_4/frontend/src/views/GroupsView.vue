<template>
  <div>
    <h1 class="text-h4 mb-4">Учебные группы</h1>
    <v-alert v-if="error" type="error" class="mb-4">{{ error }}</v-alert>
    <v-card class="mb-6">
      <v-card-title>{{ editing ? 'Изменить группу' : 'Добавить группу' }}</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="save">
          <v-row>
            <v-col cols="12" sm="6"><v-text-field v-model.trim="name" label="Название" required /></v-col>
            <v-col cols="12" sm="6"><v-text-field v-model.number="course" label="Курс" type="number" min="1" max="6" required /></v-col>
          </v-row>
          <v-btn type="submit" color="primary" :loading="saving">Сохранить</v-btn>
          <v-btn v-if="editing" class="ml-2" @click="reset">Отмена</v-btn>
        </v-form>
      </v-card-text>
    </v-card>
    <v-progress-linear v-if="loading" indeterminate />
    <v-table v-else>
      <thead><tr><th>Группа</th><th>Курс</th><th>Действия</th></tr></thead>
      <tbody>
        <tr v-for="group in groups" :key="group.id">
          <td>{{ group.name }}</td><td>{{ group.course }}</td>
          <td><v-btn size="small" variant="text" @click="edit(group)">Изменить</v-btn><v-btn size="small" variant="text" color="error" @click="remove(group)">Удалить</v-btn></td>
        </tr>
        <tr v-if="!groups.length"><td colspan="3">Групп пока нет</td></tr>
      </tbody>
    </v-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api, { apiError } from '@/services/api'

const groups = ref([])
const name = ref('')
const course = ref(1)
const editing = ref(null)
const error = ref('')
const loading = ref(false)
const saving = ref(false)

async function load() {
  loading.value = true
  try { groups.value = (await api.get('/api/groups')).data } catch (err) { error.value = apiError(err) } finally { loading.value = false }
}
function edit(group) { editing.value = group; name.value = group.name; course.value = group.course }
function reset() { editing.value = null; name.value = ''; course.value = 1 }
async function save() {
  if (!name.value || !Number.isInteger(course.value) || course.value < 1) return
  saving.value = true
  error.value = ''
  try {
    const body = { name: name.value, course: course.value }
    if (editing.value) await api.put(`/api/groups/${editing.value.id}`, body)
    else await api.post('/api/groups', body)
    reset(); await load()
  } catch (err) { error.value = apiError(err) } finally { saving.value = false }
}
async function remove(group) {
  if (!window.confirm(`Удалить группу ${group.name}?`)) return
  error.value = ''
  try { await api.delete(`/api/groups/${group.id}`); await load() } catch (err) { error.value = apiError(err) }
}
onMounted(load)
</script>
