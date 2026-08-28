<script setup lang="ts">
import { useFormContext } from '@/composables/useFormContext'

defineProps<{
  title: string
}>()

const emit = defineEmits<{
  remove: []
}>()

const { isSubmitting } = useFormContext()
</script>

<template>
  <div class="item-block">
    <div class="item-head">
      <span class="item-title">{{ title }}</span>
      <button type="button" class="btn-remove" :disabled="isSubmitting" @click="emit('remove')">
        ลบ
      </button>
    </div>
    <slot />
  </div>
</template>
