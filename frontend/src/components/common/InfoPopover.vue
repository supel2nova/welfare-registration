<script setup lang="ts">
import { onBeforeUnmount, ref, useId } from 'vue'

defineProps<{
  trigger: string
  label?: string
}>()

const open = ref(false)
const root = ref<HTMLElement | null>(null)
const panelId = useId()

function toggle() {
  open.value = !open.value
  if (open.value) {
    document.addEventListener('pointerdown', onOutside, true)
    document.addEventListener('keydown', onKeydown)
  } else {
    detach()
  }
}

function close() {
  open.value = false
  detach()
}

function detach() {
  document.removeEventListener('pointerdown', onOutside, true)
  document.removeEventListener('keydown', onKeydown)
}

function onOutside(e: Event) {
  if (!root.value?.contains(e.target as Node)) close()
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') close()
}

onBeforeUnmount(detach)
</script>

<template>
  <span ref="root" class="relative inline-flex">
    <button
      type="button"
      class="info-trigger"
      :aria-label="label ?? trigger"
      :aria-expanded="open"
      :aria-controls="panelId"
      @click="toggle"
    >
      {{ trigger }}
    </button>
    <span v-if="open" :id="panelId" role="tooltip" class="info-panel">
      <slot />
    </span>
  </span>
</template>
