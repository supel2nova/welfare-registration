<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useFormContext } from '@/composables/useFormContext'
import { useRefStore } from '@/stores/ref'
import type { AddressInput } from '@/types/api'
import { cn } from '@/utils/cn'

const props = defineProps<{
  modelValue: AddressInput
}>()

const emit = defineEmits<{
  'update:modelValue': [value: AddressInput]
}>()

const refStore = useRefStore()
const { fieldErrors, isSubmitting, clear } = useFormContext()

function errorOf(path: string) {
  return fieldErrors.value[path] ?? ''
}

onMounted(() => {
  void refStore.loadProvinces()
})

watch(
  () => props.modelValue.province_code,
  async (code, prev) => {
    if (code === prev) return
    await refStore.loadDistricts(code)
    if (prev !== undefined) {
      emit('update:modelValue', {
        ...props.modelValue,
        district_code: '',
        subdistrict_code: '',
      })
    }
  },
)

watch(
  () => props.modelValue.district_code,
  async (code, prev) => {
    if (code === prev) return
    await refStore.loadSubdistricts(code)
    if (prev !== undefined) {
      emit('update:modelValue', {
        ...props.modelValue,
        subdistrict_code: '',
      })
    }
  },
)

function patch(partial: Partial<AddressInput>, clearPath?: string) {
  emit('update:modelValue', { ...props.modelValue, ...partial })
  if (clearPath) clear(clearPath)
}
</script>

<template>
  <div class="grid-form-2">
    <label class="field">
      <span class="field-label">บ้านเลขที่</span>
      <input
        :class="cn('field-input', errorOf('personal.address.house_no') && 'field-input-error')"
        :value="modelValue.house_no"
        :disabled="isSubmitting"
        @input="patch({ house_no: ($event.target as HTMLInputElement).value }, 'personal.address.house_no')"
      />
      <span v-if="errorOf('personal.address.house_no')" class="field-error">{{
        errorOf('personal.address.house_no')
      }}</span>
    </label>

    <label class="field">
      <span class="field-label">หมู่</span>
      <input
        class="field-input"
        :value="modelValue.moo ?? ''"
        :disabled="isSubmitting"
        @input="patch({ moo: ($event.target as HTMLInputElement).value || null })"
      />
    </label>

    <label class="field">
      <span class="field-label">ถนน</span>
      <input
        class="field-input"
        :value="modelValue.road ?? ''"
        :disabled="isSubmitting"
        @input="patch({ road: ($event.target as HTMLInputElement).value || null })"
      />
    </label>

    <label class="field">
      <span class="field-label">จังหวัด</span>
      <select
        :class="cn('field-input', 'field-select', errorOf('personal.address.province_code') && 'field-input-error')"
        :value="modelValue.province_code"
        :disabled="isSubmitting"
        @change="
          patch(
            { province_code: ($event.target as HTMLSelectElement).value },
            'personal.address.province_code',
          )
        "
      >
        <option value="">เลือกจังหวัด</option>
        <option v-for="p in refStore.provinces" :key="p.code" :value="p.code">{{ p.name_th }}</option>
      </select>
      <span v-if="errorOf('personal.address.province_code')" class="field-error">{{
        errorOf('personal.address.province_code')
      }}</span>
    </label>

    <label class="field">
      <span class="field-label">อำเภอ / เขต</span>
      <select
        :class="cn('field-input', 'field-select', errorOf('personal.address.district_code') && 'field-input-error')"
        :value="modelValue.district_code"
        :disabled="isSubmitting || !modelValue.province_code"
        @change="
          patch(
            { district_code: ($event.target as HTMLSelectElement).value },
            'personal.address.district_code',
          )
        "
      >
        <option value="">เลือกอำเภอ</option>
        <option v-for="d in refStore.districts" :key="d.code" :value="d.code">{{ d.name_th }}</option>
      </select>
      <span v-if="errorOf('personal.address.district_code')" class="field-error">{{
        errorOf('personal.address.district_code')
      }}</span>
    </label>

    <label class="field">
      <span class="field-label">ตำบล / แขวง</span>
      <select
        :class="cn('field-input', 'field-select', errorOf('personal.address.subdistrict_code') && 'field-input-error')"
        :value="modelValue.subdistrict_code"
        :disabled="isSubmitting || !modelValue.district_code"
        @change="
          patch(
            { subdistrict_code: ($event.target as HTMLSelectElement).value },
            'personal.address.subdistrict_code',
          )
        "
      >
        <option value="">เลือกตำบล</option>
        <option v-for="s in refStore.subdistricts" :key="s.code" :value="s.code">{{ s.name_th }}</option>
      </select>
      <span v-if="errorOf('personal.address.subdistrict_code')" class="field-error">{{
        errorOf('personal.address.subdistrict_code')
      }}</span>
    </label>

    <label class="field">
      <span class="field-label">รหัสไปรษณีย์</span>
      <input
        :class="cn('field-input', errorOf('personal.address.postal_code') && 'field-input-error')"
        maxlength="5"
        :value="modelValue.postal_code"
        :disabled="isSubmitting"
        @input="
          patch(
            { postal_code: ($event.target as HTMLInputElement).value },
            'personal.address.postal_code',
          )
        "
      />
      <span v-if="errorOf('personal.address.postal_code')" class="field-error">{{
        errorOf('personal.address.postal_code')
      }}</span>
    </label>
  </div>
</template>
