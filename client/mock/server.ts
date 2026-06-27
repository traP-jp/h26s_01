import { createServer } from 'node:http'
import { randomUUID } from 'node:crypto'

import express from 'express'
import { Server } from 'socket.io'

import type {
  CreateRoomRequest,
  MockPingPayload,
  MockPongPayload,
  MockStatusResponse,
  Room,
  User,
} from './types'

const port = Number.parseInt(process.env.MOCK_PORT ?? '3000', 10)
const mockUserCookieName = 'mock_user_id'

const app = express()
const server = createServer(app)
const io = new Server(server, {
  cors: {
    origin: true,
  },
})

let nextMockUserIndex = 1
const rooms: Room[] = [
  {
    id: randomUUID(),
    name: '練習部屋',
    members: [],
    status: 'waiting',
  },
]

app.use(express.json())

const getCookieValue = (cookieHeader: string | undefined, name: string): string | null => {
  if (cookieHeader === undefined) {
    return null
  }

  const cookies = cookieHeader.split(';')

  for (const cookie of cookies) {
    const [rawKey, ...rawValueParts] = cookie.trim().split('=')

    if (rawKey === name) {
      try {
        return decodeURIComponent(rawValueParts.join('='))
      } catch {
        return null
      }
    }
  }

  return null
}

const createMockUserId = () => {
  const userId = `mock-user-${nextMockUserIndex}`
  nextMockUserIndex += 1
  return userId
}

const getOrCreateMockUser = (
  request: express.Request,
  response: express.Response,
): User => {
  const userId =
    getCookieValue(request.headers.cookie, mockUserCookieName) ?? createMockUserId()

  response.cookie(mockUserCookieName, userId, {
    httpOnly: true,
    path: '/',
    sameSite: 'lax',
  })

  return {
    id: userId,
  }
}

const addMemberIfNeeded = (room: Room, user: User) => {
  if (room.members.some((member) => member.id === user.id)) {
    return
  }

  room.members.push({
    ...user,
    isReady: false,
  })
}

app.get('/api/mock/status', (_request, response) => {
  const body: MockStatusResponse = {
    status: 'ok',
    service: 'mock-server',
    now: new Date().toISOString(),
  }

  response.json(body)
})

app.get('/api/me', (request, response) => {
  response.json(getOrCreateMockUser(request, response))
})

app.get('/api/rooms', (_request, response) => {
  response.json(rooms)
})

app.post('/api/rooms', (request, response) => {
  const user = getOrCreateMockUser(request, response)
  const body =
    request.body !== null && typeof request.body === 'object'
      ? (request.body as Partial<CreateRoomRequest>)
      : {}
  const name = typeof body.name === 'string' ? body.name.trim() : ''

  if (name.length === 0 || name.length > 64) {
    response.status(400).json({
      message: '部屋名は1文字以上64文字以下で入力してください',
    })
    return
  }

  const room: Room = {
    id: randomUUID(),
    name,
    members: [],
    status: 'waiting',
  }

  addMemberIfNeeded(room, user)
  rooms.push(room)

  response.status(201).json(room)
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
