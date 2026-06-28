<script lang="ts" setup>
import type { Room } from '@/types/api';

import RoomListItem from './RoomListItem.vue';

const props = defineProps<{
  rooms: Room[];
  isJoining: boolean;
  isLoading: boolean;
  error: string | null;
}>();

const emit = defineEmits<{
  join: [roomId: string];
}>();
</script>

<template>
  <section class="space-y-5 p-8">
    <h2 class="text-3xl">部屋の一覧</h2>
    <p v-if="props.isLoading" class="text-xl text-primary">読み込み中</p>
    <p v-else-if="props.error !== null" class="text-xl text-primary">{{ props.error }}</p>
    <p v-else-if="props.rooms.length === 0" class="text-xl text-primary">
      参加できる部屋がありません
    </p>
    <ul v-else class="flex flex-col gap-2">
      <RoomListItem
        v-for="room in props.rooms"
        :key="room.id"
        :disabled="props.isJoining"
        :room="room"
        @join="emit('join', $event)"
      />
    </ul>
  </section>
</template>
