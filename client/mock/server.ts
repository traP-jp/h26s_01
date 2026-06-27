import { createServer } from 'node:http'
import { randomUUID } from 'node:crypto'

import express from 'express'
import { Server } from 'socket.io'

import type {
  ClientToServerEvents,
  CreateRoomRequest,
  MockPingPayload,
  MockPongPayload,
  MockGameState,
  MockStatusResponse,
  Room,
  ServerToClientEvents,
  User,
} from './types'

const port = Number.parseInt(process.env.MOCK_PORT ?? '3000', 10)
const mockUserCookieName = 'mock_user_id'
const mockActualAnswer = '森'
const mockKanji = '森'

const app = express()
const server = createServer(app)
const io = new Server<ClientToServerEvents, ServerToClientEvents>(server, {
  cors: {
    origin: true,
  },
})

let nextMockUserIndex = 1
const socketUserIds = new Map<string, string>()
const socketRoomIds = new Map<string, string>()
const games = new Map<string, MockGameState>()
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

const getMockUserFromCookie = (cookieHeader: string | undefined): User | null => {
  const userId = getCookieValue(cookieHeader, mockUserCookieName)

  if (userId === null) {
    return null
  }

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

const findRoom = (roomId: string): Room | null =>
  rooms.find((room) => room.id === roomId) ?? null

const emitRoomUpdated = (room: Room) => {
  io.to(room.id).emit('room:updated', {
    eventType: 'room_updated',
    room,
  })
}

const emitRoomListUpdated = (room: Room) => {
  io.emit('room_list:updated', {
    eventType: 'room_updated',
    room,
  })
}

const startGameIfReady = (room: Room) => {
  if (
    room.status !== 'waiting' ||
    room.members.length < 2 ||
    !room.members.every((member) => member.isReady)
  ) {
    return
  }

  const guesser = room.members.at(-1)

  if (guesser === undefined) {
    return
  }

  const drawerIds = room.members
    .map((member) => member.id)
    .filter((memberId) => memberId !== guesser.id)
  const currentDrawerId = drawerIds[0]

  if (currentDrawerId === undefined) {
    return
  }

  room.status = 'playing'
  const game: MockGameState = {
    roomId: room.id,
    roundIndex: 1,
    turnIndex: 1,
    guesserId: guesser.id,
    drawerIds,
    currentDrawerId,
    strokes: [],
    startedAt: Date.now(),
    guesserAnswer: null,
    actualAnswer: mockActualAnswer,
  }

  games.set(room.id, game)
  emitRoomUpdated(room)
  emitRoomListUpdated(room)
  io.to(room.id).emit('round:started', {
    roundIndex: game.roundIndex,
    guesserId: game.guesserId,
    kanji: mockKanji,
  })
  io.to(room.id).emit('turn:started', {
    turnIndex: game.turnIndex,
    drawerId: game.currentDrawerId,
  })
}

const endGame = (room: Room, game: MockGameState) => {
  const totalTimeMs = Date.now() - game.startedAt

  io.to(room.id).emit('game:end', {
    cleared: true,
    totalTimeMs,
    remainingLives: 1,
    rounds: [
      {
        id: randomUUID(),
        timeMs: totalTimeMs,
        guesserId: game.guesserId,
        guesserAnswer: game.guesserAnswer ?? '',
        actualAnswer: game.actualAnswer,
        strokes: game.strokes,
      },
    ],
  })

  games.delete(room.id)
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
  io.emit('room_list:updated', {
    eventType: 'room_created',
    room,
  })

  response.status(201).json(room)
})

io.on('connection', (socket) => {
  const user = getMockUserFromCookie(socket.handshake.headers.cookie)

  if (user !== null) {
    socketUserIds.set(socket.id, user.id)
  }

  socket.on('mock:ping', (payload: MockPingPayload) => {
    const pong: MockPongPayload = {
      ok: true,
      sentAt: payload.sentAt,
      receivedAt: new Date().toISOString(),
    }

    socket.emit('mock:pong', pong)
  })

  socket.on('room:join', (payload) => {
    const userId = socketUserIds.get(socket.id)
    const room = findRoom(payload.roomId)

    if (userId === undefined || room === null || room.status !== 'waiting') {
      return
    }

    addMemberIfNeeded(room, {
      id: userId,
    })
    socket.join(room.id)
    socketRoomIds.set(socket.id, room.id)
    emitRoomUpdated(room)
    emitRoomListUpdated(room)
  })

  socket.on('game:ready', () => {
    const userId = socketUserIds.get(socket.id)
    const roomId = socketRoomIds.get(socket.id)
    const room = roomId === undefined ? null : findRoom(roomId)

    if (userId === undefined || room === null || room.status !== 'waiting') {
      return
    }

    const member = room.members.find((item) => item.id === userId)

    if (member === undefined) {
      return
    }

    member.isReady = true
    emitRoomUpdated(room)
    emitRoomListUpdated(room)
    startGameIfReady(room)
  })

  socket.on('draw:stroke', (payload) => {
    const userId = socketUserIds.get(socket.id)
    const roomId = socketRoomIds.get(socket.id)
    const room = roomId === undefined ? null : findRoom(roomId)
    const game = roomId === undefined ? undefined : games.get(roomId)

    if (
      userId === undefined ||
      room === null ||
      game === undefined ||
      userId !== game.currentDrawerId
    ) {
      return
    }

    const stroke = {
      drawerId: userId,
      x1: payload.x1,
      y1: payload.y1,
      x2: payload.x2,
      y2: payload.y2,
    }

    game.strokes.push(stroke)
    io.to(room.id).emit('draw:stroke', stroke)
  })

  socket.on('answer:submit', (payload) => {
    const userId = socketUserIds.get(socket.id)
    const roomId = socketRoomIds.get(socket.id)
    const room = roomId === undefined ? null : findRoom(roomId)
    const game = roomId === undefined ? undefined : games.get(roomId)

    if (
      userId === undefined ||
      room === null ||
      game === undefined ||
      userId !== game.guesserId
    ) {
      return
    }

    game.guesserAnswer = payload.answer
    io.to(room.id).emit('round:answer', {
      guesserAnswer: payload.answer,
      actualAnswer: game.actualAnswer,
    })
  })

  socket.on('round:end', () => {
    const roomId = socketRoomIds.get(socket.id)
    const room = roomId === undefined ? null : findRoom(roomId)
    const game = roomId === undefined ? undefined : games.get(roomId)

    if (room === null || game === undefined) {
      return
    }

    endGame(room, game)
  })

  socket.on('disconnect', () => {
    const userId = socketUserIds.get(socket.id)
    const roomId = socketRoomIds.get(socket.id)
    const room = roomId === undefined ? null : findRoom(roomId)

    if (userId !== undefined && room !== null && room.status === 'playing') {
      io.to(room.id).emit('client:disconnected', {
        userId,
      })
      games.delete(room.id)
    }

    socketUserIds.delete(socket.id)
    socketRoomIds.delete(socket.id)
  })
})

server.listen(port, () => {
  console.log(`mock server listening on http://localhost:${port}`)
})
