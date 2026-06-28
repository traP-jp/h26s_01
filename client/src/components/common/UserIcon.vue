<script setup lang="ts">
defineProps<{ traQid: string; size: 'small' | 'medium' | 'large' }>();
</script>

<template>
  <img
    :src="`https://q.trap.jp/api/v3/public/icon/${traQid}`"
    :alt="traQid + 'のアイコン'"
    :class="{
      'w-12 h-12': size == 'small',
      'w-24 h-24': size == 'medium',
      'w-36 h-36': size == 'large',
    }"
  />
</template>
