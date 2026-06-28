<script setup lang="ts">
import UserIcon from '@/components/common/UserIcon.vue';

import LeftSideSection from './LeftSideSection.vue';

const props = defineProps<{
  roundLabel: string;
  turnLabel: string;
  elapsedSeconds: number | null;
  guesserId: string | null;
  currentDrawerId: string | null;
  remainingStrokes: number;
}>();

const formatElapsed = (seconds: number | null) => {
  if (seconds === null) {
    return '-';
  }

  const minutes = Math.floor(seconds / 60);
  const restSeconds = seconds % 60;
  return `${minutes}分${restSeconds.toString().padStart(2, '0')}秒`;
};
</script>

<template>
  <aside class="h-full min-w-0 overflow-hidden border-r-4 border-primary">
    <div class="p-8 text-4xl text-center bg-primary text-background">
      <p>遊戯中</p>
    </div>
    <LeftSideSection>
      <template #title>今の状況</template>
      <template #default>
        <div class="space-y-2 text-2xl">
          <p>文字 {{ props.roundLabel || '-' }}</p>
          <p>画数 {{ props.turnLabel || '-' }}</p>
          <p>残り {{ props.remainingStrokes }}画</p>
          <p>{{ formatElapsed(props.elapsedSeconds) }}が経過</p>
        </div>
      </template>
    </LeftSideSection>
    <LeftSideSection>
      <template #title>読み手</template>
      <template #default>
        <div v-if="props.guesserId" class="flex min-w-0 items-center gap-3">
          <UserIcon :tra-qid="props.guesserId" size="small" />
          <p class="truncate text-lg" :title="props.guesserId">{{ props.guesserId }}</p>
        </div>
        <p v-else class="text-lg">-</p>
      </template>
    </LeftSideSection>
    <LeftSideSection>
      <template #title>書き手</template>
      <template #default>
        <div v-if="props.currentDrawerId" class="flex flex-col gap-3">
          <p class="underline text-lg">只今作業中</p>
          <div class="flex min-w-0 items-center gap-3">
            <UserIcon :tra-qid="props.currentDrawerId" size="small" />
            <p class="truncate text-lg" :title="props.currentDrawerId">
              {{ props.currentDrawerId }}
            </p>
          </div>
        </div>
        <p v-else class="text-lg">-</p>
      </template>
    </LeftSideSection>
  </aside>
</template>
