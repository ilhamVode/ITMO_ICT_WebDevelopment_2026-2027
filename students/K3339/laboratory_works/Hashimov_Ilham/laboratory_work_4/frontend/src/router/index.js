import { createRouter, createWebHistory } from 'vue-router'

import Hello from '@/components/Hello.vue'
import StudentsView from '@/views/StudentsView.vue'
import HealthView from '@/views/HealthView.vue'
import AuthView from '@/views/AuthView.vue'
import AccountView from '@/views/AccountView.vue'
import GroupsView from '@/views/GroupsView.vue'
import GradesView from '@/views/GradesView.vue'
import { token } from '@/services/auth'

const routes = [
  { path: '/', redirect: '/students' },
  {
    path: '/hi',
    component: Hello
  },
  {
    path: '/students',
    component: StudentsView,
    meta: { requiresAuth: true }
  },
  {
    path: '/health',
    component: HealthView
  },
  { path: '/login', component: AuthView, props: { mode: 'login' } },
  { path: '/register', component: AuthView, props: { mode: 'register' } },
  { path: '/account', component: AccountView, meta: { requiresAuth: true } },
  { path: '/groups', component: GroupsView, meta: { requiresAuth: true } },
  { path: '/grades', component: GradesView, meta: { requiresAuth: true } },
  { path: '/:pathMatch(.*)*', redirect: '/' }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !token.value) {
    return { path: '/login', query: { redirect: to.fullPath } }
  }
})

export default router
