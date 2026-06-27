import { createServer } from 'node:http'
import { randomUUID } from 'node:crypto'
import { existsSync } from 'node:fs'
import { resolve } from 'node:path'

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

const port = Number.parseInt(process.env.PORT ?? process.env.MOCK_PORT ?? '3000', 10)
const staticDirectory = resolve(process.cwd(), process.env.MOCK_STATIC_DIR ?? 'dist')
const staticIndexHtml = resolve(staticDirectory, 'index.html')
const mockUserCookieName = 'mock_user_id'
const mockCookieSecure = process.env.MOCK_COOKIE_SECURE === 'true'
const maxStrokes = 9
const mockRounds = [
  {
    kanji: '森',
    actualAnswer: '森',
  },
  {
    kanji: '箱',
    actualAnswer: '箱',
  },
] as const

const app = express()
const server = createServer(app)
const io = new Server<ClientToServerEvents, ServerToClientEvents>(server, {
  cors: {
    origin: true,
  },
})

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

const createMockUserId = () => `mock-user-${randomUUID()}`

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
    secure: mockCookieSecure,
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

const emitRoomDeleted = (roomId: string) => {
  io.emit('room_list:updated', {
    eventType: 'room_deleted',
    roomId,
  })
}

const removeRoom = (roomId: string) => {
  const roomIndex = rooms.findIndex((room) => room.id === roomId)

  if (roomIndex !== -1) {
    rooms.splice(roomIndex, 1)
  }
}

const deleteRoom = (roomId: string) => {
  games.delete(roomId)
  removeRoom(roomId)
  emitRoomDeleted(roomId)
}

const removeMember = (room: Room, userId: string) => {
  room.members = room.members.filter((member) => member.id !== userId)
}

const hasActiveSocketInRoom = (
  roomId: string,
  userId: string,
  ignoredSocketId: string,
) => {
  for (const [socketId, activeUserId] of socketUserIds) {
    if (socketId === ignoredSocketId || activeUserId !== userId) {
      continue
    }

    if (socketRoomIds.get(socketId) === roomId) {
      return true
    }
  }

  return false
}

const isFiniteStrokePayload = (payload: {
  x1: number
  y1: number
  x2: number
  y2: number
}) =>
  Number.isFinite(payload.x1) &&
  Number.isFinite(payload.y1) &&
  Number.isFinite(payload.x2) &&
  Number.isFinite(payload.y2)

const getCurrentMockRound = (game: MockGameState) => mockRounds[game.roundIndex - 1]

const emitRoundStarted = (room: Room, game: MockGameState) => {
  const mockRound = getCurrentMockRound(game)

  if (mockRound === undefined) {
    return
  }

  io.to(room.id).emit('round:started', {
    roundIndex: game.roundIndex,
    guesserId: game.guesserId,
    kanji: mockRound.kanji,
  })
  io.to(room.id).emit('turn:started', {
    turnIndex: game.turnIndex,
    drawerId: game.currentDrawerId,
  })
}

const emitTurnStarted = (room: Room, game: MockGameState) => {
  io.to(room.id).emit('turn:started', {
    turnIndex: game.turnIndex,
    drawerId: game.currentDrawerId,
  })
}

const advanceTurn = (room: Room, game: MockGameState) => {
  if (game.strokes.length >= maxStrokes) {
    return
  }

  game.turnIndex += 1
  const nextDrawerId = game.drawerIds[(game.turnIndex - 1) % game.drawerIds.length]

  if (nextDrawerId === undefined) {
    return
  }

  game.currentDrawerId = nextDrawerId
  emitTurnStarted(room, game)
}

const createRoundRoles = (room: Room, roundIndex: number) => {
  const reversedMembers = [...room.members].reverse()
  const guesser = reversedMembers[(roundIndex - 1) % reversedMembers.length]

  if (guesser === undefined) {
    return null
  }

  const drawerIds = room.members
    .map((member) => member.id)
    .filter((memberId) => memberId !== guesser.id)
  const currentDrawerId = drawerIds[0]

  if (currentDrawerId === undefined) {
    return null
  }

  return {
    currentDrawerId,
    drawerIds,
    guesserId: guesser.id,
  }
}

const completeCurrentRound = (game: MockGameState) => {
  game.completedRounds.push({
    id: randomUUID(),
    timeMs: Date.now() - game.roundStartedAt,
    guesserId: game.guesserId,
    guesserAnswer: game.guesserAnswer ?? '',
    actualAnswer: game.actualAnswer,
    strokes: game.strokes,
  })
}

