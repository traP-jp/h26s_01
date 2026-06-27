import { getSocket } from '@/api/socket';
import { useGameStore } from '@/stores/game';
import type {
  ClientDisconnectedEvent,
  DrawStrokeEvent,
  GameEndEvent,
  RoundAnswerEvent,
  RoundStartedEvent,
  Stroke,
  TurnStartedEvent,
} from '@/types/api';

let hasRegisteredListeners = false;

export const useGameSocket = () => {
  const gameStore = useGameStore();
  const socket = getSocket();

  const handleRoundStarted = (event: RoundStartedEvent) => {
    gameStore.startRound(event);
  };
  const handleTurnStarted = (event: TurnStartedEvent) => {
    gameStore.startTurn(event);
  };
  const handleStroke = (stroke: Stroke) => {
    gameStore.addStroke(stroke);
  };
  const handleRoundAnswer = (event: RoundAnswerEvent) => {
    gameStore.showRoundAnswer(event);
  };
  const handleGameEnd = (event: GameEndEvent) => {
    gameStore.endGame(event);
  };
  const handleClientDisconnected = (event: ClientDisconnectedEvent) => {
    gameStore.abortByDisconnect(event);
  };

  const register = () => {
    if (hasRegisteredListeners) {
      return;
    }

    socket.on('round:started', handleRoundStarted);
    socket.on('turn:started', handleTurnStarted);
    socket.on('draw:stroke', handleStroke);
    socket.on('round:answer', handleRoundAnswer);
    socket.on('game:end', handleGameEnd);
    socket.on('client:disconnected', handleClientDisconnected);
    hasRegisteredListeners = true;
  };

  const cleanup = () => {
    socket.off('round:started', handleRoundStarted);
    socket.off('turn:started', handleTurnStarted);
    socket.off('draw:stroke', handleStroke);
    socket.off('round:answer', handleRoundAnswer);
    socket.off('game:end', handleGameEnd);
    socket.off('client:disconnected', handleClientDisconnected);
    hasRegisteredListeners = false;
  };

  const drawStroke = (stroke: DrawStrokeEvent) => {
    if (!gameStore.canDraw) {
      return;
    }

    socket.emit('draw:stroke', stroke);
  };

  const submitAnswer = (answer: string) => {
    if (!gameStore.canSubmitAnswer) {
      return;
    }

    socket.emit('answer:submit', {
      answer,
    });
  };

  const endRound = () => {
    socket.emit('round:end');
  };

  return {
    cleanup,
    drawStroke,
    endRound,
    register,
    submitAnswer,
  };
};
