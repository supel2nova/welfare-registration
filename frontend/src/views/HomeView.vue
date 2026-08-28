<script setup lang="ts">
import { RouterLink } from 'vue-router'
import UserSwitcher from '@/components/common/UserSwitcher.vue'
</script>

<template>
  <div class="page-shell">
    <header class="topbar">
      <p class="brand">สวัสดิการแห่งรัฐ</p>
      <UserSwitcher />
    </header>

    <main class="flex min-h-[60vh] flex-col justify-center py-10 animate-rise motion-reduce:animate-none">
      <p class="m-0 text-sm font-semibold uppercase tracking-[0.08em] text-brand-600">ช่องทางเคาน์เตอร์</p>
      <h1 class="mt-1 mb-3 text-[clamp(1.75rem,3vw,2.35rem)] leading-tight">รับลงทะเบียนผู้มีสิทธิ</h1>
      <p class="mb-6 max-w-xl text-ink-muted">
        กรอกใบสมัครหน้าเคาน์เตอร์ ระบบตรวจรูปแบบและกันซ้ำก่อนบันทึก ไม่มีการพิจารณาคุณสมบัติในขั้นตอนนี้
      </p>
      <RouterLink class="btn btn-primary w-fit" to="/register">เริ่มกรอกใบสมัคร</RouterLink>
    </main>
  </div>
</template>
