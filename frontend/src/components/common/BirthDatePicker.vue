<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import InfoPopover from '@/components/common/InfoPopover.vue'
import { useFormContext } from '@/composables/useFormContext'
import { toBE, toCE } from '@/utils/buddhist'
import { cn } from '@/utils/cn'

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

const { fieldErrors, isSubmitting, clear } = useFormContext()

const paths = ['personal.birth_year', 'personal.birth_month', 'personal.birth_day']
const error = computed(() => paths.map((p) => fieldErrors.value[p]).find(Boolean) ?? '')

const monthNames = [
  'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
  'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม',
]
const monthShort = ['ม.ค.', 'ก.พ.', 'มี.ค.', 'เม.ย.', 'พ.ค.', 'มิ.ย.', 'ก.ค.', 'ส.ค.', 'ก.ย.', 'ต.ค.', 'พ.ย.', 'ธ.ค.']
const dayNames = ['อา', 'จ', 'อ', 'พ', 'พฤ', 'ศ', 'ส']

const MIN_YEAR_BE = 2443
const now = new Date()
const maxYearBE = toBE(now.getFullYear())

const open = ref(false)
const view = ref<'year' | 'month' | 'day'>('year')
const draftYear = ref<number | null>(null)
const draftMonth = ref<number | null>(null)
const gridStart = ref(maxYearBE - 11)
const panelId = useId()
const trigger = ref<HTMLButtonElement | null>(null)
const panel = ref<HTMLElement | null>(null)
const box = ref({ top: 0, left: 0, width: 320 })

const label = computed(() => {
  if (props.year == null) return ''
  if (props.month == null) return `พ.ศ. ${props.year}`
  if (props.day == null) return `${monthNames[props.month - 1]} ${props.year}`
  return `${props.day} ${monthNames[props.month - 1]} ${props.year}`
})

function daysIn(yearBE: number, month: number) {
  return new Date(toCE(yearBE), month, 0).getDate()
}

function isFuture(yearBE: number, month?: number, day?: number) {
  const d = new Date(toCE(yearBE), (month ?? 1) - 1, day ?? 1)
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  return d > today
}

const years = computed(() =>
  Array.from({ length: 12 }, (_, i) => gridStart.value + i).filter((y) => y >= MIN_YEAR_BE),
)

const leadingBlanks = computed(() => {
  if (draftYear.value == null || draftMonth.value == null) return 0
  return new Date(toCE(draftYear.value), draftMonth.value - 1, 1).getDay()
})

const days = computed(() => {
  if (draftYear.value == null || draftMonth.value == null) return []
  return Array.from({ length: daysIn(draftYear.value, draftMonth.value) }, (_, i) => i + 1)
})

function place() {
  const el = trigger.value
  if (!el) return
  const r = el.getBoundingClientRect()
  const width = 320
  const height = panel.value?.offsetHeight ?? 280
  const below = window.innerHeight - r.bottom - 12
  const up = below < height && r.top - 12 > below
  box.value = {
    top: window.scrollY + (up ? Math.max(0, r.top - 8 - height) : r.bottom + 4),
    left: window.scrollX + Math.max(12, Math.min(r.left, window.innerWidth - width - 12)),
    width,
  }
}

async function toggle() {
  open.value = !open.value
  if (!open.value) return
  draftYear.value = props.year
  draftMonth.value = props.month
  view.value = props.year == null ? 'year' : props.month == null ? 'month' : 'day'
  gridStart.value = Math.min(maxYearBE - 11, (props.year ?? maxYearBE - 30) - 6)
  await nextTick()
  place()
  await nextTick()
  place()
  window.addEventListener('resize', place)
}

function close() {
  open.value = false
  window.removeEventListener('resize', place)
}

function clearErrors() {
  for (const p of paths) clear(p)
}

function pickYear(y: number) {
  draftYear.value = y
  view.value = 'month'
}

function pickMonth(m: number) {
  draftMonth.value = m
  view.value = 'day'
}

function pickDay(d: number) {
  emit('update:year', draftYear.value)
  emit('update:month', draftMonth.value)
  emit('update:day', d)
  emit('update:precision', 'FULL')
  clearErrors()
  close()
}

function useYearOnly() {
  emit('update:year', draftYear.value)
  emit('update:month', null)
  emit('update:day', null)
  emit('update:precision', 'YEAR_ONLY')
  clearErrors()
  close()
}

function useMonthOnly() {
  emit('update:year', draftYear.value)
  emit('update:month', draftMonth.value)
  emit('update:day', null)
  emit('update:precision', 'YEAR_MONTH')
  clearErrors()
  close()
}

function reset() {
  emit('update:year', null)
  emit('update:month', null)
  emit('update:day', null)
  emit('update:precision', 'YEAR_ONLY')
  clearErrors()
  close()
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && open.value) close()
}

watch(view, async () => {
  if (!open.value) return
  await nextTick()
  place()
})

watch(open, (isOpen) => {
  if (isOpen) document.addEventListener('keydown', onKeydown)
  else document.removeEventListener('keydown', onKeydown)
})

