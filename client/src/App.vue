<script setup lang="ts">
import { onBeforeUnmount, onMounted, onUnmounted } from 'vue';

import { useGameSocket } from './composables/useGameSocket';
import { useLobbySocketEvents } from './composables/useLobbySocketEvents';
import { useRoomSocket } from './composables/useRoomSocket';
import { useSession } from './composables/useSession';

const lobbySocketEvents = useLobbySocketEvents();
const roomSocket = useRoomSocket();
const gameSocket = useGameSocket();
const session = useSession();

onMounted(() => {
  lobbySocketEvents.register();
  roomSocket.register();
  gameSocket.register();
  session.initializeSession();
});

onBeforeUnmount(() => {
  lobbySocketEvents.cleanup();
  roomSocket.cleanup();
  gameSocket.cleanup();
  session.initializeSession();
});
</script>

<template>
  <main>
    <RouterView />
  </main>
</template>

<style scoped></style>
