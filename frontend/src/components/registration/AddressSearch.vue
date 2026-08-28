<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, ref, useId, watch } from 'vue'
import { searchAddress } from '@/api/ref'
import { useFormContext } from '@/composables/useFormContext'
import type { AddressOption } from '@/types/api'
import { cn } from '@/utils/cn'

const props = defineProps<{
  selected: AddressOption | null
}>()

const emit = defineEmits<{
  select: [value: AddressOption]
  clearSelection: []
}>()

const { fieldErrors, isSubmitting, clear } = useFormContext()

const path = 'personal.address.subdistrict_code'
const error = computed(() => fieldErrors.value[path] ?? '')

const query = ref('')
const options = ref<AddressOption[]>([])
const open = ref(false)
const loading = ref(false)
const failed = ref(false)
const active = ref(-1)
const listId = useId()
const input = ref<HTMLInputElement | null>(null)
const box = ref({ top: 0, left: 0, width: 0, maxHeight: 288 })

const DEBOUNCE_MS = 450

function place() {
  const el = input.value
  if (!el) return
  const r = el.getBoundingClientRect()
  const below = window.innerHeight - r.bottom - 12
  const above = r.top - 12
  const dropUp = below < 200 && above > below
  const maxHeight = Math.min(288, Math.max(120, dropUp ? above : below))
  box.value = {
    top: window.scrollY + (dropUp ? r.top - 4 - maxHeight : r.bottom + 4),
    left: window.scrollX + r.left,
    width: r.width,
    maxHeight,
  }
}

watch(open, async (isOpen) => {
  if (!isOpen) {
    window.removeEventListener('resize', place)
    return
  }
  await nextTick()
  place()
  window.addEventListener('resize', place)
})

function nameOf(o: AddressOption) {
  return `${o.subdistrict_kind}${o.subdistrict_name}`
}

let timer: ReturnType<typeof setTimeout> | undefined
let inflight: AbortController | undefined

watch(query, (q) => {
  if (props.selected && q === nameOf(props.selected)) return

  if (props.selected) emit('clearSelection')

  clearTimeout(timer)
  inflight?.abort()
  failed.value = false

  if (q.trim().length < 2) {
    options.value = []
    open.value = false
    loading.value = false
    return
  }

  loading.value = true
  open.value = true
  timer = setTimeout(async () => {
    const controller = new AbortController()
    inflight = controller
    try {
      options.value = await searchAddress(q.trim(), controller.signal)
      active.value = options.value.length > 0 ? 0 : -1
    } catch (err) {
      if (controller.signal.aborted) return
      options.value = []
      failed.value = true
    } finally {
      if (inflight === controller) {
        inflight = undefined
        loading.value = false
      }
    }
  }, DEBOUNCE_MS)
})

watch(
  () => props.selected,
  (o) => {
    if (o) query.value = nameOf(o)
    else if (options.value.length === 0) query.value = ''
  },
)

function choose(o: AddressOption) {
  emit('select', o)
  query.value = nameOf(o)
  clear(path)
  options.value = []
  open.value = false
  active.value = -1
}

function onKeydown(e: KeyboardEvent) {
  if (!open.value || options.value.length === 0) return
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    active.value = (active.value + 1) % options.value.length
  } else if (e.key === 'ArrowUp') {
    e.preventDefault()
    active.value = (active.value - 1 + options.value.length) % options.value.length
  } else if (e.key === 'Enter') {
    e.preventDefault()
    if (active.value >= 0) choose(options.value[active.value])
  } else if (e.key === 'Escape') {
    open.value = false
  }
}

onBeforeUnmount(() => {
  clearTimeout(timer)
  inflight?.abort()
  window.removeEventListener('resize', place)
})
</script>

<template>
  <label class="field relative">
    <span class="field-label">ตำบล / แขวง</span>
    <input
      ref="input"
      :class="cn('field-input', error && 'field-input-error')"
      v-model="query"
      type="text"
      autocomplete="off"
      placeholder="พิมพ์ชื่อตำบล อำเภอ หรือรหัสไปรษณีย์"
      role="combobox"
      aria-autocomplete="list"
      :aria-expanded="open"
      :aria-controls="listId"
      :aria-invalid="error ? true : undefined"
      :disabled="isSubmitting"
      @keydown="onKeydown"
      @blur="open = false"
      @focus="open = options.length > 0"
    />

    <Teleport to="body">
      <ul
        v-if="open"
        :id="listId"
        class="combo-list"
        role="listbox"
        :style="{
          top: `${box.top}px`,
          left: `${box.left}px`,
          minWidth: `${box.width}px`,
          maxHeight: `${box.maxHeight}px`,
        }"
      >
        <li v-if="loading" class="combo-empty">กำลังค้นหา…</li>
        <li v-else-if="failed" class="combo-empty">ค้นหาไม่สำเร็จ ลองใหม่อีกครั้ง</li>
        <li v-else-if="options.length === 0" class="combo-empty">ไม่พบตำบลที่ค้นหา</li>
        <li
          v-for="(o, idx) in options"
          v-else
          :key="o.subdistrict_code + o.postal_code"
          :class="cn('combo-item', idx === active && 'combo-item-active')"
          role="option"
          :aria-selected="idx === active"
          @mousedown.prevent="choose(o)"
          @mouseenter="active = idx"
        >
          <span class="font-semibold">{{ o.subdistrict_kind }}{{ o.subdistrict_name }}</span>
          <span class="text-ink-muted">
            {{ o.district_kind }}{{ o.district_name }} · {{ o.province_name }} · {{ o.postal_code }}
          </span>
        </li>
      </ul>
    </Teleport>

    <span v-if="error" class="field-error">{{ error }}</span>
  </label>
</template>
