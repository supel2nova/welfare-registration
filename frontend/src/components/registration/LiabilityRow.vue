<script setup lang="ts">
import { computed } from 'vue'
import FieldNumber from '@/components/fields/FieldNumber.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import ItemBlock from '@/components/common/ItemBlock.vue'
import { liabilityTypeOptions } from '@/constants/options'
import type { LiabilityRow } from '@/types/form'

const props = defineProps<{
  liability: LiabilityRow
  index: number
}>()

const emit = defineEmits<{
  remove: []
}>()

const base = computed(() => `financial.liabilities[${props.index}]`)
</script>

<template>
  <ItemBlock :title="`หนี้สินรายการที่ ${index + 1}`" @remove="emit('remove')">
    <div class="grid-form-2">
      <FieldSelect
        :path="`${base}.liability_type`"
        label="หนี้สิน"
        v-model="liability.liability_type"
        :options="liabilityTypeOptions"
      />
      <FieldNumber
        :path="`${base}.credit_limit`"
        label="วงเงิน"
        :model-value="liability.credit_limit"
        format="amount"
        suffix="บาท"
        @update:model-value="(v) => (liability.credit_limit = v ?? 0)"
      />
    </div>
  </ItemBlock>
</template>
