<script setup lang="ts">
defineProps<{
  message: string;
}>();
</script>

<template>
  <div class="bg-primary text-background p-6">
    <p class="text-2xl text-center">{{ message }}</p>
  </div>
</template>
