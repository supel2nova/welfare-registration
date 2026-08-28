<script setup lang="ts">
import { computed } from 'vue'
import FieldNumber from '@/components/fields/FieldNumber.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import ItemBlock from '@/components/common/ItemBlock.vue'
import { assetTypeOptions, assetUnits, unitLabels } from '@/constants/options'
import type { AssetRow } from '@/types/form'

const props = defineProps<{
  asset: AssetRow
  index: number
}>()

const emit = defineEmits<{
  remove: []
}>()

const base = computed(() => `financial.assets[${props.index}]`)
const unitLabel = computed(() => unitLabels[props.asset.unit] ?? props.asset.unit)

function onTypeChange(value: string) {
  props.asset.asset_type = value
  props.asset.unit = assetUnits[value] ?? 'THB'
}
</script>

<template>
  <ItemBlock :title="`ทรัพย์สินรายการที่ ${index + 1}`" @remove="emit('remove')">
    <div class="grid-form-2">
      <FieldSelect
        :path="`${base}.asset_type`"
        label="ทรัพย์สิน"
        :model-value="asset.asset_type"
        :options="assetTypeOptions"
        @update:model-value="onTypeChange"
      />
      <FieldNumber
        :path="`${base}.amount`"
        label="จำนวน"
        :model-value="typeof asset.amount === 'number' ? asset.amount : Number(asset.amount)"
        format="amount"
        placeholder="0"
        :suffix="unitLabel"
        @update:model-value="(v) => (asset.amount = v ?? 0)"
      />
    </div>
  </ItemBlock>
</template>
