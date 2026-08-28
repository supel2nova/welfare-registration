<script setup lang="ts">
import { computed, ref } from 'vue'
import AddressSearch from '@/components/registration/AddressSearch.vue'
import { useFormContext } from '@/composables/useFormContext'
import type { AddressInput, AddressOption } from '@/types/api'
import { cn } from '@/utils/cn'

const props = defineProps<{
  modelValue: AddressInput
}>()

const emit = defineEmits<{
  'update:modelValue': [value: AddressInput]
}>()

const { fieldErrors, isSubmitting, clear } = useFormContext()

const selected = ref<AddressOption | null>(null)

function errorOf(path: string) {
  return fieldErrors.value[path] ?? ''
}

const houseError = computed(() => errorOf('personal.address.house_no'))

const linkedError = computed(
  () =>
    errorOf('personal.address.district_code') ||
    errorOf('personal.address.province_code') ||
    errorOf('personal.address.postal_code'),
)

function patch(partial: Partial<AddressInput>, clearPath?: string) {
  emit('update:modelValue', { ...props.modelValue, ...partial })
  if (clearPath) clear(clearPath)
}

function onSelect(o: AddressOption) {
  selected.value = o
  patch({
    province_code: o.province_code,
    district_code: o.district_code,
    subdistrict_code: o.subdistrict_code,
    postal_code: o.postal_code,
  })
  for (const p of ['district_code', 'province_code', 'postal_code']) {
    clear(`personal.address.${p}`)
  }
}

function onClearSelection() {
  selected.value = null
  patch({ province_code: '', district_code: '', subdistrict_code: '', postal_code: '' })
}
</script>

<template>
  <div>
    <div class="grid-form-3">
      <label class="field">
        <span class="field-label">บ้านเลขที่</span>
        <input
          :class="cn('field-input', houseError && 'field-input-error')"
          :value="modelValue.house_no"
          placeholder="12/34"
          :disabled="isSubmitting"
          :aria-invalid="houseError ? true : undefined"
          @input="patch({ house_no: ($event.target as HTMLInputElement).value }, 'personal.address.house_no')"
        />
        <span v-if="houseError" class="field-error">{{ houseError }}</span>
      </label>

      <label class="field">
        <span class="field-label">หมู่</span>
        <input
          class="field-input"
          :value="modelValue.moo ?? ''"
          placeholder="5"
          :disabled="isSubmitting"
          @input="patch({ moo: ($event.target as HTMLInputElement).value || null })"
        />
      </label>

      <label class="field">
        <span class="field-label">ถนน</span>
        <input
          class="field-input"
          :value="modelValue.road ?? ''"
          placeholder="ห้วยแก้ว (ถ้ามี)"
          :disabled="isSubmitting"
          @input="patch({ road: ($event.target as HTMLInputElement).value || null })"
        />
      </label>
    </div>

    <div class="mt-3.5 grid-form-2">
      <AddressSearch :selected="selected" @select="onSelect" @clear-selection="onClearSelection" />

      <div class="field">
        <span class="field-label">อำเภอ / เขต</span>
        <p :class="cn('field-static', !selected && 'field-static-empty')">
          {{ selected ? `${selected.district_kind}${selected.district_name}` : 'เลือกจากช่องตำบล' }}
        </p>
      </div>

      <div class="field">
        <span class="field-label">จังหวัด</span>
        <p :class="cn('field-static', !selected && 'field-static-empty')">
          {{ selected?.province_name ?? 'เลือกจากช่องตำบล' }}
        </p>
      </div>

      <div class="field">
        <span class="field-label">รหัสไปรษณีย์</span>
        <p :class="cn('field-static tabular-nums', !modelValue.postal_code && 'field-static-empty')">
          {{ modelValue.postal_code || 'เลือกจากช่องตำบล' }}
        </p>
      </div>
    </div>

    <p v-if="linkedError" class="field-error mt-1.5">{{ linkedError }}</p>
  </div>
</template>
