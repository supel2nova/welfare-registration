<script setup lang="ts">
import { computed, onBeforeUnmount } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import UserSwitcher from '@/components/common/UserSwitcher.vue'
import { useApplicationStore } from '@/stores/application'
import { formatThaiDateTime } from '@/utils/buddhist'

const router = useRouter()
const store = useApplicationStore()

const application = computed(() => store.submitted)

if (!application.value) void router.replace({ name: 'register' })

const submittedAt = computed(() =>
  application.value ? formatThaiDateTime(application.value.submitted_at) : '',
)

const statusLabels: Record<string, string> = {
  SUBMITTED: 'รับใบสมัครแล้ว',
  CANCELLED: 'ยกเลิกแล้ว',
}

const statusLabel = computed(() => {
  if (!application.value) return ''
  return statusLabels[application.value.status] ?? application.value.status
})

onBeforeUnmount(() => store.clear())
</script>

<template>
  <div class="page-shell">
    <header class="topbar">
      <p class="brand">สวัสดิการแห่งรัฐ</p>
      <UserSwitcher />
    </header>

    <main v-if="application" class="py-10 animate-rise motion-reduce:animate-none">
      <h1 class="mb-3 text-[clamp(1.75rem,3vw,2.35rem)] leading-tight">รับใบสมัครเรียบร้อย</h1>
      <dl class="card my-6 grid gap-3">
        <div>
          <dt class="text-sm text-ink-muted">เลขที่ใบสมัคร</dt>
          <dd class="mt-0.5 font-semibold tracking-wide tabular-nums">{{ application.application_no }}</dd>
        </div>
        <div>
          <dt class="text-sm text-ink-muted">หน่วยรับลงทะเบียน</dt>
          <dd class="mt-0.5 font-semibold">{{ application.registration_unit }}</dd>
        </div>
        <div>
          <dt class="text-sm text-ink-muted">วันที่ยื่น</dt>
          <dd class="mt-0.5 font-semibold">{{ submittedAt }}</dd>
        </div>
        <div>
          <dt class="text-sm text-ink-muted">สถานะ</dt>
          <dd class="mt-0.5 font-semibold">{{ statusLabel }}</dd>
        </div>
      </dl>
      <RouterLink class="btn btn-primary" to="/register">กรอกใบใหม่</RouterLink>
    </main>
  </div>
</template>
