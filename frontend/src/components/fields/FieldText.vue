<script setup lang="ts">
import { computed } from 'vue'
import { useFormContext } from '@/composables/useFormContext'
import {
  formatLaserId,
  formatNationalId,
  formatPhone,
  isLaserIdComplete,
  isNationalIdComplete,
  isPhoneComplete,
  parseNationalId,
  parsePhone,
} from '@/utils/idFormat'
import { cn } from '@/utils/cn'

const props = defineProps<{
  path: string
  label: string
  modelValue: string | null
  placeholder?: string
  type?: string
  maxlength?: number
  disabled?: boolean
  format?: 'national_id' | 'laser_id' | 'phone'
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const { fieldErrors, isSubmitting, clear, setError } = useFormContext()
const error = computed(() => fieldErrors.value[props.path] ?? '')

const displayValue = computed(() => {
  const raw = props.modelValue ?? ''
  if (props.format === 'national_id') return formatNationalId(raw)
  if (props.format === 'laser_id') return formatLaserId(raw)
  if (props.format === 'phone') return formatPhone(raw)
  return raw
})

const inputMode = computed(() => {
  if (props.format === 'national_id' || props.format === 'phone') return 'numeric'
  return 'text'
})

const inputMaxlength = computed(() => {
  if (props.format === 'national_id') return 17
  if (props.format === 'laser_id') return 14
  if (props.format === 'phone') return 12
  return props.maxlength
})

function onInput(e: Event) {
  const el = e.target as HTMLInputElement
  let next = el.value

  if (props.format === 'national_id') {
    next = parseNationalId(next)
    el.value = formatNationalId(next)
  } else if (props.format === 'laser_id') {
    next = formatLaserId(next)
    el.value = next
  } else if (props.format === 'phone') {
    next = parsePhone(next)
    el.value = formatPhone(next)
  }

  emit('update:modelValue', next)
  clear(props.path)
}

function onBlur() {
  const raw = props.modelValue ?? ''
  if (!raw) {
    clear(props.path)
    return
  }

  if (props.format === 'national_id') {
    if (!isNationalIdComplete(raw)) {
      setError(props.path, 'กรอกเลขประจำตัวประชาชนให้ครบ 13 หลัก')
    }
    return
  }

  if (props.format === 'laser_id') {
    if (!isLaserIdComplete(raw)) {
      setError(props.path, 'กรอกรหัสหลังบัตรให้ครบ เช่น JT8-1234567-89')
    }
    return
  }

  if (props.format === 'phone' && !isPhoneComplete(raw)) {
    setError(props.path, 'มือถือ 10 หลัก หรือเบอร์บ้าน 9 หลัก เช่น 099-119-2231 หรือ 02-123-4567')
  }
}
</script>

<template>
  <label class="field">
    <span class="field-label">{{ label }}</span>
    <input
      :class="cn('field-input', format && 'field-input-mono', error && 'field-input-error')"
      :type="type ?? 'text'"
      :inputmode="inputMode"
      :value="displayValue"
      :placeholder="placeholder"
      :maxlength="inputMaxlength"
      :disabled="disabled || isSubmitting"
      autocomplete="off"
      spellcheck="false"
      @input="onInput"
      @blur="onBlur"
    />
    <span v-if="error" class="field-error">{{ error }}</span>
  </label>
</template>
