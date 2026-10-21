<template>
  <div>
    <h1 class="text-h4 mb-4">Студенты</h1>
    <v-alert v-if="error" type="error" class="mb-4">{{ error }}</v-alert>
    <v-card class="mb-6">
      <v-card-title>{{ editing ? 'Изменить студента' : 'Добавить студента' }}</v-card-title>
      <v-card-text><StudentForm :key="formKey" :student="editing" :groups="groups" :saving="saving" @save="saveStudent" @cancel="editing = null" /></v-card-text>
    </v-card>
    <v-text-field v-model="query" label="Поиск по имени или фамилии" clearable class="mb-2" />
    <v-progress-linear v-if="loading" indeterminate />
    <StudentList v-else :students="filtered" :groups="groups" @edit="editing = $event" @remove="removeStudent" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import StudentForm from '@/components/StudentForm.vue'
import StudentList from '@/components/StudentList.vue'
import api, { apiError } from '@/services/api'

const students = ref([])
const groups = ref([])
const editing = ref(null)
const query = ref('')
const loading = ref(false)
const saving = ref(false)
const formKey = ref(0)
const error = ref('')
const filtered = computed(() => students.value.filter((s) => `${s.first_name} ${s.last_name}`.toLowerCase().includes((query.value || '').toLowerCase())))

async function load() {
  loading.value = true
  error.value = ''
  try {
    const [studentsResponse, groupsResponse] = await Promise.all([api.get('/api/students'), api.get('/api/groups')])
    students.value = studentsResponse.data
    groups.value = groupsResponse.data
  } catch (err) { error.value = apiError(err) } finally { loading.value = false }
}

async function saveStudent(value) {
  saving.value = true
  error.value = ''
  try {
    if (editing.value) await api.put(`/api/students/${editing.value.id}`, value)
    else await api.post('/api/students', { first_name: value.first_name, last_name: value.last_name, group_id: value.group_id, birth_date: value.birth_date })
    editing.value = null
    formKey.value += 1
    await load()
  } catch (err) { error.value = apiError(err) } finally { saving.value = false }
}

async function removeStudent(student) {
  if (!window.confirm(`Удалить ${student.last_name} ${student.first_name}?`)) return
  try { await api.delete(`/api/students/${student.id}`); await load() } catch (err) { error.value = apiError(err) }
}

onMounted(load)
</script>
