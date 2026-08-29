<script setup lang="ts">
import AppButton from '@/components/common/AppButton.vue'
import AppIcon from '@/components/common/AppIcon.vue'
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

const { isSubmitting, clearPrefix } = useFormContext()

function removeFrom(list: unknown[], idx: number, prefix: string) {
  list.splice(idx, 1)
  clearPrefix(prefix)
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
    <AssetRow
      v-for="(a, idx) in financial.assets"
      :key="'ast-' + idx"
      :asset="a"
      :index="idx"
      @remove="removeFrom(financial.assets, idx, 'financial.assets[')"
    />
          <AppButton variant="add" :disabled="isSubmitting" @click="financial.assets.push(emptyAsset())">
        <AppIcon name="plus" />
        เพิ่มรายการทรัพย์สิน
      </AppButton>

    <h3 class="sub-head mt-5">หนี้สิน</h3>
    <p v-if="financial.liabilities.length === 0" class="empty-note">ยังไม่ได้ระบุ</p>
    <LiabilityRow
      v-for="(l, idx) in financial.liabilities"
      :key="'liab-' + idx"
      :liability="l"
      :index="idx"
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
