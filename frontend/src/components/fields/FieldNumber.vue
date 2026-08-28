<script setup lang="ts">
import { computed } from 'vue'
import { useFormContext } from '@/composables/useFormContext'
import { cn } from '@/utils/cn'

const props = defineProps<{
  path: string
  label: string
  modelValue: number | null
  placeholder?: string
  min?: number
  max?: number
  step?: number | string
  disabled?: boolean
  format?: 'amount'
  suffix?: string
}>()

const emit = defineEmits<{
  'update:modelValue': [value: number | null]
}>()

const { fieldErrors, isSubmitting, clear } = useFormContext()
const error = computed(() => fieldErrors.value[props.path] ?? '')
const grouped = computed(() => props.format === 'amount')

function group(raw: string): string {
  if (raw === '') return ''
  const [int, dec] = raw.split('.')
  const withCommas = int.replace(/\B(?=(\d{3})+(?!\d))/g, ',')
  return dec === undefined ? withCommas : `${withCommas}.${dec}`
}

function clean(value: string): string {
  const only = value.replace(/[^\d.]/g, '')
  const [int, ...rest] = only.split('.')
  return rest.length === 0 ? int : `${int}.${rest.join('').slice(0, 2)}`
}

const displayValue = computed(() => {
  if (props.modelValue == null) return ''
  return grouped.value ? group(String(props.modelValue)) : String(props.modelValue)
})

function onInput(e: Event) {
  const el = e.target as HTMLInputElement

  if (!grouped.value) {
    const raw = el.value
    emit('update:modelValue', raw === '' ? null : Number(raw))
    clear(props.path)
    return
  }

  const caretFromEnd = el.value.length - (el.selectionStart ?? el.value.length)
  const raw = clean(el.value)
  el.value = group(raw)
  const caret = Math.max(0, el.value.length - caretFromEnd)
  el.setSelectionRange(caret, caret)

  if (raw === '' || raw === '.') {
    emit('update:modelValue', null)
  } else {
    const n = Number(raw)
    emit('update:modelValue', Number.isFinite(n) ? n : null)
  }
  clear(props.path)
}
</script>

<template>
  <label class="field">
    <span class="field-label">{{ label }}</span>
    <div v-if="suffix" class="field-affix">
      <input
        :class="cn('field-input', 'pr-14', 'text-right', 'tabular-nums', error && 'field-input-error')"
        :type="grouped ? 'text' : 'number'"
        :inputmode="grouped ? 'decimal' : undefined"
        :value="displayValue"
        :placeholder="placeholder"
        :min="grouped ? undefined : min"
        :max="grouped ? undefined : max"
        :step="grouped ? undefined : (step ?? 1)"
        :disabled="disabled || isSubmitting"
        :aria-invalid="error ? true : undefined"
        @input="onInput"
      />
      <span class="field-suffix">{{ suffix }}</span>
    </div>
    <input
      v-else
      :class="cn('field-input', grouped && 'text-right tabular-nums', error && 'field-input-error')"
      :type="grouped ? 'text' : 'number'"
      :inputmode="grouped ? 'decimal' : undefined"
      :value="displayValue"
      :placeholder="placeholder"
      :min="grouped ? undefined : min"
      :max="grouped ? undefined : max"
      :step="grouped ? undefined : (step ?? 1)"
      :disabled="disabled || isSubmitting"
      :aria-invalid="error ? true : undefined"
      @input="onInput"
    />
    <span v-if="error" class="field-error">{{ error }}</span>
  </label>
</template>
