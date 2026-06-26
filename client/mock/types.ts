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
