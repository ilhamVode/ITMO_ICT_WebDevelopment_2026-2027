<template>
  <v-table>
    <thead><tr><th>Студент</th><th>Группа</th><th>Дата рождения</th><th>Статус</th><th>Действия</th></tr></thead>
    <tbody>
      <tr v-for="student in students" :key="student.id">
        <td>{{ student.last_name }} {{ student.first_name }}</td>
        <td>{{ groups.find((g) => g.id === student.group_id)?.name || student.group_id }}</td>
        <td>{{ student.birth_date || '—' }}</td>
        <td>{{ student.is_active ? 'Обучается' : 'Отчислен' }}</td>
        <td>
          <v-btn size="small" variant="text" @click="$emit('edit', student)">Изменить</v-btn>
          <v-btn size="small" variant="text" color="error" @click="$emit('remove', student)">Удалить</v-btn>
        </td>
      </tr>
      <tr v-if="!students.length"><td colspan="5">Студенты не найдены</td></tr>
    </tbody>
  </v-table>
</template>

<script setup>
defineProps({
  students: {
    type: Array,
    required: true
  },
  groups: { type: Array, required: true }
})
defineEmits(['edit', 'remove'])
</script>
