<script setup lang="ts">
import { computed } from 'vue'
import FieldNumber from '@/components/fields/FieldNumber.vue'
import FieldSelect from '@/components/fields/FieldSelect.vue'
import FieldText from '@/components/fields/FieldText.vue'
import ItemBlock from '@/components/common/ItemBlock.vue'
import { relationOptions } from '@/constants/options'
import type { MemberRow } from '@/types/form'

const props = defineProps<{
  member: MemberRow
  index: number
}>()

const emit = defineEmits<{
  remove: []
}>()

const base = computed(() => `family.members[${props.index}]`)
</script>

<template>
  <ItemBlock :title="`สมาชิกคนที่ ${index + 1}`" @remove="emit('remove')">
    <div class="grid-form-2">
      <FieldSelect
        :path="`${base}.relation`"
        required
        label="ความสัมพันธ์"
        v-model="member.relation"
        :options="relationOptions"
      />
      <FieldText
        :path="`${base}.full_name`"
        required
        label="ชื่อ-นามสกุล"
        v-model="member.full_name"
        placeholder="สมหญิง ใจดี"
      />
      <FieldText
        :path="`${base}.national_id`"
        label="เลขบัตร"
        optional
        v-model="member.national_id"
        format="national_id"
        placeholder="1-2345-67890-12-1"
      />
      <FieldNumber
        :path="`${base}.annual_income`"
        label="รายได้ต่อปี"
        optional
        v-model="member.annual_income"
        format="amount"
        placeholder="0"
        suffix="บาท"
      />
    </div>
  </ItemBlock>
</template>
