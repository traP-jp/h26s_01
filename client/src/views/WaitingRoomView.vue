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
    <main class="mx-auto flex w-full max-w-screen-xl flex-col gap-16 px-8 py-20">
      <div class="flex flex-wrap items-baseline gap-x-24 gap-y-4">
        <p class="min-w-0 max-w-full truncate text-6xl text-primary" :title="currentRoom.name">
          {{ currentRoom.name }}
        </p>
        <p class="text-4xl text-primary">{{ toKanjiNumber(currentRoom.members.length) }}人部屋</p>
      </div>
      <div class="grid grid-cols-[repeat(auto-fit,minmax(9rem,10rem))] gap-x-12 gap-y-12">
        <div
          v-for="member in currentRoom.members"
          :key="member.id"
          class="flex min-w-0 flex-col items-center gap-4"
          :class="{ 'opacity-30': !member.isReady }"
        >
          <UserIcon :tra-qid="member.id" size="large" />
          <p class="w-full truncate text-center text-3xl text-primary" :title="member.id">
            {{ member.id }}
          </p>
        </div>
      </div>
    </main>
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
