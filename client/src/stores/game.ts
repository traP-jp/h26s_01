import { defineStore } from 'pinia';

import { MAX_ROUNDS, MAX_STROKES } from '@/constants/game';
import type {
  ClientDisconnectedEvent,
  GameEndEvent,
  RoundAnswerEvent,
  RoundStartedEvent,
  Stroke,
  TurnStartedEvent,
} from '@/types/api';

import { useRoomStore } from './room';
import { useUserStore } from './user';

export type GamePhase = 'idle' | 'playing' | 'roundResult' | 'ended' | 'aborted';

type GameState = {
  phase: GamePhase;
  roundIndex: number | null;
  maxRounds: number;
  guesserId: string | null;
  currentDrawerId: string | null;
  turnIndex: number | null;
  kanji: string | null;
  strokes: Stroke[];
  roundStartedAt: number | null;
  roundAnswer: RoundAnswerEvent | null;
  finalResult: GameEndEvent | null;
  disconnectedUserId: string | null;
  error: string | null;
};

export const useGameStore = defineStore('game', {
  state: (): GameState => ({
    phase: 'idle',
    roundIndex: null,
    maxRounds: MAX_ROUNDS,
    guesserId: null,
    currentDrawerId: null,
    turnIndex: null,
    kanji: null,
    strokes: [],
    roundStartedAt: null,
    roundAnswer: null,
    finalResult: null,
    disconnectedUserId: null,
    error: null,
  }),
  getters: {
    currentUserId(): string | null {
      return useUserStore().currentUser?.id ?? null;
    },
    isGuesser(): boolean {
      return this.currentUserId !== null && this.currentUserId === this.guesserId;
    },
    isDrawer(): boolean {
      const currentUserId = this.currentUserId;
      const isRoomMember = useRoomStore().members.some((member) => member.id === currentUserId);

      return (
        currentUserId !== null &&
        this.guesserId !== null &&
        currentUserId !== this.guesserId &&
        isRoomMember
      );
    },
    isMyTurn(): boolean {
      return this.currentUserId !== null && this.currentUserId === this.currentDrawerId;
    },
    canDraw(): boolean {
      return this.phase === 'playing' && this.isMyTurn && this.strokes.length < MAX_STROKES;
    },
    canSubmitAnswer(): boolean {
      return this.phase === 'playing' && this.isGuesser;
    },
    strokeCount: (state) => state.strokes.length,
    remainingStrokes: (state) => Math.max(MAX_STROKES - state.strokes.length, 0),
    elapsedSeconds(state): number | null {
      if (state.phase !== 'playing' || state.roundStartedAt === null) {
        return null;
      }

      return Math.floor((Date.now() - state.roundStartedAt) / 1000);
    },
    roundLabel: (state) =>
      state.roundIndex === null ? '' : `${state.roundIndex} / ${state.maxRounds}`,
    turnLabel: (state) => (state.turnIndex === null ? '' : `${state.turnIndex} / ${MAX_STROKES}`),
  },
  actions: {
    setError(error: string | null) {
      this.error = error;
    },
    startRound(event: RoundStartedEvent) {
      this.phase = 'playing';
      this.roundIndex = event.roundIndex;
      this.guesserId = event.guesserId;
      this.currentDrawerId = null;
      this.turnIndex = null;
      this.kanji = event.kanji ?? null;
      this.strokes = [];
      this.roundStartedAt = Date.now();
      this.roundAnswer = null;
      this.finalResult = null;
      this.disconnectedUserId = null;
      this.error = null;
    },
    startTurn(event: TurnStartedEvent) {
      this.turnIndex = event.turnIndex;
      this.currentDrawerId = event.drawerId;
    },
    addStroke(stroke: Stroke) {
      this.strokes.push(stroke);
    },
    showRoundAnswer(event: RoundAnswerEvent) {
      this.phase = 'roundResult';
      this.roundAnswer = event;
      this.kanji = null;
      this.currentDrawerId = null;
      this.roundStartedAt = null;
    },
    endGame(event: GameEndEvent) {
      this.phase = 'ended';
      this.finalResult = event;
      this.kanji = null;
      this.currentDrawerId = null;
      this.roundStartedAt = null;
    },
    abortByDisconnect(event: ClientDisconnectedEvent) {
      this.phase = 'aborted';
      this.disconnectedUserId = event.userId;
      this.kanji = null;
      this.currentDrawerId = event.newDrawerId ?? null;
      this.roundStartedAt = null;
    },
    resetGame() {
      this.phase = 'idle';
      this.roundIndex = null;
      this.guesserId = null;
      this.currentDrawerId = null;
      this.turnIndex = null;
      this.kanji = null;
      this.strokes = [];
      this.roundStartedAt = null;
      this.roundAnswer = null;
      this.finalResult = null;
      this.disconnectedUserId = null;
      this.error = null;
    },
  },
});