const startNextRound = (room: Room, game: MockGameState) => {
  const nextMockRound = mockRounds[game.roundIndex]

  if (nextMockRound === undefined) {
    return false
  }

  game.roundIndex += 1
  const roles = createRoundRoles(room, game.roundIndex)

  if (roles === null) {
    return false
  }

  game.turnIndex = 1
  game.guesserId = roles.guesserId
  game.drawerIds = roles.drawerIds
  game.currentDrawerId = roles.currentDrawerId
  game.strokes = []
  game.roundStartedAt = Date.now()
  game.guesserAnswer = null
  game.actualAnswer = nextMockRound.actualAnswer
  emitRoundStarted(room, game)
  return true
}

const startGameIfReady = (room: Room) => {
  if (
    room.status !== 'waiting' ||
    room.members.length < 2 ||
    !room.members.every((member) => member.isReady)
  ) {
    return
  }

  const roles = createRoundRoles(room, 1)

  if (roles === null) {
    return
  }

  room.status = 'playing'
  const firstMockRound = mockRounds[0]
  const game: MockGameState = {
    roomId: room.id,
    roundIndex: 1,
    turnIndex: 1,
    guesserId: roles.guesserId,
    drawerIds: roles.drawerIds,
    currentDrawerId: roles.currentDrawerId,
    strokes: [],
    startedAt: Date.now(),
    roundStartedAt: Date.now(),
    guesserAnswer: null,
    actualAnswer: firstMockRound.actualAnswer,
    completedRounds: [],
  }

  games.set(room.id, game)
  emitRoomUpdated(room)
  emitRoomListUpdated(room)
  emitRoundStarted(room, game)
}

const endGame = (room: Room, game: MockGameState) => {
  const totalTimeMs = Date.now() - game.startedAt

  io.to(room.id).emit('game:end', {
    cleared: true,
    totalTimeMs,
    remainingLives: 1,
    rounds: game.completedRounds,
  })

  deleteRoom(room.id)
}

const endRound = (room: Room, game: MockGameState) => {
  completeCurrentRound(game)

  if (startNextRound(room, game)) {
    return
  }

  endGame(room, game)
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

if (existsSync(staticIndexHtml)) {
  app.use(express.static(staticDirectory))
  app.use((request, response, next) => {
    if (
      request.method !== 'GET' ||
      request.path.startsWith('/api') ||
      request.path.startsWith('/socket.io') ||
      !request.accepts('html')
    ) {
      next()
      return
    }

    response.sendFile(staticIndexHtml)
  })
}

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
      room.status !== 'playing' ||
      game === undefined ||
      userId !== game.currentDrawerId ||
      game.strokes.length >= maxStrokes ||
      game.guesserAnswer !== null ||
      !isFiniteStrokePayload(payload)
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
    advanceTurn(room, game)
  })

  socket.on('answer:submit', (payload) => {
    const userId = socketUserIds.get(socket.id)
    const roomId = socketRoomIds.get(socket.id)
    const room = roomId === undefined ? null : findRoom(roomId)
    const game = roomId === undefined ? undefined : games.get(roomId)

    if (
      userId === undefined ||
      room === null ||
      room.status !== 'playing' ||
      game === undefined ||
      userId !== game.guesserId ||
      game.guesserAnswer !== null
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
    const userId = socketUserIds.get(socket.id)
    const roomId = socketRoomIds.get(socket.id)
    const room = roomId === undefined ? null : findRoom(roomId)
    const game = roomId === undefined ? undefined : games.get(roomId)

    if (
      userId === undefined ||
      room === null ||
      room.status !== 'playing' ||
      game === undefined ||
      game.guesserAnswer === null ||
      !room.members.some((member) => member.id === userId)
    ) {
      return
    }

    endRound(room, game)
  })

  socket.on('disconnect', () => {
    const userId = socketUserIds.get(socket.id)
    const roomId = socketRoomIds.get(socket.id)
    const room = roomId === undefined ? null : findRoom(roomId)

    if (
      userId !== undefined &&
      roomId !== undefined &&
      hasActiveSocketInRoom(roomId, userId, socket.id)
    ) {
      socketUserIds.delete(socket.id)
      socketRoomIds.delete(socket.id)
      return
    }

    if (userId !== undefined && room !== null && room.status === 'waiting') {
      removeMember(room, userId)

      if (room.members.length === 0) {
        deleteRoom(room.id)
      } else {
        emitRoomUpdated(room)
        emitRoomListUpdated(room)
      }
    } else if (
      userId !== undefined &&
      room !== null &&
      room.status === 'playing' &&
      games.has(room.id)
    ) {
      io.to(room.id).emit('client:disconnected', {
        userId,
      })
      deleteRoom(room.id)
    }

    socketUserIds.delete(socket.id)
    socketRoomIds.delete(socket.id)
  })
})

server.listen(port, () => {
  console.log(`mock server listening on http://localhost:${port}`)
  console.log(`mock static directory: ${staticDirectory}`)
})
