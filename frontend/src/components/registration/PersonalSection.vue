<script setup lang="ts">
import FieldNumber from '@/components/fields/FieldNumber.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import FieldText from '@/components/fields/FieldText.vue'
import FormSection from '@/components/common/FormSection.vue'
import { useFormContext } from '@/composables/useFormContext'
import { precisionOptions, titleOptions } from '@/constants/options'
import type { PersonalForm } from '@/types/form'

const props = defineProps<{
  personal: PersonalForm
  fiscalYear: number
}>()

const emit = defineEmits<{
  'update:fiscalYear': [value: number]
}>()

const { isSubmitting, clearAll } = useFormContext()

function onNoLaserChange() {
  if (props.personal.no_laser) {
    props.personal.laser_id = null
  } else {
    props.personal.id_verify_note = null
    if (props.personal.laser_id == null) props.personal.laser_id = ''
  }
  clearAll()
}
</script>

<template>
  <FormSection title="ข้อมูลส่วนตัว">
    <div class="grid-form-2">
      <FieldNumber
        path="fiscal_year"
        label="ปีงบประมาณ (ค.ศ.)"
        :model-value="fiscalYear"
        :min="2020"
        @update:model-value="(v) => emit('update:fiscalYear', v ?? 2026)"
      />
      <FieldSelect
        path="personal.title"
        label="คำนำหน้า"
        v-model="personal.title"
        :options="titleOptions"
      />
      <FieldText
        path="personal.national_id"
        label="เลขประจำตัวประชาชน"
        v-model="personal.national_id"
        format="national_id"
        placeholder="1-2345-67890-12-1"
      />
      <FieldText
        path="personal.phone"
        label="เบอร์โทร"
        v-model="personal.phone"
        format="phone"
        placeholder="099-119-2231"
      />
      <FieldText path="personal.first_name" label="ชื่อ" v-model="personal.first_name" />
      <FieldText path="personal.last_name" label="นามสกุล" v-model="personal.last_name" />
    </div>

    <div class="mt-3.5 grid-form-3">
      <FieldNumber
        path="personal.birth_year"
        label="ปีเกิด (พ.ศ.)"
        :model-value="personal.birth_year_be"
        :min="2443"
        @update:model-value="(v) => (personal.birth_year_be = v ?? 2528)"
      />
      <FieldNumber
        path="personal.birth_month"
        label="เดือน"
        v-model="personal.birth_month"
        :min="1"
        :max="12"
        :disabled="personal.birth_precision === 'YEAR_ONLY'"
      />
      <FieldNumber
        path="personal.birth_day"
        label="วัน"
        v-model="personal.birth_day"
        :min="1"
        :max="31"
        :disabled="personal.birth_precision !== 'FULL'"
      />
    </div>

    <div class="mt-3.5">
      <FieldSelect
        path="personal.birth_precision"
        label="ความละเอียดวันเกิด"
        v-model="personal.birth_precision"
        :options="precisionOptions"
      />
    </div>

    <label class="check-row">
      <input type="checkbox" class="check-box" v-model="personal.is_farmer" :disabled="isSubmitting" />
      <span>เป็นเกษตรกร</span>
    </label>

    <label class="check-row">
      <input
        type="checkbox"
        class="check-box"
        v-model="personal.no_laser"
        :disabled="isSubmitting"
        @change="onNoLaserChange"
      />
      <span>บัตรไม่มีรหัสหลังบัตร / อ่านไม่ออก (เช่น บัตรตลอดชีพ)</span>
    </label>

    <FieldText
      v-if="!personal.no_laser"
      path="personal.laser_id"
      label="รหัสหลังบัตร (Laser)"
      v-model="personal.laser_id"
      format="laser_id"
      placeholder="JT8-1234567-89"
    />
    <FieldText
      v-else
      path="personal.id_verify_note"
      label="เหตุผลที่ยืนยันด้วยตา"
      v-model="personal.id_verify_note"
      placeholder="อย่างน้อย 10 ตัวอักษร"
    />
  </FormSection>
</template>
