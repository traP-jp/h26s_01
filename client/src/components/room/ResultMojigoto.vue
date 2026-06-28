<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue';
import type { RoundResultViewData } from '@/types/result';

import Result1Moji from './Result1Moji.vue';

const props = defineProps<{
  resultData: RoundResultViewData;
  canAbort: boolean;
}>();

const emit = defineEmits<{
  abort: [];
  next: [];
}>();
</script>

<template>
  <div class="min-h-dvh bg-background">
    <div class="mx-auto flex w-full max-w-screen-xl flex-col gap-10 px-8 py-10">
      <div class="flex flex-wrap items-center justify-between gap-4">
        <p class="text-primary text-6xl">途中結果</p>
        <div class="flex flex-wrap gap-3">
          <BaseButton v-if="props.canAbort" variant="secondary" @btn-click="emit('abort')">
            途中で退出
          </BaseButton>
          <BaseButton variant="primary" @btn-click="emit('next')">次の問題へ</BaseButton>
        </div>
      </div>
      <div class="h-1 w-full bg-primary"></div>
      <Result1Moji :result-data="props.resultData" />
    </div>
  </div>
</template>
