<template>
  <div>
    <h1>Список студентов</h1>

    <input
      v-model="studentName"
      placeholder="Введите имя студента"
    />

    <button @click="addStudent">
      Добавить
    </button>

    <p v-if="students.length === 0">
      Студентов пока нет
    </p>

    <ul>
      <li v-for="student in students" :key="student.id">
        {{ student.name }}
      </li>
    </ul>
  </div>
</template>

<script setup>
defineOptions({ name: 'HelloPage' })
import { ref } from 'vue'

const studentName = ref('')
const students = ref([])
let nextId = 1

function addStudent() {
  const name = studentName.value.trim()

  if (!name) return

  students.value.push({
    id: nextId++,
    name: name
  })

  studentName.value = ''
}
</script>