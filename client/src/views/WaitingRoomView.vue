<script setup lang="ts">
import { storeToRefs } from 'pinia';

import BaseButton from '@/components/common/BaseButton.vue';
import TheHeader from '@/components/common/TheHeader.vue';
import UserIcon from '@/components/common/UserIcon.vue';
import { useReturnHomeByDisconnect } from '@/composables/useReturnHomeByDisconnect';
import { useRoomSocket } from '@/composables/useRoomSocket';
import { useRoomStore } from '@/stores/room';
import { toKanjiNumber } from '@/utils/to-kanji-number';

const roomStore = useRoomStore();
const { currentRoom } = storeToRefs(roomStore);
const { sendReady } = useRoomSocket();
const { returnHomeByDisconnect } = useReturnHomeByDisconnect();

const handleClick = () => {
  sendReady();
};

const handleAbort = async () => {
  if (!window.confirm('ゲームを中止してトップへ戻りますか？')) {
    return;
  }

  await returnHomeByDisconnect();
};
</script>

<template>
  <div v-if="currentRoom" class="min-h-screen bg-background">
    <TheHeader />
    <div>
      <div class="flex flex-row gap-24 items-baseline mt-20 ml-28">
        <p class="text-6xl text-primary">{{ currentRoom.name }}</p>
        <p class="text-4xl text-primary">{{ toKanjiNumber(currentRoom.members.length) }}人部屋</p>
      </div>
      <div class="grid grid-cols-4 gap-x-28 gap-y-12 mt-16 ml-32 w-fit">
        <div
          v-for="member in currentRoom.members"
          :key="member.id"
          class="flex flex-col items-center w-fit gap-4"
          :class="{ 'opacity-30': !member.isReady }"
        >
          <UserIcon :tra-qid="member.id" size="large" />
          <p class="text-3xl text-center text-primary">{{ member.id }}</p>
        </div>
      </div>
    </div>
    <div class="fixed bottom-16 right-24 flex flex-col items-end gap-4">
      <BaseButton variant="secondary" @btn-click="handleAbort">中止して戻る</BaseButton>
      <BaseButton
        v-if="!roomStore.isReady"
        variant="primary"
        :disabled="!roomStore.canSendReady"
        @btn-click="handleClick"
        >準備完了</BaseButton
      >
      <p v-else class="text-3xl text-primary">準備完了</p>
    </div>
  </div>
</template>
