<template>
  <v-form @submit.prevent="submitForm">
    <v-row>
      <v-col cols="12" sm="6"><v-text-field v-model.trim="firstName" label="Имя" required /></v-col>
      <v-col cols="12" sm="6"><v-text-field v-model.trim="lastName" label="Фамилия" required /></v-col>
      <v-col cols="12" sm="6"><v-select v-model="groupId" :items="groups" item-title="name" item-value="id" label="Группа" required /></v-col>
      <v-col cols="12" sm="6"><v-text-field v-model="birthDate" label="Дата рождения" type="date" /></v-col>
      <v-col v-if="student" cols="12"><v-checkbox v-model="isActive" label="Обучается" /></v-col>
    </v-row>
    <v-btn type="submit" color="primary" :loading="saving">{{ student ? 'Сохранить' : 'Добавить' }}</v-btn>
    <v-btn v-if="student" class="ml-2" @click="$emit('cancel')">Отмена</v-btn>
  </v-form>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({ student: { type: Object, default: null }, groups: { type: Array, required: true }, saving: Boolean })
const emit = defineEmits(['save', 'cancel'])
const firstName = ref('')
const lastName = ref('')
const groupId = ref(null)
const birthDate = ref('')
const isActive = ref(true)

watch(() => props.student, (student) => {
  firstName.value = student?.first_name || ''
  lastName.value = student?.last_name || ''
  groupId.value = student?.group_id || null
  birthDate.value = student?.birth_date || ''
  isActive.value = student?.is_active ?? true
}, { immediate: true })

function submitForm() {
  if (!firstName.value || !lastName.value || !groupId.value) return
  emit('save', {
    first_name: firstName.value,
    last_name: lastName.value,
    group_id: groupId.value,
    birth_date: birthDate.value || null,
    is_active: isActive.value,
  })
}
</script>
