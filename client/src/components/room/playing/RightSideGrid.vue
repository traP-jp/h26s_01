<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue';

import LifeGage from './LifeGage.vue';
</script>

<template>
  <div class="flex flex-col w-xs p-8 justify-between">
    <LifeGage />
    <div class="flex flex-col gap-3">
      <BaseButton variant="secondary">書き直し</BaseButton>
      <BaseButton variant="primary">確定</BaseButton>
    </div>
  </div>
</template>
