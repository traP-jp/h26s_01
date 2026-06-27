<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue';

import { useGameSocket } from '@/composables/useGameSocket';
import { useLobby } from '@/composables/useLobby';
import { useLobbySocketEvents } from '@/composables/useLobbySocketEvents';
import { useRoomSocket } from '@/composables/useRoomSocket';
import { useSession } from '@/composables/useSession';
import { useGameStore } from '@/stores/game';
import { useRoomStore } from '@/stores/room';
import type {
  ClientDisconnectedEvent,
  GameEndEvent,
  RoomListUpdatedEvent,
  RoomUpdatedEvent,
  RoundAnswerEvent,
  RoundStartedEvent,
  Stroke,
  TurnStartedEvent,
} from '@/types/api';

type EventLogItem = {
  id: number;
  text: string;
};

const roomName = ref('sample room');
const answer = ref('森');
const eventLog = ref<EventLogItem[]>([]);
let nextEventLogId = 1;

const { initializeSession, socketConnection, userStore } = useSession();
const { createRoom, fetchRooms, lobbyStore } = useLobby();
const lobbySocketEvents = useLobbySocketEvents();
const roomSocket = useRoomSocket();
const gameSocket = useGameSocket();
const roomStore = useRoomStore();
const gameStore = useGameStore();

const selectedRoomId = computed(() => roomStore.roomId ?? lobbyStore.rooms[0]?.id ?? '');
const latestEventLog = computed(() => eventLog.value.slice(-12).reverse());
const hasJoinedRoom = computed(() => roomStore.currentRoom !== null);
const isWaiting = computed(() => roomStore.currentRoom?.status === 'waiting');
const isPlaying = computed(() => gameStore.phase === 'playing');
const hasRoundResult = computed(() => gameStore.phase === 'roundResult');
const hasGameEnded = computed(() => gameStore.phase === 'ended' || gameStore.phase === 'aborted');
const currentStep = computed(() => {
  if (hasGameEnded.value) {
    return 'Result';
  }

  if (isPlaying.value || hasRoundResult.value) {
    return 'Game';
  }

  if (hasJoinedRoom.value) {
    return 'Waiting room';
  }

  return 'Lobby';
});

const log = (message: string) => {
  eventLog.value.push({
    id: nextEventLogId,
    text: `${new Date().toLocaleTimeString()} ${message}`,
  });
  nextEventLogId += 1;
};

const initialize = async () => {
  const user = await initializeSession();

  if (user === null) {
    log('session initialization failed');
    return;
  }

  log(`session initialized as ${user.id}`);
  await fetchRooms();
};

const logRoomListUpdated = (event: RoomListUpdatedEvent) => {
  log(`received room_list:updated (${event.eventType})`);
};

const logRoomUpdated = (event: RoomUpdatedEvent) => {
  log(`received room:updated (${event.room.name} / ${event.room.status})`);
};

const logRoundStarted = (event: RoundStartedEvent) => {
  log(
    `received round:started (${event.roundIndex}, ${event.kanji ?? 'hidden'}, guesser: ${event.guesserId})`,
  );
};

const logTurnStarted = (event: TurnStartedEvent) => {
  log(`received turn:started (drawer: ${event.drawerId})`);
};

const logStroke = (stroke: Stroke) => {
  log(`received draw:stroke (${stroke.drawerId})`);
};

const logRoundAnswer = (event: RoundAnswerEvent) => {
  log(`received round:answer (${event.guesserAnswer} / ${event.actualAnswer})`);
};

const logGameEnd = (event: GameEndEvent) => {
  log(`received game:end (${event.cleared ? 'cleared' : 'failed'})`);
};

const logClientDisconnected = (event: ClientDisconnectedEvent) => {
  log(`received client:disconnected (${event.userId})`);
};

const refreshRooms = async () => {
  const rooms = await fetchRooms();
  log(rooms === null ? 'rooms fetch failed' : `rooms fetched (${rooms.length})`);
};

const createAndJoinRoom = async () => {
  if (hasJoinedRoom.value) {
    log('already joined a room; room leave is not defined in the current API');
    return;
  }

  const room = await createRoom({
    name: roomName.value,
  });

  if (room === null) {
    log('room creation failed');
    return;
  }

  log(`room created: ${room.name}`);
  const didJoin = roomSocket.joinRoom(room.id);
  log(didJoin ? `join requested: ${room.id}` : 'join skipped');
};

const joinRoom = (roomId: string) => {
  if (hasJoinedRoom.value) {
    log('already joined a room; room leave is not defined in the current API');
    return;
  }

  if (roomId.length === 0) {
    log('no room to join');
    return;
  }

  const didJoin = roomSocket.joinRoom(roomId);
  log(didJoin ? `join requested: ${roomId}` : 'join skipped');
};

const joinSelectedRoom = () => {
  joinRoom(selectedRoomId.value);
};

