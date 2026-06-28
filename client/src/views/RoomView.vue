<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { watch } from 'vue';
import { useRouter } from 'vue-router';

import { useRoomStore } from '@/stores/room';

import PlayingRoomView from './PlayingRoomView.vue';
import WaitingRoomView from './WaitingRoomView.vue';

const router = useRouter();
const { currentRoom } = storeToRefs(useRoomStore());

watch(
  currentRoom,
  (room) => {
    if (room === null) {
      router.replace('/');
    }
  },
  { immediate: true },
);
</script>

<template>
  <WaitingRoomView v-if="currentRoom?.status === 'waiting'" />
  <PlayingRoomView v-else-if="currentRoom?.status === 'playing'" />
</template>
