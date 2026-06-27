import type { components } from '@/api/schema';

export type User = components['schemas']['User'];
export type RoomMember = components['schemas']['RoomMember'];
export type Room = components['schemas']['Room'];
export type CreateRoomRequest = components['schemas']['CreateRoomRequest'];

export type RoomJoinEvent = components['schemas']['RoomJoinEvent'];
export type DrawStrokeEvent = components['schemas']['DrawStrokeEvent'];
export type AnswerSubmitEvent = components['schemas']['AnswerSubmitEvent'];

export type RoomListUpdatedEvent = components['schemas']['RoomListUpdatedEvent'];
export type RoomCreatedEvent = components['schemas']['RoomCreatedEvent'];
export type RoomUpdatedEvent = components['schemas']['RoomUpdatedEvent'];
export type RoomDeletedEvent = components['schemas']['RoomDeletedEvent'];

export type RoundStartedEvent = components['schemas']['RoundStartedEvent'];
export type TurnStartedEvent = components['schemas']['TurnStartedEvent'];
export type ClientDisconnectedEvent = components['schemas']['ClientDisconnectedEvent'];
export type RoundAnswerEvent = components['schemas']['RoundAnswerEvent'];
export type GameEndEvent = components['schemas']['GameEndEvent'];
export type Round = components['schemas']['Round'];
export type Stroke = components['schemas']['Stroke'];
