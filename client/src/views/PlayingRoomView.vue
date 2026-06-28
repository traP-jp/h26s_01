<script setup lang="ts">
import { computed, onBeforeUnmount, ref } from 'vue';

import BaseButton from '@/components/common/BaseButton.vue';
import LastResult from '@/components/room/last-result/LastResult.vue';
import BottomGrid from '@/components/room/playing/BottomGrid.vue';
import CenterGrid from '@/components/room/playing/CenterGrid.vue';
import LeftSideGrid from '@/components/room/playing/LeftSideGrid.vue';
import RightSideGrid from '@/components/room/playing/RightSideGrid.vue';
import StrokeCanvas from '@/components/room/playing/StrokeCanvas.vue';
import ResultMojigoto from '@/components/room/ResultMojigoto.vue';
import { useGameSocket } from '@/composables/useGameSocket';
import { useReturnHomeByDisconnect } from '@/composables/useReturnHomeByDisconnect';
import { useGameStore } from '@/stores/game';
import { useRoomStore } from '@/stores/room';
import type { DrawStrokeEvent } from '@/types/api';
import type { RoundResultViewData } from '@/types/result';

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

const canAbort = computed(() => ['playing', 'roundResult'].includes(gameStore.phase));

const roundResultData = computed<RoundResultViewData | null>(() => {
  if (
    gameStore.roundIndex === null ||
    gameStore.roundAnswer === null ||
    gameStore.guesserId === null
  ) {
    return null;
  }

  return {
    count: gameStore.roundIndex,
    actualAnswer: gameStore.roundAnswer.actualAnswer,
    guesserAnswer: gameStore.roundAnswer.guesserAnswer,
    guesserId: gameStore.guesserId,
    strokes: gameStore.strokes,
  };
});

const finalResultData = computed<RoundResultViewData[]>(() =>
  gameStore.finalResult === null
    ? []
    : gameStore.finalResult.rounds.map((round, index) => ({
        count: index + 1,
        actualAnswer: round.actualAnswer,
        guesserAnswer: round.guesserAnswer,
        guesserId: round.guesserId,
        strokes: round.strokes,
      })),
);

const bottomMessage = computed(() => {
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

const returnHome = async () => {
  await returnHomeByDisconnect();
};

const handleAbort = async () => {
  if (!window.confirm('ゲームを中止してトップへ戻りますか？')) {
    return;
  }

  await returnHome();
};

onBeforeUnmount(() => {
  window.clearInterval(intervalId);
});
</script>

<template>
  <div class="min-h-dvh bg-background text-primary">
    <div v-if="gameStore.phase === 'playing'" class="grid h-dvh grid-rows-[1fr_auto]">
      <div
        class="grid min-h-0 grid-cols-[minmax(14rem,18rem)_minmax(0,1fr)_minmax(16rem,22rem)] overflow-hidden"
      >
        <LeftSideGrid
          :current-drawer-id="gameStore.currentDrawerId"
          :elapsed-seconds="elapsedSeconds"
          :guesser-id="gameStore.guesserId"
          :remaining-strokes="gameStore.remainingStrokes"
          :round-label="gameStore.roundLabel"
          :turn-label="gameStore.turnLabel"
        />
        <CenterGrid :is-guesser="gameStore.isGuesser" :kanji="gameStore.kanji">
          <div class="mx-auto w-full max-w-3xl">
            <StrokeCanvas
              v-model:draft-stroke="draftStroke"
              :can-draw="gameStore.canDraw"
              :disabled="gameStore.phase !== 'playing'"
              :strokes="gameStore.strokes"
            />
          </div>
          <p v-if="roomStore.currentRoom === null" class="text-xl">
            部屋情報を取得できませんでした。
          </p>
        </CenterGrid>
        <RightSideGrid
          v-model:answer="answer"
          :can-abort="canAbort"
          :can-clear-stroke="draftStroke !== null"
          :can-end-round="false"
          :can-submit-answer="canSubmitAnswer"
          :can-submit-stroke="canSubmitDraftStroke"
          :is-guesser="gameStore.isGuesser"
          :phase="gameStore.phase"
          @abort="handleAbort"
          @clear-stroke="clearDraftStroke"
          @end-round="endRound"
          @submit-answer="submitAnswer"
          @submit-stroke="submitDraftStroke"
        />
      </div>
      <BottomGrid :message="bottomMessage" />
    </div>

    <ResultMojigoto
      v-else-if="gameStore.phase === 'roundResult' && roundResultData !== null"
      :can-abort="canAbort"
      :result-data="roundResultData"
      @abort="handleAbort"
      @next="endRound"
    />

    <div
      v-else-if="gameStore.phase === 'roundResult'"
      class="min-h-dvh flex flex-col items-center justify-center gap-8 text-3xl"
    >
      <p>結果データを取得できませんでした。</p>
      <BaseButton variant="primary" @btn-click="endRound">次の問題へ</BaseButton>
    </div>

    <LastResult
      v-else-if="gameStore.phase === 'ended'"
      :results="finalResultData"
      @return-home="returnHome"
    />

    <div
      v-else-if="gameStore.phase === 'aborted'"
      class="min-h-dvh flex flex-col items-center justify-center gap-8 text-3xl"
    >
      <p>切断が発生したためゲームは中止されました。</p>
      <p>切断したユーザー: {{ gameStore.disconnectedUserId ?? '-' }}</p>
      <BaseButton variant="primary" @btn-click="returnHome">トップへ戻る</BaseButton>
    </div>

    <div v-else class="min-h-dvh flex items-center justify-center text-3xl">
      ゲーム情報を待機しています。
    </div>
  </div>
</template>
