import type { Stroke } from '@/types/api';

export type RoundResultViewData = {
  count: number;
  actualAnswer: string;
  guesserAnswer: string;
  guesserId: string;
  strokes: Stroke[];
};
