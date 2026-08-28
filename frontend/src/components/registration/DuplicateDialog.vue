<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { DuplicateInfo } from '@/types/api'

const props = defineProps<{
  open: boolean
  info: DuplicateInfo | null
}>()

const emit = defineEmits<{
  close: []
}>()

const closeButton = ref<HTMLButtonElement | null>(null)

watch(
  () => props.open,
  async (open) => {
    if (!open) return
    await nextTick()
    closeButton.value?.focus()
  },
)

function onKeydown(e: KeyboardEvent) {
  if (props.open && e.key === 'Escape') emit('close')
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => window.removeEventListener('keydown', onKeydown))
</script>

<template>
  <Teleport to="body">
    <div v-if="open && info" class="dialog-backdrop" @click.self="emit('close')">
      <div class="dialog-panel" role="dialog" aria-modal="true" aria-labelledby="dup-title">
        <h2 id="dup-title" class="mb-1 text-lg">ลงทะเบียนแล้ว</h2>
        <p class="m-0 text-ink-muted">เลขประจำตัวประชาชนนี้มีใบสมัครในปีงบประมาณนี้อยู่แล้ว</p>
        <dl class="card my-6 grid gap-3">
          <div>
            <dt class="text-sm text-ink-muted">เลขที่ใบสมัคร</dt>
            <dd class="mt-0.5 font-semibold">{{ info.application_no }}</dd>
          </div>
          <div>
            <dt class="text-sm text-ink-muted">หน่วยรับลงทะเบียน</dt>
            <dd class="mt-0.5 font-semibold">{{ info.registered_unit }}</dd>
          </div>
          <div>
            <dt class="text-sm text-ink-muted">วันที่ยื่น</dt>
            <dd class="mt-0.5 font-semibold">{{ info.registered_at }}</dd>
          </div>
        </dl>
        <button ref="closeButton" type="button" class="btn btn-primary" @click="emit('close')">ปิด</button>
      </div>
    </div>
  </Teleport>
</template>