const sendReady = () => {
  const didSend = roomSocket.sendReady();
  log(didSend ? 'ready requested' : 'ready skipped');
};

const drawSampleStroke = () => {
  const didSend = gameSocket.drawStroke({
    x1: 0.1,
    y1: 0.1,
    x2: 0.9,
    y2: 0.9,
  });
  log(didSend ? 'stroke requested' : 'stroke skipped');
};

const submitSampleAnswer = () => {
  const didSend = gameSocket.submitAnswer(answer.value);
  log(didSend ? `answer submitted: ${answer.value}` : 'answer skipped');
};

const endRound = () => {
  const didSend = gameSocket.endRound();
  log(didSend ? 'round end requested' : 'round end skipped');
};

onMounted(() => {
  lobbySocketEvents.register();
  roomSocket.register();
  gameSocket.register();
  lobbySocketEvents.onRoomListUpdated(logRoomListUpdated);
  roomSocket.onRoomUpdated(logRoomUpdated);
  gameSocket.onRoundStarted(logRoundStarted);
  gameSocket.onTurnStarted(logTurnStarted);
  gameSocket.onStroke(logStroke);
  gameSocket.onRoundAnswer(logRoundAnswer);
  gameSocket.onGameEnd(logGameEnd);
  gameSocket.onClientDisconnected(logClientDisconnected);
  void initialize();
});

onBeforeUnmount(() => {
  lobbySocketEvents.offRoomListUpdated(logRoomListUpdated);
  roomSocket.offRoomUpdated(logRoomUpdated);
  gameSocket.offRoundStarted(logRoundStarted);
  gameSocket.offTurnStarted(logTurnStarted);
  gameSocket.offStroke(logStroke);
  gameSocket.offRoundAnswer(logRoundAnswer);
  gameSocket.offGameEnd(logGameEnd);
  gameSocket.offClientDisconnected(logClientDisconnected);
  lobbySocketEvents.cleanup();
  roomSocket.cleanup();
  gameSocket.cleanup();
});
</script>

