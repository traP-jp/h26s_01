<script setup lang="ts">
import { computed } from 'vue';

import BaseButton from '@/components/common/BaseButton.vue';
import TheHeader from '@/components/common/TheHeader.vue';
import Result1Moji from '@/components/room/Result1Moji.vue';
import { MAX_ROUNDS } from '@/constants/game';
import type { RoundResultViewData } from '@/types/result';
import { toKanjiNumber } from '@/utils/to-kanji-number';

const props = defineProps<{
  cleared: boolean;
  results: RoundResultViewData[];
}>();

const emit = defineEmits<{
  returnHome: [];
}>();

const successCount = computed(
  () => props.results.filter((result) => result.actualAnswer === result.guesserAnswer).length,
);
const loseCount = computed(() => props.results.length - successCount.value);
</script>

<template>
  <div class="min-h-dvh bg-background">
    <TheHeader />
    <main class="mx-auto flex w-full max-w-screen-xl flex-col gap-16 px-8 py-10">
      <div class="flex justify-end">
        <BaseButton variant="primary" @btn-click="emit('returnHome')">トップへ戻る</BaseButton>
      </div>

      <section class="grid items-end gap-10 text-primary lg:grid-cols-[auto_minmax(0,1fr)]">
        <div class="grid gap-8 lg:justify-items-end">
          <p class="text-6xl">結果</p>
          <div class="grid gap-6 text-5xl">
            <p>推測成功</p>
            <p>推測失敗</p>
          </div>
        </div>
        <div class="grid gap-6">
          <p class="text-8xl">{{ props.cleared ? '成功' : '不成功' }}</p>
          <p class="text-5xl">
            {{ toKanjiNumber(successCount) }}文字 / {{ toKanjiNumber(MAX_ROUNDS) }}文字
          </p>
          <p class="text-5xl">
            {{ toKanjiNumber(loseCount) }}文字 / {{ toKanjiNumber(MAX_ROUNDS) }}文字
          </p>
        </div>
      </section>

      <div v-if="props.results.length === 0" class="text-center text-primary text-4xl">
        結果がありません
      </div>
      <div v-else class="flex flex-col gap-16">
        <Result1Moji v-for="result in props.results" :key="result.count" :result-data="result" />
      </div>
    </main>
  </div>
</template>
