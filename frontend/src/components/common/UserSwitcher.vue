<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()

onMounted(() => {
  void auth.loadUsers().catch(() => undefined)
})
</script>

<template>
  <label class="flex w-full min-w-0 items-center gap-2 sm:w-auto">
    <span class="field-label whitespace-nowrap">เจ้าหน้าที่</span>
    <select
      class="field-input field-select"
      :value="auth.userId"
      :disabled="!auth.loaded || auth.users.length === 0"
      @change="auth.selectUser(($event.target as HTMLSelectElement).value)"
    >
      <option v-if="!auth.loaded" value="">กำลังโหลด…</option>
      <option v-for="u in auth.users" :key="u.id" :value="u.id">
        {{ u.username }} · {{ u.org_name }}
      </option>
    </select>
  </label>
</template>
