<script setup lang="ts">
import AppIcon from '@/components/common/AppIcon.vue'
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
      <button
        type="button"
        class="btn-remove"
        :disabled="isSubmitting"
        :aria-label="`ลบ${title}`"
        :title="`ลบ${title}`"
        @click="emit('remove')"
      >
        <AppIcon name="trash" />
      </button>
    </div>
    <slot />
  </div>
</template>
