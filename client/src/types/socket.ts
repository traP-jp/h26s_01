import type { Socket } from 'socket.io-client';

import type {
  AnswerSubmitEvent,
  ClientDisconnectedEvent,
  DrawStrokeEvent,
  GameEndEvent,
  RoomJoinEvent,
  RoomListUpdatedEvent,
  RoomUpdatedEvent,
  RoundAnswerEvent,
  RoundStartedEvent,
  Stroke,
  TurnStartedEvent,
} from './api';

export type ClientToServerEvents = {
  'room:join': (payload: RoomJoinEvent) => void;
  // This event is agreed in the game flow but is not represented by a schema yet.
  'game:ready': () => void;
  'draw:stroke': (payload: DrawStrokeEvent) => void;
  'answer:submit': (payload: AnswerSubmitEvent) => void;
  'round:end': () => void;
};

export type ServerToClientEvents = {
  'room_list:updated': (payload: RoomListUpdatedEvent) => void;
  'room:updated': (payload: RoomUpdatedEvent) => void;
  'round:started': (payload: RoundStartedEvent) => void;
  'turn:started': (payload: TurnStartedEvent) => void;
  'draw:stroke': (payload: Stroke) => void;
  'client:disconnected': (payload: ClientDisconnectedEvent) => void;
  'round:answer': (payload: RoundAnswerEvent) => void;
  'game:end': (payload: GameEndEvent) => void;
};

export type AppSocket = Socket<ServerToClientEvents, ClientToServerEvents>;
