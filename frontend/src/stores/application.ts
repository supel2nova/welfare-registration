import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { CreateApplicationResponse } from '@/types/api'

export const useApplicationStore = defineStore('application', () => {
  const submitted = ref<CreateApplicationResponse | null>(null)

  function keep(application: CreateApplicationResponse) {
    submitted.value = application
  }

  function clear() {
    submitted.value = null
  }

  return { submitted, keep, clear }
})
