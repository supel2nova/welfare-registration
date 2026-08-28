<script setup lang="ts">
defineProps<{
  title: string
}>()
</script>

<template>
  <section class="card-section">
    <div class="mb-4 flex flex-col items-start justify-between gap-3 sm:flex-row sm:items-center">
      <h2 class="m-0 text-lg">{{ title }}</h2>
      <slot name="header" />
    </div>
    <slot />
  </section>
</template>
