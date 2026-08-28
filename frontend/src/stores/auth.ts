import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { fetchDevUsers } from '@/api/ref'
import type { DevUser } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const users = ref<DevUser[]>([])
  const userId = ref<string>('')
  const loaded = ref(false)

  const currentUser = computed(() => users.value.find((u) => u.id === userId.value) ?? null)

  async function loadUsers() {
    users.value = await fetchDevUsers()
    loaded.value = true
    if (!userId.value && users.value.length > 0) {
      const preferred = users.value.find((u) => u.username === 'somying.baac')
      userId.value = preferred?.id ?? users.value[0].id
    }
  }

  function selectUser(id: string) {
    userId.value = id
  }

  return { users, userId, loaded, currentUser, loadUsers, selectUser }
})
