import type { components } from '../src/api/schema'

export type MockStatusResponse = {
  status: 'ok'
  service: 'mock-server'
  now: string
}

export type MockPingPayload = {
  sentAt: string
}

export type MockPongPayload = MockPingPayload & {
  ok: true
  receivedAt: string
}

export type User = components['schemas']['User']
export type Room = components['schemas']['Room']
export type CreateRoomRequest = components['schemas']['CreateRoomRequest']
export type RoomJoinEvent = components['schemas']['RoomJoinEvent']
export type DrawStrokeEvent = components['schemas']['DrawStrokeEvent']
export type AnswerSubmitEvent = components['schemas']['AnswerSubmitEvent']
export type RoomListUpdatedEvent = components['schemas']['RoomListUpdatedEvent']
export type RoomUpdatedEvent = components['schemas']['RoomUpdatedEvent']
export type RoundStartedEvent = components['schemas']['RoundStartedEvent']
export type TurnStartedEvent = components['schemas']['TurnStartedEvent']
export type ClientDisconnectedEvent = components['schemas']['ClientDisconnectedEvent']
export type RoundAnswerEvent = components['schemas']['RoundAnswerEvent']
export type GameEndEvent = components['schemas']['GameEndEvent']
export type Stroke = components['schemas']['Stroke']

export type MockGameState = {
  roomId: string
  roundIndex: number
  turnIndex: number
  guesserId: string
  drawerIds: string[]
  currentDrawerId: string
  strokes: Stroke[]
  startedAt: number
  guesserAnswer: string | null
  actualAnswer: string
}

export type ClientToServerEvents = {
  'room:join': (payload: RoomJoinEvent) => void
  'game:ready': () => void
  'draw:stroke': (payload: DrawStrokeEvent) => void
  'answer:submit': (payload: AnswerSubmitEvent) => void
  'round:end': () => void
  'mock:ping': (payload: MockPingPayload) => void
}

export type ServerToClientEvents = {
  'room_list:updated': (payload: RoomListUpdatedEvent) => void
  'room:updated': (payload: RoomUpdatedEvent) => void
  'round:started': (payload: RoundStartedEvent) => void
  'turn:started': (payload: TurnStartedEvent) => void
  'draw:stroke': (payload: Stroke) => void
  'client:disconnected': (payload: ClientDisconnectedEvent) => void
  'round:answer': (payload: RoundAnswerEvent) => void
  'game:end': (payload: GameEndEvent) => void
  'mock:pong': (payload: MockPongPayload) => void
}
