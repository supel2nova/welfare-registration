import { createRouter, createWebHistory } from 'vue-router'
import RegisterFormView from '@/views/RegisterFormView.vue'
import SuccessView from '@/views/SuccessView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', redirect: { name: 'register' } },
    { path: '/register', name: 'register', component: RegisterFormView },
    { path: '/success', name: 'success', component: SuccessView },
  ],
  scrollBehavior() {
    return { top: 0 }
  },
})

export default router
