import { createServer } from 'node:http'

import express from 'express'
import { Server } from 'socket.io'

import type { MockPingPayload, MockPongPayload, MockStatusResponse } from './types'

const port = Number.parseInt(process.env.MOCK_PORT ?? '3000', 10)

const app = express()
const server = createServer(app)
const io = new Server(server, {
  cors: {
    origin: true,
  },
})

app.get('/api/mock/status', (_request, response) => {
  const body: MockStatusResponse = {
    status: 'ok',
    service: 'mock-server',
    now: new Date().toISOString(),
  }

  response.json(body)
})

io.on('connection', (socket) => {
  socket.on('mock:ping', (payload: MockPingPayload) => {
    const pong: MockPongPayload = {
      ok: true,
      sentAt: payload.sentAt,
      receivedAt: new Date().toISOString(),
    }

    socket.emit('mock:pong', pong)
  })
})

server.listen(port, () => {
  console.log(`mock server listening on http://localhost:${port}`)
})
