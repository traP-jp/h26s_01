<script lang="ts" setup>
import type { Room } from '@/types/api';
import { toKanjiNumber } from '@/utils/to-kanji-number';

import BaseButton from '../common/BaseButton.vue';
import UserIcon from '../common/UserIcon.vue';

const props = withDefaults(
  defineProps<{
    room: Room;
    disabled?: boolean;
  }>(),
  {
    disabled: false,
  },
);

const emit = defineEmits<{
  join: [roomId: string];
}>();

const handleClick = () => {
  emit('join', props.room.id);
};
</script>

<template>
  <li class="flex items-center justify-between text-2xl px-8 border-l-4 gap-5">
    <p>{{ props.room.name }}</p>
    <div class="grid grid-cols-3 items-center gap-3 w-xl">
      <div class="flex justify-center gap-1">
        <UserIcon
          v-for="member in props.room.members"
          :key="member.id"
          :tra-qid="member.id"
          size="small"
        />
      </div>
      <p class="text-center">{{ toKanjiNumber(props.room.members.length) }}人が参加中</p>
      <BaseButton
        v-if="props.room.status === 'waiting'"
        variant="primary"
        :disabled="props.disabled"
        @btn-click="handleClick"
        >参加</BaseButton
      >
      <BaseButton v-else-if="props.room.status === 'playing'" variant="secondary" disabled
        >遊戯中</BaseButton
      >
    </div>
  </li>
</template>
