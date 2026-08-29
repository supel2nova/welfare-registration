<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/utils/cn'

const props = withDefaults(
  defineProps<{
    variant?: 'primary' | 'ghost' | 'danger' | 'soft' | 'add'
    type?: 'button' | 'submit'
    disabled?: boolean
  }>(),
  { variant: 'ghost', type: 'button', disabled: false },
)

const classes = computed(() => cn('btn', `btn-${props.variant}`))
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled">
    <slot />
  </button>
</template>
