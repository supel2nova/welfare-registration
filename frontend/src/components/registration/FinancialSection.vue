<script setup lang="ts">
import AppButton from '@/components/common/AppButton.vue'
import AppIcon from '@/components/common/AppIcon.vue'
import FieldNumber from '@/components/fields/FieldNumber.vue'
import FormSection from '@/components/common/FormSection.vue'
import TypeAmountRow from '@/components/common/TypeAmountRow.vue'
import { useFormContext } from '@/composables/useFormContext'
import {
  assetTypeOptions,
  assetUnits,
  emptyAsset,
  emptyIncome,
  emptyLiability,
  incomeTypeOptions,
  liabilityTypeOptions,
  unitLabels,
} from '@/constants/options'
import type { AssetRow, FinancialForm } from '@/types/form'

defineProps<{
  financial: FinancialForm
}>()

const { isSubmitting, clearPrefix } = useFormContext()

function removeFrom(list: unknown[], idx: number, prefix: string) {
  list.splice(idx, 1)
  clearPrefix(prefix)
}

function setAssetType(asset: AssetRow, value: string) {
  asset.asset_type = value
  asset.unit = assetUnits[value] ?? 'THB'
}
</script>

<template>
  <FormSection title="รายได้และทรัพย์สิน">
    <h3 class="sub-head">แหล่งรายได้</h3>
    <p v-if="financial.income_sources.length === 0" class="empty-note">
      ยังไม่ได้ระบุ — ไม่มีรายได้ก็ส่งใบสมัครได้
    </p>
    <TypeAmountRow
      v-for="(s, idx) in financial.income_sources"
      :key="'inc-' + idx"
      :title="`แหล่งรายได้ที่ ${idx + 1}`"
      :type-path="`financial.income_sources[${idx}].source_type`"
      type-label="แหล่งรายได้"
      :type-options="incomeTypeOptions"
      :type="s.source_type"
      :amount-path="`financial.income_sources[${idx}].annual_amount`"
      amount-label="รายได้ต่อปี"
      :amount="s.annual_amount"
      suffix="บาท"
      @update:type="(v) => (s.source_type = v)"
      @update:amount="(v) => (s.annual_amount = v)"
      @remove="removeFrom(financial.income_sources, idx, 'financial.income_sources[')"
    />
    <AppButton variant="add" :disabled="isSubmitting" @click="financial.income_sources.push(emptyIncome())">
      <AppIcon name="plus" />
      เพิ่มแหล่งรายได้
    </AppButton>

    <div class="mt-5">
      <FieldNumber
        path="financial.expense_to_others"
        required
        label="ค่าเลี้ยงดูผู้อื่นต่อปี"
        :model-value="financial.expense_to_others"
        format="amount"
        placeholder="0"
        suffix="บาท"
        @update:model-value="(v) => (financial.expense_to_others = v ?? 0)"
      />
    </div>

    <h3 class="sub-head mt-5">ทรัพย์สิน</h3>
    <p v-if="financial.assets.length === 0" class="empty-note">ยังไม่ได้ระบุ</p>
    <TypeAmountRow
      v-for="(a, idx) in financial.assets"
      :key="'ast-' + idx"
      :title="`ทรัพย์สินรายการที่ ${idx + 1}`"
      :type-path="`financial.assets[${idx}].asset_type`"
      type-label="ทรัพย์สิน"
      :type-options="assetTypeOptions"
      :type="a.asset_type"
      :amount-path="`financial.assets[${idx}].amount`"
      amount-label="จำนวน"
      :amount="Number(a.amount)"
      :suffix="unitLabels[a.unit] ?? a.unit"
      @update:type="(v) => setAssetType(a, v)"
      @update:amount="(v) => (a.amount = v)"
      @remove="removeFrom(financial.assets, idx, 'financial.assets[')"
    />
    <AppButton variant="add" :disabled="isSubmitting" @click="financial.assets.push(emptyAsset())">
      <AppIcon name="plus" />
      เพิ่มรายการทรัพย์สิน
    </AppButton>

    <h3 class="sub-head mt-5">หนี้สิน</h3>
    <p v-if="financial.liabilities.length === 0" class="empty-note">ยังไม่ได้ระบุ</p>
    <TypeAmountRow
      v-for="(l, idx) in financial.liabilities"
      :key="'liab-' + idx"
      :title="`หนี้สินรายการที่ ${idx + 1}`"
      :type-path="`financial.liabilities[${idx}].liability_type`"
      type-label="หนี้สิน"
      :type-options="liabilityTypeOptions"
      :type="l.liability_type"
      :amount-path="`financial.liabilities[${idx}].credit_limit`"
      amount-label="วงเงิน"
      :amount="l.credit_limit"
      suffix="บาท"
      @update:type="(v) => (l.liability_type = v)"
      @update:amount="(v) => (l.credit_limit = v)"
      @remove="removeFrom(financial.liabilities, idx, 'financial.liabilities[')"
    />
    <AppButton variant="add" :disabled="isSubmitting" @click="financial.liabilities.push(emptyLiability())">
      <AppIcon name="plus" />
      เพิ่มรายการหนี้สิน
    </AppButton>

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
