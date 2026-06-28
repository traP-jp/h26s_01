<script lang="ts" setup>
import { useRoomSocket } from '@/composables/useRoomSocket.ts';
import router from '@/router/index.ts';
import type { Room } from '@/types/api.ts';
import { toKanjiNumber } from '@/utils/to-kanji-number.ts';

import BaseButton from '../common/BaseButton.vue';
import UserIcon from '../common/UserIcon.vue';

const { joinRoom } = useRoomSocket();

const { room } = defineProps<{
  room: Room;
}>();

const handleClick = async () => {
  joinRoom(room.id);
  router.push('/game');
};
</script>

<template>
  <li class="flex items-center justify-between text-2xl px-8 border-l-4 gap-5">
    <p>{{ room.name }}</p>
    <div class="grid grid-cols-3 items-center gap-3 w-xl">
      <div class="flex justify-center gap-1">
        <UserIcon v-for="member in room.members" :tra-qid="member.id" size="medium" />
      </div>
      <p class="text-center">{{ toKanjiNumber(room.members.length) }}人が参加中</p>
      <BaseButton v-if="room.status == 'waiting'" variant="primary" v-on:btn-click="handleClick"
        >参加</BaseButton
      >
      <BaseButton v-if="room.status == 'playing'" variant="secondary" disabled>遊戯中</BaseButton>
    </div>
  </li>
</template>
