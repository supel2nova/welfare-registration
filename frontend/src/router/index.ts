import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'
import RegisterFormView from '@/views/RegisterFormView.vue'
import SuccessView from '@/views/SuccessView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', name: 'home', component: HomeView },
    { path: '/register', name: 'register', component: RegisterFormView },
    { path: '/success', name: 'success', component: SuccessView },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

export default router
