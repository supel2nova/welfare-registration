<script setup lang="ts">
const paths = {
  plus: 'M10 4v12M4 10h12',
  trash: 'M4 6h12M8.5 6V4.5h3V6M6 6l.6 9.5h6.8L14 6M8.5 9v4M11.5 9v4',
} as const

defineProps<{
  name: keyof typeof paths
}>()
</script>

<template>
  <svg
    class="app-icon"
    viewBox="0 0 20 20"
    fill="none"
    stroke="currentColor"
    stroke-width="1.6"
    stroke-linecap="round"
    stroke-linejoin="round"
    aria-hidden="true"
    focusable="false"
  >
    <path :d="paths[name]" />
  </svg>
</template>
