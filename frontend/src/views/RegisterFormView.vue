<script setup lang="ts">
import AddressSelector from '@/components/registration/AddressSelector.vue'
import AppButton from '@/components/common/AppButton.vue'
import ConfirmDialog from '@/components/common/ConfirmDialog.vue'
import DuplicateDialog from '@/components/registration/DuplicateDialog.vue'
import FamilySection from '@/components/registration/FamilySection.vue'
import FinancialSection from '@/components/registration/FinancialSection.vue'
import FormSection from '@/components/common/FormSection.vue'
import PersonalSection from '@/components/registration/PersonalSection.vue'
import UserSwitcher from '@/components/common/UserSwitcher.vue'
import { useRegistrationForm } from '@/composables/useRegistrationForm'
import { ref } from 'vue'

const { form, isSubmitting, duplicateOpen, duplicateInfo, formError, isDirty, reset, submit } =
  useRegistrationForm()

const cancelOpen = ref(false)

function onCancel() {
  if (!isDirty.value) return
  cancelOpen.value = true
}

function confirmCancel() {
  cancelOpen.value = false
  reset()
  window.scrollTo({ top: 0 })
}
</script>

<template>
  <div class="page-shell">
    <header class="topbar">
      <p class="brand">สวัสดิการแห่งรัฐ</p>
      <UserSwitcher />
    </header>

    <main>
      <div class="mb-6 animate-rise motion-reduce:animate-none">
        <h1 class="mt-1 mb-3 text-[clamp(1.75rem,3vw,2.35rem)] leading-tight">ลงทะเบียนผู้มีสิทธิ</h1>
      </div>

      <form @submit.prevent="submit">
        <PersonalSection :personal="form.personal" />

        <FormSection title="ที่อยู่">
          <AddressSelector v-model="form.personal.address" />
        </FormSection>

        <FamilySection v-model:include="form.includeFamily" :family="form.family" />

        <FinancialSection :financial="form.financial" />

        <p v-if="formError" class="banner banner-error">{{ formError }}</p>

        <div class="form-actions">
          <AppButton :disabled="isSubmitting" @click="onCancel">ยกเลิก</AppButton>
          <AppButton type="submit" variant="primary" :disabled="isSubmitting">
            {{ isSubmitting ? 'กำลังบันทึก…' : 'บันทึกใบสมัคร' }}
          </AppButton>
        </div>
      </form>
    </main>

    <DuplicateDialog :open="duplicateOpen" :info="duplicateInfo" @close="duplicateOpen = false" />

    <ConfirmDialog
      :open="cancelOpen"
      title="ล้างข้อมูลในฟอร์ม"
      message="ข้อมูลที่กรอกไว้ทั้งหมดจะหายไป ต้องการล้างฟอร์มใช่ไหม"
      confirm-label="ล้างข้อมูล"
      cancel-label="กรอกต่อ"
      danger
      @confirm="confirmCancel"
      @close="cancelOpen = false"
    />
  </div>
</template>