<template>
  <section class="min-h-screen bg-background px-8 py-6 text-primary">
    <div class="mx-auto max-w-5xl space-y-5">
      <div class="space-y-1 border-b border-primary pb-3">
        <h1 class="text-3xl leading-tight">Socket Sample</h1>
        <p class="text-sm">
          Current step: <span class="font-bold">{{ currentStep }}</span>
          <span class="ml-3">mock socket と store の状態遷移確認用ページ</span>
        </p>
      </div>

      <div class="space-y-4">
        <section class="border border-primary p-4">
          <h2 class="mb-3 text-xl">1. Session</h2>
          <dl class="grid grid-cols-[7rem_minmax(0,1fr)] gap-x-4 gap-y-2 text-sm">
            <dt class="text-sm opacity-70">User</dt>
            <dd class="break-all">{{ userStore.currentUser?.id ?? 'not loaded' }}</dd>

            <dt class="text-sm opacity-70">Socket</dt>
            <dd>
              <span class="font-bold">
                {{ socketConnection.isConnected ? 'connected' : 'disconnected' }}
              </span>
              <span v-if="socketConnection.connectionError">
                ({{ socketConnection.connectionError }})
              </span>
            </dd>
          </dl>
          <button
            type="button"
            class="mt-3 border border-primary px-3 py-1 text-sm hover:bg-primary hover:text-background"
            @click="initialize"
          >
            Initialize
          </button>
        </section>

        <section class="border border-primary p-4">
          <h2 class="mb-3 text-xl">2. Lobby</h2>
          <p v-if="hasJoinedRoom" class="mb-3 text-sm">
            You are already in a room. The current API has no room leave event, so joining another
            room is disabled in this sample.
          </p>
          <div class="flex flex-wrap items-end gap-2">
            <label class="flex min-w-64 flex-col gap-1 text-sm">
              Room name
              <input
                v-model="roomName"
                type="text"
                class="border border-primary bg-background px-2 py-1 text-sm outline-none focus:bg-white"
              />
            </label>
            <button
              type="button"
              class="border border-primary px-3 py-1 text-sm hover:bg-primary hover:text-background"
              @click="refreshRooms"
            >
              Refresh rooms
            </button>
            <button
              type="button"
              :disabled="hasJoinedRoom"
              class="border border-primary px-3 py-1 text-sm hover:bg-primary hover:text-background disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-background disabled:hover:text-primary"
              @click="createAndJoinRoom"
            >
              Create and join
            </button>
          </div>

          <ul class="mt-3 space-y-1">
            <li
              v-for="room in lobbyStore.rooms"
              :key="room.id"
              class="flex flex-wrap items-center gap-2 text-sm"
            >
              <button
                type="button"
                :disabled="hasJoinedRoom"
                class="border border-primary px-2 py-0.5 text-xs hover:bg-primary hover:text-background disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-background disabled:hover:text-primary"
                @click="joinRoom(room.id)"
              >
                Join
              </button>
              <span class="text-base">{{ room.name }}</span>
              <span class="opacity-70">{{ room.status }} / {{ room.members.length }} members</span>
            </li>
          </ul>
        </section>

        <section class="border border-primary p-4">
          <h2 class="mb-3 text-xl">3. Waiting room</h2>
          <dl class="grid grid-cols-[7rem_minmax(0,1fr)] gap-x-4 gap-y-2 text-sm">
            <dt class="text-sm opacity-70">Current room</dt>
            <dd>{{ roomStore.currentRoom?.name ?? 'not joined' }}</dd>

            <dt class="text-sm opacity-70">Members</dt>
            <dd class="flex flex-wrap gap-2">
              <span
                v-for="member in roomStore.members"
                :key="member.id"
                class="border border-primary/40 px-2 py-0.5 text-xs"
              >
                {{ member.id }}{{ member.isReady ? ' ready' : '' }}
              </span>
            </dd>
          </dl>
          <div class="mt-3 flex flex-wrap gap-2">
            <button
              type="button"
              :disabled="hasJoinedRoom"
              class="border border-primary px-3 py-1 text-sm hover:bg-primary hover:text-background disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-background disabled:hover:text-primary"
              @click="joinSelectedRoom"
            >
              Join first room
            </button>
            <button
              type="button"
              :disabled="!isWaiting || roomStore.isReady"
              class="border border-primary px-3 py-1 text-sm hover:bg-primary hover:text-background disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-background disabled:hover:text-primary"
              @click="sendReady"
            >
              Ready
            </button>
          </div>
          <div class="mt-3 space-y-1 text-sm opacity-80">
            <p>Room leave is not implemented because the current API has no leave event.</p>
            <p>
              In the mock server, closing the tab while waiting removes that member from the room.
              Closing the tab while playing aborts the game for the remaining clients.
            </p>
          </div>
        </section>

        <section class="border border-primary p-4">
          <h2 class="mb-3 text-xl">4. Game</h2>
          <dl class="grid grid-cols-[7rem_minmax(0,1fr)] gap-x-4 gap-y-2 text-sm">
            <dt class="text-sm opacity-70">Phase</dt>
            <dd>{{ gameStore.phase }}</dd>

            <dt class="text-sm opacity-70">Round</dt>
            <dd>{{ gameStore.roundLabel || '-' }}</dd>

            <dt class="text-sm opacity-70">Turn</dt>
            <dd>{{ gameStore.turnLabel || '-' }}</dd>

            <dt class="text-sm opacity-70">Role</dt>
            <dd>
              guesser: {{ gameStore.isGuesser }} / drawer: {{ gameStore.isDrawer }} / my turn:
              {{ gameStore.isMyTurn }}
            </dd>

            <dt class="text-sm opacity-70">Strokes</dt>
            <dd>{{ gameStore.strokeCount }}</dd>

            <dt class="text-sm opacity-70">Round answer</dt>
            <dd>
              {{ gameStore.roundAnswer?.guesserAnswer ?? '-' }} /
              {{ gameStore.roundAnswer?.actualAnswer ?? '-' }}
            </dd>
          </dl>

          <div class="mt-3 flex flex-wrap items-end gap-2">
            <button
              type="button"
              :disabled="!gameStore.canDraw"
              class="border border-primary px-3 py-1 text-sm hover:bg-primary hover:text-background disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-background disabled:hover:text-primary"
              @click="drawSampleStroke"
            >
              Draw sample stroke
            </button>
            <label class="flex min-w-40 flex-col gap-1 text-sm">
              Answer
              <input
                v-model="answer"
                type="text"
                class="border border-primary bg-background px-2 py-1 text-sm outline-none focus:bg-white"
              />
            </label>
            <button
              type="button"
              :disabled="!gameStore.canSubmitAnswer"
              class="border border-primary px-3 py-1 text-sm hover:bg-primary hover:text-background disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-background disabled:hover:text-primary"
              @click="submitSampleAnswer"
            >
              Submit answer
            </button>
            <button
              type="button"
              :disabled="!hasRoundResult"
              class="border border-primary px-3 py-1 text-sm hover:bg-primary hover:text-background disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-background disabled:hover:text-primary"
              @click="endRound"
            >
              End round
            </button>
          </div>
          <p class="mt-3 text-sm opacity-80">
            The mock game runs two fixed rounds, 森 then 箱. The first round end starts the next
            round, and the second round end sends game:end.
          </p>
        </section>

        <section class="border border-primary p-4">
          <h2 class="mb-3 text-xl">5. Event log</h2>
          <ol class="max-h-80 space-y-1 overflow-y-auto text-sm">
            <li v-for="item in latestEventLog" :key="item.id">
              {{ item.text }}
            </li>
          </ol>
        </section>
      </div>
    </div>
  </section>
</template>

<style scoped></style>
