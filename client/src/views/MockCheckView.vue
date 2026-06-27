<script setup lang="ts">
import { io, type Socket } from 'socket.io-client';
import { onBeforeUnmount, onMounted, ref } from 'vue';

type MockStatusResponse = {
  status: 'ok';
  service: 'mock-server';
  now: string;
};
type MockPongPayload = {
  ok: boolean;
  sentAt: string;
  receivedAt: string;
};

const restStatus = ref('not checked');
const socketStatus = ref('not connected');
const pongStatus = ref('not checked');

let socket: Socket | undefined;

async function checkRest() {
  restStatus.value = 'checking';

  try {
    const response = await fetch('/api/mock/status');

    if (!response.ok) {
      restStatus.value = 'request failed';
      return;
    }

    const data = (await response.json()) as MockStatusResponse;
    restStatus.value = `${data.status} (${data.service}, ${data.now})`;
  } catch (error) {
    restStatus.value = error instanceof Error ? error.message : 'request failed';
  }
}

function connectSocket() {
  if (socket) {
    return;
  }

  socketStatus.value = 'connecting';
  socket = io({ path: '/socket.io' });

  socket.on('connect', () => {
    socketStatus.value = `connected (${socket?.id ?? 'unknown'})`;
  });

  socket.on('disconnect', (reason) => {
    socketStatus.value = `disconnected (${reason})`;
  });

  socket.on('connect_error', (error) => {
    socketStatus.value = error.message;
  });

  socket.on('mock:pong', (payload: MockPongPayload) => {
    pongStatus.value = `${payload.ok ? 'ok' : 'ng'} (${payload.receivedAt})`;
  });
}

function sendPing() {
  if (!socket?.connected) {
    pongStatus.value = 'socket is not connected';
    return;
  }

  pongStatus.value = 'waiting';
  socket.emit('mock:ping', { sentAt: new Date().toISOString() });
}

onMounted(() => {
  void checkRest();
  connectSocket();
});

onBeforeUnmount(() => {
  socket?.disconnect();
  socket = undefined;
});
</script>

<template>
  <section>
    <h1>Mock Check</h1>

    <dl>
      <dt>REST</dt>
      <dd>{{ restStatus }}</dd>

      <dt>Socket.IO</dt>
      <dd>{{ socketStatus }}</dd>

      <dt>Ping</dt>
      <dd>{{ pongStatus }}</dd>
    </dl>

    <button type="button" @click="checkRest">Check REST</button>
    <button type="button" @click="sendPing">Send Ping</button>
  </section>
</template>

<style scoped></style>
