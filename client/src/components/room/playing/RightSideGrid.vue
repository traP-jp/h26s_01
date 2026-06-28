<script setup lang="ts">
import BaseButton from '@/components/common/BaseButton.vue';
import type { GamePhase } from '@/stores/game';

import LifeGage from './LifeGage.vue';

const props = defineProps<{
  phase: GamePhase;
  answer: string;
  canClearStroke: boolean;
  canSubmitStroke: boolean;
  canSubmitAnswer: boolean;
  canEndRound: boolean;
  isGuesser: boolean;
}>();

const emit = defineEmits<{
  'update:answer': [answer: string];
  clearStroke: [];
  submitStroke: [];
  submitAnswer: [];
  endRound: [];
}>();
</script>

<template>
  <div class="flex flex-col w-xs p-8 justify-between">
    <LifeGage />
    <div v-if="props.phase === 'playing' && !props.isGuesser" class="flex flex-col gap-3">
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
    <div v-else-if="props.phase === 'playing' && props.isGuesser" class="flex flex-col gap-3">
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
    <div v-else-if="props.phase === 'roundResult'" class="flex flex-col gap-3">
      <BaseButton variant="primary" :disabled="!props.canEndRound" @btn-click="emit('endRound')"
        >次の問題へ</BaseButton
      >
    </div>
  </div>
</template>
