<script setup lang="ts">
import BirthDatePicker from '@/components/common/BirthDatePicker.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import FieldText from '@/components/fields/FieldText.vue'
import FormSection from '@/components/common/FormSection.vue'
import { useFormContext } from '@/composables/useFormContext'
import { idVerifyReasonOptions, titleOptions } from '@/constants/options'
import type { PersonalForm } from '@/types/form'

const props = defineProps<{
  personal: PersonalForm
}>()

const { isSubmitting, clear } = useFormContext()

function onNoLaserChange() {
  if (props.personal.no_laser) {
    props.personal.laser_id = null
  } else {
    props.personal.id_verify_note = null
    if (props.personal.laser_id == null) props.personal.laser_id = ''
  }
  props.personal.id_verify_reason = ''
  clear('personal.laser_id')
  clear('personal.id_verify_reason')
  clear('personal.id_verify_note')

}

function onReasonChange(value: string) {
  props.personal.id_verify_reason = value
  if (value !== 'OTHER') {
    props.personal.id_verify_note = null
    clear('personal.id_verify_note')
  }
}
</script>

<template>
  <FormSection title="ข้อมูลส่วนตัว">
    <div class="grid-form-3">
      <FieldSelect
        path="personal.title"
        required
        label="คำนำหน้า"
        v-model="personal.title"
        :options="titleOptions"
      />
      <FieldText
        path="personal.first_name"
        required
        label="ชื่อ"
        v-model="personal.first_name"
        placeholder="สมชาย"
      />
      <FieldText
        path="personal.last_name"
        required
        label="นามสกุล"
        v-model="personal.last_name"
        placeholder="ใจดี"
      />
    </div>

    <div class="grid-form-2 mt-3.5">
      <FieldText
        path="personal.national_id"
        required
        label="เลขประจำตัวประชาชน"
        v-model="personal.national_id"
        format="national_id"
        placeholder="1-2345-67890-12-1"
      />
      <BirthDatePicker
        :year="personal.birth_year_be"
        :month="personal.birth_month"
        :day="personal.birth_day"
        @update:year="(v) => (personal.birth_year_be = v)"
        @update:month="(v) => (personal.birth_month = v)"
        @update:day="(v) => (personal.birth_day = v)"
        @update:precision="(v) => (personal.birth_precision = v)"
      />
      <FieldText
        path="personal.phone"
        required
        label="เบอร์โทร"
        v-model="personal.phone"
        format="phone"
        placeholder="099-119-2231"
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
        required
      label="รหัสหลังบัตร (Laser)"
      v-model="personal.laser_id"
      format="laser_id"
      placeholder="JT8-1234567-89"
    />
    <template v-else>
      <FieldSelect
        path="personal.id_verify_reason"
        required
        label="เหตุผลที่ไม่ใช้รหัสหลังบัตร"
        :model-value="personal.id_verify_reason"
        :options="idVerifyReasonOptions"
        @update:model-value="onReasonChange"
      />
      <FieldText
        v-if="personal.id_verify_reason === 'OTHER'"
        path="personal.id_verify_note"
        required
        label="ระบุเหตุผล"
        v-model="personal.id_verify_note"
        placeholder="อย่างน้อย 10 ตัวอักษร"
      />
    </template>
  </FormSection>
</template>
