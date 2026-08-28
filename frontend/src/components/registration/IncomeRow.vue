<script setup lang="ts">
import { computed } from 'vue'
import FieldNumber from '@/components/fields/FieldNumber.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import ItemBlock from '@/components/common/ItemBlock.vue'
import { incomeTypeOptions } from '@/constants/options'
import type { IncomeRow } from '@/types/form'

const props = defineProps<{
  income: IncomeRow
  index: number
}>()

const emit = defineEmits<{
  remove: []
}>()

const base = computed(() => `financial.income_sources[${props.index}]`)
</script>

<template>
  <ItemBlock :title="`แหล่งรายได้ที่ ${index + 1}`" @remove="emit('remove')">
    <div class="grid-form-2">
      <FieldSelect
        :path="`${base}.source_type`"
        label="แหล่งรายได้"
        v-model="income.source_type"
        :options="incomeTypeOptions"
      />
      <FieldNumber
        :path="`${base}.annual_amount`"
        label="รายได้ต่อปี"
        :model-value="income.annual_amount"
        format="amount"
        placeholder="0"
        suffix="บาท"
        @update:model-value="(v) => (income.annual_amount = v ?? 0)"
      />
    </div>
  </ItemBlock>
</template>
