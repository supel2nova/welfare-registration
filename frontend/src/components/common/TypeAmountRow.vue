<script setup lang="ts">
import FieldNumber from '@/components/fields/FieldNumber.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import ItemBlock from '@/components/common/ItemBlock.vue'
import type { SelectOption } from '@/types/form'

defineProps<{
  title: string
  typePath: string
  typeLabel: string
  typeOptions: SelectOption[]
  type: string
  amountPath: string
  amountLabel: string
  amount: number
  suffix: string
}>()

const emit = defineEmits<{
  'update:type': [value: string]
  'update:amount': [value: number]
  remove: []
}>()
</script>

<template>
  <ItemBlock :title="title" @remove="emit('remove')">
    <div class="grid-form-2">
      <FieldSelect
        :path="typePath"
        required
        :label="typeLabel"
        :model-value="type"
        :options="typeOptions"
        @update:model-value="(v) => emit('update:type', v)"
      />
      <FieldNumber
        :path="amountPath"
        required
        :label="amountLabel"
        :model-value="amount"
        format="amount"
        placeholder="0"
        :suffix="suffix"
        @update:model-value="(v) => emit('update:amount', v ?? 0)"
      />
    </div>
  </ItemBlock>
</template>
