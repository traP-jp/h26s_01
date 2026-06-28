<script setup lang="ts">
import { ref } from 'vue';

import { useLobby } from '@/composables/useLobby.ts';
import { useRoomSocket } from '@/composables/useRoomSocket.ts';
import router from '@/router/index.ts';

import BaseButton from '../common/BaseButton.vue';

const { createRoom } = useLobby();
const { joinRoom } = useRoomSocket();

const inputtedRoomName = ref('');
const handleClick = async () => {
  const createdRoom = await createRoom({ name: inputtedRoomName.value });
  if (createdRoom) {
    joinRoom(createdRoom.id);
    router.push('/game');
  }
};
</script>

<template>
  <section class="space-y-5 p-8">
    <h2 class="text-3xl">新たに作成</h2>
    <form class="flex gap-5 pr-8">
      <input
        v-model="inputtedRoomName"
        class="outline-none text-2xl p-4 border-4 border-primary flex-1"
        type="text"
        placeholder="あなたの部屋の名前"
      />
      <BaseButton variant="primary" @btn-click="handleClick">作成</BaseButton>
    </form>
  </section>
</template>
