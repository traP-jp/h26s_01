<script lang="ts" setup>
import { storeToRefs } from 'pinia';
import { onMounted } from 'vue';

import { useLobby } from '@/composables/useLobby.ts';
import { useLobbyStore } from '@/stores/lobby.ts';

import RoomListItem from './RoomListItem.vue';

const { fetchRooms } = useLobby();
const { waitingRooms } = storeToRefs(useLobbyStore());

onMounted(async () => {
  console.log('DEBUG1');
  const r = await fetchRooms();
  console.log('DEBUG2');
  console.log(waitingRooms);
  console.log(r);
});
</script>

<template>
  <section class="space-y-5 p-8">
    <h2 class="text-3xl">部屋の一覧</h2>
    <ul class="flex flex-col gap-2">
      <RoomListItem v-for="room in waitingRooms" :key="room.id" :room="room" />
    </ul>
  </section>
</template>
