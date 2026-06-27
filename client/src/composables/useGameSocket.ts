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
const roundStartedHandlers = new Set<(event: RoundStartedEvent) => void>();
const turnStartedHandlers = new Set<(event: TurnStartedEvent) => void>();
const strokeHandlers = new Set<(stroke: Stroke) => void>();
const roundAnswerHandlers = new Set<(event: RoundAnswerEvent) => void>();
const gameEndHandlers = new Set<(event: GameEndEvent) => void>();
const clientDisconnectedHandlers = new Set<(event: ClientDisconnectedEvent) => void>();

export const useGameSocket = () => {
  const gameStore = useGameStore();
  const socket = getSocket();

  const handleRoundStarted = (event: RoundStartedEvent) => {
    gameStore.startRound(event);
    roundStartedHandlers.forEach((handler) => {
      handler(event);
    });
  };
  const handleTurnStarted = (event: TurnStartedEvent) => {
    gameStore.startTurn(event);
    turnStartedHandlers.forEach((handler) => {
      handler(event);
    });
  };
  const handleStroke = (stroke: Stroke) => {
    gameStore.addStroke(stroke);
    strokeHandlers.forEach((handler) => {
      handler(stroke);
    });
  };
  const handleRoundAnswer = (event: RoundAnswerEvent) => {
    gameStore.showRoundAnswer(event);
    roundAnswerHandlers.forEach((handler) => {
      handler(event);
    });
  };
  const handleGameEnd = (event: GameEndEvent) => {
    gameStore.endGame(event);
    gameEndHandlers.forEach((handler) => {
      handler(event);
    });
  };
  const handleClientDisconnected = (event: ClientDisconnectedEvent) => {
    gameStore.abortByDisconnect(event);
    clientDisconnectedHandlers.forEach((handler) => {
      handler(event);
    });
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
      return false;
    }

    socket.emit('draw:stroke', stroke);
    return true;
  };

  const submitAnswer = (answer: string) => {
    if (!gameStore.canSubmitAnswer) {
      return false;
    }

    socket.emit('answer:submit', {
      answer,
    });
    return true;
  };

  const endRound = () => {
    if (gameStore.phase !== 'roundResult') {
      return false;
    }

    socket.emit('round:end');
    return true;
  };

  return {
    cleanup,
    drawStroke,
    endRound,
    offClientDisconnected: (handler: (event: ClientDisconnectedEvent) => void) =>
      clientDisconnectedHandlers.delete(handler),
    offGameEnd: (handler: (event: GameEndEvent) => void) => gameEndHandlers.delete(handler),
    offRoundAnswer: (handler: (event: RoundAnswerEvent) => void) =>
      roundAnswerHandlers.delete(handler),
    offRoundStarted: (handler: (event: RoundStartedEvent) => void) =>
      roundStartedHandlers.delete(handler),
    offStroke: (handler: (stroke: Stroke) => void) => strokeHandlers.delete(handler),
    offTurnStarted: (handler: (event: TurnStartedEvent) => void) =>
      turnStartedHandlers.delete(handler),
    onClientDisconnected: (handler: (event: ClientDisconnectedEvent) => void) =>
      clientDisconnectedHandlers.add(handler),
    onGameEnd: (handler: (event: GameEndEvent) => void) => gameEndHandlers.add(handler),
    onRoundAnswer: (handler: (event: RoundAnswerEvent) => void) => roundAnswerHandlers.add(handler),
    onRoundStarted: (handler: (event: RoundStartedEvent) => void) =>
      roundStartedHandlers.add(handler),
    onStroke: (handler: (stroke: Stroke) => void) => strokeHandlers.add(handler),
    onTurnStarted: (handler: (event: TurnStartedEvent) => void) => turnStartedHandlers.add(handler),
    register,
    submitAnswer,
  };
};
