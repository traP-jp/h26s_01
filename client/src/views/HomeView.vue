<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

import Header from '@/components/common/TheHeader.vue';
import CreateRoomSection from '@/components/home/CreateRoomSection.vue';
import RoomListSection from '@/components/home/RoomListSection.vue';
import { useLobby } from '@/composables/useLobby';
import { useRoomSocket } from '@/composables/useRoomSocket';
import { useLobbyStore } from '@/stores/lobby';
import { useRoomStore } from '@/stores/room';
import { useUserStore } from '@/stores/user';
import type { RoomUpdatedEvent } from '@/types/api';

const router = useRouter();
const { createRoom, fetchRooms } = useLobby();
const roomSocket = useRoomSocket();
const lobbyStore = useLobbyStore();
const roomStore = useRoomStore();
const userStore = useUserStore();
const { waitingRooms, error, isLoading } = storeToRefs(lobbyStore);
const { currentRoom, isJoining } = storeToRefs(roomStore);
const { userId } = storeToRefs(userStore);

const joiningRoomId = ref<string | null>(null);
const defaultRoomName = computed(() => (userId.value === null ? '' : `${userId.value} の部屋`));

watch(
  currentRoom,
  (room) => {
    if (room !== null) {
      router.replace('/game');
    }
  },
  { immediate: true },
);

const requestJoin = (roomId: string) => {
  joiningRoomId.value = roomId;

  const didJoin = roomSocket.joinRoom(roomId);
  if (!didJoin) {
    joiningRoomId.value = null;
  }
};

const handleCreateRoom = async (roomName: string) => {
  const createdRoom = await createRoom({ name: roomName });
  if (createdRoom === null) {
    return;
  }

  requestJoin(createdRoom.id);
};

const handleJoinRoom = (roomId: string) => {
  requestJoin(roomId);
};

const handleRoomUpdated = (event: RoomUpdatedEvent) => {
  if (joiningRoomId.value === null || event.room.id !== joiningRoomId.value) {
    return;
  }

  joiningRoomId.value = null;
  router.replace('/game');
};

onMounted(() => {
  roomSocket.onRoomUpdated(handleRoomUpdated);
  fetchRooms();
});

onBeforeUnmount(() => {
  roomSocket.offRoomUpdated(handleRoomUpdated);
});
</script>

<template>
  <Header />
  <RoomListSection
    :error="error"
    :is-joining="isJoining || joiningRoomId !== null"
    :is-loading="isLoading"
    :rooms="waitingRooms"
    @join="handleJoinRoom"
  />
  <CreateRoomSection
    :disabled="isJoining || joiningRoomId !== null"
    :initial-room-name="defaultRoomName"
    @create="handleCreateRoom"
  />
</template>

<style scoped></style>
