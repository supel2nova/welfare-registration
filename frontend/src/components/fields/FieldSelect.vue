<script setup lang="ts">
import { computed } from 'vue'
import { useFormContext } from '@/composables/useFormContext'
import type { SelectOption } from '@/types/form'
import { cn } from '@/utils/cn'

const props = defineProps<{
  path: string
  label: string
  modelValue: string | null
  options: SelectOption[]
  disabled?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { fieldErrors, isSubmitting, clear } = useFormContext()
const error = computed(() => fieldErrors.value[props.path] ?? '')

function onChange(e: Event) {
  emit('update:modelValue', (e.target as HTMLSelectElement).value)
  clear(props.path)
}
</script>

<template>
  <label class="field">
    <span class="field-label">{{ label }}</span>
    <select
      :class="cn('field-input', 'field-select', error && 'field-input-error')"
      :value="modelValue ?? ''"
      :disabled="disabled || isSubmitting"
      @change="onChange"
    >
      <option v-for="o in options" :key="o.value" :value="o.value">{{ o.label }}</option>
    </select>
    <span v-if="error" class="field-error">{{ error }}</span>
  </label>
</template>
