<script setup lang="ts">
import AppButton from '@/components/common/AppButton.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import FormSection from '@/components/common/FormSection.vue'
import MemberRow from '@/components/registration/MemberRow.vue'
import { useFormContext } from '@/composables/useFormContext'
import { emptyMember, maritalOptions } from '@/constants/options'
import type { FamilyForm } from '@/types/form'

const props = defineProps<{
  family: FamilyForm
  include: boolean
}>()

const emit = defineEmits<{
  'update:include': [value: boolean]
}>()

const { isSubmitting, clearAll } = useFormContext()

function removeMember(idx: number) {
  props.family.members.splice(idx, 1)
  clearAll()
}
</script>

<template>
  <FormSection title="ครอบครัว">
    <template #header>
      <label class="flex items-start gap-2.5 text-[0.95rem]">
        <input
          type="checkbox"
          class="check-box"
          :checked="include"
          :disabled="isSubmitting"
          @change="emit('update:include', ($event.target as HTMLInputElement).checked)"
        />
        <span>ระบุข้อมูลครอบครัว</span>
      </label>
    </template>

    <template v-if="include">
      <FieldSelect
        path="family.marital_status"
        label="สถานภาพสมรส"
        v-model="family.marital_status"
        :options="maritalOptions"
      />
      <p v-if="family.members.length === 0" class="empty-note">ยังไม่ได้ระบุสมาชิก</p>
      <MemberRow
        v-for="(m, idx) in family.members"
        :key="idx"
        :member="m"
        :index="idx"
        @remove="removeMember(idx)"
      />
      <div class="btn-row">
        <AppButton :disabled="isSubmitting" @click="family.members.push(emptyMember())">
          เพิ่มสมาชิก
        </AppButton>
      </div>
    </template>
  </FormSection>
</template>
