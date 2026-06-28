<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue';

import BaseButton from '@/components/common/BaseButton.vue';
import BottomGrid from '@/components/room/playing/BottomGrid.vue';
import CenterGrid from '@/components/room/playing/CenterGrid.vue';
import LeftSideGrid from '@/components/room/playing/LeftSideGrid.vue';
import RightSideGrid from '@/components/room/playing/RightSideGrid.vue';
import StrokeCanvas from '@/components/room/playing/StrokeCanvas.vue';
import { useGameSocket } from '@/composables/useGameSocket';
import { useReturnHomeByDisconnect } from '@/composables/useReturnHomeByDisconnect';
import { useGameStore } from '@/stores/game';
import { useRoomStore } from '@/stores/room';
import type { DrawStrokeEvent } from '@/types/api';

const MIN_STROKE_LENGTH = 0.01;

const gameStore = useGameStore();
const roomStore = useRoomStore();
const gameSocket = useGameSocket();
const { returnHomeByDisconnect } = useReturnHomeByDisconnect();
const draftStroke = ref<DrawStrokeEvent | null>(null);
const answer = ref('');
const now = ref(Date.now());

const intervalId = window.setInterval(() => {
  now.value = Date.now();
}, 1000);

const elapsedSeconds = computed(() => {
  if (gameStore.phase !== 'playing' || gameStore.roundStartedAt === null) {
    return null;
  }

  return Math.floor((now.value - gameStore.roundStartedAt) / 1000);
});

const draftStrokeLength = computed(() => {
  if (draftStroke.value === null) {
    return 0;
  }

  return Math.hypot(
    draftStroke.value.x2 - draftStroke.value.x1,
    draftStroke.value.y2 - draftStroke.value.y1,
  );
});

const canSubmitDraftStroke = computed(
  () =>
    gameStore.canDraw && draftStroke.value !== null && draftStrokeLength.value >= MIN_STROKE_LENGTH,
);

const canSubmitAnswer = computed(() => gameStore.canSubmitAnswer && answer.value.trim().length > 0);

const bottomMessage = computed(() => {
  if (gameStore.phase === 'roundResult') {
    return '回答結果を確認してください。';
  }

  if (gameStore.phase === 'ended') {
    return 'ゲームが終了しました。';
  }

  if (gameStore.phase === 'aborted') {
    return '切断が発生したためゲームは中止されました。';
  }

  if (gameStore.isMyTurn) {
    return 'あなたの順番です。直線を一つ書き加えたら、確定を押してください。';
  }

  if (gameStore.isGuesser) {
    return '読み手です。推測したお題を入力してください。';
  }

  return '他の人が書いています。';
});

const clearDraftStroke = () => {
  draftStroke.value = null;
};

const submitDraftStroke = () => {
  if (!canSubmitDraftStroke.value || draftStroke.value === null) {
    return;
  }

  const didSend = gameSocket.drawStroke(draftStroke.value);
  if (didSend) {
    draftStroke.value = null;
  }
};

const submitAnswer = () => {
  const trimmedAnswer = answer.value.trim();
  if (trimmedAnswer.length === 0) {
    return;
  }

  if (gameSocket.submitAnswer(trimmedAnswer)) {
    answer.value = '';
  }
};

const endRound = () => {
  gameSocket.endRound();
};

const handleAbort = async () => {
  if (!window.confirm('ゲームを中止してトップへ戻りますか？')) {
    return;
  }

  await returnHomeByDisconnect();
};

onBeforeUnmount(() => {
  window.clearInterval(intervalId);
});
</script>

<template>
  <div class="h-dvh flex flex-col bg-background text-primary">
    <div class="fixed top-4 right-4 z-10">
      <BaseButton variant="secondary" @btn-click="handleAbort">中止して戻る</BaseButton>
    </div>
    <div class="flex-1 flex justify-between overflow-hidden">
      <LeftSideGrid
        :current-drawer-id="gameStore.currentDrawerId"
        :elapsed-seconds="elapsedSeconds"
        :guesser-id="gameStore.guesserId"
        :remaining-strokes="gameStore.remainingStrokes"
        :round-label="gameStore.roundLabel"
        :turn-label="gameStore.turnLabel"
      />
      <CenterGrid :is-guesser="gameStore.isGuesser" :kanji="gameStore.kanji">
        <div class="w-[min(60vh,42rem)] max-w-full">
          <StrokeCanvas
            v-model:draft-stroke="draftStroke"
            :can-draw="gameStore.canDraw"
            :disabled="gameStore.phase !== 'playing'"
            :strokes="gameStore.strokes"
          />
        </div>
        <div v-if="gameStore.phase === 'roundResult' && gameStore.roundAnswer" class="text-3xl">
          <p>回答: {{ gameStore.roundAnswer.guesserAnswer }}</p>
          <p>正解: {{ gameStore.roundAnswer.actualAnswer }}</p>
        </div>
        <div v-else-if="gameStore.phase === 'ended' && gameStore.finalResult" class="text-3xl">
          <p>{{ gameStore.finalResult.cleared ? '成功' : '失敗' }}</p>
          <p>残機: {{ gameStore.finalResult.remainingLives }}</p>
        </div>
        <div v-else-if="gameStore.phase === 'aborted'" class="text-3xl">
          <p>切断: {{ gameStore.disconnectedUserId ?? '-' }}</p>
        </div>
        <p v-if="roomStore.currentRoom === null" class="text-xl">
          部屋情報を取得できませんでした。
        </p>
      </CenterGrid>
      <RightSideGrid
        v-model:answer="answer"
        :can-clear-stroke="draftStroke !== null"
        :can-end-round="gameStore.phase === 'roundResult'"
        :can-submit-answer="canSubmitAnswer"
        :can-submit-stroke="canSubmitDraftStroke"
        :is-guesser="gameStore.isGuesser"
        :phase="gameStore.phase"
        @clear-stroke="clearDraftStroke"
        @end-round="endRound"
        @submit-answer="submitAnswer"
        @submit-stroke="submitDraftStroke"
      />
    </div>
    <BottomGrid :message="bottomMessage" />
  </div>
</template>
