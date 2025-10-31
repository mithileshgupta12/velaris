import { createRouter, createWebHistory } from 'vue-router'
import LoginView from './LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [{ path: '/login', component: LoginView }],
})

export default router
