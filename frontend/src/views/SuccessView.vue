<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import UserSwitcher from '@/components/common/UserSwitcher.vue'
import type { CreateApplicationResponse } from '@/types/api'
import { formatThaiDateTime } from '@/utils/buddhist'

const route = useRoute()

const application = computed(() => {
  const q = route.query
  return {
    application_no: String(q.application_no ?? ''),
    registration_unit: String(q.registration_unit ?? ''),
    submitted_at: String(q.submitted_at ?? ''),
    status: String(q.status ?? 'SUBMITTED'),
  } satisfies Partial<CreateApplicationResponse>
})

const submittedAt = computed(() => formatThaiDateTime(application.value.submitted_at))

const statusLabels: Record<string, string> = {
  SUBMITTED: 'รับใบสมัครแล้ว',
  CANCELLED: 'ยกเลิกแล้ว',
}

const statusLabel = computed(() => statusLabels[application.value.status] ?? application.value.status)
</script>

<template>
  <div class="page-shell">
    <header class="topbar">
      <p class="brand">สวัสดิการแห่งรัฐ</p>
      <UserSwitcher />
    </header>

    <main class="py-10 animate-rise motion-reduce:animate-none">
      <p class="m-0 text-sm font-semibold uppercase tracking-[0.08em] text-brand-600">บันทึกแล้ว</p>
      <h1 class="mt-1 mb-3 text-[clamp(1.75rem,3vw,2.35rem)] leading-tight">รับใบสมัครเรียบร้อย</h1>
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
      <div class="mt-4 flex flex-wrap gap-3">
        <RouterLink class="btn btn-primary" to="/register">กรอกใบใหม่</RouterLink>
        <RouterLink class="btn btn-ghost" to="/">กลับหน้าแรก</RouterLink>
      </div>
    </main>
  </div>
</template>
