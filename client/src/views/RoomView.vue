<script setup lang="ts">
import { storeToRefs } from 'pinia';
import { computed, onBeforeUnmount, onMounted, watch } from 'vue';
import { onBeforeRouteLeave, useRouter } from 'vue-router';

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

const shouldConfirmLeave = computed(() => {
  if (['ended', 'aborted'].includes(gameStore.phase)) {
    return false;
  }

  return currentRoom.value !== null;
});

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

const handleBeforeUnload = (event: BeforeUnloadEvent) => {
  if (!shouldConfirmLeave.value) {
    return;
  }

  event.preventDefault();
  event.returnValue = '';
};

onMounted(() => {
  window.addEventListener('beforeunload', handleBeforeUnload);
});

onBeforeUnmount(() => {
  window.removeEventListener('beforeunload', handleBeforeUnload);
});

onBeforeRouteLeave(() => {
  if (!shouldConfirmLeave.value) {
    return true;
  }

  return window.confirm('ゲームを中断してこのページを離れますか？');
});
</script>

<template>
  <WaitingRoomView v-if="shouldShowWaitingRoom" />
  <PlayingRoomView v-else-if="shouldShowPlayingRoom" />
</template>
