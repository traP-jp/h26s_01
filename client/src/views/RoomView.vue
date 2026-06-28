<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { computed, watch } from 'vue';
import { useRouter } from 'vue-router';

import { useGameStore } from '@/stores/game';
import { useRoomStore } from '@/stores/room';

import PlayingRoomView from './PlayingRoomView.vue';
import WaitingRoomView from './WaitingRoomView.vue';

const router = useRouter();
const gameStore = useGameStore();
const { currentRoom } = storeToRefs(useRoomStore());

const canShowGamePhase = computed(() =>
  ['roundResult', 'ended', 'aborted'].includes(gameStore.phase),
);

const shouldShowWaitingRoom = computed(() => currentRoom.value?.status === 'waiting');
const shouldShowPlayingRoom = computed(
  () => currentRoom.value?.status === 'playing' || canShowGamePhase.value,
);

watch(
  [currentRoom, canShowGamePhase],
  ([room, canShowResult]) => {
    if (room === null && !canShowResult) {
      router.replace('/');
    }
  },
  { immediate: true },
);
</script>

<template>
  <WaitingRoomView v-if="shouldShowWaitingRoom" />
  <PlayingRoomView v-else-if="shouldShowPlayingRoom" />
</template>
