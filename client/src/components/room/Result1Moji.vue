<script setup lang="ts">
import UserIcon from '@/components/common/UserIcon.vue';
import StrokeCanvasView from '@/components/room/playing/StrokeCanvasView.vue';
import type { RoundResultViewData } from '@/types/result';
import { toKanjiNumber } from '@/utils/to-kanji-number';

withDefaults(
  defineProps<{
    resultData: RoundResultViewData;
    fitViewport?: boolean;
  }>(),
  {
    fitViewport: false,
  },
);
</script>

<template>
  <section
    class="w-full bg-background"
    :class="{ 'flex min-h-0 flex-col overflow-y-auto': fitViewport }"
  >
    <p class="text-primary text-5xl">{{ toKanjiNumber(resultData.count) }}文字目</p>
    <div
      class="mt-8 grid items-start gap-10 lg:grid-cols-2"
      :class="{ 'min-h-0 flex-1': fitViewport }"
    >
      <div
        class="result-canvas-stage grid min-h-0 place-items-center"
        :class="{ 'result-canvas-stage--fitted h-full': fitViewport }"
      >
        <StrokeCanvasView
          class="result-canvas"
          :class="{ 'result-canvas--fitted': fitViewport }"
          :strokes="resultData.strokes"
        />
      </div>
      <div class="flex min-w-0 flex-col gap-8">
        <p class="text-primary text-7xl">
          推測{{ resultData.actualAnswer === resultData.guesserAnswer ? '成功' : '失敗' }}
        </p>
        <div class="grid gap-4 text-primary">
          <div class="grid grid-cols-[auto_minmax(0,1fr)] items-center gap-6">
            <p class="text-4xl">お題</p>
            <div class="aspect-square w-full max-w-32 bg-white flex items-center justify-center">
              <p class="text-7xl">{{ resultData.actualAnswer }}</p>
            </div>
          </div>
          <div class="grid grid-cols-[auto_minmax(0,1fr)] items-center gap-6">
            <p class="text-4xl">回答</p>
            <div class="aspect-square w-full max-w-32 bg-white flex items-center justify-center">
              <p class="text-7xl">{{ resultData.guesserAnswer }}</p>
            </div>
          </div>
          <div class="grid grid-cols-[auto_minmax(0,1fr)] items-center gap-6">
            <p class="text-4xl">回答者</p>
            <UserIcon :tra-qid="resultData.guesserId" size="small" />
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.result-canvas-stage--fitted {
  container-type: size;
}

.result-canvas--fitted {
  width: min(100cqw, 100cqh);
  height: min(100cqw, 100cqh);
}
</style>
