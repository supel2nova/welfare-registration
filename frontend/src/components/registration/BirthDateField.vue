<script setup lang="ts">
import { computed } from 'vue'
import FieldNumber from '@/components/fields/FieldNumber.vue'
import InfoPopover from '@/components/common/InfoPopover.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import { useFormContext } from '@/composables/useFormContext'
import type { SelectOption } from '@/types/form'
import { toBE, toCE } from '@/utils/buddhist'

const props = defineProps<{
  year: number | null
  month: number | null
  day: number | null
}>()

const emit = defineEmits<{
  'update:year': [value: number | null]
  'update:month': [value: number | null]
  'update:day': [value: number | null]
  'update:precision': [value: string]
}>()

const { fieldErrors } = useFormContext()

const thaiMonths = [
  'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
  'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม',
]

const unknown = { value: '', label: 'ไม่ทราบ' }

const monthOptions = computed<SelectOption[]>(() => [
  unknown,
  ...thaiMonths.map((label, i) => ({ value: String(i + 1), label })),
])

const daysInMonth = computed(() => {
  if (props.month == null) return 31
  const year = props.year == null ? 2024 : toCE(props.year)
  return new Date(year, props.month, 0).getDate()
})

const dayOptions = computed<SelectOption[]>(() => [
  unknown,
  ...Array.from({ length: daysInMonth.value }, (_, i) => ({
    value: String(i + 1),
    label: String(i + 1),
  })),
])

const precision = computed(() => {
  if (props.month == null) return 'YEAR_ONLY'
  if (props.day == null) return 'YEAR_MONTH'
  return 'FULL'
})

const summary = computed(() => {
  if (props.year == null) return ''
  const year = String(props.year)
  if (props.month == null) return `ทราบแต่ปีเกิด — พ.ศ. ${year}`
  const month = thaiMonths[props.month - 1]
  if (props.day == null) return `ทราบแต่เดือนและปีเกิด — ${month} ${year}`
  return `${props.day} ${month} ${year}`
})

const hasError = computed(() =>
  ['birth_year', 'birth_month', 'birth_day', 'birth_precision'].some(
    (f) => fieldErrors.value[`personal.${f}`],
  ),
)

const maxYear = toBE(new Date().getFullYear())

function setDay(value: string) {
  emit('update:day', value === '' ? null : Number(value))
  emit('update:precision', value === '' ? (props.month == null ? 'YEAR_ONLY' : 'YEAR_MONTH') : 'FULL')
}

function setMonth(value: string) {
  const month = value === '' ? null : Number(value)
  emit('update:month', month)
  if (month == null) {
    emit('update:day', null)
    emit('update:precision', 'YEAR_ONLY')
    return
  }
  if (props.day != null && props.day > new Date(toCE(props.year ?? maxYear), month, 0).getDate()) {
    emit('update:day', null)
    emit('update:precision', 'YEAR_MONTH')
    return
  }
  emit('update:precision', props.day == null ? 'YEAR_MONTH' : 'FULL')
}

function setYear(value: number | null) {
  emit('update:year', value)
  emit('update:precision', precision.value)
}
</script>

<template>
  <div>
    <div class="field-group-label">
      <InfoPopover trigger="ไม่ทราบวันเกิด?" label="อธิบายกรณีไม่ทราบวันเกิด">
        บัตรประชาชนรุ่นเก่าและทะเบียนราษฎรบางฉบับระบุแต่ปีเกิด
        เลือก <strong>“ไม่ทราบ”</strong> ที่ช่องวันหรือเดือนได้ ระบบรับใบสมัครตามปกติ
        และคิดอายุจากวันที่ 1 มกราคมของปีนั้นซึ่งเป็นผลดีกับผู้ยื่น
      </InfoPopover>
    </div>

    <div class="grid-form-3">
      <FieldSelect
        path="personal.birth_day"
        label="วัน"
        :model-value="day == null ? '' : String(day)"
        :options="dayOptions"
        :disabled="month == null"
        @update:model-value="setDay"
      />
      <FieldSelect
        path="personal.birth_month"
        label="เดือน"
        :model-value="month == null ? '' : String(month)"
        :options="monthOptions"
        @update:model-value="setMonth"
      />
      <FieldNumber
        path="personal.birth_year"
        label="ปี (พ.ศ.)"
        :model-value="year"
        :min="2443"
        :max="maxYear"
        placeholder="2528"
        @update:model-value="setYear"
      />
    </div>
    <p v-if="summary && !hasError" class="mt-1.5 text-sm text-ink-muted">
      วันเกิดที่บันทึก: {{ summary }}
    </p>
  </div>
</template>
