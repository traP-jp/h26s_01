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
