<script setup lang="ts">
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import AppButton from '@/components/common/AppButton.vue'

const props = withDefaults(
  defineProps<{
    open: boolean
    title: string
    message: string
    confirmLabel?: string
    cancelLabel?: string
    danger?: boolean
  }>(),
  { confirmLabel: 'ยืนยัน', cancelLabel: 'ย้อนกลับ', danger: false },
)

const emit = defineEmits<{
  confirm: []
  close: []
}>()

const cancelButton = ref<InstanceType<typeof AppButton> | null>(null)

watch(
  () => props.open,
  async (open) => {
    if (!open) return
    await nextTick()
    ;(cancelButton.value?.$el as HTMLButtonElement | undefined)?.focus()
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
    <div v-if="open" class="dialog-backdrop" @click.self="emit('close')">
      <div class="dialog-panel" role="dialog" aria-modal="true" aria-labelledby="confirm-title">
        <h2 id="confirm-title" class="mb-1 text-lg">{{ title }}</h2>
        <p class="m-0 text-ink-muted">{{ message }}</p>
        <div class="btn-row mt-6">
          <AppButton ref="cancelButton" @click="emit('close')">{{ cancelLabel }}</AppButton>
          <AppButton :variant="danger ? 'danger' : 'primary'" @click="emit('confirm')">
            {{ confirmLabel }}
          </AppButton>
        </div>
      </div>
    </div>
  </Teleport>
</template>
