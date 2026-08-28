<script setup lang="ts">
import AppButton from '@/components/common/AppButton.vue'
import AssetRow from '@/components/registration/AssetRow.vue'
import FieldNumber from '@/components/fields/FieldNumber.vue'
import FormSection from '@/components/common/FormSection.vue'
import IncomeRow from '@/components/registration/IncomeRow.vue'
import LiabilityRow from '@/components/registration/LiabilityRow.vue'
import { useFormContext } from '@/composables/useFormContext'
import { emptyAsset, emptyIncome, emptyLiability } from '@/constants/options'
import type { FinancialForm } from '@/types/form'

const props = defineProps<{
  financial: FinancialForm
}>()

const { isSubmitting, clearAll } = useFormContext()

function removeFrom(list: unknown[], idx: number) {
  list.splice(idx, 1)
  clearAll()
}
</script>

<template>
  <FormSection title="รายได้และทรัพย์สิน">
    <h3 class="sub-head">แหล่งรายได้</h3>
    <p v-if="financial.income_sources.length === 0" class="empty-note">
      ยังไม่ได้ระบุ — ไม่มีรายได้ก็ส่งใบสมัครได้
    </p>
    <IncomeRow
      v-for="(s, idx) in financial.income_sources"
      :key="'inc-' + idx"
      :income="s"
      :index="idx"
      @remove="removeFrom(financial.income_sources, idx)"
    />
    <div class="btn-row">
      <AppButton :disabled="isSubmitting" @click="financial.income_sources.push(emptyIncome())">
        เพิ่มแหล่งรายได้
      </AppButton>
    </div>

    <div class="mt-5">
      <FieldNumber
        path="financial.expense_to_others"
        label="เงินส่งคนอื่นต่อปี"
        :model-value="financial.expense_to_others"
        format="amount"
        suffix="บาท"
        @update:model-value="(v) => (financial.expense_to_others = v ?? 0)"
      />
    </div>

    <h3 class="sub-head mt-5">ทรัพย์สิน</h3>
    <p v-if="financial.assets.length === 0" class="empty-note">ยังไม่ได้ระบุ</p>
    <AssetRow
      v-for="(a, idx) in financial.assets"
      :key="'ast-' + idx"
      :asset="a"
      :index="idx"
      @remove="removeFrom(financial.assets, idx)"
    />
    <div class="btn-row">
      <AppButton :disabled="isSubmitting" @click="financial.assets.push(emptyAsset())">
        เพิ่มทรัพย์สิน
      </AppButton>
    </div>

    <h3 class="sub-head mt-5">หนี้สิน</h3>
    <p v-if="financial.liabilities.length === 0" class="empty-note">ยังไม่ได้ระบุ</p>
    <LiabilityRow
      v-for="(l, idx) in financial.liabilities"
      :key="'liab-' + idx"
      :liability="l"
      :index="idx"
      @remove="removeFrom(financial.liabilities, idx)"
    />
    <div class="btn-row">
      <AppButton :disabled="isSubmitting" @click="financial.liabilities.push(emptyLiability())">
        เพิ่มหนี้สิน
      </AppButton>
    </div>

    <label class="check-row">
      <input
        type="checkbox"
        class="check-box"
        v-model="financial.has_credit_card"
        :disabled="isSubmitting"
      />
      <span>มีบัตรเครดิต</span>
    </label>
  </FormSection>
</template>
