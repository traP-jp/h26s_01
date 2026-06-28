<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue';
import type { GamePhase } from '@/stores/game';

import LifeGage from './LifeGage.vue';

const props = defineProps<{
  phase: GamePhase;
  answer: string;
  canAbort: boolean;
  canClearStroke: boolean;
  canSubmitStroke: boolean;
  canSubmitAnswer: boolean;
  canEndRound: boolean;
  isGuesser: boolean;
}>();

const emit = defineEmits<{
  'update:answer': [answer: string];
  abort: [];
  clearStroke: [];
  submitStroke: [];
  submitAnswer: [];
  endRound: [];
}>();
</script>

<template>
  <aside class="flex h-full min-w-0 flex-col gap-8 border-l-4 border-primary p-8">
    <div class="flex flex-col gap-6">
      <BaseButton v-if="props.canAbort" variant="secondary" @btn-click="emit('abort')">
        途中で退出
      </BaseButton>
      <LifeGage />
    </div>

    <div v-if="props.phase === 'playing' && !props.isGuesser" class="mt-auto flex flex-col gap-3">
      <BaseButton
        variant="secondary"
        :disabled="!props.canClearStroke"
        @btn-click="emit('clearStroke')"
        >書き直し</BaseButton
      >
      <BaseButton
        variant="primary"
        :disabled="!props.canSubmitStroke"
        @btn-click="emit('submitStroke')"
        >確定</BaseButton
      >
    </div>
    <div
      v-else-if="props.phase === 'playing' && props.isGuesser"
      class="mt-auto flex flex-col gap-3"
    >
      <input
        :value="props.answer"
        class="outline-none text-2xl p-3 border-4 border-primary bg-background text-primary"
        type="text"
        placeholder="回答"
        @input="emit('update:answer', ($event.target as HTMLInputElement).value)"
      />
      <BaseButton
        variant="primary"
        :disabled="!props.canSubmitAnswer"
        @btn-click="emit('submitAnswer')"
        >回答</BaseButton
      >
    </div>
    <div v-else-if="props.phase === 'roundResult'" class="mt-auto flex flex-col gap-3">
      <BaseButton variant="primary" :disabled="!props.canEndRound" @btn-click="emit('endRound')"
        >次の問題へ</BaseButton
      >
    </div>
  </aside>
</template>
