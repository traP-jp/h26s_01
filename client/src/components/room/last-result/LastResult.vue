<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue';
import TheHeader from '@/components/common/TheHeader.vue';

import Result1Moji from '../Result1Moji.vue';
import { type ResultProps } from '../Result1Moji.vue';

export interface LastResultProps {
  successGame: number;
  loseGame: number;
}
defineProps<{
  lastData: LastResultProps;
}>();

const test: ResultProps = {
  count: 3,
  actualAnswer: '詰',
  guesserAnswer: '読',
  traqId: 'yadorigi',
};
</script>

<template>
  <div class="min-h-screen bg-background">
    <TheHeader />
    <div class="mt-16 mx-32 flex justify-end">
      <BaseButton variant="primary">退出</BaseButton>
    </div>
    <div class="flex flex-row justify-center gap-32 mt-12">
      <div class="flex flex-col items-end">
        <p class="text-primary text-6xl mt-8">結果</p>
        <div class="flex flex-col items-end">
          <p class="text-primary text-6xl mt-48">推測成功</p>
          <p class="text-primary text-6xl mt-32">推測失敗</p>
        </div>
      </div>
      <div class="flex flex-col">
        <p class="text-primary text-9xl">成功</p>
        <p class="text-primary text-6xl mt-40">
          {{
            lastData.successGame.toLocaleString('ja-JP-u-nu-hanidec', { useGrouping: false })
          }}文字 / 九文字
        </p>
        <p class="text-primary text-6xl mt-32">
          {{ lastData.loseGame.toLocaleString('ja-JP-u-nu-hanidec', { useGrouping: false }) }}文字 /
          九文字
        </p>
      </div>
    </div>
    <div class="flex flex-col gap-48 mt-48 justify-start ml-64">
      <div v-for="count in 9" :key="count">
        <Result1Moji :result-data="test" />
      </div>
    </div>
  </div>
</template>
