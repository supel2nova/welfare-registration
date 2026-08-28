import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from '@/App.vue'
import { setRequestHeaders } from '@/api/client'
import router from '@/router'
import { useAuthStore } from '@/stores/auth'
import '@/style.css'

const app = createApp(App)
const pinia = createPinia()
app.use(pinia)
app.use(router)

setRequestHeaders((): Record<string, string> => {
  const auth = useAuthStore(pinia)
  if (!auth.userId) return {}
  return { 'X-Debug-User-Id': auth.userId }
})

app.mount('#app')
