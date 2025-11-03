import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '@/views/auth/LoginView.vue'
import RegisterView from '@/views/auth/RegisterView.vue'
import DashboardView from '@/views/DashboardView.vue'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: HomeView, name: 'home' },
    {
      path: '/auth',
      children: [
        { path: 'register', component: RegisterView, name: 'auth.register' },
        { path: 'login', component: LoginView, name: 'auth.login' },
      ],
    },
    { path: '/dashboard', component: DashboardView, name: 'dashboard' },
  ],
})

router.beforeEach(() => {
  console.log('Hello, World!')
})

export default router