onBeforeUnmount(() => {
  document.removeEventListener('keydown', onKeydown)
  window.removeEventListener('resize', place)
})
</script>

<template>
  <div class="field">
    <div class="flex items-center justify-between gap-2">
      <span class="field-label">วัน/เดือน/ปีเกิด</span>
      <InfoPopover trigger="ไม่ทราบวันเกิด?" label="อธิบายกรณีไม่ทราบวันเกิด">
        บัตรประชาชนรุ่นเก่าและทะเบียนราษฎรบางฉบับระบุแต่ปีเกิด
        เลือกปีแล้วกด <strong>“ทราบแค่ปี…”</strong> หรือเลือกถึงเดือนแล้วกด
        <strong>“ทราบแค่เดือน…”</strong> ได้เลย ระบบรับใบสมัครตามปกติ
        และคิดอายุจากวันที่ 1 มกราคมของปีนั้นซึ่งเป็นผลดีกับผู้ยื่น
      </InfoPopover>
    </div>
    <button
      ref="trigger"
      type="button"
      :class="cn('field-input text-left', error && 'field-input-error', !label && 'text-neutral-400')"
      :disabled="isSubmitting"
      :aria-expanded="open"
      :aria-controls="panelId"
      :aria-invalid="error ? true : undefined"
      @click="toggle"
    >
      {{ label || 'เลือกวันเกิด' }}
    </button>
    <span v-if="error" class="field-error">{{ error }}</span>

    <Teleport to="body">
      <div v-if="open" class="picker-backdrop" @click="close" />
      <div
        v-if="open"
        ref="panel"
        :id="panelId"
        class="picker"
        role="dialog"
        :style="{ top: `${box.top}px`, left: `${box.left}px`, width: `${box.width}px` }"
      >
        <div class="picker-head">
          <button
            v-if="view !== 'year'"
            type="button"
            class="picker-nav"
            @click="view = view === 'day' ? 'month' : 'year'"
          >
            ‹ ย้อนกลับ
          </button>
          <span v-else class="picker-nav-static">เลือกปีเกิด (พ.ศ.)</span>
          <span class="picker-title">
            {{ view === 'year' ? `${gridStart}–${gridStart + 11}` : '' }}
            {{ view === 'month' ? draftYear : '' }}
            {{ view === 'day' && draftMonth ? `${monthNames[draftMonth - 1]} ${draftYear}` : '' }}
          </span>
        </div>

        <div v-if="view === 'year'" class="picker-grid picker-grid-3">
          <button
            type="button"
            class="picker-nav"
            :disabled="gridStart <= MIN_YEAR_BE"
            @click="gridStart = Math.max(MIN_YEAR_BE, gridStart - 12)"
          >
            ‹ ก่อนหน้า
          </button>
          <span />
          <button
            type="button"
            class="picker-nav"
            :disabled="gridStart + 12 > maxYearBE"
            @click="gridStart = Math.min(maxYearBE - 11, gridStart + 12)"
          >
            ถัดไป ›
          </button>
          <button
            v-for="y in years"
            :key="y"
            type="button"
            :class="cn('picker-cell', y === year && 'picker-cell-on')"
            :disabled="y > maxYearBE"
            @click="pickYear(y)"
          >
            {{ y }}
          </button>
        </div>

        <div v-else-if="view === 'month'" class="picker-grid picker-grid-3">
          <button
            v-for="(m, idx) in monthShort"
            :key="m"
            type="button"
            :class="cn('picker-cell', idx + 1 === month && draftYear === year && 'picker-cell-on')"
            :disabled="draftYear != null && isFuture(draftYear, idx + 1)"
            @click="pickMonth(idx + 1)"
          >
            {{ m }}
          </button>
        </div>

        <div v-else class="picker-grid picker-grid-7">
          <span v-for="d in dayNames" :key="d" class="picker-dayname">{{ d }}</span>
          <span v-for="b in leadingBlanks" :key="'b' + b" />
          <button
            v-for="d in days"
            :key="d"
            type="button"
            :class="cn('picker-cell', d === day && draftMonth === month && draftYear === year && 'picker-cell-on')"
            :disabled="draftYear != null && draftMonth != null && isFuture(draftYear, draftMonth, d)"
            @click="pickDay(d)"
          >
            {{ d }}
          </button>
        </div>

        <div class="picker-foot">
          <button v-if="year != null || month != null" type="button" class="picker-nav" @click="reset">
            ล้าง
          </button>
          <span v-else />
          <button
            v-if="view === 'month' && draftYear != null"
            type="button"
            class="picker-confirm"
            @click="useYearOnly"
          >
            ทราบแค่ปี {{ draftYear }}
          </button>
          <button
            v-else-if="view === 'day' && draftMonth != null"
            type="button"
            class="picker-confirm"
            @click="useMonthOnly"
          >
            ทราบแค่ {{ monthNames[draftMonth - 1] }} {{ draftYear }}
          </button>
          <span v-else />
        </div>
      </div>
    </Teleport>
  </div>
</template>
