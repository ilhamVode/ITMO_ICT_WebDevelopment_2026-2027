<template>
  <div>
    <h1 class="text-h4 mb-4">Оценки</h1>
    <v-alert v-if="error" type="error" class="mb-4">{{ error }}</v-alert>
    <v-card class="mb-6">
      <v-card-title>{{ editing ? 'Изменить оценку' : 'Добавить оценку' }}</v-card-title>
      <v-card-text>
        <v-form @submit.prevent="save">
          <v-row>
            <v-col cols="12" sm="6"><v-select v-model="studentId" :items="students" :item-title="s => `${s.last_name} ${s.first_name}`" item-value="id" label="Студент" required /></v-col>
            <v-col cols="12" sm="6"><v-select v-model="subjectId" :items="subjects" item-title="name" item-value="id" label="Предмет" required /></v-col>
            <v-col cols="12" sm="6"><v-select v-model="semesterId" :items="semesters" :item-title="s => `${s.academic_year}, семестр ${s.number}`" item-value="id" label="Семестр" required /></v-col>
            <v-col cols="12" sm="6"><v-text-field v-model.number="value" label="Оценка" type="number" min="2" max="5" required /></v-col>
          </v-row>
          <v-btn type="submit" color="primary" :loading="saving">Сохранить</v-btn>
          <v-btn v-if="editing" class="ml-2" @click="reset">Отмена</v-btn>
        </v-form>
      </v-card-text>
    </v-card>
    <v-progress-linear v-if="loading" indeterminate />
    <v-table v-else>
      <thead><tr><th>Студент</th><th>Предмет</th><th>Семестр</th><th>Оценка</th><th>Действия</th></tr></thead>
      <tbody>
        <tr v-for="grade in grades" :key="grade.id">
          <td>{{ studentName(grade.student_id) }}</td>
          <td>{{ subjects.find(s => s.id === grade.subject_id)?.name || grade.subject_id }}</td>
          <td>{{ semesterName(grade.semester_id) }}</td>
          <td>{{ grade.value }}</td>
          <td><v-btn size="small" variant="text" @click="edit(grade)">Изменить</v-btn><v-btn size="small" variant="text" color="error" @click="remove(grade)">Удалить</v-btn></td>
        </tr>
        <tr v-if="!grades.length"><td colspan="5">Оценок пока нет</td></tr>
      </tbody>
    </v-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api, { apiError } from '@/services/api'

const grades = ref([])
const students = ref([])
const subjects = ref([])
const semesters = ref([])
const studentId = ref(null)
const subjectId = ref(null)
const semesterId = ref(null)
const value = ref(5)
const editing = ref(null)
const error = ref('')
const loading = ref(false)
const saving = ref(false)

function studentName(id) { const s = students.value.find(s => s.id === id); return s ? `${s.last_name} ${s.first_name}` : id }
function semesterName(id) { const s = semesters.value.find(s => s.id === id); return s ? `${s.academic_year}, ${s.number}` : id }
async function load() {
  loading.value = true
  error.value = ''
  try {
    const responses = await Promise.all(['/api/grades', '/api/students', '/api/subjects', '/api/semesters'].map(path => api.get(path)))
    grades.value = responses[0].data
    students.value = responses[1].data
    subjects.value = responses[2].data
    semesters.value = responses[3].data
  } catch (err) { error.value = apiError(err) } finally { loading.value = false }
}
function edit(grade) { editing.value = grade; studentId.value = grade.student_id; subjectId.value = grade.subject_id; semesterId.value = grade.semester_id; value.value = grade.value }
function reset() { editing.value = null; studentId.value = null; subjectId.value = null; semesterId.value = null; value.value = 5 }
async function save() {
  if (!studentId.value || !subjectId.value || !semesterId.value || !Number.isInteger(value.value)) return
  saving.value = true
  error.value = ''
  try {
    const body = { student_id: studentId.value, subject_id: subjectId.value, semester_id: semesterId.value, value: value.value }
    if (editing.value) await api.put(`/api/grades/${editing.value.id}`, body)
    else await api.post('/api/grades', body)
    reset(); await load()
  } catch (err) { error.value = apiError(err) } finally { saving.value = false }
}
async function remove(grade) {
  if (!window.confirm('Удалить оценку?')) return
  error.value = ''
  try { await api.delete(`/api/grades/${grade.id}`); await load() } catch (err) { error.value = apiError(err) }
}
onMounted(load)
</script>
